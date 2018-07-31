package fantype

import "steve/room/majong/utils"

// checkKanZhang 坎张 胡 2 张牌之间的牌,4556 胡 5 也为坎张,手中有 45567 胡 6 不算坎张
func checkKanZhang(tc *typeCalculator) bool {
	huCard := tc.getHuCard()

	huValue := utils.ServerCard2Number(huCard.Card)

	for _, combine := range tc.combines {
		for _, shun := range combine.shuns {
			if shun+1 == huValue {
				return true
			}
		}
	}
	return false
}
