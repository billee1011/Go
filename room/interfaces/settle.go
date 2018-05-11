package interfaces

import (
	server_pb "steve/server_pb/majong"
)

// DeskSettler 牌桌结算
type DeskSettler interface {
	// Settle
	Settle(desk Desk, mjContext server_pb.MajongContext)
	// RoundSettle 单局结算
	RoundSettle(desk Desk, mjContext server_pb.MajongContext)
}

// DeskSettlerFactory 牌桌结算工厂
type DeskSettlerFactory interface {
	CreateDeskSettler(gameID int) DeskSettler
}
