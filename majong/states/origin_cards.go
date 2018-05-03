package states

import (
	majongpb "steve/server_pb/majong"
)

func getOriginCards(gameID int) []*majongpb.Card {
	cards := []struct {
		card  majongpb.Card
		count int
	}{
		// 万
		{card: Card1W, count: 4},
		{card: Card2W, count: 4},
		{card: Card3W, count: 4},
		{card: Card4W, count: 4},
		{card: Card5W, count: 4},
		{card: Card6W, count: 4},
		{card: Card7W, count: 4},
		{card: Card8W, count: 4},
		{card: Card9W, count: 4},
		// 条
		{card: Card1T, count: 4},
		{card: Card2T, count: 4},
		{card: Card3T, count: 4},
		{card: Card4T, count: 4},
		{card: Card5T, count: 4},
		{card: Card6T, count: 4},
		{card: Card7T, count: 4},
		{card: Card8T, count: 4},
		{card: Card9T, count: 4},
		// 筒
		{card: Card1B, count: 4},
		{card: Card2B, count: 4},
		{card: Card3B, count: 4},
		{card: Card4B, count: 4},
		{card: Card5B, count: 4},
		{card: Card6B, count: 4},
		{card: Card7B, count: 4},
		{card: Card8B, count: 4},
		{card: Card9B, count: 4},
	}

	result := make([]*majongpb.Card, 0, 200)
	for index, cardx := range cards {
		for i := 0; i < cardx.count; i++ {
			result = append(result, &cards[index].card)
		}
	}
	return result
}
