package fantype

import (
	"steve/room/majong/utils"
)

// checkDanDiaoJiang 单钓将:钓单张牌作将成胡,1112 胡 2 算单钓将,1234 胡 1、4 不算单钓将
func checkDanDiaoJiang(tc *typeCalculator) bool {
	huCard := tc.getHuCard()
	huValue := utils.ServerCard2Number(huCard.Card)
	isJiang := false
	for _, combine := range tc.combines {
		if combine.jiang == huValue {
			isJiang = true
			break
		}
	}
	return isJiang
}
