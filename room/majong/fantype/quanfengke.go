package fantype

import (
	"steve/gutils"
)

// checkQuanFengKe 检测圈风刻 与本圈风相同的刻子
func checkQuanFengKe(tc *typeCalculator) bool {
	renNum := len(tc.mjContext.GetPlayers())
	feng := gutils.GetPlayerSeat(renNum, int(tc.mjContext.GetZhuangjiaIndex()))
	for _, combine := range tc.combines {
		for _, ke := range combine.kes {
			if ke == feng {
				return true
			}
		}
	}
	for _, peng := range tc.getPengCards() {
		pengValue := gutils.ServerCard2Number(peng.GetCard())
		if pengValue == uint32(feng) {
			return true
		}
	}
	return false
}
