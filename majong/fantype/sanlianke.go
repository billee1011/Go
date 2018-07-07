package fantype

// checkSanLianKe 三连刻:胡牌时,含有一种花色 3 副依次递增一位数字的刻子
func checkSanLianKe(tc *typeCalculator) bool {
	pengCount := len(tc.getPengCards())
	keCombines := make([]Combine, 0)
	for _, combine := range tc.combines {
		if len(combine.kes)+pengCount >= 3 {
			keCombines = append(keCombines, combine)
		}
	}
	for _, keCombine := range keCombines {
		colorCount, cardCount, _ := getPengCardsDetails(tc.getPengCards())
		for _, ke := range keCombine.kes {
			cardCount[ke] = cardCount[ke] + 1
			kcolor := ke / 10
			colorCount[kcolor] = colorCount[kcolor] + 1
		}
		has3LianKe := false
		for card := range cardCount {
			for i := 0; i < 3; i++ {
				if cardCount[card+i] == 0 {
					break
				}
				if i == 2 {
					has3LianKe = true
				}
			}
		}
		if has3LianKe {
			return true
		}
	}
	return false
}
