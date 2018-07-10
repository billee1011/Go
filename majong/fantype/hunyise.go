package fantype

import majongpb "steve/server_pb/majong"

// checkHunYiSe 检测混一色
func checkHunYiSe(tc *typeCalculator) bool {
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
		checkCards = append(checkCards, chiCards.OprCard)
	}
	for _, handCard := range handCards {
		checkCards = append(checkCards, handCard)
	}
	checkCards = append(checkCards, huCard.Card)

	cardColor := majongpb.CardColor(-1)
	for _, card := range checkCards {
		if IsNotFlowerCard(card) {
			continue
		}
		if cardColor == -1 {
			cardColor = card.Color
		} else if cardColor != card.Color {
			return false
		}
	}
	return true
}
