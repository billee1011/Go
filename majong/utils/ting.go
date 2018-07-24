package utils

import (
	"steve/majong/global"
	majongpb "steve/entity/majong"
)

// GetTingCardNum 获取听牌数量
func GetTingCardNum(mjContext *majongpb.MajongContext, playerID uint64, leftCards []*majongpb.Card,
	c2nMap map[int]uint32, laizis map[Card]bool) (num uint32) {
	tingCards, _ := GetTingCards(leftCards, laizis)
	for _, card := range tingCards {
		num += c2nMap[int(card)]
	}
	return
}

// GetAnCardAndNum 获取未亮牌和数量card2num map
func GetAnCardAndNum(mjContext *majongpb.MajongContext, playerID uint64, lenCard int) map[int]uint32 {
	c2nMap := make(map[int]uint32, len(global.GetOriginCards(mjContext))-lenCard)
	for _, player := range mjContext.GetPlayers() {
		if player.GetPalyerId() == playerID {
			continue
		}
		for _, card := range player.GetHandCards() {
			c2nMap[ServerCard2Number(card)]++
		}
	}

	for _, card := range mjContext.GetWallCards() {
		c2nMap[ServerCard2Number(card)]++
	}

	return c2nMap
}

// CalcTianHuCardNum 计算天胡时胡牌
func CalcTianHuCardNum(mjContext *majongpb.MajongContext, playerID uint64) (tingMax uint32, huCard *majongpb.Card) {
	player := GetPlayerByID(mjContext.GetPlayers(), playerID)
	handCards := player.GetHandCards()
	c2nMap := GetAnCardAndNum(mjContext, playerID, len(handCards))

	var leftCards = make([]*majongpb.Card, len(handCards)-1)
	hcMap := make(map[int]bool, len(handCards))
	for index, card := range handCards {
		if hcMap[ServerCard2Number(card)] == true {
			continue
		}
		hcMap[ServerCard2Number(card)] = true
		copy(leftCards, handCards[0:index])
		copy(leftCards[index:], handCards[index+1:])
		num := GetTingCardNum(mjContext, playerID, leftCards, c2nMap, nil)
		if num > tingMax {
			tingMax = num
			huCard = card
		}
	}
	return
}
