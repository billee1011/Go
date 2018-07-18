package interfaces

import (
	majongpb "steve/server_pb/majong"
)

// MajongState 牌局状态
type MajongState interface {
	ProcessEvent(eventID majongpb.EventID, eventContext []byte, flow MajongFlow) (newState majongpb.StateID, err error)
	OnEntry(flow MajongFlow)
	OnExit(flow MajongFlow)
}

// MajongStateFactory 麻将状态工厂
type MajongStateFactory interface {
	// CreateState 根据 gameID 和 stateID 创建麻将状态
	CreateState(gameID int32, stateID majongpb.StateID) MajongState
}
