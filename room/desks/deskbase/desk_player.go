package deskbase

import (
	msgid "steve/client_pb/msgid"
	"steve/client_pb/room"
	"steve/room/interfaces"
	"steve/room/interfaces/facade"
	"sync"

	"github.com/golang/protobuf/proto"
)

type deskPlayer struct {
	playerID    uint64
	seat        uint32 // 座号
	ecoin       uint64 // 进牌桌金币数
	quit        bool   // 是否已经退出牌桌
	overTime    int    // 超时计数
	maxOverTime int    // 最大超时次数
	tuoguan     bool   // 是否在托管中
	robotLv     int    // 机器人等级

	mu sync.RWMutex
}

// CreateDeskPlayer 创建牌桌玩家
// maxOverTime : 最大超时次数，超过此次数将会被自动托管
func CreateDeskPlayer(playerID uint64, seat uint32, coin uint64, maxOverTime int, robotLv int) interfaces.DeskPlayer {
	return &deskPlayer{
		playerID:    playerID,
		seat:        seat,
		ecoin:       coin,
		maxOverTime: maxOverTime,
		robotLv:     robotLv,
	}
}

// GetRobotLv 获取机器人等级
func (dp *deskPlayer) GetRobotLv() int {
	return dp.robotLv
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
	defer dp.mu.Unlock()
	dp.quit = true
	dp.tuoguan = true // 退出后自动托管
}

// EnterDesk 进入牌桌
func (dp *deskPlayer) EnterDesk() {
	dp.mu.Lock()
	defer dp.mu.Unlock()
	dp.quit = false
}

// OnPlayerOverTime 玩家超时
func (dp *deskPlayer) OnPlayerOverTime() {
	dp.mu.Lock()
	defer dp.mu.Unlock()
	dp.overTime++

	if dp.overTime >= dp.maxOverTime && !dp.tuoguan {
		dp.tuoguan = true
		dp.notifyTuoguan(dp.playerID, true)
	}
}

// IsTuoguan 玩家是否在托管中
func (dp *deskPlayer) IsTuoguan() bool {
	dp.mu.RLock()
	defer dp.mu.RUnlock()
	return dp.tuoguan
}

// SetTuoguan 设置托管状态
func (dp *deskPlayer) SetTuoguan(tuoguan bool, notify bool) {
	dp.mu.Lock()
	defer dp.mu.Unlock()
	dp.tuoguan = tuoguan
	if notify {
		dp.notifyTuoguan(dp.playerID, tuoguan)
	}
}

// notifyTuoguan 通知玩家托管状态
func (dp *deskPlayer) notifyTuoguan(playerID uint64, tuoguan bool) {
	facade.SendMessageToPlayer(playerID, msgid.MsgID_ROOM_TUOGUAN_NTF, &room.RoomTuoGuanNtf{
		Tuoguan: proto.Bool(tuoguan),
	})
}
