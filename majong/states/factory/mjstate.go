package factory

import (
	"steve/majong/interfaces"
	"steve/majong/states/common"
	majongpb "steve/server_pb/majong"
)

func createState(stateID majongpb.StateID) interfaces.MajongState {
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
		return new(common.ZimoState)
	case majongpb.StateID_state_zimo_settle:
		return new(common.ZiMoSettleState)
	case majongpb.StateID_state_hu:
		return new(common.HuState)
	case majongpb.StateID_state_hu_settle:
		return new(common.HuSettleState)
	case majongpb.StateID_state_qiangganghu:
		return new(common.QiangganghuState)
	case majongpb.StateID_state_qiangganghu_settle:
		return new(common.QiangGangHuSettleState)
	case majongpb.StateID_state_angang:
		return new(common.AnGangState)
	case majongpb.StateID_state_gang_settle:
		return new(common.GangSettleState)
	case majongpb.StateID_state_bugang:
		return new(common.BuGangState)
	case majongpb.StateID_state_gang:
		return new(common.MingGangState)
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
	case majongpb.StateID_state_gamestart_buhua:
		return new(common.GameStartBuhuaState)
	case majongpb.StateID_state_xingpai_buhua:
		return new(common.XingPaiBuhuaState)
	case majongpb.StateID_state_chi:
		return new(common.ChiState)
	default:
		return nil
	}
}
