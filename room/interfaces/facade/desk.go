package facade

import "steve/room/interfaces"

// GetDeskPlayerByID 根据玩家 ID 获取牌桌玩家对象
func GetDeskPlayerByID(d interfaces.DeskPlayerMgr, playerID uint64) interfaces.DeskPlayer {
	players := d.GetDeskPlayers()
	for _, deskPlayer := range players {
		if deskPlayer.GetPlayerID() == playerID {
			return deskPlayer
		}
	}
	return nil
}

// GetDeskPlayerIDs 获取牌桌玩家 ID 列表， 座号作为索引
func GetDeskPlayerIDs(d interfaces.DeskPlayerMgr) []uint64 {
	players := d.GetDeskPlayers()
	result := make([]uint64, len(players))
	for _, player := range players {
		result[player.GetSeat()] = player.GetPlayerID()
	}
	return result
}

// GetTuoguanPlayers 获取牌桌所有托管玩家
func GetTuoguanPlayers(desk interfaces.DeskPlayerMgr) []uint64 {
	players := desk.GetDeskPlayers()
	result := make([]uint64, 0, len(players))
	for _, player := range players {
		if player.IsTuoguan() {
			result = append(result, player.GetPlayerID())
		}
	}
	return result
}
