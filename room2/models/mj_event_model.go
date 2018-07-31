package models

import (
	"context"
	"runtime/debug"
	"steve/client_pb/msgid"
	"steve/client_pb/room"
	"steve/gutils"
	majong_process "steve/majong/export/process"
	"steve/room2/ai"
	context2 "steve/room2/contexts"
	"steve/room2/desk"
	"steve/room2/fixed"
	"steve/room2/player"
	server_pb "steve/server_pb/majong"
	"steve/structs/proto/gate_rpc"
	"sync"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
)

type MjEventModel struct {
	event chan desk.DeskEvent // 牌桌事件通道
	BaseModel
	closed bool
	mu     sync.Mutex
}

func NewMjEventModel(desk *desk.Desk) DeskModel {
	result := &MjEventModel{}
	result.SetDesk(desk)
	return result
}

func (model *MjEventModel) GetName() string {
	return fixed.Event
}

func (model *MjEventModel) Start() {
	model.event = make(chan desk.DeskEvent, 16)

	go func() {
		model.processEvents(context.Background())
		GetModelManager().StopDeskModel(model.GetDesk().GetUid())
	}()
	go func() {
		model.timerTask(context.Background())
	}()

	event := desk.NewDeskEvent(int(server_pb.EventID_event_start_game), fixed.NormalEvent, model.GetDesk(), desk.CreateEventParams(
		model.GetDesk().GetConfig().Context.(*context2.MjContext).StateNumber,
		[]byte{},
		0,
	))
	model.PushEvent(event)
}

// Stop 停止
func (model *MjEventModel) Stop() {
	model.mu.Lock()
	if model.closed {
		model.mu.Unlock()
		return
	}
	model.closed = true
	close(model.event)
	model.mu.Unlock()
}

// PushEvent 压入事件
func (model *MjEventModel) PushEvent(event desk.DeskEvent) {
	model.mu.Lock()
	if model.closed {
		model.mu.Unlock()
		return
	}
	model.event <- event
	model.mu.Unlock()
}

// pushAutoEvent 一段时间后压入自动事件
func (model *MjEventModel) pushAutoEvent(autoEvent *server_pb.AutoEvent, stateNumber int) {
	time.Sleep(time.Millisecond * time.Duration(autoEvent.GetWaitTime()))
	if model.GetDesk().GetConfig().Context.(*context2.MjContext).StateNumber != stateNumber {
		return
	}

	event := desk.NewDeskEvent(int(autoEvent.EventId), fixed.NormalEvent, model.GetDesk(), desk.CreateEventParams(model.GetDesk().GetConfig().Context.(*context2.MjContext).StateNumber, autoEvent.EventContext, 0))

	model.PushEvent(event)
}

// PushRequest 压入玩家请求
func (model *MjEventModel) PushRequest(playerID uint64, head *steve_proto_gaterpc.Header, bodyData []byte) {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name":  "desk.PushRequest",
		"desk_uid":   model.GetDesk().GetUid(),
		"game_id":    model.GetDesk().GetGameId(),
		"player_id":  playerID,
		"message_id": head.GetMsgId(),
	})

	trans := GetTranslator()
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

	event := desk.NewDeskEvent(int(server_pb.EventID(eventID)),
		fixed.NormalEvent,
		model.GetDesk(),
		desk.CreateEventParams(model.GetDesk().GetConfig().Context.(*context2.MjContext).StateNumber, eventConetxtByte, playerID))

	model.PushEvent(event)
}

