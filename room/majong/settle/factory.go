package settle
/*
功能:
     结算工厂类：实现所有麻将的结算实现。
     定义单独的子游戏目录(比如scxl, scxz,ermj)，然后在新建目录内实现自己的gangSettle.go(杠结算), huSettle.go(胡结算), roundSettle.go(单局结算)。
作者: SkyWang
日期: 2018-7-16
 */

import (
	"steve/room/majong/interfaces"
	"steve/room/majong/settle/majong"
)
// map[游戏Id]*结算管理器
// 不同子游戏定义不同的结算管理器
var mapSettle =  map[int32] *settlerMgr{
	0 : {gangSettle: &majong.GangSettle{}, huSettle: &majong.HuSettle{}, roundSettle: &majong.RoundSettle{}},
	1 : {gangSettle: &majong.GangSettle{}, huSettle: &majong.HuSettle{}, roundSettle: &majong.RoundSettle{}},
	2: {gangSettle: &majong.GangSettle{}, huSettle: &majong.HuSettle{}, roundSettle: &majong.RoundSettle{}},
}

func init() {


}

// 结算管理器
type settlerMgr struct {
	gangSettle interfaces.GangSettle 		// 杠结算
	huSettle interfaces.HuSettle			// 胡结算
	roundSettle interfaces.RoundSettle		// 单局结算
}

// 结算工厂类
type SettlerFactory struct{}


// CreateGangSettler 创建杠结算器
func (f *SettlerFactory) CreateGangSettler(gameId int32) interfaces.GangSettle {
	if t, ok := mapSettle[gameId]; ok {
		return t.gangSettle
	}
	return mapSettle[0].gangSettle
}

// CreateHuSettler 创建胡结算器
func (f *SettlerFactory) CreateHuSettler(gameId int32) interfaces.HuSettle {
	if t, ok := mapSettle[gameId]; ok {
		return t.huSettle
	}
	return mapSettle[0].huSettle

}

// CreateRoundSettle 创建单局结算器
func (f *SettlerFactory) CreateRoundSettle(gameId int32) interfaces.RoundSettle {
	if t, ok := mapSettle[gameId]; ok {
		return t.roundSettle
	}
	return mapSettle[0].roundSettle

}
