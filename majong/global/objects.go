package global

import (
	"steve/majong/interfaces"
)

var gCardTypeCalculator interfaces.CardTypeCalculator
var gStateFactory interfaces.MajongStateFactory

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
