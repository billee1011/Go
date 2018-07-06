package gutils

import (
	"steve/common/mjoption"
	majongpb "steve/server_pb/majong"

	"github.com/Sirupsen/logrus"
)

// GetMajongPlayer 从 MajongContext 中根据玩家 ID 获取玩家
func GetMajongPlayer(playerID uint64, mjContext *majongpb.MajongContext) *majongpb.Player {
	for _, player := range mjContext.GetPlayers() {
		if player.GetPalyerId() == playerID {
			return player
		}
	}
	return nil
}

// GetPlayerIndex 获取玩家索引
func GetPlayerIndex(playerID uint64, players []*majongpb.Player) int {
	for index, player := range players {
		if player.GetPalyerId() == playerID {
			return index
		}
	}
	return -1
}

// GetPlayerAndIndex 获取玩家索引
func GetPlayerAndIndex(playerID uint64, players []*majongpb.Player) (int, *majongpb.Player) {
	for index, player := range players {
		if player.GetPalyerId() == playerID {
			return index, player
		}
	}
	return -1, nil
}

// IsPlayerContinue   玩家的状态在麻将可行牌数组中包含则返回true
func IsPlayerContinue(playerStater majongpb.XingPaiState, mjContext *majongpb.MajongContext) bool {
	// 麻将不可行牌数组
	xpOption := mjoption.GetXingpaiOption(int(mjContext.GetXingpaiOptionId()))
	flag := true
	for _, state := range xpOption.PlayerStates {
		if int(state) == int(playerStater) {
			flag = false
			break
		}
	}
	logrus.WithFields(logrus.Fields{"playerStater": playerStater,
		"canNotXpStates": xpOption.PlayerStates, "isCanXp": flag}).Info("判断玩家是否可以继续")
	return flag
}
