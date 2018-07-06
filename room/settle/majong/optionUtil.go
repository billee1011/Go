package majong

import (
	"steve/common/mjoption"
	majongpb "steve/server_pb/majong"
)

// GetSettleOption 获取游戏的结算配置
func GetSettleOption(gameID int) *mjoption.SettleOption {
	return mjoption.GetSettleOption(mjoption.GetGameOptions(gameID).SettleOptionID)
}

// IsGangSettle 是否是杠结算方式
func IsGangSettle(settleType majongpb.SettleType) bool {
	return map[majongpb.SettleType]bool{
		majongpb.SettleType_settle_angang:   true,
		majongpb.SettleType_settle_bugang:   true,
		majongpb.SettleType_settle_minggang: true,
	}[settleType]
}

// IsHuSettle 是否是胡结算方式
func IsHuSettle(settleType majongpb.SettleType) bool {
	return map[majongpb.SettleType]bool{
		majongpb.SettleType_settle_dianpao: true,
		majongpb.SettleType_settle_zimo:    true,
	}[settleType]
}

// IsRoundSettle 是否是单局结算方式
func IsRoundSettle(settleType majongpb.SettleType) bool {
	return map[majongpb.SettleType]bool{
		majongpb.SettleType_settle_yell:      true,
		majongpb.SettleType_settle_flowerpig: true,
		majongpb.SettleType_settle_taxrebeat: true,
	}[settleType]
}

// CanInstantSettle 能否立即结算
func CanInstantSettle(settleType majongpb.SettleType, settleOption *mjoption.SettleOption) bool {
	if IsGangSettle(settleType) {
		return settleOption.GangInstantSettle
	} else if IsHuSettle(settleType) {
		return settleOption.HuInstantSettle
	}
	return true
}

// CanRoundSettle 玩家是否可以单局结算
func CanRoundSettle(playerID uint64, huQuitPlayers map[uint64]bool, settleOption *mjoption.SettleOption) bool {
	if huQuitPlayers[playerID] {
		if _, ok := settleOption.HuQuitPlayerCanSettle["huPlayer_can_round_settele"]; !ok {
			return true
		}
		return settleOption.HuQuitPlayerCanSettle["huPlayer_can_round_settele"]
	}
	return true
}
