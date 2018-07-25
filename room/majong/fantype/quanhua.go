package fantype

// checkQuanHua 检测全花
func checkQuanHua(tc *typeCalculator) bool {
	return len(tc.getHuaCards()) == 8
}
