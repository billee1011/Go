package fantype

func checkQingyise(tc *typeCalculator) bool {
	handCards := tc.getHandCards()
	color := handCards[0].Color
	for _, card := range handCards {
		if card.Color != color {
			return false
		}
		if !IsXuShuCard(card) {
			return false
		}
	}
	huCard := tc.getHuCard()
	if huCard != nil && huCard.GetCard().GetColor() != color {
		return false
	}
	if !IsXuShuCard(huCard.GetCard()) {
		return false
	}
	gangCards := tc.getGangCards()
	for _, card := range gangCards {
		if card.GetCard().GetColor() != color {
			return false
		}
		if !IsXuShuCard(card.GetCard()) {
			return false
		}
	}

	pengCards := tc.getPengCards()
	for _, card := range pengCards {
		if card.GetCard().GetColor() != color {
			return false
		}
		if !IsXuShuCard(card.GetCard()) {
			return false
		}
	}

	chiCards := tc.getChiCards()
	for _, chiCard := range chiCards {
		if chiCard.Card.GetColor() != color {
			return false
		}
		if !IsXuShuCard(chiCard.Card) {
			return false
		}
	}
	return true
}
