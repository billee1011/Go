package null

import (
	"steve/room/majong/interfaces"
	majongpb "steve/entity/majong"
)


// HuSettler 空的胡结算器
type HuSettle struct{}

// Settle 结算
func (s *HuSettle) Settle(params interfaces.HuSettleParams) []*majongpb.SettleInfo {
	return []*majongpb.SettleInfo{}
}

// GangSettler 空的杠结算器
type GangSettle struct{}

// Settle 结算
func (s *GangSettle) Settle(params interfaces.GangSettleParams) *majongpb.SettleInfo {
	return nil
}

// RoundSettler 单局结算器
type RoundSettle struct{}

// Settle 结算
func (s *RoundSettle) Settle(params interfaces.RoundSettleParams) ([]*majongpb.SettleInfo, []uint64) {
	return []*majongpb.SettleInfo{}, []uint64{}
}
