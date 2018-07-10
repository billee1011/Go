package fantype

import (
	"steve/majong/utils"
)

// checkSiGuiYi 4 张相同的牌归于一家的顺、刻子、对、将牌中(不包括杠牌)
func checkSiGuiYi(tc *typeCalculator) bool {
	pengCards := tc.getPengCards()
	chiCards := tc.getChiCards()
	handCards := tc.getHandCards()
	huCard := tc.getHuCard()

	cardCount := make(map[int]int)
	cardValue := 0
	for _, pengCard := range pengCards {
		cardValue = utils.ServerCard2Number(pengCard.Card)
		cardCount[cardValue] = cardCount[cardValue] + 3
	}
	for _, chiCard := range chiCards {
		cardValue = utils.ServerCard2Number(chiCard.Card)
		cardCount[cardValue] = cardCount[cardValue] + 1
		cardCount[cardValue+1] = cardCount[cardValue+1] + 1
		cardCount[cardValue+2] = cardCount[cardValue+2] + 13
	}
	for _, handCard := range handCards {
		cardValue = utils.ServerCard2Number(handCard)
		cardCount[cardValue] = cardCount[cardValue] + 1
	}
	cardValue = utils.ServerCard2Number(huCard.Card)
	cardCount[cardValue] = cardCount[cardValue] + 1
	for _, count := range cardCount {
		if count == 4 {
			return true
		}
	}
	return false

}
