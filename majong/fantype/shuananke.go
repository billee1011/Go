package fantype

// checkShuanAnKe 胡牌时,含有 2 个暗刻
func checkShuanAnKe(tc *typeCalculator) bool {
	keCombines := make([]combine, 0)
	for _, combine := range tc.combines {
		if len(combine.kes) >= 2 {
			keCombines = append(keCombines, combine)
		}
	}
	for _, keCombine := range keCombines {
		anKeCount := 0
		for _, ke := range keCombine.kes {
			if contains(tc.handCards, ke) {
				anKeCount++
			}
			if anKeCount >= 2 {
				return true
			}
		}

	}
	return false
}
