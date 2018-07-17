package room2

import (
	"sync"
	"steve/room/interfaces/facade"
	"steve/client_pb/msgId"
	"steve/client_pb/room"
	"github.com/golang/protobuf/proto"
	playerdata "steve/common/data/player"
	"steve/room2/desk"
)

type RoomPlayer struct {
	PlayerID    uint64
	seat        uint32 // 座号
	ecoin       uint64 // 进牌桌金币数
	quit        bool   // 是否已经退出牌桌
	overTime    int    // 超时计数
	maxOverTime int    // 最大超时次数
	tuoguan     bool   // 是否在托管中
	desk        *desk.Desk

	mu sync.RWMutex
}

func (dp *RoomPlayer) GetDesk() *desk.Desk{
	return dp.desk
}

// GetPlayerID 获取玩家 ID
func (dp *RoomPlayer) GetPlayerID() uint64 {
	dp.mu.RLock()
	defer dp.mu.RUnlock()
	return dp.PlayerID
}

// GetSeat 获取座号
func (dp *RoomPlayer) GetSeat() int {
	dp.mu.RLock()
	defer dp.mu.RUnlock()
	return int(dp.seat)
}
func (dp *RoomPlayer) SetSeat(seat uint32) {
	dp.mu.RLock()
	defer dp.mu.RUnlock()
	dp.seat = seat
}

// GetEcoin 获取进入时金币数
func (dp *RoomPlayer) GetEcoin() int {
	dp.mu.RLock()
	defer dp.mu.RUnlock()
	return int(dp.ecoin)
}

// IsQuit 是否已经退出
func (dp *RoomPlayer) IsQuit() bool {
	dp.mu.RLock()
	defer dp.mu.RUnlock()
	return dp.quit
}

// QuitDesk 退出房间
func (dp *RoomPlayer) QuitDesk(desk *desk.Desk) {
	dp.mu.Lock()
	defer dp.mu.Unlock()
	dp.quit = true
	dp.tuoguan = true // 退出后自动托管
	dp.desk = nil
}

// EnterDesk 进入房间
func (dp *RoomPlayer) EnterDesk(desk *desk.Desk) {
	dp.mu.Lock()
	defer dp.mu.Unlock()
	dp.quit = false
	dp.desk = desk
	dp.ecoin = dp.GetCoin()
}

// OnPlayerOverTime 玩家超时
func (dp *RoomPlayer) OnPlayerOverTime() {
	dp.mu.Lock()
	defer dp.mu.Unlock()
	dp.overTime++

	if dp.overTime >= dp.maxOverTime && !dp.tuoguan {
		dp.tuoguan = true
		dp.notifyTuoguan(dp.PlayerID, true)
	}
}

// IsTuoguan 玩家是否在托管中
func (dp *RoomPlayer) IsTuoguan() bool {
	dp.mu.RLock()
	defer dp.mu.RUnlock()
	return dp.tuoguan
}

// SetTuoguan 设置托管
func (dp *RoomPlayer) SetTuoguan(tuoguan bool, notify bool) {
	dp.mu.Lock()
	defer dp.mu.Unlock()
	dp.tuoguan = tuoguan
	if notify {
		dp.notifyTuoguan(dp.PlayerID, tuoguan)
	}
}

func (p *RoomPlayer) GetCoin() uint64 {
	return playerdata.GetPlayerCoin(p.PlayerID)
}

func (p *RoomPlayer) SetCoin(coin uint64) {
	playerdata.SetPlayerCoin(p.PlayerID, coin)
}

// GetUserName() string

// 判断玩家是否在线
func (p *RoomPlayer) IsOnline() bool {
	return playerdata.GetPlayerGateAddr(p.PlayerID) != ""
}

func (dp *RoomPlayer) notifyTuoguan(playerID uint64, tuoguan bool) {
	facade.SendMessageToPlayer(playerID, msgid.MsgID_ROOM_TUOGUAN_NTF, &room.RoomTuoGuanNtf{
		Tuoguan: proto.Bool(tuoguan),
	})
}
