package utils

import (
	majongpb "steve/server_pb/majong"

	"github.com/Sirupsen/logrus"
)

//GetNormalPlayersAll 获取正常玩家数组
func GetNormalPlayersAll(players []*majongpb.Player) []*majongpb.Player {
	newPlalyers := make([]*majongpb.Player, 0)
	for _, player := range players {
		// 不是正常行牌的玩家，不能检查胡，碰，杠，摸牌。。。
		if player.GetXpState() != majongpb.XingPaiState_normal {
			continue
		}
		newPlalyers = append(newPlalyers, player)
	}
	// 正常的玩家ID
	playerID := make([]uint64, 0)
	for _, player := range newPlalyers {
		playerID = append(playerID, player.GetPalyerId())
	}
	logrus.WithFields(logrus.Fields{"PlayerIDs": playerID}).Info("获取正常玩家数组")
	return newPlalyers
}

//IsNormalPlayer 判断玩家是否是正常玩家，根据游戏ID(非SCXZ返回true，在SCXZ下是正常状态下返回true)
func IsNormalPlayer(player *majongpb.Player) bool {
	flag := player.GetXpState() == majongpb.XingPaiState_normal
	logrus.WithFields(logrus.Fields{"playerStatus": player.GetXpState(),
		"isNormal": flag}).Info("玩家是否是正常")
	return flag
}

//GetNextNormalPlayerByID 获取下个正常状态的玩家
func GetNextNormalPlayerByID(players []*majongpb.Player, srcPlayerID uint64) (nextPalyer *majongpb.Player) {
	curPlayerID, i := srcPlayerID, 0
	for i < 4 {
		nextPalyer = GetNextPlayerByID(players, curPlayerID)
		if nextPalyer.GetXpState() == majongpb.XingPaiState_normal {
			break
		}
		curPlayerID = nextPalyer.GetPalyerId()
	}
	logrus.WithFields(logrus.Fields{"playerID": nextPalyer.GetPalyerId(),
		"playerStatus": nextPalyer.GetXpState()}).Info("获取下个正常状态的玩家")
	return nextPalyer
}

//IsNormalPlayerInsufficient 正常状态的玩家是否人数不够
func IsNormalPlayerInsufficient(players []*majongpb.Player) bool {
	conut := 0
	for _, player := range players {
		if player.GetXpState() == majongpb.XingPaiState_normal {
			conut++
		}
	}
	logrus.WithFields(logrus.Fields{"NormalPlayerConut": conut}).Infoln("正常状态玩家数量")
	return conut <= 1
}
