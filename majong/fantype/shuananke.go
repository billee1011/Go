package fantype

// checkShuanAnKe 胡牌时,含有 2 个暗刻
func checkShuanAnKe(tc *typeCalculator) bool {
	return getPlayerMaxAnKeNum(tc.combines, 2)
}
