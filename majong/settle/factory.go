package settle
/*
功能: 结算工厂类：实现所有麻将的结算实现
作者: Sky
日期: 2018-7-16
 */

import (
	"steve/majong/interfaces"
	"steve/majong/settle/null"
	"steve/majong/settle/majong"
)
var mapSettle  map[int32] *settlerMgr

func init() {
	// 不同子游戏定义不同的结算管理器
	mapSettle = map[int32] *settlerMgr{
		0 : {gangSettle: &null.GangSettle{}, huSettle: &null.HuSettle{}, roundSettle: &null.RoundSettle{}},
		1 : {gangSettle: &majong.GangSettle{}, huSettle: &majong.HuSettle{}, roundSettle: &majong.RoundSettle{}},
		2: {gangSettle: &majong.GangSettle{}, huSettle: &majong.HuSettle{}, roundSettle: &majong.RoundSettle{}},
	}
}


// 结算管理器
type settlerMgr struct {
	gangSettle interfaces.GangSettle 		// 杠结算
	huSettle interfaces.HuSettle			// 胡结算
	roundSettle interfaces.RoundSettle		// 单局结算
}

// Settle
type SettlerFactory struct{}


// CreateGangSettler 创建杠结算器
func (f *SettlerFactory) CreateGangSettler(gameId int32) interfaces.GangSettle {
	return mapSettle[gameId].gangSettle
}

// CreateHuSettler 创建胡结算器
func (f *SettlerFactory) CreateHuSettler(gameId int32) interfaces.HuSettle {
	return mapSettle[gameId].huSettle
}

// CreateRoundSettle 创建单局结算器
func (f *SettlerFactory) CreateRoundSettle(gameId int32) interfaces.RoundSettle {
	return mapSettle[gameId].roundSettle
}
