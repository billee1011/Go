package gutils

import (
	"steve/client_pb/room"
	majongpb "steve/server_pb/majong"
)

//

// IsTing 玩家是否是听的状态
func IsTing(player *majongpb.Player) bool {
	return player.GetTingStateInfo().GetIsTing()
}

// GetTingType 获取玩家听的类型
func GetTingType(player *majongpb.Player) (tingType room.TingType) {
	tingState := player.GetTingStateInfo()
	if tingState.GetIsTing() {
		tingType = room.TingType_TT_NORMAL_TING
	}
	if tingState.GetIsTianting() {
		tingType = room.TingType_TT_TIAN_TING
	}
	return tingType
}

// IsHu 玩家是否时胡的状态
func IsHu(player *majongpb.Player) bool {
	if len(player.GetHuCards()) > 0 {
		return true
	}
	return false
}

// GetZixunPlayer 获取当前自询的玩家
func GetZixunPlayer(mjContext *majongpb.MajongContext) uint64 {
	zxType := mjContext.GetZixunType()
	switch zxType {
	case majongpb.ZixunType_ZXT_PENG:
		return mjContext.GetLastPengPlayer()
	case majongpb.ZixunType_ZXT_CHI:
		return mjContext.GetLastChiPlayer()
	default:
		return mjContext.GetLastMopaiPlayer()
	}
}

// SetNextZhuangIndex 设置续局庄家Index
func SetNextZhuangIndex(huPlayerID []uint64, lostPlayerID uint64, mjContext *majongpb.MajongContext) {
	huPlayerCount := len(huPlayerID)
	if !mjContext.GetFixNextBankerSeat() {
		if huPlayerCount == 1 {
			mjContext.NextBankerSeat = uint32(GetPlayerIndex(huPlayerID[0], mjContext.GetPlayers()))
		} else if huPlayerCount > 1 {
			mjContext.NextBankerSeat = uint32(GetPlayerIndex(lostPlayerID, mjContext.GetPlayers()))
		}
	}
}
