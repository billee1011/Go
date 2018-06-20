package utils

import (
	"steve/gutils"
	majongpb "steve/server_pb/majong"

	"github.com/Sirupsen/logrus"
)

//GetNormalPlayersByGameID 获取正常玩家数组
func GetNormalPlayersByGameID(gameID int32, players []*majongpb.Player) []*majongpb.Player {
	// 不是血战直接返回所有玩家
	if gameID != gutils.SCXZGameID {
		return players
	}
	playersID := make([]uint64, 0)
	newPlalyers := make([]*majongpb.Player, 0)
	for _, player := range players {
		// 不等与正常行牌的，不能检查胡，碰，杠，摸牌。。。
		if player.GetXpState() != majongpb.XingPaiState_normal {
			playersID = append(playersID, player.GetPalyerId())
			continue
		}
		newPlalyers = append(newPlalyers, player)
	}
	defer func() {
		logrus.WithFields(logrus.Fields{
			"func_name":            "GetNormalPlayersByGameID",
			"gameID":               gameID,
			"off_normal_playersID": playersID,
		}).Info("判断玩家是否是正常玩家")
	}()
	return newPlalyers
}

//IsNormalPlayerByGameID 判断玩家是否是正常玩家，根据游戏ID(非SCXZ返回fale，在SCXZ下是正常状态下返回fale)
func IsNormalPlayerByGameID(gameID int32, player *majongpb.Player) bool {
	if gameID != gutils.SCXZGameID || player.GetXpState() == majongpb.XingPaiState_normal {
		return true
	}
	defer func() {
		logrus.WithFields(logrus.Fields{
			"func_name":    "IsNormalPlayerByGameID",
			"gameID":       gameID,
			"playerStatus": player.GetXpState(),
		}).Info("判断玩家是否是正常玩家")
	}()
	return false
}

//GetNextNormalPlayerByID 获取下个正常状态的玩家
func GetNextNormalPlayerByID(gameID int32, srcPlayerID uint64, players []*majongpb.Player) uint64 {
	log := logrus.WithFields(logrus.Fields{
		"func_name":   "GetNextNormalPlayerByID",
		"gameID":      gameID,
		"srcPlayerID": srcPlayerID,
	})
	curPlayerID, i := srcPlayerID, 0
	for i < 4 {
		palyer := GetNextPlayerByID(players, curPlayerID)
		if gameID != gutils.SCXZGameID || palyer.GetXpState() == majongpb.XingPaiState_normal {
			log.WithFields(logrus.Fields{
				"nextPlayerID": palyer.GetPalyerId(),
			}).Infoln("获取下个正常状态的玩家")
			return palyer.GetPalyerId()
		}
		curPlayerID = palyer.GetPalyerId()
	}
	log.Errorln("获取下个正常状态的玩家失败！")
	return srcPlayerID
}
