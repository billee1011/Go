package robotservice

import (
	"steve/server_pb/robot"

	"github.com/Sirupsen/logrus"
)

var robotState = map[robot.RobotPlayerState]bool{
	robot.RobotPlayerState_RPS_IDIE:     true,
	robot.RobotPlayerState_RPS_MATCHING: true,
	robot.RobotPlayerState_RPS_GAMEING:  true,
}

// 检验金币和胜率
func checkGetLeisureRobotArgs(coinsRange *robot.CoinsRange, winRateRange *robot.WinRateRange, newState int) bool {
	switch {
	case coinsRange.High < coinsRange.Low:
		logrus.Warningln("coinsRange:High Must be greater than or equal to Low")
		return false
	case coinsRange.High < 0 || coinsRange.Low < 0:
		logrus.Warningln("coinsRange:Both High and Low must be greater than or equal to 0")
		return false
	case winRateRange.High < winRateRange.Low:
		logrus.Warningln("winRateRange:High Must be greater than or equal to Low")
		return false
	case winRateRange.High < 0 || winRateRange.Low < 0:
		logrus.Warningln("winRateRange:Both High and Low must be greater than or equal to 0")
		return false
	}
	if robot.RobotPlayerState(newState) == robot.RobotPlayerState_RPS_IDIE || !robotState[robot.RobotPlayerState(newState)] {
		logrus.Warningln("Robot Player_state is incorrect, newState:%d", newState)
		return false
	}
	return true
}

// 检验状态请求过来的参数
func checkSateArgs(playerID int64, newState, oldState int) bool {
	if playerID < 0 {
		logrus.Warningln("Robot Player ID cannot be less than 0:%v", playerID)
		return false
	}
	if oldState != newState && (!robotState[robot.RobotPlayerState(oldState)] || !robotState[robot.RobotPlayerState(newState)]) {
		logrus.Warningln("Robot Player_state is incorrect, oldState:%d,newState:%d", oldState, newState)
		return false
	}
	return true
}
