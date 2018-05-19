package global

import (
	majongpb "steve/server_pb/majong"
)

var (
	// Card1W 1 万
	Card1W = majongpb.Card{Color: majongpb.CardColor_ColorWan, Point: 1}
	// Card2W 2 万
	Card2W = majongpb.Card{Color: majongpb.CardColor_ColorWan, Point: 2}
	// Card3W 3 万
	Card3W = majongpb.Card{Color: majongpb.CardColor_ColorWan, Point: 3}
	// Card4W 4 万
	Card4W = majongpb.Card{Color: majongpb.CardColor_ColorWan, Point: 4}
	// Card5W 5 万
	Card5W = majongpb.Card{Color: majongpb.CardColor_ColorWan, Point: 5}
	// Card6W 6 万
	Card6W = majongpb.Card{Color: majongpb.CardColor_ColorWan, Point: 6}
	// Card7W 7 万
	Card7W = majongpb.Card{Color: majongpb.CardColor_ColorWan, Point: 7}
	// Card8W 8 万
	Card8W = majongpb.Card{Color: majongpb.CardColor_ColorWan, Point: 8}
	// Card9W 9 万
	Card9W = majongpb.Card{Color: majongpb.CardColor_ColorWan, Point: 9}

	// Card1T 1 条
	Card1T = majongpb.Card{Color: majongpb.CardColor_ColorTiao, Point: 1}
	// Card2T 2 条
	Card2T = majongpb.Card{Color: majongpb.CardColor_ColorTiao, Point: 2}
	// Card3T 3 条
	Card3T = majongpb.Card{Color: majongpb.CardColor_ColorTiao, Point: 3}
	// Card4T 4 条
	Card4T = majongpb.Card{Color: majongpb.CardColor_ColorTiao, Point: 4}
	// Card5T 5 条
	Card5T = majongpb.Card{Color: majongpb.CardColor_ColorTiao, Point: 5}
	// Card6T 6 条
	Card6T = majongpb.Card{Color: majongpb.CardColor_ColorTiao, Point: 6}
	// Card7T 7 条
	Card7T = majongpb.Card{Color: majongpb.CardColor_ColorTiao, Point: 7}
	// Card8T 8 条
	Card8T = majongpb.Card{Color: majongpb.CardColor_ColorTiao, Point: 8}
	// Card9T 9 条
	Card9T = majongpb.Card{Color: majongpb.CardColor_ColorTiao, Point: 9}

	// Card1B 1 筒
	Card1B = majongpb.Card{Color: majongpb.CardColor_ColorTong, Point: 1}
	// Card2B 2 筒
	Card2B = majongpb.Card{Color: majongpb.CardColor_ColorTong, Point: 2}
	// Card3B 3 筒
	Card3B = majongpb.Card{Color: majongpb.CardColor_ColorTong, Point: 3}
	// Card4B 4 筒
	Card4B = majongpb.Card{Color: majongpb.CardColor_ColorTong, Point: 4}
	// Card5B 5 筒
	Card5B = majongpb.Card{Color: majongpb.CardColor_ColorTong, Point: 5}
	// Card6B 6 筒
	Card6B = majongpb.Card{Color: majongpb.CardColor_ColorTong, Point: 6}
	// Card7B 7 筒
	Card7B = majongpb.Card{Color: majongpb.CardColor_ColorTong, Point: 7}
	// Card8B 8 筒
	Card8B = majongpb.Card{Color: majongpb.CardColor_ColorTong, Point: 8}
	// Card9B 9 筒
	Card9B = majongpb.Card{Color: majongpb.CardColor_ColorTong, Point: 9}
)

// GetOriginCards 获取gameID游戏的所有牌
func GetOriginCards(gameID int) []*majongpb.Card {
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
