package states

import (
	"steve/room/desks/ddzdesk/flow/machine"
	"steve/server_pb/ddz"

	"github.com/Sirupsen/logrus"
)

type initState struct{}

func (s *initState) OnEnter(m machine.Machine) {
	logrus.WithField("context", getDDZContext(m)).Debugln("进入初始状态")
}

func (s *initState) OnExit(m machine.Machine) {
	logrus.WithField("context", getDDZContext(m)).Debugln("离开初始状态")
}

func (s *initState) OnEvent(m machine.Machine, event machine.Event) (int, error) {
	if event.EventID == int(ddz.EventID_event_start_game) {
		return s.onStartGame(m)
	}
	return int(ddz.StateID_state_init), nil
}

func (s *initState) onStartGame(m machine.Machine) (int, error) {
	logrus.WithField("context", getDDZContext(m)).Debugln("开始游戏")
	return int(ddz.StateID_state_deal), nil
}
