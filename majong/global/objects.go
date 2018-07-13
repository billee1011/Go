package global

import (
	"steve/majong/interfaces"
)

var gFanTypeCalculator interfaces.FantypeCalculator
var gStateFactory interfaces.MajongStateFactory

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
