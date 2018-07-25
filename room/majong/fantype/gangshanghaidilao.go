package fantype

//checkGangShangHaiDiLao 检测杠上海底捞 杠上开花最后一张牌
func checkGangShangHaiDiLao(tc *typeCalculator) bool {
	return tc.callCheckFunc(gangshangkaihuaFuncID) && tc.callCheckFunc(miaoshouhuichunFuncID)
}
