package global

import (
	"steve/common/mjoption"
	majongpb "steve/entity/majong"
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

	// Card1Z 东
	Card1Z = majongpb.Card{Color: majongpb.CardColor_ColorZi, Point: 1}
	// Card2Z 南
	Card2Z = majongpb.Card{Color: majongpb.CardColor_ColorZi, Point: 2}
	// Card3Z 西
	Card3Z = majongpb.Card{Color: majongpb.CardColor_ColorZi, Point: 3}
	// Card4Z 北
	Card4Z = majongpb.Card{Color: majongpb.CardColor_ColorZi, Point: 4}
	// Card5Z 中
	Card5Z = majongpb.Card{Color: majongpb.CardColor_ColorZi, Point: 5}
	// Card6Z 发
	Card6Z = majongpb.Card{Color: majongpb.CardColor_ColorZi, Point: 6}
	// Card7Z 白
	Card7Z = majongpb.Card{Color: majongpb.CardColor_ColorZi, Point: 7}

	// Card1H 春
	Card1H = majongpb.Card{Color: majongpb.CardColor_ColorHua, Point: 1}
	// Card2H 夏
	Card2H = majongpb.Card{Color: majongpb.CardColor_ColorHua, Point: 2}
	// Card3H 秋
	Card3H = majongpb.Card{Color: majongpb.CardColor_ColorHua, Point: 3}
	// Card4H 东
	Card4H = majongpb.Card{Color: majongpb.CardColor_ColorHua, Point: 4}
	// Card5H 梅
	Card5H = majongpb.Card{Color: majongpb.CardColor_ColorHua, Point: 5}
	// Card6H 兰
	Card6H = majongpb.Card{Color: majongpb.CardColor_ColorHua, Point: 6}
	// Card7H 竹
	Card7H = majongpb.Card{Color: majongpb.CardColor_ColorHua, Point: 7}
	// Card8H 菊
	Card8H = majongpb.Card{Color: majongpb.CardColor_ColorHua, Point: 8}
)

func getMjCard(v int) *majongpb.Card {
	switch v {
	case 11:
		return &Card1W
	case 12:
		return &Card2W
	case 13:
		return &Card3W
	case 14:
		return &Card4W
	case 15:
		return &Card5W
	case 16:
		return &Card6W
	case 17:
		return &Card7W
	case 18:
		return &Card8W
	case 19:
		return &Card9W
	case 21:
		return &Card1T
	case 22:
		return &Card2T
	case 23:
		return &Card3T
	case 24:
		return &Card4T
	case 25:
		return &Card5T
	case 26:
		return &Card6T
	case 27:
		return &Card7T
	case 28:
		return &Card8T
	case 29:
		return &Card9T
	case 31:
		return &Card1B
	case 32:
		return &Card2B
	case 33:
		return &Card3B
	case 34:
		return &Card4B
	case 35:
		return &Card5B
	case 36:
		return &Card6B
	case 37:
		return &Card7B
	case 38:
		return &Card8B
	case 39:
		return &Card9B
	case 41:
		return &Card1Z
	case 42:
		return &Card2Z
	case 43:
		return &Card3Z
	case 44:
		return &Card4Z
	case 45:
		return &Card5Z
	case 46:
		return &Card6Z
	case 47:
		return &Card7Z
	case 51:
		return &Card1H
	case 52:
		return &Card2H
	case 53:
		return &Card3H
	case 54:
		return &Card4H
	case 55:
		return &Card5H
	case 56:
		return &Card6H
	case 57:
		return &Card7H
	case 58:
		return &Card8H
	}
	return &majongpb.Card{}
}

// GetOriginCards 获取gameID游戏的所有牌
func GetOriginCards(mjContext *majongpb.MajongContext) []*majongpb.Card {
	xpOption := mjoption.GetXingpaiOption(int(mjContext.GetXingpaiOptionId()))
	result := make([]*majongpb.Card, 0, 200)
	for _, v := range xpOption.WallCards {
		card := getMjCard(v)
		result = append(result, card)
	}
	return result
}
