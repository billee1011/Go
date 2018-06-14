package null

import (
	"steve/majong/interfaces"
	majongpb "steve/server_pb/majong"
)

// SettlerFactory 空结算器工厂
type SettlerFactory struct{}

// CreateGangSettler 创建杠结算器
func (f *SettlerFactory) CreateGangSettler() interfaces.GangSettle {
	return &GangSettler{}
}

// CreateHuSettler 创建胡结算器
func (f *SettlerFactory) CreateHuSettler() interfaces.HuSettle {
	return &HuSettler{}
}

// CreateRoundSettle 创建单局结算器
func (f *SettlerFactory) CreateRoundSettle() interfaces.RoundSettle {
	return &RoundSettler{}
}

// HuSettler 空的胡结算器
type HuSettler struct{}

// Settle 结算
func (s *HuSettler) Settle(params interfaces.HuSettleParams) []*majongpb.SettleInfo {
	return []*majongpb.SettleInfo{}
}

// GangSettler 空的杠结算器
type GangSettler struct{}

// Settle 结算
func (s *GangSettler) Settle(params interfaces.GangSettleParams) *majongpb.SettleInfo {
	return nil
}

// RoundSettler 单局结算器
type RoundSettler struct{}

// Settle 结算
func (s *RoundSettler) Settle(params interfaces.RoundSettleParams) ([]*majongpb.SettleInfo, []uint64) {
	return []*majongpb.SettleInfo{}, []uint64{}
}
