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
	for _, chiCards := range chiCards {
		checkCards = append(checkCards, chiCards.Card)
	}
	for _, handCard := range handCards {
		checkCards = append(checkCards, handCard)
	}
	checkCards = append(checkCards, huCard.Card)

	for _, card := range checkCards {
		if card.Point == 1 || card.Point == 9 {
			return false
		}
		if card.GetColor() == majongpb.CardColor_ColorFeng {
			return false

		}
	}
	return true

}
