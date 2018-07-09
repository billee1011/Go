package fantype

// checkSanGang 三杠
func checkSanGang(tc *typeCalculator) bool {
	gangCount := len(tc.getGangCards())
	if gangCount == 3 {
		return true
	}
	return false
}
