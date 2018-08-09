package data

import (
	"encoding/json"
	"strconv"

	"github.com/Sirupsen/logrus"
)

//InterToUint64 接口转uint64
func InterToUint64(param interface{}) uint64 {
	if param == nil {
		return 0
	}
	str := param.(string)
	result, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		logrus.WithFields(logrus.Fields{"func_name": "InterToUint64",
			"param": param}).Infoln("InterToUint64失败")
		return 0
	}
	return result
}

//FmtRobotPlayer 格式化RobotPlayer
func FmtRobotPlayer(robotPlayer *RobotPlayer) map[string]interface{} {
	robotPlayerMap := make(map[string]interface{})
	if robotPlayer.PlayerID > 0 {
		robotPlayerMap[RobotPlayerIDField] = robotPlayer.PlayerID
	}
	if len(robotPlayer.GameIDWinRate) > 0 {
		robotPlayerMap[RobotPlayerGameIDWinRate] = GameIDWinRateToJSON(robotPlayer.GameIDWinRate)
	}
	robotPlayerMap[RobotPlayerStateField] = robotPlayer.State //默认是空闲
	return robotPlayerMap
}

// GameIDWinRateToJSON 游戏对应胜率 转 JSON
func GameIDWinRateToJSON(gameIDWinRate map[uint64]uint64) string {
	if len(gameIDWinRate) == 0 || gameIDWinRate == nil {
		return ""
	}
	str, err := json.Marshal(gameIDWinRate)
	if err != nil {
		logrus.WithFields(logrus.Fields{"func_name": "GameIDWinRateToJSON",
			"gameIDWinRate": gameIDWinRate}).Infoln("游戏对应胜率 转 JSON失败")
	}
	return string(str)
}

// JSONToGameIDWinRate JSON 转 游戏对应胜率
func JSONToGameIDWinRate(gameIDWinRateJSON string) map[uint64]uint64 {
	gameIDWinRate := make(map[uint64]uint64)
	if gameIDWinRateJSON == "" {
		return gameIDWinRate
	}
	giwrbyte := []byte(gameIDWinRateJSON)
	if err := json.Unmarshal(giwrbyte, &gameIDWinRate); err != nil {
		logrus.WithFields(logrus.Fields{"func_name": "JSONToGameIDWinRate",
			"gameIDWinRateJSON": gameIDWinRateJSON}).Infoln("JSON 转 游戏对应胜率失败")
	}
	return gameIDWinRate
}
