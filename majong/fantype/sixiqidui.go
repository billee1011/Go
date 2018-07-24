package fantype

import (
	"steve/majong/global"
	"steve/majong/utils"
	majongpb "steve/entity/majong"
)

// checkSiXiQiDui 四喜七对:胡牌为七对,并且包含“东南西北”
func checkSiXiQiDui(tc *typeCalculator) bool {
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
	fengCards := []majongpb.Card{global.Card1Z, global.Card2Z, global.Card3Z, global.Card4Z}
	for _, fengCard := range fengCards {
		cardValue := utils.ServerCard2Number(&fengCard)
		if countMap[cardValue] == 0 {
			return false
		}
	}
	return true
}
