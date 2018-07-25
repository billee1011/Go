package fantype

import majongpb "steve/entity/majong"

func checkQingyise(tc *typeCalculator) bool {
	handCards := tc.getHandCards()
	color := handCards[0].Color

	checkCards := make([]*majongpb.Card, 0)

	checkCards = append(checkCards, handCards...)
	checkCards = append(checkCards, tc.getHuCard().GetCard())

	for _, gangCard := range tc.getGangCards() {
		checkCards = append(checkCards, gangCard.Card)
	}
	for _, pengCard := range tc.getPengCards() {
		checkCards = append(checkCards, pengCard.Card)
	}
	for _, chiCard := range tc.getChiCards() {
		checkCards = append(checkCards, chiCard.Card)
	}
	for _, card := range checkCards {
		if card.Color != color {
			return false
		}
		if !IsXuShuCard(card) {
			return false
		}
	}
	return true
}
