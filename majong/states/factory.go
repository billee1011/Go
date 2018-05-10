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
	case majongpb.StateID_state_angang:
		return new(AnGangState)
	case majongpb.StateID_state_bugang:
		return new(BuGangState)
	case majongpb.StateID_state_gang:
		return new(MingGangState)
	case majongpb.StateID_state_peng:
		return new(PengState)
	case majongpb.StateID_state_dingque:
		return new(DingqueState)
	case majongpb.StateID_state_waitqiangganghu:
		return new(WaitQiangganghuState)
	default:
		return nil
	}
}
