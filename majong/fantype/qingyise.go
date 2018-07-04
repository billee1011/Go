package fantype

func checkQingyise(tc *typeCalculator) bool {
	handCards := tc.getHandCards()
	color := handCards[0].Color
	for _, card := range handCards {
		if card.Color != color {
			return false
		}
	}
	huCard := tc.getHuCard()
	if huCard != nil && huCard.GetCard().GetColor() != color {
		return false
	}

	gangCards := tc.getGangCards()
	for _, card := range gangCards {
		if card.GetCard().GetColor() != color {
			return false
		}
	}

	pengCards := tc.getPengCards()
	for _, card := range pengCards {
		if card.GetCard().GetColor() != color {
			return false
		}
	}

	chiCards := tc.getChiCards()
	for _, card := range chiCards {
		if card.GetCard().GetColor() != color {
			return false
		}
	}
	return true
}