func (model *MjEventModel) processEvents(ctx context.Context) {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "desk.processEvent",
		"desk_uid":  model.GetDesk().GetUid(),
		"game_id":   model.GetDesk().GetGameId(),
	})
	defer func() {
		if x := recover(); x != nil {
			logEntry.Errorln(x)
			debug.PrintStack()
		}
	}()

	playerModel := GetModelManager().GetPlayerModel(model.GetDesk().GetUid())
	playerEnterChannel := playerModel.getEnterChannel()
	playerLeaveChannel := playerModel.getLeaveChannel()

	for {
		select {
		case <-ctx.Done():
			{
				logEntry.Infoln("done")
				return
			}
		case enterInfo := <-playerEnterChannel:
			{
				model.handlePlayerEnter(enterInfo)
			}
		case leaveInfo := <-playerLeaveChannel:
			{
				model.handlePlayerLeave(leaveInfo)
			}
		case event := <-model.event:
			{
				mjContext := model.GetDesk().GetConfig().Context.(*context2.MjContext)
				stateNumber := event.Params.Params[0].(int)
				context := event.Params.Params[1].([]byte)
				if needCompareStateNumber(&event) && stateNumber != mjContext.StateNumber {
					continue
				}
				model.processEvent(event.EventID, context)
				if model.checkGameOver(logEntry) { // 游戏结束
					return
				}
				model.recordTuoguanOverTimeCount(event)
			}
		}
	}
}

func (model *MjEventModel) recoverGameForPlayer(playerID uint64) {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "recoverGameForPlayer",
		"playerID":  playerID,
	})
	ctx := model.GetDesk().GetConfig().Context.(*context2.MjContext)
	mjContext := &ctx.MjContext
	bankerSeat := mjContext.GetZhuangjiaIndex()
	totalCardsNum := mjContext.GetCardTotalNum()
	gameStage := GetGameStage(mjContext.GetCurState())
	gameID := gutils.GameIDServer2Client(int(mjContext.GetGameId()))
	gameDeskInfo := room.GameDeskInfo{
		GameId:            &gameID,
		GameStage:         &gameStage,
		Players:           GetRecoverPlayerInfo(playerID, model.GetDesk()),
		Dices:             mjContext.GetDices(),
		BankerSeat:        &bankerSeat,
		EastSeat:          &bankerSeat,
		TotalCards:        &totalCardsNum,
		RemainCards:       proto.Uint32(uint32(len(mjContext.GetWallCards()))),
		CostTime:          proto.Uint32(GetStateCostTime(ctx.StateTime.Unix())),
		OperatePid:        GetOperatePlayerID(mjContext),
		NeedHsz:           proto.Bool(gutils.GameHasHszState(mjContext)),
		LastOutCard:       proto.Uint32(getLastOutCard(mjContext.GetLastOutCard())),
		LastOutCardPlayer: proto.Uint64(mjContext.GetLastChupaiPlayer()),
	}
	gameDeskInfo.HasZixun, gameDeskInfo.ZixunInfo = GetZixunInfo(playerID, mjContext)
	gameDeskInfo.HasWenxun, gameDeskInfo.WenxunInfo = GetWenxunInfo(playerID, mjContext)
	gameDeskInfo.HasQgh, gameDeskInfo.QghInfo = GetQghInfo(playerID, mjContext)

	_, gameDeskInfo.HuansanzhangInfo = getHuansanzhangInfo(playerID, mjContext)
	_, gameDeskInfo.DingqueInfo = getDingqueInfo(playerID, mjContext)
	if gameDeskInfo.GetHasZixun() {
		gameDeskInfo.DoorCard = GetDoorCard(mjContext)
	}
	rsp, err := proto.Marshal(&room.RoomResumeGameRsp{
		ResumeRes: room.RoomError_SUCCESS.Enum(),
		GameInfo:  &gameDeskInfo,
	})
	logEntry.WithField("desk_info", gameDeskInfo).Infoln("恢复数据")
	if err != nil {
		logEntry.WithError(err).Errorln("序列化失败")
		return
	}
	GetModelManager().GetMjEventModel(model.GetDesk().GetUid()).Reply([]server_pb.ReplyClientMessage{
		server_pb.ReplyClientMessage{
			Players: []uint64{playerID},
			MsgId:   int32(msgid.MsgID_ROOM_RESUME_GAME_RSP),
			Msg:     rsp,
		},
	})
}

// getContextPlayer 获取context玩家
func (model *MjEventModel) getContextPlayer(playerID uint64) *server_pb.Player {
	mjDeskContext := model.GetGameContext().(*context2.MjContext)
	for _, contextPlayer := range mjDeskContext.MjContext.GetPlayers() {
		if contextPlayer.GetPalyerId() == playerID {
			return contextPlayer
		}
	}
	return nil
}

