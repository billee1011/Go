package ddzmachine

import (
	"fmt"
	msgid "steve/client_pb/msgId"
	"steve/room/desks/ddzdesk/flow/machine"
	"steve/server_pb/ddz"
	"time"

	"github.com/golang/protobuf/proto"
)

// MessageSender 状态机的消息发送器
type MessageSender func(players []uint64, msgID msgid.MsgID, body proto.Message) error

// DDZMachine 斗地主状态机
type DDZMachine struct {
	factory    machine.StateFactory
	ddzContext *ddz.DDZContext
	sender     MessageSender
}

// CreateDDZMachine 创建斗地主状态机
func CreateDDZMachine(ddzContext *ddz.DDZContext, stateFactory machine.StateFactory, sender MessageSender) *DDZMachine {
	return &DDZMachine{
		ddzContext: ddzContext,
		factory:    stateFactory,
		sender:     sender,
	}
}

// ProcessEvent 处理事件
func (m *DDZMachine) ProcessEvent(event machine.Event) error {
	return machine.DefaultProcessor(m, m.factory, event)
}

// GetStateID 获取状态 ID
func (m *DDZMachine) GetStateID() int {
	return int(m.ddzContext.GetCurState())
}

// SetStateID 设置状态 ID
func (m *DDZMachine) SetStateID(state int) {
	m.ddzContext.CurState = ddz.StateID(state)
}

// GetDDZContext 获取牌局现场
func (m *DDZMachine) GetDDZContext() *ddz.DDZContext {
	return m.ddzContext
}

// SendMessage 发送消息
func (m *DDZMachine) SendMessage(players []uint64, msgID msgid.MsgID, body proto.Message) error {
	if m.sender == nil {
		return fmt.Errorf("没有设置消息发送器")
	}
	return m.sender(players, msgID, body)
}

// SetAutoEvent 设置自动事件， duration ： 多久后触发
func (m *DDZMachine) SetAutoEvent(event machine.Event, duration time.Duration) {

}
