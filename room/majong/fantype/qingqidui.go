package fantype

// checkQingqidui 检测清七对
func checkQingqidui(tc *typeCalculator) bool {
	return tc.callCheckFunc(qiduiFuncID) && tc.callCheckFunc(qingyiseFuncID)
}
