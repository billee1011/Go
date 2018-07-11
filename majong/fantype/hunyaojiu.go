package fantype

import "steve/gutils"

// checkHunYaoJiu 检测混幺九 由序数牌1,9和字牌的刻子，将牌组成
func checkHunYaoJiu(tc *typeCalculator) bool {
	// 吃，杠数量
	chiGangNum := len(tc.getChiCards()) + len(tc.getGangCards())
	if chiGangNum != 0 {
		return false
	}
	// 幺九只能是刻子和将
	for _, combine := range tc.combines {
		if len(combine.shuns) != 0 {
			continue
		}
		jiang := combine.jiang
		if !isYaoJiuByInt(jiang) {
			return false
		}
		for _, ke := range combine.kes {
			if !isYaoJiuByInt(ke) {
				return false
			}
		}
		return true
	}
	return false
}

//isYaoJiuByInt 判断是否是幺九(1,9,字)
func isYaoJiuByInt(card int) bool {
	if card < gutils.Dong {
		cardValue := card % 10
		if cardValue > 1 && cardValue < 9 {
			return false
		}
	}
	return true
}
