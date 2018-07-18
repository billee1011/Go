package mjdesk

import (
	"context"
	"errors"
	"runtime/debug"
	"steve/client_pb/msgid"
	"steve/client_pb/room"
	"steve/gutils"
	majong_initial "steve/majong/export/initial"
	majong_process "steve/majong/export/process"
	"steve/room/config"
	"steve/room/desks/deskbase"
	"steve/room/interfaces"
	"steve/room/interfaces/facade"
	"steve/room/interfaces/global"
	"steve/room/peipai/handle"
	server_pb "steve/server_pb/majong"
	"steve/structs/proto/gate_rpc"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
	"github.com/spf13/viper"

	_ "steve/room/ai" // 加载 AI 包
)

var errInitMajongContext = errors.New("初始化麻将现场失败")
var errAllocDeskIDFailed = errors.New("分配牌桌 ID 失败")
var errPlayerNotExist = errors.New("玩家不存在")
var errPlayerNeedXingPai = errors.New("玩家需要参与行牌")

// deskEvent 房间事件
type deskEvent struct {
	event       interfaces.Event
	stateNumber int
}

// deskContext 牌桌现场
type deskContext struct {
	mjContext   server_pb.MajongContext // 牌局现场
	stateNumber int                     // 状态序号
	stateTime   time.Time               // 状态时间
}

type autoEvent struct {
	aevent      *server_pb.AutoEvent // 自动事件
	stateNumber int                  // 对应的状态序号
	createTime  time.Time            // 创建时间
}

type desk struct {
	*deskbase.DeskBase
	dContext     *deskContext                 // 牌桌现场
	settler      interfaces.DeskSettler       // 结算器
	event        chan deskEvent               // 牌桌事件通道
	cancel       context.CancelFunc           // 取消事件处理
	createOption interfaces.CreateDeskOptions // 创建选项
}

// Start 启动牌桌逻辑
// finish : 当牌桌逻辑完成时调用
// step 1. 初始化牌桌现场
// step 2. 启动发送事件的 goroutine
// step 3. 写入开始游戏事件
func (d *desk) Start(finish func()) error {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "desk.Start",
		"desk_uid":  d.GetUID(),
		"game_id":   d.GetGameID(),
	})

	if err := d.initMajongContext(); err != nil {
		logEntry.WithError(err).Errorln(errInitMajongContext)
		return errInitMajongContext
	}
	var ctx context.Context
	ctx, d.cancel = context.WithCancel(context.Background())

	go func() {
		d.processEvents(ctx)
		logEntry.Infoln("处理事件完成")
		finish()
	}()
	go func() {
		d.timerTask(ctx)
	}()

	d.event <- deskEvent{
		event: interfaces.Event{
			ID:        int32(server_pb.EventID_event_start_game),
			Context:   []byte{},
			EventType: interfaces.NormalEvent,
		},
		stateNumber: d.dContext.stateNumber,
	}
	return nil
}

// Stop 停止桌面
// step1，桌面解散开始
// step2，广播桌面解散通知
func (d *desk) Stop() error {
	d.cancel()
	ntf := room.RoomDeskDismissNtf{}
	facade.BroadCastDeskMessage(d, nil, msgid.MsgID_ROOM_DESK_DISMISS_NTF, &ntf, true)
	return nil
}

// PushRequest 压入玩家请求
func (d *desk) PushRequest(playerID uint64, head *steve_proto_gaterpc.Header, bodyData []byte) {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name":  "desk.PushRequest",
		"desk_uid":   d.GetUID(),
		"game_id":    d.GetGameID(),
		"player_id":  playerID,
		"message_id": head.GetMsgId(),
	})

	trans := global.GetReqEventTranslator()
	eventID, eventContext, err := trans.Translate(playerID, head, bodyData)
	if err != nil {
		logEntry.WithError(err).Errorln("消息转事件失败")
		return
	}
	eventMessage, ok := eventContext.(proto.Message)
	if !ok {
		logEntry.Errorln("转换事件函数返回值类型错误")
		return
	}
	eventConetxtByte, err := proto.Marshal(eventMessage)
	if err != nil {
		logEntry.WithError(err).Errorln("序列化事件现场失败")
	}

	d.PushEvent(interfaces.Event{
		ID:        int32(eventID),
		Context:   eventConetxtByte,
		EventType: interfaces.NormalEvent,
		PlayerID:  playerID,
	})
}

