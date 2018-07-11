package fantype

import majongpb "steve/server_pb/majong"

// checkDuanYao 断幺
func checkDuanYao(tc *typeCalculator) bool {
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
		checkCards = append(checkCards, chiCard.OprCard)
	}
	for _, handCard := range handCards {
		checkCards = append(checkCards, handCard)
	}
	checkCards = append(checkCards, huCard.Card)

	for _, card := range checkCards {
		if card.Point == 1 || card.Point == 9 {
			return false
		}
		if !IsFlowerCard(card) {
			return false

		}
	}
	return true

}
