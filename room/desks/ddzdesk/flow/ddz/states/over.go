package states

import (
	"steve/room/desks/ddzdesk/flow/machine"
	"github.com/Sirupsen/logrus"
	"steve/server_pb/ddz"
	"steve/client_pb/msgId"
	"steve/client_pb/room"
	"github.com/gogo/protobuf/proto"
)

type overState struct{}

func (s *overState) OnEnter(m machine.Machine) {
	logrus.WithField("context", getDDZContext(m)).Debugln("进入Over状态")

	context := getDDZContext(m)

	antiSpring := !context.Spring && context.AntiSpring
	broadcast(m, msgid.MsgID_ROOM_DDZ_GAME_OVER_NTF,
		&room.DDZGameOverNtf{
			WinnerId: &context.WinnerId,
			ShowHandTime: proto.Uint32(4),
			Spring: &context.Spring,
			AntiSpring: &antiSpring,
		},
	)
}

func (s *overState) OnExit(m machine.Machine) {
	logrus.WithField("context", getDDZContext(m)).Debugln("离开Over状态")
}

func (s *overState) OnEvent(m machine.Machine, event machine.Event) (int, error) {
	return int(ddz.StateID_state_over), nil
}
