package factory

import (
	"steve/majong/global"
	"steve/majong/interfaces"
	majongpb "steve/server_pb/majong"
)

type createGameStateFunc func(majongpb.StateID) interfaces.MajongState

type factory struct {
	creators createGameStateFunc
}

var _ interfaces.MajongStateFactory = new(factory)

// newFactory 创建状态工厂
func newFactory() interfaces.MajongStateFactory {
	return &factory{
		creators: createState,
	}
}

func (f *factory) CreateState(gameID int, stateID majongpb.StateID) interfaces.MajongState {
	return f.creators(stateID)
}

func init() {
	f := newFactory()
	global.SetMajongStateFacotry(f)
}
