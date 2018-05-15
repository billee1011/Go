package interfaces

import (
	majongpb "steve/server_pb/majong"
)

// CardTypeCalculator 牌型计算器
type CardTypeCalculator interface {
	Calculate(params CardCalcParams) (cardTypes []CardType, gengCount uint32)
	// CardTypeValue 牌型的倍数
	CardTypeValue(cardTypes []CardType, gengCount uint32) uint32
	// CardGenSum 牌的根数量
	CardGenSum(params CardCalcParams) uint32
}

// CardType 卡牌类型
type CardType int

// CardCalcParams 计算牌型的参数
type CardCalcParams struct {
	HandCard []*majongpb.Card
	PengCard []*majongpb.Card
	GangCard []*majongpb.Card
	HuCard   *majongpb.Card
	GameID   int
}