func (model *MjEventModel) setMjPlayerQuitDesk(playerID uint64, isQuit bool) {
	mjPlayer := model.getContextPlayer(playerID)
	if mjPlayer != nil {
		mjPlayer.IsQuit = isQuit
	}
}

// handlePlayerEnter 处理玩家进入牌桌
func (model *MjEventModel) handlePlayerEnter(enterInfo playerIDWithChannel) {
	model.setMjPlayerQuitDesk(enterInfo.playerID, false)
	modelMgr := GetModelManager()
	modelMgr.GetPlayerModel(model.GetDesk().GetUid()).handlePlayerEnter(enterInfo.playerID)
	model.recoverGameForPlayer(enterInfo.playerID)
	close(enterInfo.finishChannel)
}

// handlePlayerLeave 处理玩家离开牌桌
func (model *MjEventModel) handlePlayerLeave(leaveInfo playerIDWithChannel) {
	modelMgr := GetModelManager()
	playerID := leaveInfo.playerID

	modelMgr.GetPlayerModel(model.GetDesk().GetUid()).handlePlayerLeave(playerID)
	model.setMjPlayerQuitDesk(playerID, true)
	mjPlayer := model.getContextPlayer(playerID)
	ctx := model.GetDesk().GetConfig().Context.(*context2.MjContext)
	mjContext := &ctx.MjContext
	if !gutils.IsPlayerContinue(mjPlayer.GetXpState(), mjContext) {
		playerMgr := player.GetPlayerMgr()
		playerMgr.GetPlayer(playerID).SetDesk(nil)
		playerMgr.UnbindPlayerRoomAddr([]uint64{playerID})
	}
	logrus.WithField("player_id", playerID).Debugln("玩家退出")
	close(leaveInfo.finishChannel)
}

// recordTuoguanOverTimeCount 记录托管超时计数
func (model *MjEventModel) recordTuoguanOverTimeCount(event desk.DeskEvent) {
	if event.EventType != fixed.OverTimeEvent {
		return
	}
	playerID := event.Params.Params[2].(uint64)
	if playerID == 0 {
		return
	}
	id := event.EventID
	if id == int(server_pb.EventID_event_huansanzhang_request) || id == int(server_pb.EventID_event_dingque_request) {
		return
	}
	deskPlayer := player.GetPlayerMgr().GetPlayer(playerID)
	if deskPlayer != nil {
		deskPlayer.OnPlayerOverTime()
	}
}

// processEvent 处理单个事件
// step 1. 调用麻将逻辑的接口来处理事件(返回最新麻将现场, 自动事件， 发送给玩家的消息)， 并且更新 mjContext
// step 2. 将消息发送给玩家
// step 3. 调用 room 的结算逻辑来处理结算
// step 4. 如果有自动事件， 将自动事件写入自动事件通道
// step 5. 如果当前状态是游戏结束状态， 调用 cancel 终止游戏
func (model *MjEventModel) processEvent(eventID int, eventContext []byte) {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "desk.ProcessEvent",
		"event_id":  eventID,
	})
	result, succ := model.callEventHandler(logEntry, eventID, eventContext)
	if !succ {
		return
	}

	// 发送消息给玩家
	model.Reply(result.ReplyMsgs)
	model.GetDesk().GetConfig().Settle.(*MajongSettle).Settle(model.GetDesk(), model.GetDesk().GetConfig())

	// 自动事件不为空，继续处理事件
	if result.AutoEvent != nil {
		if result.AutoEvent.GetWaitTime() == 0 {
			model.processEvent(int(result.AutoEvent.GetEventId()), result.AutoEvent.GetEventContext())
		} else {
			go model.pushAutoEvent(result.AutoEvent, model.GetDesk().GetConfig().Context.(*context2.MjContext).StateNumber)
		}
	}
}

