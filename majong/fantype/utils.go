package fantype

import majongpb "steve/server_pb/majong"

// int 转 Card
func intToCard(cardInt int) *majongpb.Card {
	var color majongpb.CardColor
	switch cardInt / 10 {
	case 1:
		color = majongpb.CardColor_ColorWan
	case 2:
		color = majongpb.CardColor_ColorTiao
	case 3:
		color = majongpb.CardColor_ColorTong
	case 4:
		color = majongpb.CardColor_ColorFeng
	}
	point := int32(cardInt % 10)
	card := &majongpb.Card{
		Color: color,
		Point: point,
	}
	return card
}

func intsToCards(cardInts []int) []*majongpb.Card {
	newCard := make([]*majongpb.Card, 0, len(cardInts))
	for _, card := range cardInts {
		newCard = append(newCard, intToCard(card))
	}
	return newCard
}
