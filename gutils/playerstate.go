package gutils

import (
	"math"
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
	if !FixNextBankerSeat(mjContext) {
		huPlayerCount := len(huPlayerID)
		if huPlayerCount == 1 {
			mjContext.NextBankerSeat = uint32(GetPlayerIndex(huPlayerID[0], mjContext.GetPlayers()))
		} else if huPlayerCount > 1 {
			mjContext.NextBankerSeat = uint32(GetPlayerIndex(lostPlayerID, mjContext.GetPlayers()))
		} else if huPlayerCount == 0 {
			mjContext.NextBankerSeat = 0
		}
		// mjContext.FixNextBankerSeat = true
	}
}

// FixNextBankerSeat 是否填充了下个庄家
func FixNextBankerSeat(mjContext *majongpb.MajongContext) bool {
	if mjContext.GetNextBankerSeat() == math.MaxUint32 {
		return false
	}
	return true
}
