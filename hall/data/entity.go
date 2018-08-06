package data

import (
	"steve/entity/cache"
	"strconv"
)

// PlayerState 玩家状态
type PlayerState struct {
	PlayerID  uint64
	State     uint32
	GameID    uint32
	IPAddr    string
	GateAddr  string
	MatchAddr string
	RoomAddr  string
}

func (pState *PlayerState) generatePlayerState(info map[string]string) {
	// 游戏状态
	state, _ := strconv.ParseUint(info[cache.GameState], 10, 64)
	pState.State = uint32(state)
	// 游戏状态
	gameID, _ := strconv.ParseUint(info[cache.GameID], 10, 64)
	pState.GameID = uint32(gameID)
	// ip地址
	pState.IPAddr = info[cache.IPAddr]
	// 网关服地址
	pState.GateAddr = info[cache.GateAddr]
	// 匹配服地址
	pState.MatchAddr = info[cache.MatchAddr]
	// 房间服地址
	pState.RoomAddr = info[cache.RoomAddr]
}
