package fantype

import (
	"steve/gutils"
)

// checkJianKe 检测箭刻 由“中发白”三张相同的牌组成的刻子
func checkJianKe(tc *typeCalculator) bool {
	// 碰
	for _, peng := range tc.getPengCards() {
		cardValue := gutils.ServerCard2Number(peng.GetCard())
		if cardValue >= gutils.Zhong && cardValue <= gutils.Bai {
			return true
		}
	}
	// 刻
	for _, combine := range tc.combines {
		for _, ke := range combine.kes {
			if ke >= gutils.Zhong && ke <= gutils.Bai {
				return true
			}
		}
	}
	return false
}
