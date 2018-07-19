package fantype

import (
	"steve/gutils"
)

// checkJianKe 检测箭刻 由“中发白”三张相同的牌组成的刻子,包括杠
func checkJianKe(tc *typeCalculator) bool {
	cardsAll := getPlayerCardAll(tc)
	cardMap := make(map[int]int)
	for _, card := range cardsAll {
		cardValue := int(gutils.ServerCard2Number(card))
		if cardValue >= gutils.Zhong && cardValue <= gutils.Bai {
			cardMap[cardValue] = cardMap[cardValue] + 1
			if cardMap[cardValue] >= 3 {
				return true
			}
		}
	}
	return false
}
