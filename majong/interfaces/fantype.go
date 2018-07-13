package interfaces

import majongpb "steve/server_pb/majong"

// FantypeCalculator 番型计算器
type FantypeCalculator interface {
	Calculate(params FantypeParams) (fanTypes []int, gengCount int, huaCount int)
	// CardTypeValue 牌型的倍数,根数
	CardTypeValue(mjContext *majongpb.MajongContext, fanTypes []int, gengCount int, huaCount int) uint64
}

// FantypeParams 计算番型的参数
type FantypeParams struct {
	PlayerID  uint64
	MjContext *majongpb.MajongContext
	HandCard  []*majongpb.Card
	PengCard  []*majongpb.Card
	GangCard  []*majongpb.GangCard
	HuCard    *majongpb.HuCard
	GameID    int
}