func (d *desk) selectZhuang() uint32 {
	if !d.createOption.FixBankerSeat {
		return 0
	}
	return uint32(d.createOption.BankerSeat)
}

func (d *desk) initMajongContext() error {
	players := facade.GetDeskPlayerIDs(d)
	param := server_pb.InitMajongContextParams{
		GameId:  int32(d.GetGameID()),
		Players: players,
		Option: &server_pb.MajongCommonOption{
			MaxFapaiCartoonTime:        uint32(viper.GetInt(config.MaxFapaiCartoonTime)),
			MaxHuansanzhangCartoonTime: uint32(viper.GetInt(config.MaxHuansanzhangCartoonTime)),
			HasHuansanzhang:            handle.GetHsz(d.GetGameID()),                     //设置玩家是否开启换三张
			Cards:                      handle.GetPeiPai(d.GetGameID()),                  //设置是否配置墙牌
			WallcardsLength:            uint32(handle.GetLensOfWallCards(d.GetGameID())), //设置墙牌长度
			HszFx: &server_pb.Huansanzhangfx{
				NeedDeployFx:   handle.GetHSZFangXiang(d.GetGameID()) != -1,
				HuansanzhangFx: int32(handle.GetHSZFangXiang(d.GetGameID())),
			}, //设置换三张方向
		}, //设置庄家
		ZhuangIndex:  d.selectZhuang(), // TODO
		MajongOption: []byte{},
	}
	var mjContext server_pb.MajongContext
	var err error
	if mjContext, err = majong_initial.InitMajongContext(param); err != nil {
		return err
	}
	if err := fillContextOptions(d.GetGameID(), &mjContext); err != nil {
		return err
	}
	d.dContext = &deskContext{
		mjContext:   mjContext,
		stateNumber: 0,
		stateTime:   time.Now(),
	}
	return nil
}

// genTimerEvent 生成计时事件
func (d *desk) genTimerEvent() {
	g := global.GetDeskAutoEventGenerator()
	// 先将 context 指针读出来拷贝， 后面的 context 修改都会分配一块新的内存
	dContext := d.dContext
	tuoGuanPlayers := facade.GetTuoguanPlayers(d)
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name":       "desk.genTimerEvent",
		"state_number":    dContext.stateNumber,
		"tuoguan_players": tuoGuanPlayers,
	})
	deskPlayers := d.GetDeskPlayers()
	robotLvs := make(map[uint64]int, len(deskPlayers))
	for _, deskPlayer := range deskPlayers {
		robotLv := deskPlayer.GetRobotLv()
		if robotLv != 0 {
			robotLvs[deskPlayer.GetPlayerID()] = robotLv
		}
	}
	result := g.GenerateV2(&interfaces.AutoEventGenerateParams{
		Desk:          d,
		MajongContext: &dContext.mjContext,
		StartTime:     dContext.stateTime,
		RobotLv:        robotLvs,
	})
	for _, event := range result.Events {
		logEntry.WithFields(logrus.Fields{
			"event_id":     event.ID,
			"event_player": event.PlayerID,
			"event_type":   event.EventType,
		}).Debugln("注入计时事件")
		d.event <- deskEvent{
			event:       event,
			stateNumber: dContext.stateNumber,
		}
	}
}

// timerTask 定时任务，产生自动事件
func (d *desk) timerTask(ctx context.Context) {
	defer func() {
		if x := recover(); x != nil {
			debug.PrintStack()
		}
	}()

	t := time.NewTicker(time.Millisecond * 200)
	defer t.Stop()
	for {
		select {
		case <-t.C:
			{
				d.genTimerEvent()
			}
		case <-ctx.Done():
			{
				return
			}
		}
	}
}

// needCompareStateNumber 判断事件是否需要比较 stateNumber
func (d *desk) needCompareStateNumber(event *deskEvent) bool {
	if event.event.ID == int32(server_pb.EventID_event_huansanzhang_request) ||
		event.event.ID == int32(server_pb.EventID_event_dingque_request) {
		return false
	}
	return true
}

