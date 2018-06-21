package factory

import (
	"steve/majong/interfaces"
	"steve/majong/states/common"
	"steve/majong/states/scxz"
	majongpb "steve/server_pb/majong"
)

func createSCXZState(stateID majongpb.StateID) interfaces.MajongState {
	switch stateID {
	case majongpb.StateID_state_init:
		return new(common.InitState)
	case majongpb.StateID_state_xipai:
		return new(common.XipaiState)
	case majongpb.StateID_state_fapai:
		return new(common.FapaiState)
	case majongpb.StateID_state_huansanzhang:
		return new(common.HuansanzhangState)
	case majongpb.StateID_state_zixun:
		return new(common.ZiXunState)
	case majongpb.StateID_state_chupai:
		return new(common.ChupaiState)
	case majongpb.StateID_state_zimo:
		return new(scxz.ZimoState)
	case majongpb.StateID_state_zimo_settle:
		return new(scxz.ZiMoSettleState)
	case majongpb.StateID_state_hu:
		return new(scxz.HuState)
	case majongpb.StateID_state_hu_settle:
		return new(scxz.HuSettleState)
	case majongpb.StateID_state_qiangganghu:
		return new(scxz.QiangganghuState)
	case majongpb.StateID_state_qiangganghu_settle:
		return new(scxz.QiangGangHuSettleState)
	case majongpb.StateID_state_angang:
		return new(scxz.AnGangState)
	case majongpb.StateID_state_gang_settle:
		return new(scxz.GangSettleState)
	case majongpb.StateID_state_bugang:
		return new(scxz.BuGangState)
	case majongpb.StateID_state_gang:
		return new(scxz.MingGangState)
	case majongpb.StateID_state_peng:
		return new(common.PengState)
	case majongpb.StateID_state_dingque:
		return new(common.DingqueState)
	case majongpb.StateID_state_waitqiangganghu:
		return new(common.WaitQiangganghuState)
	case majongpb.StateID_state_chupaiwenxun:
		return new(common.ChupaiwenxunState)
	case majongpb.StateID_state_mopai:
		return new(common.MoPaiState)
	case majongpb.StateID_state_gameover:
		return new(common.GameOverState)
	default:
		return nil
	}
}
