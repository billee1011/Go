package desks

import (
	"context"
	"errors"
	"runtime/debug"
	msgid "steve/client_pb/msgId"
	"steve/client_pb/room"
	"steve/gutils"
	majong_initial "steve/majong/export/initial"
	majong_process "steve/majong/export/process"
	"steve/room/interfaces"
	"steve/room/interfaces/facade"
	"steve/room/interfaces/global"
	server_pb "steve/server_pb/majong"
	"steve/structs/proto/gate_rpc"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"

	_ "steve/room/ai" // 加载 AI 包
)

var errInitMajongContext = errors.New("初始化麻将现场失败")
var errAllocDeskIDFailed = errors.New("分配牌桌 ID 失败")
var errPlayerNotExist = errors.New("玩家不存在")

// deskEvent 房间事件
type deskEvent struct {
	event       interfaces.Event
	stateNumber int
}

// enterQuitInfo 退出以及进入信息
type enterQuitInfo struct {
	playerID uint64
	quit     bool // true 为退出， false 为进入
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
	deskUID      uint64                       // 牌桌唯一 ID
	gameID       int                          // 游戏 ID
	createOption interfaces.CreateDeskOptions // 创建选项
	dContext     *deskContext                 // 牌桌现场
	settler      interfaces.DeskSettler       // 结算器
	players      map[uint32]*deskPlayer       // Seat -> player
	event        chan deskEvent               // 牌桌事件通道
	enterQuits   chan enterQuitInfo           // 退出以及进入信息
	cancel       context.CancelFunc           // 取消事件处理
	tuoGuanMgr   interfaces.TuoGuanMgr        // 托管管理器
}

func makeDeskPlayers(logEntry *logrus.Entry, players []uint64) (map[uint32]*deskPlayer, error) {
	playerMgr := global.GetPlayerMgr()
	deskPlayers := make(map[uint32]*deskPlayer, 4)
	seat := uint32(0)
	for _, playerID := range players {
		player := playerMgr.GetPlayer(playerID)
		if player == nil {
			logEntry.WithField("player_id", playerID).Errorln(errPlayerNotExist)
			return nil, errPlayerNotExist
		}
		deskPlayers[seat] = newDeskPlayer(playerID, seat)
		seat++
	}
	return deskPlayers, nil
}

func newDesk(players []uint64, gameID int, opt interfaces.CreateDeskOptions) (result interfaces.CreateDeskResult, err error) {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "newDesk",
		"game_id":   gameID,
		"players":   players,
	})
	alloc := global.GetDeskIDAllocator()
	id, err := alloc.AllocDeskID()
	if err != nil {
		logEntry.Errorln(errAllocDeskIDFailed)
		err = errAllocDeskIDFailed
		return
	}
	logEntry = logEntry.WithField("desk_uid", id)
	deskPlayers, err := makeDeskPlayers(logEntry, players)
	if err != nil {
		return
	}
	return interfaces.CreateDeskResult{
		Desk: &desk{
			deskUID:      id,
			gameID:       gameID,
			createOption: opt,
			settler:      global.GetDeskSettleFactory().CreateDeskSettler(gameID),
			players:      deskPlayers,
			event:        make(chan deskEvent, 16),
			tuoGuanMgr:   newTuoGuanMgr(),
			enterQuits:   make(chan enterQuitInfo),
		},
	}, nil
}

// GetUID 获取牌桌 UID
func (d *desk) GetUID() uint64 {
	return d.deskUID
}

// GetGameID 获取游戏 ID
func (d *desk) GetGameID() int {
	return d.gameID
}

// GetPlayers 获取牌桌玩家数据
func (d *desk) GetPlayers() []*room.RoomPlayerInfo {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "desk.GetPlayers",
		"desk_uid":  d.deskUID,
		"game_id":   d.gameID,
	})
	result := []*room.RoomPlayerInfo{}

	playerMgr := global.GetPlayerMgr()

	for seat, deskPlayer := range d.players {
		player := playerMgr.GetPlayer(deskPlayer.playerID)
		if player == nil {
			logEntry.WithField("player_id", deskPlayer.playerID).Errorln("玩家不存在")
			return nil
		}
		result = append(result, &room.RoomPlayerInfo{
			PlayerId: proto.Uint64(deskPlayer.GetPlayerID()),
			Coin:     proto.Uint64(player.GetCoin()),
			Seat:     proto.Uint32(uint32(seat)),
		})
	}
	return result
}

