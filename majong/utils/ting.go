package utils

import (
	"fmt"
	"steve/majong/global"
	majongpb "steve/server_pb/majong"
)

// GetTingCardNum 获取听牌数量
func GetTingCardNum(mjContext *majongpb.MajongContext, playerID uint64, leftCards []Card,
	c2nMap map[Card]uint32, laizis map[Card]bool) (num uint32) {
	tingCards := FastCheckTingV2(leftCards, laizis)
	for _, card := range tingCards {
		num += c2nMap[card]
	}
	return
}

// GetAnCardAndNum 获取未亮牌和数量card2num map
func GetAnCardAndNum(mjContext *majongpb.MajongContext, playerID uint64, lenCard int) map[Card]uint32 {
	c2nMap := make(map[Card]uint32, len(global.GetOriginCards(int(mjContext.GetGameId())))-lenCard)
	for _, player := range mjContext.GetPlayers() {
		if player.GetPalyerId() == playerID {
			continue
		}
		cards := CardsToUtilCards(player.GetHandCards())
		for _, card := range cards {
			c2nMap[card]++
		}
	}

	cards := CardsToUtilCards(mjContext.GetWallCards())
	for _, card := range cards {
		c2nMap[card]++
	}

	return c2nMap
}

// CalcTianHuCardNum 计算天胡时胡牌
func CalcTianHuCardNum(mjContext *majongpb.MajongContext, playerID uint64) (tingMax uint32, huCard *majongpb.Card) {
	player := GetPlayerByID(mjContext.GetPlayers(), playerID)
	handCards := CardsToUtilCards(player.GetHandCards())
	c2nMap := GetAnCardAndNum(mjContext, playerID, len(handCards))

	var huCardInt Card
	var leftCards = make([]Card, len(handCards)-1)
	hcMap := make(map[Card]bool, len(handCards))
	for index, card := range handCards {
		if hcMap[card] == true {
			continue
		}
		hcMap[card] = true
		copy(leftCards, handCards[0:index])
		copy(leftCards[index:], handCards[index+1:])
		num := GetTingCardNum(mjContext, playerID, leftCards, c2nMap, nil)
		if num > tingMax {
			tingMax = num
			huCardInt = card
		}
	}

	huCard, err := IntToCard(int32(huCardInt))
	if err != nil {
		fmt.Println("转换成花色失败", huCardInt)
		return 0, nil
	}
	return
}
