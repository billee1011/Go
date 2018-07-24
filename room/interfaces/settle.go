package interfaces

import (
	server_pb "steve/entity/majong"
)

// DeskSettler 牌桌结算
type DeskSettler interface {
	// Settle
	Settle(desk Desk, mjContext server_pb.MajongContext)
	// RoundSettle 单局结算
	RoundSettle(desk Desk, mjContext server_pb.MajongContext)
	// GetStatistics 获取当前的统计数据，每个玩家赢了多少钱
	GetStatistics() map[uint64]int64
}

// DeskSettlerFactory 牌桌结算工厂
type DeskSettlerFactory interface {
	CreateDeskSettler(gameID int) DeskSettler
}
