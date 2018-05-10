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
	case majongpb.StateID_state_xipai:
		return new(XipaiState)
	case majongpb.StateID_state_fapai:
		return new(FapaiState)
	case majongpb.StateID_state_huansanzhang:
		return new(HuansanzhangState)
	case majongpb.StateID_state_zixun:
		return new(ZiXunState)
	case majongpb.StateID_state_chupai:
		return new(ChupaiState)
	case majongpb.StateID_state_zimo:
		return new(ZimoState)
	case majongpb.StateID_state_hu:
		return new(HuState)
	case majongpb.StateID_state_qiangganghu:
		return new(QiangganghuState)
	default:
		return nil
	}
}
