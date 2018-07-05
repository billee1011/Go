package machine

import (
	"fmt"
)

// Machine 状态机
type Machine interface {
	ProcessEvent(event Event) error
	GetStateID() int
	SetStateID(state int)
}

// DefaultProcessor 默认事件处理器
func DefaultProcessor(m Machine, stateFactory StateFactory, event Event) error {
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
		return nil
	}

	newState := stateFactory.NewState(newStateID)
	if newState == nil {
		return fmt.Errorf("创建状态失败 %v", newStateID)
	}

	curState.OnEnter(m)
	m.SetStateID(newStateID)
	newState.OnExit(m)
	return nil
}
