package deskplayer

import (
	"steve/room/interfaces"
)

// CreateDeskPlayer 创建牌桌玩家
// maxOverTime : 最大超时次数，超过此次数将会被自动托管
func CreateDeskPlayer(playerID uint64, seat uint32, coin uint64, maxOverTime int) interfaces.DeskPlayer {
	return &deskPlayer{
		playerID:    playerID,
		seat:        seat,
		ecoin:       coin,
		maxOverTime: maxOverTime,
	}
}
