package states

import (
	"steve/room/desks/ddzdesk/flow/machine"
	"steve/server_pb/ddz"

	"github.com/Sirupsen/logrus"
)

type dealState struct{}

func (s *dealState) OnEnter(m machine.Machine) {
	logrus.WithField("context", getDDZContext(m)).Debugln("进入发牌状态")
}

func (s *dealState) OnExit(m machine.Machine) {
	logrus.WithField("context", getDDZContext(m)).Debugln("离开发牌状态")
}

func (s *dealState) OnEvent(m machine.Machine, event machine.Event) (int, error) {
	logrus.WithFields(logrus.Fields{
		"context": getDDZContext(m),
		"event":   event,
	}).Debugln("发牌处理事件")
	return int(ddz.StateID_state_deal), nil
}