// GetDeskPlayers 获取牌桌玩家
func (d *desk) GetDeskPlayers() []interfaces.DeskPlayer {
	result := []interfaces.DeskPlayer{}
	for _, deskPlayer := range d.players {
		result = append(result, deskPlayer)
	}
	return result
}

// Start 启动牌桌逻辑
// finish : 当牌桌逻辑完成时调用
// step 1. 初始化牌桌现场
// step 2. 启动发送事件的 goroutine
// step 3. 写入开始游戏事件
func (d *desk) Start(finish func()) error {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "desk.Start",
		"desk_uid":  d.deskUID,
		"game_id":   d.gameID,
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
			ID:        server_pb.EventID_event_start_game,
			Context:   []byte{},
			EventType: interfaces.NormalEvent,
		},
		stateNumber: d.dContext.stateNumber,
	}
	return nil
}

// PlayerQuit 玩家退出
func (d *desk) PlayerQuit(playerID uint64) {
	d.enterQuits <- enterQuitInfo{
		playerID: playerID,
		quit:     true,
	}
}

// PlayerEnter 玩家进入
func (d *desk) PlayerEnter(playerID uint64) {
	d.enterQuits <- enterQuitInfo{
		playerID: playerID,
		quit:     false,
	}
}

// Stop 停止桌面
// step1，桌面解散开始
// step2，广播桌面解散通知
func (d *desk) Stop() error {
	d.cancel()

	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "desk.Stop",
		"desk_uid":  d.deskUID,
		"game_id":   d.gameID,
	})
	players := d.GetDeskPlayers()
	clientIDs := []uint64{}

	playerMgr := global.GetPlayerMgr()
	for _, player := range players {
		playerID := player.GetPlayerID()
		p := playerMgr.GetPlayer(playerID)
		if p != nil {
			clientIDs = append(clientIDs, p.GetClientID())
		}
	}

	ntf := room.RoomDeskDismissNtf{}
	head := &steve_proto_gaterpc.Header{
		MsgId: uint32(msgid.MsgID_ROOM_DESK_DISMISS_NTF)}
	ms := global.GetMessageSender()
	err := ms.BroadcastPackage(clientIDs, head, &ntf)
	if err != nil {
		logEntry.WithError(err)
	}

	return err
}

// PushRequest 压入玩家请求
func (d *desk) PushRequest(playerID uint64, head *steve_proto_gaterpc.Header, bodyData []byte) {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name":  "desk.PushRequest",
		"desk_uid":   d.deskUID,
		"game_id":    d.gameID,
		"player_id":  playerID,
		"message_id": head.GetMsgId(),
	})

	trans := global.GetReqEventTranslator()
	eventID, eventContext, err := trans.Translate(playerID, head, bodyData)
	if err != nil {
		logEntry.WithError(err).Errorln("消息转事件失败")
		return
	}
	eventConetxtByte, err := proto.Marshal(eventContext)
	if err != nil {
		logEntry.WithError(err).Errorln("序列化事件现场失败")
	}

	d.event <- deskEvent{
		event: interfaces.Event{
			ID:        eventID,
			Context:   eventConetxtByte,
			EventType: interfaces.NormalEvent,
			PlayerID:  playerID,
		},
		stateNumber: d.dContext.stateNumber,
	}
}

// GetTuoGuanMgr 获取托管管理器
func (d *desk) GetTuoGuanMgr() interfaces.TuoGuanMgr {
	return d.tuoGuanMgr
}

func (d *desk) initMajongContext() error {
	players := make([]uint64, len(d.players))

	for seat, player := range d.players {
		players[seat] = player.playerID
	}

	param := server_pb.InitMajongContextParams{
		GameId:  int32(d.gameID),
		Players: players,
		Option: &server_pb.MajongCommonOption{
			MaxFapaiCartoonTime:        10 * 1000,
			MaxHuansanzhangCartoonTime: 10 * 1000,
		},
		MajongOption: []byte{},
	}
	var mjContext server_pb.MajongContext
	var err error
	if mjContext, err = majong_initial.InitMajongContext(param); err != nil {
		return err
	}
	d.dContext = &deskContext{
		mjContext:   mjContext,
		stateNumber: 0,
		stateTime:   time.Now(),
	}
	return nil
}

func (d *desk) getTuoguanPlayers() []uint64 {
	return d.tuoGuanMgr.GetTuoGuanPlayers()
}

