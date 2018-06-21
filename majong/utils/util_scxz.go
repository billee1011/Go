package utils

import (
	majongpb "steve/server_pb/majong"

	"github.com/Sirupsen/logrus"
)

//GetNormalPlayersAll 获取正常玩家数组
func GetNormalPlayersAll(players []*majongpb.Player) []*majongpb.Player {
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
			"func_name":            "GetNormalPlayersAll",
			"off_normal_playersID": playersID,
		}).Info("判断玩家是否是正常玩家")
	}()
	return newPlalyers
}

//IsNormalPlayer 判断玩家是否是正常玩家，根据游戏ID(非SCXZ返回true，在SCXZ下是正常状态下返回true)
func IsNormalPlayer(player *majongpb.Player) bool {
	if player.GetXpState() == majongpb.XingPaiState_normal {
		return true
	}
	defer func() {
		logrus.WithFields(logrus.Fields{
			"func_name":    "IsNormalPlayer",
			"playerStatus": player.GetXpState(),
		}).Info("判断玩家是否是正常玩家")
	}()
	return false
}

//GetNextNormalPlayerByID 获取下个正常状态的玩家
func GetNextNormalPlayerByID(players []*majongpb.Player, srcPlayerID uint64) *majongpb.Player {
	log := logrus.WithFields(logrus.Fields{
		"func_name":   "GetNextNormalPlayerByID",
		"srcPlayerID": srcPlayerID,
	})
	curPlayerID, i := srcPlayerID, 0
	for i < 4 {
		palyer := GetNextPlayerByID(players, curPlayerID)
		if palyer.GetXpState() == majongpb.XingPaiState_normal {
			log.WithFields(logrus.Fields{
				"nextPlayerID": palyer.GetPalyerId(),
			}).Infoln("获取下个正常状态的玩家")
			return palyer
		}
		curPlayerID = palyer.GetPalyerId()
	}
	log.Errorln("获取下个正常状态的玩家失败！")
	return nil
}
