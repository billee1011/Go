package ddzdesk

import (
	msgid "steve/client_pb/msgId"
	"steve/room/desks/ddzdesk/flow/ddz/ddzmachine"
	"steve/room/desks/ddzdesk/flow/ddz/procedure"
	"steve/room/desks/deskbase"
	"steve/room/interfaces"
	"steve/room/interfaces/facade"
	"steve/room/interfaces/global"
	"steve/server_pb/ddz"
	"steve/structs/proto/gate_rpc"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
)

// deskEvent 牌桌事件
type deskEvent struct {
	eventID      int
	eventContext []byte
}

// desk 斗地主牌桌
type desk struct {
	*deskbase.DeskBase
	eventChannel   chan deskEvent
	closingChannel chan struct{}
	ddzContext     *ddz.DDZContext
}

// initDDZContext 初始化斗地主现场
func (d *desk) initDDZContext() {
	d.ddzContext = procedure.CreateInitDDZContext(facade.GetDeskPlayerIDs(d))
}

// Start 启动牌桌逻辑
// finish : 当牌桌逻辑完成时调用
func (d *desk) Start(finish func()) error {
	d.eventChannel = make(chan deskEvent, 4)
	d.closingChannel = make(chan struct{})

	d.initDDZContext()
	go func() {
		d.run()
		finish()
	}()
	d.pushEvent(&deskEvent{
		eventID: int(ddz.EventID_event_start_game),
	})
	return nil
}

// Stop 停止牌桌
func (d *desk) Stop() error {
	d.closingChannel <- struct{}{}
	return nil
}

// PushRequest 压入玩家请求
func (d *desk) PushRequest(playerID uint64, head *steve_proto_gaterpc.Header, bodyData []byte) {
	entry := logrus.WithFields(logrus.Fields{
		"func_name": "desk.PushRequest",
		"player_id": playerID,
		"msg_id":    head.GetMsgId(),
	})

	translator := global.GetReqEventTranslator()
	eventID, eventData, err := translator.Translate(playerID, head, bodyData)
	if err != nil {
		entry.WithError(err).Errorln("事件转换失败")
		return
	}
	if eventID == 0 {
		entry.Warningln("没有对应事件")
		return
	}
	eventMessage, ok := eventData.(proto.Message)
	if !ok {
		entry.Errorln("事件数据不是 proto.Message 类型")
		return
	}
	eventContext, err := proto.Marshal(eventMessage)
	if err != nil {
		entry.WithError(err).Errorln("事件消息序列化失败")
		return
	}
	d.pushEvent(&deskEvent{
		eventID:      eventID,
		eventContext: eventContext,
	})
}

func (d *desk) pushEvent(e *deskEvent) {
	d.eventChannel <- *e
}

// PushEvent 压入事件
func (d *desk) PushEvent(event interfaces.Event) {
	return
}

// run 执行牌桌逻辑
func (d *desk) run() {

forstart:
	for {
		select {
		case event := <-d.eventChannel:
			{
				d.processEvent(&event)
			}
		case <-d.closingChannel:
			{
				break forstart
			}
		}
	}
}

func (d *desk) processEvent(e *deskEvent) {
	entry := logrus.WithFields(logrus.Fields{
		"func_name": "desk.processEvent",
		"event_id":  e.eventID,
	})

	params := procedure.Params{
		Context:      *d.ddzContext,
		Sender:       d.getMessageSender(),
		EventID:      e.eventID,
		EventContext: e.eventContext,
	}
	result := procedure.HandleEvent(params)
	if !result.Succeed {
		entry.Errorln("处理事件失败")
		return
	}

	d.ddzContext = &result.Context
	// 游戏结束
	if d.ddzContext.GetCurState() == ddz.StateID_state_over {
		go func() { d.closingChannel <- struct{}{} }()
		return
	}
	if result.HasAutoEvent {
		if result.AutoEventDuration == time.Duration(0) {
			d.pushEvent(&deskEvent{
				eventID:      result.AutoEventID,
				eventContext: result.AutoEventContext,
			})
		} else {
			go func() {
				timer := time.NewTimer(result.AutoEventDuration)
				<-timer.C
				d.pushEvent(&deskEvent{
					eventID:      result.AutoEventID,
					eventContext: result.AutoEventContext,
				})
			}()
		}
	}
}

func (d *desk) getMessageSender() ddzmachine.MessageSender {
	return func(players []uint64, msgID msgid.MsgID, body proto.Message) error {
		return facade.BroadCastDeskMessage(d, players, msgID, body, true)
	}
}
