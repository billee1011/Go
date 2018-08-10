package fantype

import (
	"steve/gutils"
)

// checkMenFengKe 检测门风刻 与本门风相同的刻子
func checkMenFengKe(tc *typeCalculator) bool {
	playerAll := tc.mjContext.GetPlayers()
	renNum, index := len(playerAll), gutils.GetPlayerIndex(tc.getPlayer().GetPlayerId(), playerAll)
	feng := gutils.GetPlayerSeat(renNum, index)

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
