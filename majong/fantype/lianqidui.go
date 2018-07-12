package fantype

import (
	"steve/majong/utils"
	majongpb "steve/server_pb/majong"
)

// checkLianQiDui 连七对:由一种花色序数牌组成序数相连的 7 个对子组成的胡牌;
func checkLianQiDui(tc *typeCalculator) bool {
	countMap := map[int]int{}
	handCards := tc.getHandCards()
	cardColor := handCards[0].GetColor()
	if len(handCards) != 13 {
		return false
	}
	huCard := tc.getHuCard()
	if huCard == nil {
		return false
	}
	minCard := utils.ServerCard2Number(handCards[0])
	for _, card := range handCards {
		if card.GetColor() != cardColor {
			return false
		}
		if card.GetColor() != majongpb.CardColor_ColorWan ||
			card.GetColor() != majongpb.CardColor_ColorTiao || card.GetColor() != majongpb.CardColor_ColorTong {
			return false
		}
		c := utils.ServerCard2Number(card)
		if c <= minCard {
			minCard = c
		}
		countMap[c] = countMap[c] + 1
	}
	c := utils.ServerCard2Number(huCard.GetCard())
	if huCard.GetCard().GetColor() != cardColor {
		return false
	}
	if c <= minCard {
		minCard = c
	}
	countMap[c] = countMap[c] + 1

	for _, count := range countMap {
		if count%2 != 0 {
			return false
		}
	}
	for i := 0; i < 7; i++ {
		if countMap[minCard+i] != 2 {
			return false
		}
	}
	return true
}