// checkGameOver 检查游戏结束
func (model *MjEventModel) checkGameOver(logEntry *logrus.Entry) bool {
	mjContext := model.GetDesk().GetConfig().Context.(*context2.MjContext).MjContext
	// 游戏结束
	if mjContext.GetCurState() == server_pb.StateID_state_gameover {
		continueModel := GetContinueModel(model.GetDesk().GetUid())
		continueModel.ContinueDesk(mjContext.GetFixNextBankerSeat(), int(mjContext.GetNextBankerSeat()))
		model.GetDesk().GetConfig().Settle.(*MajongSettle).RoundSettle(model.GetDesk(), model.GetDesk().GetConfig())
		logEntry.Infoln("游戏结束状态")
		return true
	}
	return false
}

func (model *MjEventModel) cancelTuoguanGameOver() {
	playerModel := GetModelManager().GetPlayerModel(model.GetDesk().GetUid())
	for _, player := range playerModel.GetDeskPlayers() {
		if player.IsTuoguan() {
			player.SetTuoguan(false, true)
		}
	}
}

func (model *MjEventModel) Reply(replyMsgs []server_pb.ReplyClientMessage) {
	if replyMsgs == nil {
		return
	}
	for _, msg := range replyMsgs {
		GetModelManager().GetMessageModel(model.GetDesk().GetUid()).BroadcastMessage(msg.GetPlayers(), msgid.MsgID(msg.GetMsgId()), msg.GetMsg(), true)
	}
}

// callEventHandler 调用事件处理器
func (model *MjEventModel) callEventHandler(logEntry *logrus.Entry, eventID int, eventContext []byte) (result majong_process.HandleMajongEventResult, succ bool) {
	succ = false
	conte := model.GetDesk().GetConfig().Context.(*context2.MjContext)
	stateNumber, mjContext, stateTime := conte.StateNumber, conte.MjContext, conte.StateTime
	oldState := mjContext.GetCurState()
	if result = majong_process.HandleMajongEvent(majong_process.HandleMajongEventParams{
		MajongContext: mjContext,
		EventID:       server_pb.EventID(eventID),
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
	model.GetDesk().GetConfig().Context = &context2.MjContext{
		MjContext: newContext,
		//StateNumber: stateNumber,
		StateTime: stateTime,
	}
	model.GetDesk().GetConfig().Context.(*context2.MjContext).SetStateNumber(stateNumber)
	println("更新桌子状体 old ", model.GetDesk().GetConfig().Context.(*context2.MjContext).StateNumber)
	succ = true
	return
}

// needCompareStateNumber 判断事件是否需要比较 stateNumber
func needCompareStateNumber(event *desk.DeskEvent) bool {
	if event.EventID == int(server_pb.EventID_event_huansanzhang_request) ||
		event.EventID == int(server_pb.EventID_event_dingque_request) {
		return false
	}
	return true
}

// timerTask 定时任务，产生自动事件
func (model *MjEventModel) timerTask(ctx context.Context) {
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
				if model.closed {
					return
				}
				model.genTimerEvent()
			}
		case <-ctx.Done():
			{
				return
			}
		}
	}
}

// genTimerEvent 生成计时事件
func (model *MjEventModel) genTimerEvent() {
	// 先将 context 指针读出来拷贝， 后面的 context 修改都会分配一块新的内存
	dContext := model.GetDesk().GetConfig().Context.(*context2.MjContext)
	tuoGuanPlayers := GetModelManager().GetPlayerModel(model.GetDesk().GetUid()).GetTuoguanPlayers()

	deskPlayers := GetModelManager().GetPlayerModel(model.GetDesk().GetUid()).GetDeskPlayers()
	robotLvs := make(map[uint64]int, len(deskPlayers))
	for _, deskPlayer := range deskPlayers {
		robotLv := deskPlayer.GetRobotLv()
		if robotLv != 0 {
			robotLvs[deskPlayer.GetPlayerID()] = robotLv
		}
	}
	result := ai.GetAtEvent().GenerateV2(&ai.AutoEventGenerateParams{
		Desk:           model.GetDesk(),
		CurTime:        time.Now(),
		StartTime:      dContext.StateTime,
		RobotLv:        robotLvs,
		TuoGuanPlayers: tuoGuanPlayers,
	})
	for _, event := range result.Events {
		logrus.WithField("event_id", event.EventID).Debugln("写入自动事件")
		model.PushEvent(event)
	}
}
