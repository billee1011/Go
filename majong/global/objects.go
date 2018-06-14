package global

import (
	"steve/majong/interfaces"
)

var gCardTypeCalculator interfaces.CardTypeCalculator
var gStateFactory interfaces.MajongStateFactory
var gGameSettlerFactory interfaces.GameSettlerFactory

// SetCardTypeCalculator set global card type calculator
func SetCardTypeCalculator(ctc interfaces.CardTypeCalculator) {
	gCardTypeCalculator = ctc
}

// GetCardTypeCalculator get global card type calc
func GetCardTypeCalculator() interfaces.CardTypeCalculator {
	return gCardTypeCalculator
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
