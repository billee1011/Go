package fantype

//checkErrRenPingHe 检测平和 由4副顺子及序数牌作将牌组成的胡牌
func checkErrRenPingHe(tc *typeCalculator) bool {
	for _, combine := range tc.combines {
		jiangCard := intToCard(combine.jiang)
		if !IsFlowerCard(jiangCard) {
			continue
		}
		if len(combine.kes) == 0 && len(combine.shuns) == 4 {
			return true
		}
	}
	return false
}
