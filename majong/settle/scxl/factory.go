package scxl

import "steve/majong/interfaces"

// SettlerFactory 四川血流结算器工厂
type SettlerFactory struct{}

// CreateGangSettler 创建杠结算器
func (f *SettlerFactory) CreateGangSettler() interfaces.GangSettle {
	return &GangSettle{}
}

// CreateHuSettler 创建胡结算器
func (f *SettlerFactory) CreateHuSettler() interfaces.HuSettle {
	return &HuSettle{}
}

// CreateRoundSettle 创建单局结算器
func (f *SettlerFactory) CreateRoundSettle() interfaces.RoundSettle {
	return &RoundSettle{}
}