// genTimerEvent 生成计时事件
func (d *desk) genTimerEvent() {
	g := global.GetDeskAutoEventGenerator()
	// 先将 context 指针读出来拷贝， 后面的 context 修改都会分配一块新的内存
	dContext := d.dContext
	tuoGuanPlayers := d.getTuoguanPlayers()
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name":       "desk.genTimerEvent",
		"state_number":    dContext.stateNumber,
		"tuoguan_players": tuoGuanPlayers,
	})
	result := g.GenerateV2(&interfaces.AutoEventGenerateParams{
		MajongContext:  &dContext.mjContext,
		CurTime:        time.Now(),
		StateTime:      dContext.stateTime,
		RobotLv:        map[uint64]int{},
		TuoGuanPlayers: tuoGuanPlayers,
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
	if event.event.ID == server_pb.EventID_event_huansanzhang_request ||
		event.event.ID == server_pb.EventID_event_dingque_request {
		return false
	}
	return true
}

func (d *desk) processEvents(ctx context.Context) {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "desk.processEvent",
		"desk_uid":  d.deskUID,
		"game_id":   d.gameID,
	})
	defer func() {
		if x := recover(); x != nil {
			logEntry.Errorln(x)
			debug.PrintStack()
		}
	}()

	for {
		select {
		case <-ctx.Done():
			{
				logEntry.Infoln("done")
				return
			}
		case enterQuitInfo := <-d.enterQuits:
			{
				d.handleEnterQuit(enterQuitInfo)
			}
		case event := <-d.event:
			{
				if d.needCompareStateNumber(&event) && event.stateNumber != d.dContext.stateNumber {
					continue
				}
				d.processEvent(event.event.ID, event.event.Context)
				if event.event.EventType == interfaces.OverTimeEvent && event.event.PlayerID != 0 {
					d.tuoGuanMgr.OnPlayerTimeOut(event.event.PlayerID)
				}
			}
		}
	}
}

// getDeskPlayer 获取牌桌玩家
func (d *desk) getDeskPlayer(playerID uint64) *deskPlayer {
	for _, deskPlayer := range d.players {
		if deskPlayer.GetPlayerID() == playerID {
			return deskPlayer
		}
	}
	return nil
}

// handleEnterQuit 处理退出进入信息
func (d *desk) handleEnterQuit(eqi enterQuitInfo) {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "handleEnterQuit",
		"player_id": eqi.playerID,
		"quit":      eqi.quit,
	})
	var msgs []server_pb.ReplyClientMessage
	deskPlayer := d.getDeskPlayer(eqi.playerID)
	if deskPlayer == nil {
		logEntry.Errorln("玩家不在牌桌上")
		return
	}
	if eqi.quit {
		deskPlayer.quitDesk()
		d.tuoGuanMgr.SetTuoGuan(eqi.playerID, true, false) // 退出后自动托管
		logEntry.Debugln("玩家退出")
	} else {
		deskPlayer.enterDesk()
		msgs = d.recoverGameForPlayer(eqi.playerID)
		d.tuoGuanMgr.SetTuoGuan(eqi.playerID, false, false) // 进入后取消托管
		logEntry.Debugln("玩家进入")
	}
	// 发送消息给玩家
	d.reply(msgs)
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

