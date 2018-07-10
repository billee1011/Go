package fantype

import (
	"steve/gutils"
)

//checkSiTongShun 检测四同顺 含有一种花色4副相同的顺子,吃也算
func checkSiTongShun(tc *typeCalculator) bool {
	// 不能有碰杠
	if len(tc.getGangCards())+len(tc.getPengCards()) > 0 {
		return false
	}
	// 顺子
	for _, combine := range tc.combines {
		// 刻子为0
		if len(combine.kes) == 0 {
			cardMap := make(map[uint32]int)
			// 吃
			for _, chi := range tc.getChiCards() {
				cardValue := gutils.ServerCard2Number(chi.GetOprCard())
				cardMap[cardValue] = cardMap[cardValue] + 1
			}
			for _, shun := range combine.shuns {
				shunValue := uint32(shun)
				cardMap[shunValue] = cardMap[shunValue] + 1
			}
			if len(cardMap) == 1 {
				return true
			}
		}
	}
	return false
}
