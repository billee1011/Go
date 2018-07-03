package deskbase

import (
	"steve/room/desks/deskplayer"
	"steve/room/interfaces"
	"steve/room/interfaces/global"
)

// EnterQuitInfo 退出以及进入信息
type EnterQuitInfo struct {
	PlayerID uint64
	Quit     bool // true 为退出， false 为进入
}

// DeskPlayerMgr 牌桌玩家管理器
type DeskPlayerMgr struct {
	players    map[uint32]interfaces.DeskPlayer
	enterQuits chan EnterQuitInfo // 退出以及进入信息
}

// CreateDeskPlayerMgr 创建牌桌玩家管理器
func CreateDeskPlayerMgr() *DeskPlayerMgr {
	return &DeskPlayerMgr{
		enterQuits: make(chan EnterQuitInfo),
	}
}

// SetPlayers 设置玩家列表
func (dpm *DeskPlayerMgr) SetPlayers(players []uint64) {
	playerMgr := global.GetPlayerMgr()
	dpm.players = make(map[uint32]interfaces.DeskPlayer, len(players))
	var seat uint32
	for _, playerID := range players {
		player := playerMgr.GetPlayer(playerID)
		var coin uint64
		if player == nil {
			coin = player.GetCoin()
		}
		dpm.players[seat] = deskplayer.CreateDeskPlayer(playerID, seat, coin, 2) // TODO， 最大超时次数
		seat++
	}
}

// GetDeskPlayers 获取牌桌玩家列表
func (dpm *DeskPlayerMgr) GetDeskPlayers() []interfaces.DeskPlayer {
	result := []interfaces.DeskPlayer{}
	for _, deskPlayer := range dpm.players {
		result = append(result, deskPlayer)
	}
	return result
}

// PlayerQuit 玩家退出
func (dpm *DeskPlayerMgr) PlayerQuit(playerID uint64) {
	dpm.enterQuits <- EnterQuitInfo{
		PlayerID: playerID,
		Quit:     true,
	}
}

// PlayerEnter 玩家进入
func (dpm *DeskPlayerMgr) PlayerEnter(playerID uint64) {
	dpm.enterQuits <- EnterQuitInfo{
		PlayerID: playerID,
		Quit:     false,
	}
}

// GetEnterQuitChan 获取玩家进入退出信息通道
func (dpm *DeskPlayerMgr) GetEnterQuitChan() <-chan EnterQuitInfo {
	return dpm.enterQuits
}
