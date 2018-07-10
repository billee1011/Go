package fantype

// checkSanAnKe 三暗刻:胡牌时,含有 3 个暗刻
func checkSanAnKe(tc *typeCalculator) bool {
	return getPlayerMaxAnKeNum(tc.combines, 3)
}
