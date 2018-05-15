package interfaces

import (
	majongpb "steve/server_pb/majong"
)

// CardTypeCalculator 牌型计算器
type CardTypeCalculator interface {
	Calculate(params CardCalcParams) (cardTypes []majongpb.CardType, gengCount uint32)
	// CardTypeValue 牌型的倍数,根数
	CardTypeValue(cardTypes []majongpb.CardType, gengCount uint32) (uint32, uint32)
}

// CardCalcParams 计算牌型的参数
type CardCalcParams struct {
	HandCard []*majongpb.Card
	PengCard []*majongpb.Card
	GangCard []*majongpb.Card
	HuCard   *majongpb.Card
	GameID   int
}
