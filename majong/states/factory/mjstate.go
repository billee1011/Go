package factory
/*
功能： 麻将通用状态机定义

 */
import (
	"steve/majong/interfaces"
	"steve/majong/states/common"
	majongpb "steve/entity/majong"
)

// 状态列表：在下面列表中增加自己添加的新状态.
var mapState =  map[majongpb.StateID] interfaces.MajongState{
	majongpb.StateID_state_init:new(common.InitState),
	majongpb.StateID_state_xipai:new(common.XipaiState),
	majongpb.StateID_state_fapai:new(common.FapaiState),
	majongpb.StateID_state_huansanzhang:new(common.HuansanzhangState),
	majongpb.StateID_state_zixun:new(common.ZiXunState),
	majongpb.StateID_state_chupai:new(common.ChupaiState),
	majongpb.StateID_state_zimo:new(common.ZimoState),
	majongpb.StateID_state_zimo_settle:new(common.ZiMoSettleState),
	majongpb.StateID_state_hu:new(common.HuState),
	majongpb.StateID_state_hu_settle:new(common.HuSettleState),
	majongpb.StateID_state_qiangganghu:new(common.QiangganghuState),
	majongpb.StateID_state_qiangganghu_settle:new(common.QiangGangHuSettleState),
	majongpb.StateID_state_angang:new(common.AnGangState),
	majongpb.StateID_state_gang_settle:new(common.GangSettleState),
	majongpb.StateID_state_bugang:new(common.BuGangState),
	majongpb.StateID_state_gang:new(common.MingGangState),
	majongpb.StateID_state_peng:new(common.PengState),
	majongpb.StateID_state_dingque:new(common.DingqueState),
	majongpb.StateID_state_waitqiangganghu:new(common.WaitQiangganghuState),
	majongpb.StateID_state_chupaiwenxun:new(common.ChupaiwenxunState),
	majongpb.StateID_state_mopai:new(common.MoPaiState),
	majongpb.StateID_state_gameover:new(common.GameOverState),
	majongpb.StateID_state_gamestart_buhua:new(common.GameStartBuhuaState),
	majongpb.StateID_state_xingpai_buhua:new(common.XingPaiBuhuaState),
	majongpb.StateID_state_chi:new(common.ChiState),
}

type mjStateMgr struct {
}

func (m*mjStateMgr) NewState(stateID majongpb.StateID) interfaces.MajongState {
	return mapState[stateID]
}

/*
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
*/