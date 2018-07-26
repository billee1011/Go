package fantype

//checkQingPeng 检测清碰 满足清一色，碰碰胡条件
func checkQingPeng(tc *typeCalculator) bool {
	return tc.callCheckFunc(qingyiseFuncID) && tc.callCheckFunc(pengpenghuFuncID)
}
