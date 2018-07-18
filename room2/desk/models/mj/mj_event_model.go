package mj

import (
	"runtime/debug"
	"github.com/Sirupsen/logrus"
	"context"
	server_pb "steve/server_pb/majong"
	majong_process "steve/majong/export/process"
	"steve/room2/desk"
	"steve/room2/common"
	"steve/room2/desk/models"
	"steve/room2/desk/models/public"
	context2 "steve/room2/desk/contexts"
	"time"
	"steve/client_pb/msgid"
	"github.com/golang/protobuf/proto"
	"steve/structs/proto/gate_rpc"
	"steve/room2/util"
)

type MjEventModel struct {
	event chan desk.DeskEvent // 牌桌事件通道
	public.BaseModel
}

func NewMjEventModel(desk *desk.Desk) MjEventModel {
	result := MjEventModel{}
	result.SetDesk(desk)
	return result
}

func (model MjEventModel) GetName() string {
	return models.Event
}

func (model MjEventModel) Start() {
	model.event = make(chan desk.DeskEvent, 16)
}

func (model MjEventModel) Stop() {
	close(model.event)
}


// PushEvent 压入事件
func (model MjEventModel) PushEvent(event desk.DeskEvent) {
	model.event <- event
}

// pushAutoEvent 一段时间后压入自动事件
func (model MjEventModel) pushAutoEvent(autoEvent *server_pb.AutoEvent, stateNumber int) {
	time.Sleep(time.Millisecond * time.Duration(autoEvent.GetWaitTime()))
	if model.GetDesk().GetConfig().Context.(context2.MjContext).StateNumber != stateNumber {
		return
	}

	event := desk.NewDeskEvent(int(autoEvent.EventId),common.NormalEvent,model.GetDesk(),common.CreateEventParams(model.GetDesk().GetConfig().Context.(context2.MjContext).StateNumber,autoEvent.EventContext,0))

	model.PushEvent(event)
}



// PushRequest 压入玩家请求
func (model MjEventModel) PushRequest(playerID uint64, head *steve_proto_gaterpc.Header, bodyData []byte) {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name":  "desk.PushRequest",
		"desk_uid":   model.GetDesk().GetUid(),
		"game_id":    model.GetDesk().GetGameId(),
		"player_id":  playerID,
		"message_id": head.GetMsgId(),
	})

	trans := util.GetTranslator()
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
		common.NormalEvent,
		model.GetDesk(),
		common.CreateEventParams(model.GetDesk().GetConfig().Context.(context2.MjContext).StateNumber,eventConetxtByte,playerID))

	model.PushEvent(event)
}



func (model MjEventModel) processEvents(ctx context.Context) {
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
	mjContext := model.GetDesk().GetConfig().Context.(*context2.MjContext)
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
				stateNumber := event.Params.Params[0].(int)
				context := event.Params.Params[1].([]byte)
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
func (model MjEventModel) recordTuoguanOverTimeCount(event desk.DeskEvent) {
	if event.EventType != common.OverTimeEvent {
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
	deskPlayer := model.GetDesk().GetPlayer(playerID)
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
func (model MjEventModel) processEvent(eventID int, eventContext []byte) {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "desk.ProcessEvent",
		"event_id":  eventID,
	})
	result, succ := model.callEventHandler(logEntry, eventID, eventContext)
	if !succ {
		return
	}

	// 发送消息给玩家
	model.reply(result.ReplyMsgs)
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
func (model *MjEventModel) checkGameOver(logEntry *logrus.Entry) bool {
	mjContext := model.GetDesk().GetConfig().Context.(context2.MjContext).MjContext
	// 游戏结束
	if mjContext.GetCurState() == server_pb.StateID_state_gameover {
		d.settler.RoundSettle(d, mjContext)
		logEntry.Infoln("游戏结束状态")
		d.cancel()
		return true
	}
	return false
}

func (model MjEventModel) Reply(replyMsgs []server_pb.ReplyClientMessage) {
	if replyMsgs == nil {
		return
	}
	for _, msg := range replyMsgs {
		model.GetDesk().GetModel(models.Message).(public.MessageModel).BroadcastMessage(msg.GetPlayers(), msgid.MsgID(msg.GetMsgId()), msg.GetMsg(), true)
	}
}



// callEventHandler 调用事件处理器
func (model MjEventModel) callEventHandler(logEntry *logrus.Entry, eventID int, eventContext []byte) (result majong_process.HandleMajongEventResult, succ bool) {
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
		StateNumber: stateNumber,
		StateTime:   stateTime,
	}
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