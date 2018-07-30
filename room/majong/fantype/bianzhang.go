package fantype

import (
	"steve/gutils"
	"steve/room/majong/utils"
)

// checkBianZhang 单胡 123 的 3 及 789 的 7 或 1233 胡 3、77879 胡 7 都为张;手中有 12345胡 3,56789 胡 6 不算边张
func checkBianZhang(tc *typeCalculator) bool {
	huCard := tc.getHuCard()
	if huCard == nil {
		return false
	}
	tingCardInfo, _ := utils.GetTingCards(tc.getHandCards(), nil)
	if len(tingCardInfo) != 1 {
		return false
	}
	huValue := gutils.ServerCard2Number(huCard.GetCard())
	if huValue != uint32(tingCardInfo[0]) {
		return false
	}
	for _, com := range tc.combines {
		for _, shun := range com.shuns {
			currShun := uint32(shun)
			if currShun == huValue || currShun+2 == huValue {
				return true
			}
		}
	}
	return false
}
