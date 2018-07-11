package deskbase

import (
	"steve/client_pb/msgid"
	"steve/room/interfaces"
	"steve/room/interfaces/facade"
	"steve/room/interfaces/global"

	"github.com/Sirupsen/logrus"
)

// deskPlayerMgr 牌桌玩家管理器
type deskPlayerMgr struct {
	players    map[uint32]interfaces.DeskPlayer
	enterQuits chan interfaces.PlayerEnterQuitInfo // 退出以及进入信息
}

// CreateDeskPlayerMgr 创建牌桌玩家管理器
func createDeskPlayerMgr() *deskPlayerMgr {
	return &deskPlayerMgr{
		enterQuits: make(chan interfaces.PlayerEnterQuitInfo),
	}
}

// SetPlayers 设置玩家列表
func (dpm *deskPlayerMgr) setPlayers(players []uint64) {
	playerMgr := global.GetPlayerMgr()
	dpm.players = make(map[uint32]interfaces.DeskPlayer, len(players))
	var seat uint32
	for _, playerID := range players {
		player := playerMgr.GetPlayer(playerID)
		var coin uint64
		if player == nil {
			coin = player.GetCoin()
		}
		dpm.players[seat] = createDeskPlayer(playerID, seat, coin, 2) // TODO， 最大超时次数
		seat++
	}
}

// GetDeskPlayers 获取牌桌玩家列表
func (dpm *deskPlayerMgr) GetDeskPlayers() []interfaces.DeskPlayer {
	result := []interfaces.DeskPlayer{}
	for _, deskPlayer := range dpm.players {
		result = append(result, deskPlayer)
	}
	return result
}

// PlayerQuit 玩家退出
func (dpm *deskPlayerMgr) PlayerQuit(playerID uint64) chan struct{} {
	finishChannel := make(chan struct{})
	dpm.enterQuits <- interfaces.PlayerEnterQuitInfo{
		PlayerID:      playerID,
		Quit:          true,
		FinishChannel: finishChannel,
	}
	return finishChannel
}

// PlayerEnter 玩家进入
func (dpm *deskPlayerMgr) PlayerEnter(playerID uint64) chan struct{} {
	finishChannel := make(chan struct{})
	dpm.enterQuits <- interfaces.PlayerEnterQuitInfo{
		PlayerID:      playerID,
		Quit:          false,
		FinishChannel: finishChannel,
	}
	return finishChannel
}

// PlayerEnterQuitChannel 获取玩家进入退出信息通道
func (dpm *deskPlayerMgr) PlayerEnterQuitChannel() <-chan interfaces.PlayerEnterQuitInfo {
	return dpm.enterQuits
}

// BroadcastMessage 向玩家广播消息
func (dpm *deskPlayerMgr) BroadcastMessage(playerIDs []uint64, msgID msgid.MsgID, body []byte, exceptQuit bool) {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name":       "deskPlayerMgr.BroadcastMessage",
		"dest_player_ids": playerIDs,
		"msg_id":          msgID,
	})
	// 是否针对所有玩家
	if playerIDs == nil || len(playerIDs) == 0 {
		playerIDs = facade.GetDeskPlayerIDs(dpm)
		logEntry = logEntry.WithField("all_player_ids", playerIDs)
	}
	playerIDs = dpm.removeQuit(playerIDs)
	logEntry = logEntry.WithField("real_dest_player_ids", playerIDs)

	if len(playerIDs) == 0 {
		return
	}
	facade.BroadCastMessageBare(playerIDs, msgID, body)
	logEntry.Debugln("广播消息")
}

// removeQuit 移除已经退出的玩家
func (dpm *deskPlayerMgr) removeQuit(playerIDs []uint64) []uint64 {
	deskPlayerIDs := map[uint64]bool{}
	deskPlayers := dpm.GetDeskPlayers()
	for _, deskPlayer := range deskPlayers {
		playerID := deskPlayer.GetPlayerID()
		deskPlayerIDs[playerID] = deskPlayer.IsQuit()
	}
	result := []uint64{}
	for _, playerID := range playerIDs {
		if quited, _ := deskPlayerIDs[playerID]; !quited {
			result = append(result, playerID)
		}
	}
	return result
}
