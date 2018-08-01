package robotservice

import (
	"steve/server_pb/robot"

	"github.com/Sirupsen/logrus"
)

// 检验金币和胜率
func checkCoinsWinRtaeRange(coinsRange *robot.CoinsRange, winRateRange *robot.WinRateRange) bool {
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
	return true
}

// 检验状态请求过来的参数
func checkSateArgs(playerID uint64, newState, oldState, serverType int, serverAddr string) bool {
	if playerID < 0 {
		logrus.Warningln("Robot Player ID cannot be less than 0:%v", playerID)
		return false
	}
	robotState := map[robot.RobotPlayerState]bool{
		robot.RobotPlayerState_RPS_IDIE:     true,
		robot.RobotPlayerState_RPS_MATCHING: true,
		robot.RobotPlayerState_RPS_GAMEING:  true,
	}
	if oldState != newState && (!robotState[robot.RobotPlayerState(oldState)] || !robotState[robot.RobotPlayerState(newState)]) {
		logrus.Warningln("Robot Player_state is incorrect, oldState:%d,newState:%d", oldState, newState)
		return false
	}

	robotServerType := map[robot.ServerType]bool{
		robot.ServerType_ST_GATE:  true,
		robot.ServerType_ST_MATCH: true,
		robot.ServerType_ST_ROOM:  true,
	}

	if !robotServerType[robot.ServerType(serverType)] {
		logrus.Warningln("server_type is incorrect, server_type:%d", serverType)
		return false
	}
	if len(serverAddr) == 0 {
		logrus.Warningln("server_addr is empty, server_addr:%d", serverAddr)
		return false
	}
	return true
}