// pushAutoEvent 一段时间后压入自动事件
func (d *desk) pushAutoEvent(autoEvent *server_pb.AutoEvent, stateNumber int) {
	time.Sleep(time.Millisecond * time.Duration(autoEvent.GetWaitTime()))
	if d.dContext.stateNumber != stateNumber {
		return
	}
	d.event <- deskEvent{
		event: interfaces.Event{
			ID:        autoEvent.EventId,
			Context:   autoEvent.EventContext,
			EventType: interfaces.NormalEvent,
			PlayerID:  0,
		},
		stateNumber: stateNumber,
	}
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

// checkGameOver 检查游戏结束
func (d *desk) checkGameOver(logEntry *logrus.Entry) bool {
	mjContext := d.dContext.mjContext
	// 游戏结束
	if mjContext.GetCurState() == server_pb.StateID_state_gameover {
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

// removeQuit 移除已经退出的玩家
func (d *desk) removeQuit(playerIDs []uint64) []uint64 {
	deskPlayerIDs := map[uint64]bool{}
	deskPlayers := d.GetDeskPlayers()
	for _, deskPlayer := range deskPlayers {
		playerID := deskPlayer.GetPlayerID()
		deskPlayerIDs[playerID] = deskPlayer.IsQuit()
	}
	result := []uint64{}
	for _, playerID := range playerIDs {
		if quited, _ := deskPlayerIDs[playerID]; !quited {
			result = append(result, playerID)
		}
	}
	return result
}

// allPlayerIDs 获取所有玩家的 ID
func (d *desk) allPlayerIDs() []uint64 {
	result := []uint64{}
	deskPlayers := d.GetDeskPlayers()
	for _, deskPlayer := range deskPlayers {
		playerID := deskPlayer.GetPlayerID()
		result = append(result, playerID)
	}
	return result
}

// BroadcastMessage 向玩家广播消息
func (d *desk) BroadcastMessage(playerIDs []uint64, msgID msgid.MsgID, body []byte, exceptQuit bool) {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name":       "BroadcastMessage",
		"dest_player_ids": playerIDs,
		"msg_id":          msgID,
	})
	// 是否针对所有玩家
	if playerIDs == nil || len(playerIDs) == 0 {
		playerIDs = d.allPlayerIDs()
		logEntry = logEntry.WithField("all_player_ids", playerIDs)
	}
	playerIDs = d.removeQuit(playerIDs)
	logEntry = logEntry.WithField("real_dest_player_ids", playerIDs)

	if len(playerIDs) == 0 {
		return
	}
	facade.BroadCastMessageBare(playerIDs, msgID, body)
	logEntry.Debugln("广播消息")
}

func (d *desk) recoverGameForPlayer(playerID uint64) []server_pb.ReplyClientMessage {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "recoverGameForPlayer",
		"playerID":  playerID,
	})

	mjContext := d.dContext.mjContext
	bankerSeat := mjContext.GetZhuangjiaIndex()
	totalCardsNum := uint32(108) //global.GetOriginCards(mjContext.GetGameId()),该函数在麻将里面，room调用不到
	remainCardsNum := uint32(len(mjContext.GetWallCards()))

	// 获取本状态已经停留的时间
	var costTime uint32
	nowTime := time.Now().Unix()
	stateEntryTime := d.dContext.stateTime.Unix()
	if nowTime > stateEntryTime {
		costTime = uint32(nowTime - stateEntryTime)
	}

	gameStage := d.getGameStage()

	gameDeskInfo := room.GameDeskInfo{
		GameStage:   &gameStage,               // 需要根据状态确定
		Players:     d.getRecoverPlayerInfo(), // 需要构造
		Dices:       mjContext.GetDices(),
		BankerSeat:  &bankerSeat,
		EastSeat:    &bankerSeat,
		TotalCards:  &totalCardsNum,
		RemainCards: &remainCardsNum,
		CostTime:    &costTime,
		OperatePid:  d.getOperatePlayerID(),
		DoorCard:    d.getDoorCard(),
		// ZixunInfo        *RoomZixunNtf
		// WenxunInfo       *RoomChupaiWenxunNtf
		// QghInfo          *RoomWaitQianggangHuNtf
	}
	rsp := room.RoomResumeGameRsp{
		ResumeRes: room.RoomError_SUCCESS.Enum(),
		GameInfo:  &gameDeskInfo,
	}
	body, err := proto.Marshal(&rsp)
	if err != nil {
		logEntry.WithError(err).Errorln("序列化失败")
		return nil
	}
	msg := server_pb.ReplyClientMessage{
		Players: []uint64{playerID},
		MsgId:   int32(msgid.MsgID_ROOM_RESUME_GAME_RSP),
		Msg:     body,
	}
	return []server_pb.ReplyClientMessage{msg}
}

func (d *desk) getOperatePlayerID() *uint64 {
	mjContext := d.dContext.mjContext
	state := mjContext.GetCurState()
	var playerID uint64
	// TODO 状态的定义在mjpb
	switch state {
	case 1: // 自询状态 刚开局是庄家，其他情况是最近摸牌者
		playerID = mjContext.GetLastMopaiPlayer()
	case 2: // 他询状态， 出牌者
		playerID = mjContext.GetLastChupaiPlayer()
	case 3: // 等待抢杠胡 杠牌玩家
		playerID = mjContext.GetLastGangPlayer()
	default:
		return nil
	}
	return &playerID
}

func (d *desk) getGameStage() (stage room.GameStage) {
	switch d.dContext.mjContext.GetCurState() {
	case server_pb.StateID_state_huansanzhang:
		stage = room.GameStage_GAMESTAGE_HUANSANZHANG
	case server_pb.StateID_state_dingque:
		stage = room.GameStage_GAMESTAGE_DINGQUE
	default:
		stage = room.GameStage_GAMESTAGE_PLAYCARD
	}
	return
}

func (d *desk) getDoorCard() *uint32 {
	mjContext := &d.dContext.mjContext
	if mjContext.GetCurState() == server_pb.StateID_state_zixun {
		DoorCard := uint32(mjContext.GetLastMopaiCard().GetPoint())
		return &DoorCard
	}
	return nil
}

