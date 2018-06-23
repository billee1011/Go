package utils

import (
	"errors"
	"steve/gutils"
	majongpb "steve/server_pb/majong"

	"github.com/Sirupsen/logrus"
)

var errGameIDNoExist = errors.New("游戏ID不存在")

// GetPlayersByGameID 根据游戏ID，获取玩家数组
func GetPlayersByGameID(gameID int32, players []*majongpb.Player) []*majongpb.Player {
	switch gameID {
	case gutils.SCXLGameID:
		return players
	case gutils.SCXZGameID:
		return GetNormalPlayersAll(players)
	}
	logrus.WithFields(logrus.Fields{
		"func_name": "GetPlayersByGameID",
		"err":       errGameIDNoExist,
	}).Infoln("获取玩家数组失败！")
	return players 
}

// IsPlayerContinueByGameID  根据游戏ID，判断玩家是否可以继续,可以继续返回true
func IsPlayerContinueByGameID(gameID int32, player *majongpb.Player) bool {
	switch gameID {
	case gutils.SCXLGameID:
		return true
	case gutils.SCXZGameID:
		return IsNormalPlayer(player)
	}
	logrus.WithFields(logrus.Fields{
		"func_name": "IsPlayerContinueByGameID",
		"err":       errGameIDNoExist,
	}).Infoln("判断玩家是否可以失败！")
	return true
}

//GetNextPlayerByGameID  根据游戏ID，获取下一个玩家
func GetNextPlayerByGameID(gameID int32, srcPlayerID uint64, players []*majongpb.Player) *majongpb.Player {
	switch gameID {
	case gutils.SCXLGameID:
		return GetNextPlayerByID(players, srcPlayerID)
	case gutils.SCXZGameID:
		return GetNextNormalPlayerByID(players, srcPlayerID)
	}
	logrus.WithFields(logrus.Fields{
		"func_name": "GetNextPlayerByGameID",
		"err":       errGameIDNoExist,
	}).Infoln("获取下一个玩家失败！")
	return GetNextPlayerByID(players, srcPlayerID)
}
