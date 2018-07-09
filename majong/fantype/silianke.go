package fantype

// checkSiLianKe 四连刻:胡牌时,含有一种花色 4 副依次递增一位数的刻子;
func checkSiLianKe(tc *typeCalculator) bool {
	pengCount := len(tc.getPengCards())
	keCombines := make([]Combine, 0)
	for _, combine := range tc.combines {
		if len(combine.kes)+pengCount >= 4 {
			keCombines = append(keCombines, combine)
		}
	}
	for _, keCombine := range keCombines {
		colorCount, cardCount := getChiCardsDetails(tc.getChiCards())
		for _, ke := range keCombine.kes {
			keCard := intToCard(ke)
			cardCount[ke] = cardCount[ke] + 1
			colorCount[keCard.Color] = colorCount[keCard.Color] + 1
		}
		hasColor := false
		for _, count := range colorCount {
			if count >= 4 {
				hasColor = true
				break
			}
		}
		hasValue := false
		for _, count := range cardCount {
			if cardCount[count+1] != 0 && cardCount[count+2] != 0 {
				hasValue = true
			}
		}
		if !hasColor || !hasValue {
			return false
		}
	}
	return false
}
