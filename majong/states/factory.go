package states

import (
	"steve/majong/interfaces"
	majongpb "steve/server_pb/majong"
)

type factory struct{}

var _ interfaces.MajongStateFactory = new(factory)

// NewFactory 创建状态工厂
func NewFactory() interfaces.MajongStateFactory {
	return new(factory)
}

func (f *factory) CreateState(gameID int, stateID majongpb.StateID) interfaces.MajongState {
	switch stateID {
	case majongpb.StateID_state_init:
		return new(InitState)
	default:
		return nil
	}
}
