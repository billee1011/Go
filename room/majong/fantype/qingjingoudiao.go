package fantype

//checkQingJinGouDiao 检测清金钩钓 满足金钩钓，清一色条件
func checkQingJinGouDiao(tc *typeCalculator) bool {
	return tc.callCheckFunc(qingyiseFuncID) && tc.callCheckFunc(jingoudiaoFuncID)
}
