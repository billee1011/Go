package fantype

import "steve/majong/utils"

// checkRenHu 人胡
func checkRenHu(tc *typeCalculator) bool {
	mjContext := tc.mjContext
	zhuangJia := mjContext.Players[mjContext.ZhuangjiaIndex]

	if len(zhuangJia.GangCards) != 0 {
		return false
	}

	if len(zhuangJia.OutCards) != 1 {
		return false
	}

	for _, player := range mjContext.Players {
		if player.PalyerId == tc.playerID || player.PalyerId == zhuangJia.PalyerId {
			continue
		}
		if len(player.OutCards) != 0 {
			return false
		}
	}
	outValue := utils.ServerCard2Number(zhuangJia.OutCards[0])
	huValue := utils.ServerCard2Number(tc.getHuCard().Card)

	if outValue != huValue {
		return false
	}

	return true
}
