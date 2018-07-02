package deskplayer

import (
	"sync"
)

type deskPlayer struct {
	playerID uint64
	seat     uint32 // 座号
	ecoin    uint64 // 进牌桌金币数
	quit     bool   // 是否已经退出牌桌

	mu sync.RWMutex
}

// GetPlayerID 获取玩家 ID
func (dp *deskPlayer) GetPlayerID() uint64 {
	dp.mu.RLock()
	defer dp.mu.RUnlock()
	return dp.playerID
}

// GetSeat 获取座号
func (dp *deskPlayer) GetSeat() int {
	dp.mu.RLock()
	defer dp.mu.RUnlock()
	return int(dp.seat)
}

// GetEcoin 获取进牌桌金币数
func (dp *deskPlayer) GetEcoin() int {
	dp.mu.RLock()
	defer dp.mu.RUnlock()
	return int(dp.ecoin)
}

// IsQuit 是否已经退出
func (dp *deskPlayer) IsQuit() bool {
	dp.mu.RLock()
	defer dp.mu.RUnlock()
	return dp.quit
}

// QuitDesk 退出牌桌
func (dp *deskPlayer) QuitDesk() {
	dp.mu.Lock()
	dp.mu.Unlock()
	dp.quit = true
}

// EnterDesk 进入牌桌
func (dp *deskPlayer) EnterDesk() {
	dp.mu.Lock()
	dp.mu.Unlock()
	dp.quit = false
}
