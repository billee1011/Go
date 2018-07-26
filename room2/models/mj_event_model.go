package models

import (
	"runtime/debug"
	"github.com/Sirupsen/logrus"
	"context"
	server_pb "steve/server_pb/majong"
	majong_process "steve/majong/export/process"
	context2 "steve/room2/contexts"
	"time"
	"steve/client_pb/msgid"
	"github.com/golang/protobuf/proto"
	"steve/structs/proto/gate_rpc"
	"steve/room2/ai"
	"steve/room2/player"
	"steve/room2/desk"
	"steve/room2/fixed"
)

type MjEventModel struct {
	event chan desk.DeskEvent // 牌桌事件通道
	BaseModel
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
		model.processEvents(model.GetDesk().Context)
	}()
	go func() {
		model.timerTask(model.GetDesk().Context)
	}()

	event := desk.NewDeskEvent(int(server_pb.EventID_event_start_game), fixed.NormalEvent, model.GetDesk(), desk.CreateEventParams(
		model.GetDesk().GetConfig().Context.(*context2.MjContext).StateNumber,
		[]byte{},
		0,
	))
	model.PushEvent(event)
}

func (model *MjEventModel) Stop() {
	close(model.event)
}

// PushEvent 压入事件
func (model *MjEventModel) PushEvent(event desk.DeskEvent) {
	model.event <- event
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

	/*interfaces.Event{
		ID:        server_pb.EventID(eventID),
		Context:   eventConetxtByte,
		EventType: interfaces.NormalEvent,
		PlayerID:  playerID,
	}
*/
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
	for {

		select {
		case <-ctx.Done():
			{
				logEntry.Infoln("done")
				return
			}
			/*case enterQuitInfo := <-model.GetDesk().PlayerEnterQuitChannel():
				{
					d.handleEnterQuit(enterQuitInfo)
				}*/
		case event := <-model.event:
			{
				println("收到事件 : ", event.EventID)
				if event.EventID == 3 {
					println("1111")
				}
				mjContext := model.GetDesk().GetConfig().Context.(*context2.MjContext)
				stateNumber := event.Params.Params[0].(int)
				context := event.Params.Params[1].([]byte)
				println("event state:",stateNumber,"-----context state:",mjContext.StateNumber)
				if needCompareStateNumber(&event) && stateNumber != mjContext.StateNumber {
					continue
				}
				model.processEvent(event.EventID, context)
				model.recordTuoguanOverTimeCount(event)
			}
		}
	}
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
	if eventID == 4 || eventID == 3 {
		println(1111)
	}
	result, succ := model.callEventHandler(logEntry, eventID, eventContext)
	if !succ {
		return
	}

	// 发送消息给玩家
	model.Reply(result.ReplyMsgs)
	model.GetDesk().GetConfig().Settle.(*MajongSettle).Settle(model.GetDesk(), model.GetDesk().GetConfig())

	if model.checkGameOver(logEntry) {
		return
	}
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
		model.GetDesk().GetConfig().Settle.(*MajongSettle).RoundSettle(model.GetDesk(), model.GetDesk().GetConfig())
		logEntry.Infoln("游戏结束状态")
		model.GetDesk().Cancel()
		return true
	}
	return false
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
		MjContext:   newContext,
		//StateNumber: stateNumber,
		StateTime:   stateTime,
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
		MajongContext:  &dContext.MjContext,
		CurTime:        time.Now(),
		StateTime:      dContext.StateTime,
		RobotLv:        robotLvs,
		TuoGuanPlayers: tuoGuanPlayers,
	})
	for _, event := range result.Events {
		GetModelManager().GetMjEventModel(model.GetDesk().GetUid()).PushEvent(event)
	}
}
