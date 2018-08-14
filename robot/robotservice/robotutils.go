package robotservice

import (
	"steve/server_pb/robot"

	"github.com/Sirupsen/logrus"
)

// 检验金币和胜率
func checkGetLeisureRobotArgs(coinsRange *robot.CoinsRange, winRateRange *robot.WinRateRange) bool {
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
