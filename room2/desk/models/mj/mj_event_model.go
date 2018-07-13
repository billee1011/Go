package mj

import (
	"runtime/debug"
	"github.com/Sirupsen/logrus"
	"context"
	server_pb "steve/server_pb/majong"
	"steve/room2/desk"
	"steve/room2/desk/models/public"
	"steve/room2/desk/models"
	context2 "steve/room2/desk/contexts"
	"time"
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
	mjContext := model.GetDesk().GetConfig().Context.(context2.MjContext)
	for {
		select {
		case <-ctx.Done():
			{
				logEntry.Infoln("done")
				return
			}
		/*case enterQuitInfo := <-d.PlayerEnterQuitChannel():
			{
				d.handleEnterQuit(enterQuitInfo)
			}*/
		case event := <-model.event:
			{
				if needCompareStateNumber(&event) && event.stateNumber != mjContext.StateNumber {
					continue
				}
				model.processEvent(event.EventID, event.event.Context)
				model.recordTuoguanOverTimeCount(event.event)
			}
		}
	}
}



// processEvent 处理单个事件
// step 1. 调用麻将逻辑的接口来处理事件(返回最新麻将现场, 自动事件， 发送给玩家的消息)， 并且更新 mjContext
// step 2. 将消息发送给玩家
// step 3. 调用 room 的结算逻辑来处理结算
// step 4. 如果有自动事件， 将自动事件写入自动事件通道
// step 5. 如果当前状态是游戏结束状态， 调用 cancel 终止游戏
func (model MjEventModel) processEvent(eventID server_pb.EventID, eventContext []byte) {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "desk.ProcessEvent",
		"event_id":  eventID,
	})
	result, succ := model.callEventHandler(logEntry, eventID, eventContext)
	if !succ {
		return
	}

	// 发送消息给玩家
	d.reply(result.ReplyMsgs)

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



// callEventHandler 调用事件处理器
func (model MjEventModel) callEventHandler(logEntry *logrus.Entry, eventID server_pb.EventID, eventContext []byte) (result majong_process.HandleMajongEventResult, succ bool) {
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

// needCompareStateNumber 判断事件是否需要比较 stateNumber
func needCompareStateNumber(event *desk.DeskEvent) bool {
	if event.EventID == int(server_pb.EventID_event_huansanzhang_request) ||
		event.EventID == int(server_pb.EventID_event_dingque_request) {
		return false
	}
	return true
}