package desk

// DeskSettler 牌桌结算器
type DeskSettler interface {
	// Settle
	Settle(desk *Desk, config *DeskConfig)
	// RoundSettle 单局结算
	RoundSettle(desk *Desk, config *DeskConfig)
	// 获取结算统计信息
	GetStatistics() map[uint64]int64
}