// recordTuoguanOverTimeCount 记录托管超时计数
func (d *desk) recordTuoguanOverTimeCount(event interfaces.Event) {
	if event.EventType != interfaces.OverTimeEvent {
		return
	}
	playerID := event.PlayerID
	if playerID == 0 {
		return
	}
	id := event.ID
	if id == int32(server_pb.EventID_event_huansanzhang_request) || id == int32(server_pb.EventID_event_dingque_request) {
		return
	}
	deskPlayer := facade.GetDeskPlayerByID(d, playerID)
	if deskPlayer != nil {
		deskPlayer.OnPlayerOverTime()
	}
}

func (d *desk) processEvents(ctx context.Context) {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "desk.processEvent",
		"desk_uid":  d.GetUID(),
		"game_id":   d.GetGameID(),
	})
	defer func() {
		if x := recover(); x != nil {
			logEntry.Errorln(x)
			debug.PrintStack()
		}
	}()
	defer d.consumeAllEnterQuit() // 消费完所有的退出进入数据

	for {
		select {
		case <-ctx.Done():
			{

				logEntry.Infoln("done")
				return
			}
		case enterQuitInfo := <-d.PlayerEnterQuitChannel():
			{
				d.handleEnterQuit(enterQuitInfo)
			}
		case event := <-d.event:
			{
				if d.needCompareStateNumber(&event) && event.stateNumber != d.dContext.stateNumber {
					continue
				}
				d.processEvent(server_pb.EventID(event.event.ID), event.event.Context)
				d.recordTuoguanOverTimeCount(event.event)
			}
		}
	}
}

func (d *desk) consumeAllEnterQuit() {
	for {
		select {
		case enterQuitInfo := <-d.PlayerEnterQuitChannel():
			{
				d.handleEnterQuit(enterQuitInfo)
			}
		default:
			return
		}
	}
}

// getContextPlayer 获取context玩家
func (d *desk) getContextPlayer(playerID uint64) *server_pb.Player {
	for _, contextPlayer := range d.dContext.mjContext.GetPlayers() {
		if contextPlayer.GetPalyerId() == playerID {
			return contextPlayer
		}
	}
	return nil
}

// handleEnterQuit 处理退出进入信息
func (d *desk) handleEnterQuit(eqi interfaces.PlayerEnterQuitInfo) {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "handleEnterQuit",
		"player_id": eqi.PlayerID,
		"quit":      eqi.Quit,
	})
	deskPlayer := facade.GetDeskPlayerByID(d, eqi.PlayerID)
	defer close(eqi.FinishChannel)

	if deskPlayer == nil {
		logEntry.Errorln("玩家不在牌桌上")
		return
	}
	if eqi.Quit {
		deskPlayer.QuitDesk()
		d.setMjPlayerQuitDesk(eqi.PlayerID, true)
		d.handleQuitByPlayerState(eqi.PlayerID)
		d.playerQuitEnterDeskNtf(eqi.PlayerID, room.QuitEnterType_QET_QUIT)
		logEntry.Debugln("玩家退出")
	} else {
		// 判断行牌状态, 选项化后需修改
		mjPlayer := gutils.GetMajongPlayer(eqi.PlayerID, &d.dContext.mjContext)
		// 非主动退出，再进入后取消托管；主动退出再进入不取消托管
		// 胡牌后没有托管，但是在客户端退出时，需要托管来自动胡牌,重新进入后把托管取消
		if !deskPlayer.IsQuit() || mjPlayer.GetXpState() != server_pb.XingPaiState_normal {
			deskPlayer.SetTuoguan(false, false)
		}
		deskPlayer.EnterDesk()
		d.recoverGameForPlayer(eqi.PlayerID)
		d.setMjPlayerQuitDesk(eqi.PlayerID, false)
		d.playerQuitEnterDeskNtf(eqi.PlayerID, room.QuitEnterType_QET_ENTER)
		logEntry.Debugln("玩家进入")
	}
}

func (d *desk) handleQuitByPlayerState(playerID uint64) {
	mjContext := d.dContext.mjContext
	player := gutils.GetMajongPlayer(playerID, &mjContext)

	if !gutils.IsPlayerContinue(player.GetXpState(), &mjContext) {
		deskMgr := global.GetDeskMgr()
		deskMgr.RemoveDeskPlayerByPlayerID(playerID)
	}
	logrus.WithFields(logrus.Fields{
		"funcName":    "handleQuitByPlayerState",
		"gameID":      mjContext.GetGameId(),
		"playerState": player.GetXpState(),
	}).Infof("玩家:%v退出后的相关处理", playerID)
}

