package bus
/*
 功能： 系统管理总线， 包括程序内部所有的模块管理器接口【单件模式】。
 作者： SkyWang
 日期： 2018-7-18
 */
import (
	"steve/room/majong/interfaces"
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
