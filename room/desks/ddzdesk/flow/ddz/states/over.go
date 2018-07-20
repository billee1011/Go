package states

import (
	"github.com/Sirupsen/logrus"
	"steve/room/desks/ddzdesk/flow/machine"
	"steve/server_pb/ddz"
)

type overState struct{}

func (s *overState) OnEnter(m machine.Machine) {
	logrus.WithField("context", getDDZContext(m)).Debugln("进入Over状态")
}

func (s *overState) OnExit(m machine.Machine) {
	logrus.WithField("context", getDDZContext(m)).Debugln("离开Over状态")
}

func (s *overState) OnEvent(m machine.Machine, event machine.Event) (int, error) {
	return int(ddz.StateID_state_over), nil
}