// callEventHandler 调用事件处理器
func (d *desk) callEventHandler(logEntry *logrus.Entry, eventID server_pb.EventID, eventContext []byte) (result majong_process.HandleMajongEventResult, succ bool) {
	succ = false
	stateNumber, mjContext, stateTime := d.dContext.stateNumber, d.dContext.mjContext, d.dContext.stateTime
	oldState := mjContext.GetCurState()

	if result = majong_process.HandleMajongEvent(majong_process.HandleMajongEventParams{
		MajongContext: mjContext,
		EventID:       eventID,
		EventContext:  eventContext,
	}); !result.Succeed {
		logEntry.Debugln("处理事件不成功")
		return
	}
	newContext := result.NewContext
	newState := newContext.GetCurState()
	if newState != oldState {
		stateNumber++
		stateTime = time.Now()
	}
	// dContext 的每次修改都是一块新内存，用来确保并发安全。
	d.dContext = &deskContext{
		mjContext:   newContext,
		stateNumber: stateNumber,
		stateTime:   stateTime,
	}
	succ = true
	return
}

// PushEvent 压入事件
func (d *desk) PushEvent(event interfaces.Event) {
	d.event <- deskEvent{
		event:       event,
		stateNumber: d.dContext.stateNumber,
	}
}

// pushAutoEvent 一段时间后压入自动事件
func (d *desk) pushAutoEvent(autoEvent *server_pb.AutoEvent, stateNumber int) {
	time.Sleep(time.Millisecond * time.Duration(autoEvent.GetWaitTime()))
	if d.dContext.stateNumber != stateNumber {
		return
	}
	d.PushEvent(interfaces.Event{
		ID:        int32(autoEvent.EventId),
		Context:   autoEvent.EventContext,
		EventType: interfaces.NormalEvent,
		PlayerID:  0,
	})
}

// processEvent 处理单个事件
// step 1. 调用麻将逻辑的接口来处理事件(返回最新麻将现场, 自动事件， 发送给玩家的消息)， 并且更新 mjContext
// step 2. 将消息发送给玩家
// step 3. 调用 room 的结算逻辑来处理结算
// step 4. 如果有自动事件， 将自动事件写入自动事件通道
// step 5. 如果当前状态是游戏结束状态， 调用 cancel 终止游戏
func (d *desk) processEvent(eventID server_pb.EventID, eventContext []byte) {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "desk.ProcessEvent",
		"event_id":  eventID,
	})
	result, succ := d.callEventHandler(logEntry, eventID, eventContext)
	if !succ {
		return
	}

	// 发送消息给玩家
	d.reply(result.ReplyMsgs)
	d.settler.Settle(d, d.dContext.mjContext)

	if d.checkGameOver(logEntry) {
		return
	}
	// 自动事件不为空，继续处理事件
	if result.AutoEvent != nil {
		if result.AutoEvent.GetWaitTime() == 0 {
			d.processEvent(result.AutoEvent.GetEventId(), result.AutoEvent.GetEventContext())
		} else {
			go d.pushAutoEvent(result.AutoEvent, d.dContext.stateNumber)
		}
	}
}

func (d *desk) getWinners() []uint64 {
	statistics := d.settler.GetStatistics()
	winners := make([]uint64, 0, len(statistics))
	players := d.GetDeskPlayers()

	for _, player := range players {
		playerID := player.GetPlayerID()
		val, ok := statistics[playerID]
		if !ok || val >= 0 {
			winners = append(winners, playerID)
		}
	}
	return winners
}

// checkGameOver 检查游戏结束
func (d *desk) checkGameOver(logEntry *logrus.Entry) bool {
	mjContext := d.dContext.mjContext
	// 游戏结束
	if mjContext.GetCurState() == server_pb.StateID_state_gameover {
		d.ContinueDesk(true, int(mjContext.GetNextBankerSeat()), d.getWinners())
		d.settler.RoundSettle(d, mjContext)
		logEntry.Infoln("游戏结束状态")
		d.cancel()
		return true
	}
	return false
}

func (d *desk) reply(replyMsgs []server_pb.ReplyClientMessage) {
	if replyMsgs == nil {
		return
	}
	for _, msg := range replyMsgs {
		d.BroadcastMessage(msg.GetPlayers(), msgid.MsgID(msg.GetMsgId()), msg.GetMsg(), true)
	}
}

