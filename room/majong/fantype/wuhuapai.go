package fantype

// checkWuHuaPai 无花牌
func checkWuHuaPai(tc *typeCalculator) bool {
	return len(tc.getHuaCards()) == 0
}
