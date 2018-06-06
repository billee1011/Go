package desks

import "steve/room/interfaces"

type tuoGuanPlayer struct {
	overTimerCount int  // 超时计数
	tuoGuaning     bool // 是否在托管中
}

type tuoGuanMgr struct {
	players      map[uint64]*tuoGuanPlayer
	maxOverTimer int //	最大超时次数，超过此次数则进入托管状态
}

// getTuoGuanPlayers 获取托管玩家
func (tg *tuoGuanMgr) getTuoGuanPlayers() []uint64 {
	result := []uint64{}
	for playerID, player := range tg.players {
		if player.tuoGuaning {
			result = append(result, playerID)
		}
	}
	return result
}

// afterPlayerEvent 处理完成玩家事件
func (tg *tuoGuanMgr) afterPlayerEvent(playerID uint64, eventType interfaces.EventType) {
	if eventType == interfaces.TuoGuanEvent {
		tg.afterOverTimerEvent(playerID)
		return
	}
	tg.afterNormalEvent(playerID)
}

// cancelTuoguan 取消托管
func (tg *tuoGuanMgr) cancelTuoguan(playerID uint64) {
	player, exist := tg.players[playerID]
	if !exist {
		return
	}
	player.tuoGuaning = false
	player.overTimerCount = 0 // TODO: 是否要清零
}

// afterOverTimerEvent 处理完成超时事件
func (tg *tuoGuanMgr) afterOverTimerEvent(playerID uint64) {
	player, exist := tg.players[playerID]
	if !exist {
		player = &tuoGuanPlayer{}
		tg.players[playerID] = player
	}
	player.overTimerCount++
	if player.overTimerCount >= tg.maxOverTimer && !player.tuoGuaning {
		player.tuoGuaning = true
	}
}

// afterNormalEvent 处理完成其他事件
func (tg *tuoGuanMgr) afterNormalEvent(playerID uint64) {
	player, exist := tg.players[playerID]
	if !exist {
		return
	}
	player.overTimerCount = 0 // TODO 是否要清零
	if player.tuoGuaning {
		player.tuoGuaning = false
	}
}
