package fantype

import (
	"steve/room/majong/utils"
)

// checkQidui 检测七对
func checkQidui(tc *typeCalculator) bool {
	countMap := map[int]int{}
	handCards := tc.getHandCards()
	if len(handCards) != 13 {
		return false
	}
	huCard := tc.getHuCard()
	if huCard == nil {
		return false
	}
	for _, card := range handCards {
		c := utils.ServerCard2Number(card)
		countMap[c] = countMap[c] + 1
	}
	c := utils.ServerCard2Number(huCard.GetCard())
	countMap[c] = countMap[c] + 1

	for _, count := range countMap {
		if count%2 != 0 {
			return false
		}
	}
	return true
}
