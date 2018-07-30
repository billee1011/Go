package fantype

// checkPengpenghu 检测碰碰胡
func checkPengpenghu(tc *typeCalculator) bool {
	if len(tc.getChiCards()) != 0 {
		return false
	}
	for _, combine := range tc.combines {
		if len(combine.shuns) == 0 {
			return true
		}
	}
	return false
}
