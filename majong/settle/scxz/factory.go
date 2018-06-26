package scxz

import (
	"steve/majong/interfaces"
	"steve/majong/settle/scxl"
)

// SettlerFactory 四川血战结算器工厂
type SettlerFactory struct{}

// CreateGangSettler 创建杠结算器
func (f *SettlerFactory) CreateGangSettler() interfaces.GangSettle {
	return &scxl.GangSettle{}
}

// CreateHuSettler 创建胡结算器
func (f *SettlerFactory) CreateHuSettler() interfaces.HuSettle {
	return &scxl.HuSettle{}
}

// CreateRoundSettle 创建单局结算器
func (f *SettlerFactory) CreateRoundSettle() interfaces.RoundSettle {
	return &RoundSettle{}
}
