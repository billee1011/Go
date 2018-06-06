package desks

import (
	"sync"
)

type deskPlayer struct {
	playerID uint64
	seat     uint32 // 座号
	quit     bool   // 是否已经退出牌桌

	mu sync.RWMutex
}

func newDeskPlayer(playerID uint64, seat uint32) *deskPlayer {
	return &deskPlayer{
		playerID: playerID,
		seat:     seat,
	}
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

// IsQuit 是否已经退出
func (dp *deskPlayer) IsQuit() bool {
	dp.mu.RLock()
	defer dp.mu.RUnlock()
	return dp.quit
}

// quitDesk 退出牌桌
func (dp *deskPlayer) quitDesk() {
	dp.mu.Lock()
	dp.mu.Unlock()
	dp.quit = true
}

// enterDesk 进入牌桌
func (dp *deskPlayer) enterDesk() {
	dp.mu.Lock()
	dp.mu.Unlock()
	dp.quit = false
}
