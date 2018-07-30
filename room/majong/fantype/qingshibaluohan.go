package fantype

//checkQingShiBaLuoHan 检测清十八罗汉 满足十八罗汉，清一色条件
func checkQingShiBaLuoHan(tc *typeCalculator) bool {
	return tc.callCheckFunc(qingyiseFuncID) && tc.callCheckFunc(shibaluohanFuncID)
}