func (d *desk) getRecoverPlayerInfo() (recoverPlayerInfo []*room.GamePlayerInfo) {
	mjContext := &d.dContext.mjContext
	roomPlayerInfos := d.GetPlayers()
	for _, roomPlayerInfo := range roomPlayerInfos {
		var player *server_pb.Player
		for _, player = range mjContext.GetPlayers() {
			if player.GetPalyerId() == roomPlayerInfo.GetPlayerId() {
				break
			}
		}
		playerID := player.GetPalyerId()
		svrHandCard := player.GetHandCards()
		handCardCount := uint32(len(svrHandCard))
		gamePlayerInfo := &room.GamePlayerInfo{
			PlayerInfo:    roomPlayerInfo,
			Color:         gutils.ServerColor2ClientColor(player.DingqueColor).Enum(),
			HandCardCount: &handCardCount,
		}

		// 手牌组
		cltHandCard := gutils.ServerCards2Numbers(svrHandCard)
		handCardGroup := &room.CardsGroup{
			Cards: cltHandCard,
			Type:  room.CardsGroupType_CGT_HAND.Enum(),
		}
		gamePlayerInfo.CardsGroup = append(gamePlayerInfo.CardsGroup, handCardGroup)
		// 吃牌组
		//chiCardGroup

		// 碰牌组,每一次碰牌填1张还是三张
		var pengCardGroups []*room.CardsGroup
		for _, pengCard := range player.GetPengCards() {
			srcPlayerID := pengCard.GetSrcPlayer()
			pengCardGroup := &room.CardsGroup{
				Cards: []uint32{gutils.ServerCard2Number(pengCard.GetCard())},
				Type:  room.CardsGroupType_CGT_PENG.Enum(),
				Pid:   &srcPlayerID,
			}
			pengCardGroups = append(pengCardGroups, pengCardGroup)
		}
		gamePlayerInfo.CardsGroup = append(gamePlayerInfo.CardsGroup, pengCardGroups...)
		// 杠牌组
		var gangCardGroups []*room.CardsGroup
		for _, gangCard := range player.GetGangCards() {
			groupType := gutils.SvrGangType2CardGroupType(gangCard.GetType())
			srcPlayerID := gangCard.GetSrcPlayer()
			gangCardGroup := &room.CardsGroup{
				Cards: []uint32{gutils.ServerCard2Number(gangCard.GetCard())},
				Type:  &groupType,
				Pid:   &srcPlayerID,
			}
			gangCardGroups = append(gangCardGroups, gangCardGroup)
		}
		gamePlayerInfo.CardsGroup = append(gamePlayerInfo.CardsGroup, gangCardGroups...)
		// 胡牌组
		var huCardGroups []*room.CardsGroup
		for _, huCard := range player.GetHuCards() {
			srcPlayerID := huCard.GetSrcPlayer()
			huCardGroup := &room.CardsGroup{
				Cards: []uint32{gutils.ServerCard2Number(huCard.GetCard())},
				Type:  room.CardsGroupType_CGT_HU.Enum(),
				Pid:   &srcPlayerID,
			}
			huCardGroups = append(huCardGroups, huCardGroup)
		}
		gamePlayerInfo.CardsGroup = append(gamePlayerInfo.CardsGroup, huCardGroups...)
		// 花牌组

		// 出牌组
		outCardGroup := &room.CardsGroup{
			Cards: gutils.ServerCards2Numbers(player.GetOutCards()),
			Type:  room.CardsGroupType_CGT_OUT.Enum(),
			Pid:   &playerID,
		}
		gamePlayerInfo.CardsGroup = append(gamePlayerInfo.CardsGroup, outCardGroup)
	}
	return
}

// // 麻将组类型
// enum CardsGroupType{
//     CGT_HAND = 0;     // 手牌
//     CGT_CHI = 1;      // 吃牌
//     CGT_PENG = 2;     // 碰牌
//     CGT_MINGGANG = 3; // 明杠牌
//     CGT_ANGANG = 4;   // 暗杠牌
//     CGT_BUGANG = 5;   // 补杠牌
//     CGT_HU = 6;       // 胡的牌
//     CGT_HUA = 7;      // 花牌
// 	CGT_OUT = 8;      // 出去的牌
// }

// // 麻将组
// message CardsGroup {
//     repeated uint32 cards = 1;           // 牌
//     optional CardsGroupType type = 2;    // 麻将组类型
//     optional uint64 pid = 3;             // 吃碰杠谁的牌
// }
