package interfaces

import (
	majongpb "steve/server_pb/majong"
)

// TransitionValidator 状态转移验证器
type TransitionValidator interface {
	Valid(oldState majongpb.StateID, newState majongpb.StateID, event majongpb.EventID) error
}
