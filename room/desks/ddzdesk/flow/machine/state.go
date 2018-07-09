package machine

// State 状态
type State interface {
	OnEnter(m Machine)
	OnExit(m Machine)
	OnEvent(m Machine, event Event) (int, error)
}

// StateFactory 状态工厂
type StateFactory interface {
	NewState(stateID int) State
}
