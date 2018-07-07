package fantype

import (
	"steve/majong/utils"
)

// checkSanAnKe 三暗刻:胡牌时,含有 3 个暗刻
func checkSanAnKe(tc *typeCalculator) bool {
	keCombines := make([]Combine, 0)
	for _, combine := range tc.combines {
		if len(combine.kes) >= 3 {
			keCombines = append(keCombines, combine)
		}
	}
	for _, keCombine := range keCombines {
		anKeCount := 0
		for _, ke := range keCombine.kes {
			for _, handCard := range tc.handCards {
				if ke == utils.ServerCard2Number(handCard) {
					anKeCount++
				}
			}
			if anKeCount >= 3 {
				return true
			}
		}

	}
	return false
}
