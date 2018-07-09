package fantype

import (
	"steve/majong/utils"
	majongpb "steve/server_pb/majong"
)

// checkSiGuiYi 4 张相同的牌归于一家的顺、刻子、对、将牌中(不包括杠牌)
func checkSiGuiYi(tc *typeCalculator) bool {
	pengCards := tc.getPengCards()
	chiCards := tc.getChiCards()
	handCards := tc.getHandCards()
	huCard := tc.getHuCard()

	checkCards := make([]*majongpb.Card, 0)

	for _, pengCard := range pengCards {
		checkCards = append(checkCards, pengCard.Card)
	}
	for _, chiCards := range chiCards {
		checkCards = append(checkCards, chiCards.Card)
	}
	for _, handCard := range handCards {
		checkCards = append(checkCards, handCard)
	}
	checkCards = append(checkCards, huCard.Card)

	cardCount := make(map[int]int)
	for _, card := range checkCards {
		cardValue := utils.ServerCard2Number(card)
		cardCount[cardValue] = cardCount[cardValue] + 1
		if cardCount[cardValue] == 4 {
			return true
		}
	}
	return false

}
