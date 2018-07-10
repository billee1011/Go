package fantype

import (
	majongpb "steve/server_pb/majong"
)

// checkZiYiSe 字一色:由字牌组成的胡牌;
func checkZiYiSe(tc *typeCalculator) bool {
	gangCards := tc.getGangCards()
	pengCards := tc.getPengCards()
	chiCards := tc.getChiCards()
	handCards := tc.getHandCards()
	huCard := tc.getHuCard()

	checkCards := make([]*majongpb.Card, 0)

	for _, gangCard := range gangCards {
		checkCards = append(checkCards, gangCard.Card)
	}
	for _, pengCard := range pengCards {
		checkCards = append(checkCards, pengCard.Card)
	}
	for _, chiCard := range chiCards {
		checkCards = append(checkCards, chiCard.Card)
	}
	for _, handCard := range handCards {
		checkCards = append(checkCards, handCard)
	}
	checkCards = append(checkCards, huCard.Card)

	for _, checkCard := range checkCards {
		if !IsNotFlowerCard(checkCard) {
			return false
		}
	}
	return true
}
