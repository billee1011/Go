package interfaces

import (
	majongpb "steve/server_pb/majong"
)

// SettlerType 结算类型
type SettlerType uint32

// Settler 结算借口
type Settler interface {
}

// HuSettleParams 胡结算参数
type HuSettleParams struct {
	HuPlayers     []uint64 //胡玩家
	SrcPlayer     uint64   //点炮胡为放炮者的玩家id，自摸为玩家自己
	AllPlayers    []uint64 //所有玩家
	SettleType    majongpb.SettleType
	CardTypes     []majongpb.CardType // 牌型
	CardTypeValue int                 // 牌型倍数
	SettleID      uint64              // 结算信息id
}

// HuSettle 胡结算
type HuSettle interface {
	Settle(params HuSettleParams) []*majongpb.SettleInfo
}

// GangSettleParams 杠结算参数
type GangSettleParams struct {
	GangPlayer uint64            // 杠的玩家
	SrcPlayer  uint64            // 放杠者玩家
	AllPlayers []uint64          // 所有玩家
	GangType   majongpb.GangType // 杠的类型
	SettleID   uint64            // 结算信息id
}

// GangSettle 杠结算
type GangSettle interface {
	Settle(params GangSettleParams) []*majongpb.SettleInfo
}

// RoundSettleParams 单局结算参数
type RoundSettleParams struct {
	FlowerPigPlayers []uint64               // 花猪玩家
	HuPlayers        []uint64               // 胡牌玩家
	TingPlayersInfo  map[uint64]int         // 听玩家及胡牌最大倍数
	NotTingPlayers   []uint64               // 未听玩家,排除花猪玩家
	SettleInfos      []*majongpb.SettleInfo // 历史结算信息
	SettleID         uint64                 // 结算信息id
}

// RoundSettle 单局结算
type RoundSettle interface {
	Settle(params RoundSettleParams) ([]*majongpb.SettleInfo, []uint64)
}
