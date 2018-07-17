package desk

type DeskSettler interface {
	// Settle
	Settle(desk Desk, config DeskConfig)
	// RoundSettle 单局结算
	RoundSettle(desk Desk, config DeskConfig)
}