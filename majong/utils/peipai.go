package utils

import (
	"steve/peipai"
	majongpb "steve/server_pb/majong"
	"strconv"
)

// PeiPai 配牌工具
func PeiPai(wallCards []*majongpb.Card, game string) ([]*majongpb.Card, error) {
	value, err := peipai.GetPeiPai(game)
	if err != nil {
		return nil, err
	}
	var cards []*majongpb.Card
	for i := 0; i < len(value); i = i + 3 {
		ca, err := strconv.Atoi(value[i : i+2])
		if err != nil {
			return nil, err
		}
		card, err := IntToCard(int32(ca))
		if err != nil {
			return nil, err
		}
		cards = append(cards, card)
	}
	for i := 0; i < len(cards); i++ {
		for j := len(wallCards) - 1; j >= 0; j-- {
			if cards[i].Point == wallCards[j].Point && cards[i].Color == wallCards[j].Color {
				wallCards[i], wallCards[j] = wallCards[j], wallCards[i]
			}
		}
	}
	return wallCards, nil
}
