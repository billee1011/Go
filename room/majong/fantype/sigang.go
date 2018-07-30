package fantype

// checkSiGang 四杠:胡牌时,含有 4 个杠(明杠、暗杠);
func checkSiGang(tc *typeCalculator) bool {
	gangCount := len(tc.getGangCards())
	if gangCount == 4 {
		return true
	}
	return false
}
