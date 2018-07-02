package deskplayer

import (
	"steve/room/interfaces"
	"steve/room/interfaces/global"
)

// CreateDeskPlayer 创建牌桌玩家
func CreateDeskPlayer(playerID uint64, seat uint32) interfaces.DeskPlayer {
	return &deskPlayer{
		playerID: playerID,
		seat:     seat,
		ecoin:    global.GetPlayerMgr().GetPlayer(playerID).GetCoin(),
	}
}
