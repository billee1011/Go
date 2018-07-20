package machine

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"time"
)

// Machine 状态机
type Machine interface {
	ProcessEvent(event Event) error
	GetStateID() int
	SetStateID(state int)
}

// DefaultProcessor 默认事件处理器
func DefaultProcessor(m Machine, stateFactory StateFactory, event Event) error {
	start := time.Now()
	logEntry := logrus.WithFields(logrus.Fields{"eventId": event.EventID,
		"start": start})
	curStateID := m.GetStateID()

	curState := stateFactory.NewState(curStateID)
	if curState == nil {
		return fmt.Errorf("创建状态失败 %v", curStateID)
	}

	newStateID, err := curState.OnEvent(m, event)
	if err != nil {
		return err
	}

	if curStateID == newStateID {
		end := time.Now()
		logEntry.WithFields(logrus.Fields{"end": end, "duration": end.Sub(start)}).Debug("状态机退出")
		return nil
	}

	newState := stateFactory.NewState(newStateID)
	if newState == nil {
		return fmt.Errorf("创建状态失败 %v", newStateID)
	}

	curState.OnExit(m)
	m.SetStateID(newStateID)
	newState.OnEnter(m)

	end := time.Now()
	logEntry.WithFields(logrus.Fields{"end": end, "duration": end.Sub(start)}).Debug("状态机退出")
	return nil
}
