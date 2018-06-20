package factory

import (
	"steve/gutils"
	"steve/majong/global"
	"steve/majong/interfaces"
	majongpb "steve/server_pb/majong"
)

type createGameStateFunc func(majongpb.StateID) interfaces.MajongState

type factory struct {
	creators map[int]createGameStateFunc
}

var _ interfaces.MajongStateFactory = new(factory)

// newFactory 创建状态工厂
func newFactory() interfaces.MajongStateFactory {
	creators := map[int]createGameStateFunc{
		gutils.SCXLGameID: createSCXLState,
		gutils.SCXZGameID: createSCXZState,
	}
	return &factory{
		creators: creators,
	}
}

func (f *factory) CreateState(gameID int, stateID majongpb.StateID) interfaces.MajongState {
	creator, exists := f.creators[gameID]
	if !exists {
		return nil
	}
	return creator(stateID)
}

func init() {
	f := newFactory()
	global.SetMajongStateFacotry(f)
}
