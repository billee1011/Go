package interfaces

import (
	majongpb "steve/server_pb/majong"
)

// CardTypeCalculator 牌型计算器
type CardTypeCalculator interface {
	Calculate(params CardCalcParams) (cardTypes []CardType, gengCount int)
	// CardTypeValue 牌型的倍数
	CardTypeValue(cardTypes []CardType, gengCount int) int
}

// CardType 卡牌类型
type CardType int

// CardCalcParams 计算牌型的参数
type CardCalcParams struct {
	handCard []*majongpb.Card
	pengCard []*majongpb.Card
	gangCard []*majongpb.Card
	huCard   *majongpb.Card
	gameID   int
}