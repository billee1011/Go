package fantype

// checkQuanQiuRen 全求人:全靠吃牌、碰牌、单钓别人打出的牌胡牌
func checkQuanQiuRen(tc *typeCalculator) bool {
	if len(tc.getGangCards()) != 0 {
		return false
	}
	if len(tc.handCards) == 1 && tc.huCard != nil {
		return true
	}
	return false
}
