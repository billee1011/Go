package global

import (
	"steve/majong/interfaces"
)

var gFanTypeCalculator interfaces.FantypeCalculator
var gStateFactory interfaces.MajongStateFactory
var gGameSettlerFactory interfaces.GameSettlerFactory

// SetFanTypeCalculator set global fan type calculator
func SetFanTypeCalculator(ctc interfaces.FantypeCalculator) {
	gFanTypeCalculator = ctc
}

// GetFanTypeCalculator get global fan type calc
func GetFanTypeCalculator() interfaces.FantypeCalculator {
	return gFanTypeCalculator
}

// SetMajongStateFacotry set majong state factory
func SetMajongStateFacotry(f interfaces.MajongStateFactory) {
	gStateFactory = f
}

// GetMajongStateFactory get majong state factory
func GetMajongStateFactory() interfaces.MajongStateFactory {
	return gStateFactory
}

// SetGameSettlerFactory 设置游戏结算器工厂
func SetGameSettlerFactory(f interfaces.GameSettlerFactory) {
	gGameSettlerFactory = f
}

func GetGameSettlerFactory() interfaces.GameSettlerFactory {
	return gGameSettlerFactory
}
