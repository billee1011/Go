package fantype

import "steve/majong/utils"

// checkBianZhang 单胡 123 的 3 及 789 的 7 或 1233 胡 3、77879 胡 7 都为张;手中有 12345胡 3,56789 胡 6 不算边张
func checkBianZhang(tc *typeCalculator) bool {
	// 胡牌只在一个顺子里
	huCard := tc.getHuCard()

	huValue := utils.ServerCard2Number(huCard.Card)

	for _, combine := range tc.combines {
		contain := 0
		for _, shun := range combine.shuns {
			if shun == huValue || shun+1 == huValue || shun+2 == huValue {
				contain = contain + 1
			}
			if contain > 1 || contain < 1 {
				return false
			}
		}
		if contain == 1 {
			return true
		}
	}
	return false
}
