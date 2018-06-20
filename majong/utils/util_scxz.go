package utils

import (
	"steve/gutils"
	majongpb "steve/server_pb/majong"

	"github.com/Sirupsen/logrus"
)

//GetCanCheckPlayerByGameID 获取能检测胡玩家，根据游戏ID
func GetCanCheckPlayerByGameID(gameID int32, players []*majongpb.Player) []*majongpb.Player {
	if gameID != gutils.SCXZGameID {
		return players
	}
	playersID := make([]uint64, 0)
	newPlalyers := make([]*majongpb.Player, 0)
	for _, player := range players {
		// 不等与正常行牌的，不能检查胡，碰，杠，摸牌。。。
		if player.GetPlayerState() != majongpb.PlayerState_normal {
			playersID = append(playersID, player.GetPalyerId())
			continue
		}
		newPlalyers = append(newPlalyers, player)
	}
	defer func() {
		logrus.WithFields(logrus.Fields{
			"func_name":            GetCanCheckPlayerByGameID,
			"gameID":               gameID,
			"off_normal_playersID": playersID,
		}).Info("判断玩家是否是正常玩家")
	}()
	return newPlalyers
}

//IsNormalPlayerByGameID 判断玩家是否是正常玩家，根据游戏ID(非SCXZ返回fale，在SCXZ下是正常状态下返回fale)
func IsNormalPlayerByGameID(gameID int32, player *majongpb.Player) bool {
	if gameID != gutils.SCXZGameID || player.GetPlayerState() == majongpb.PlayerState_normal {
		return false
	}
	defer func() {
		logrus.WithFields(logrus.Fields{
			"func_name":    IsNormalPlayerByGameID,
			"gameID":       gameID,
			"playerStatus": player.GetPlayerState(),
		}).Info("判断玩家是否是正常玩家")
	}()
	return true
}

//GetNextNormalPlayerByID 获取下个正常状态的玩家
func GetNextNormalPlayerByID(gameID int32, srcPlayerID uint64, players []*majongpb.Player) uint64 {
	palyer := GetNextPlayerByID(players, srcPlayerID)
	if gameID != gutils.SCXZGameID {
		return palyer.GetPalyerId()
	}
	for palyer.GetPlayerState() != majongpb.PlayerState_normal {
		palyer = GetNextPlayerByID(players, srcPlayerID)
	}
	defer func() {
		logrus.WithFields(logrus.Fields{
			"func_name":   GetNextNormalPlayerByID,
			"gameID":      gameID,
			"srcPlayerID": srcPlayerID,
			"nextPlayer":  palyer.GetPalyerId(),
		}).Info("获取下个正常状态的玩家")
	}()
	return palyer.GetPalyerId()
}