func (d *desk) recoverGameForPlayer(playerID uint64) {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "recoverGameForPlayer",
		"playerID":  playerID,
	})

	mjContext := &d.dContext.mjContext
	bankerSeat := mjContext.GetZhuangjiaIndex()
	totalCardsNum := mjContext.GetCardTotalNum()
	gameStage := getGameStage(mjContext.GetCurState())
	gameID := gutils.GameIDServer2Client(int(mjContext.GetGameId()))
	gameDeskInfo := room.GameDeskInfo{
		GameId:      &gameID,
		GameStage:   &gameStage,
		Players:     getRecoverPlayerInfo(playerID, d),
		Dices:       mjContext.GetDices(),
		BankerSeat:  &bankerSeat,
		EastSeat:    &bankerSeat,
		TotalCards:  &totalCardsNum,
		RemainCards: proto.Uint32(uint32(len(mjContext.GetWallCards()))),
		CostTime:    proto.Uint32(getStateCostTime(d.dContext.stateTime.Unix())),
		OperatePid:  getOperatePlayerID(mjContext),
		DoorCard:    getDoorCard(mjContext),
		NeedHsz:     proto.Bool(gutils.GameHasHszState(mjContext)),
	}
	gameDeskInfo.HasZixun, gameDeskInfo.ZixunInfo = getZixunInfo(playerID, mjContext)
	gameDeskInfo.HasWenxun, gameDeskInfo.WenxunInfo = getWenxunInfo(playerID, mjContext)
	gameDeskInfo.HasQgh, gameDeskInfo.QghInfo = getQghInfo(playerID, mjContext)
	rsp, err := proto.Marshal(&room.RoomResumeGameRsp{
		ResumeRes: room.RoomError_SUCCESS.Enum(),
		GameInfo:  &gameDeskInfo,
	})
	logEntry.Infoln("恢复数据")
	logEntry.Infoln(gameDeskInfo)
	if err != nil {
		logEntry.WithError(err).Errorln("序列化失败")
		return
	}
	d.reply([]server_pb.ReplyClientMessage{
		server_pb.ReplyClientMessage{
			Players: []uint64{playerID},
			MsgId:   int32(msgid.MsgID_ROOM_RESUME_GAME_RSP),
			Msg:     rsp,
		},
	})
}

func (d *desk) playerQuitEnterDeskNtf(playerID uint64, qeType room.QuitEnterType) {
	deskPlayer := facade.GetDeskPlayerByID(d, playerID)
	if deskPlayer == nil {
		return
	}
	roomPlayer := translateToRoomPlayer(deskPlayer)
	ntf := room.RoomDeskQuitEnterNtf{
		PlayerId:   &playerID,
		Type:       &qeType,
		PlayerInfo: &roomPlayer,
	}
	facade.BroadCastDeskMessageExcept(d, []uint64{playerID}, true, msgid.MsgID_ROOM_DESK_QUIT_ENTER_NTF, &ntf)
}

func (d *desk) setMjPlayerQuitDesk(playerID uint64, isQuit bool) {
	mjPlayer := d.getContextPlayer(playerID)
	if mjPlayer != nil {
		mjPlayer.IsQuit = isQuit
	}
}

// ChangePlayer 换对手
func (d *desk) ChangePlayer(playerID uint64) error {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "d.ChangePlayer",
		"playerID":  playerID,
	})
	mjContext := &d.dContext.mjContext
	player := gutils.GetMajongPlayer(playerID, mjContext)

	if gutils.IsPlayerContinue(player.GetXpState(), mjContext) {
		logEntry.WithFields(logrus.Fields{
			"XpState": player.GetXpState(),
		}).WithError(errPlayerNeedXingPai).Errorln("不能换对手")
		return errPlayerNeedXingPai
	}
	deskMgr := global.GetDeskMgr()
	deskPlayer := facade.GetDeskPlayerByID(d, playerID)
	d.playerQuitEnterDeskNtf(playerID, room.QuitEnterType_QET_QUIT)
	deskPlayer.QuitDesk()
	deskMgr.RemoveDeskPlayerByPlayerID(playerID)
	// getJoinApplyMgr().joinPlayer(playerID, room.GameId(mjContext.GetGameId()))
	return nil
}
