package fantype

//checkQingLongQiDui 检查清龙七对，满足清一色，1根（4张相同的牌），七对条件
func checkQingLongQiDui(tc *typeCalculator) bool {
	return tc.callCheckFunc(qingyiseFuncID) && tc.callCheckFunc(longqiduiFuncID)
}
