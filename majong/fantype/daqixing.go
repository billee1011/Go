package fantype

// checkDaQiXing 大七星:胡牌为七对,并且由“东南西北中发白”其中的字牌构成
func checkDaQiXing(tc *typeCalculator) bool {
	if !tc.callCheckFunc(qiduiFuncID) {
		return false
	}
	handCards := tc.getHandCards()
	if len(handCards) != 13 {
		return false
	}
	huCard := tc.getHuCard()
	if huCard == nil {
		return false
	}
	for _, card := range handCards {
		if IsXuShuCard(card) {
			return false
		}
	}
	if IsXuShuCard(huCard.GetCard()) {
		return false
	}
	return true
}
