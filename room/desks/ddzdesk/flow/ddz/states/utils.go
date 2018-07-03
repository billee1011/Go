package states

import (
	"fmt"
	msgid "steve/client_pb/msgId"
	"steve/room/desks/ddzdesk/flow/ddz/ddzmachine"
	"steve/room/desks/ddzdesk/flow/machine"
	"steve/server_pb/ddz"
	"time"

	"github.com/golang/protobuf/proto"
)

// getDDZContext 从状态机中获取斗地主现场
func getDDZContext(m machine.Machine) *ddz.DDZContext {
	dm, ok := m.(*ddzmachine.DDZMachine)
	if !ok {
		return nil
	}
	return dm.GetDDZContext()
}

// sendMessage 向玩家发送消息
func sendMessage(m machine.Machine, players []uint64, msgID msgid.MsgID, body proto.Message) error {
	dm, ok := m.(*ddzmachine.DDZMachine)
	if !ok {
		return fmt.Errorf("不是斗地主状态机")
	}
	return dm.SendMessage(players, msgID, body)
}

// setMachineAutoEvent 设置状态机自动事件
func setMachineAutoEvent(m machine.Machine, event machine.Event, duration time.Duration) {
	dm, ok := m.(*ddzmachine.DDZMachine)
	if !ok {
		return
	}
	dm.SetAutoEvent(event, duration)
}
