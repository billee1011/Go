package player

import (
	"sync"
	"steve/room/interfaces/facade"
	"steve/client_pb/msgId"
	"steve/client_pb/room"
	"github.com/golang/protobuf/proto"
	playerdata "steve/common/data/player"
	"steve/room2/desk"
)

type Player struct {
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

func (dp *Player) GetDesk() *desk.Desk{
	return dp.desk
}

// GetPlayerID 获取玩家 ID
func (dp *Player) GetPlayerID() uint64 {
	dp.mu.RLock()
	defer dp.mu.RUnlock()
	return dp.PlayerID
}

// GetSeat 获取座号
func (dp *Player) GetSeat() int {
	dp.mu.RLock()
	defer dp.mu.RUnlock()
	return int(dp.seat)
}
func (dp *Player) SetSeat(seat uint32) {
	dp.mu.RLock()
	defer dp.mu.RUnlock()
	dp.seat = seat
}

// GetEcoin 获取进入时金币数
func (dp *Player) GetEcoin() int {
	dp.mu.RLock()
	defer dp.mu.RUnlock()
	return int(dp.ecoin)
}

// IsQuit 是否已经退出
func (dp *Player) IsQuit() bool {
	dp.mu.RLock()
	defer dp.mu.RUnlock()
	return dp.quit
}

// QuitDesk 退出房间
func (dp *Player) QuitDesk(desk *desk.Desk) {
	dp.mu.Lock()
	defer dp.mu.Unlock()
	dp.quit = true
	dp.tuoguan = true // 退出后自动托管
	dp.desk = nil
}

// EnterDesk 进入房间
func (dp *Player) EnterDesk(desk *desk.Desk) {
	dp.mu.Lock()
	defer dp.mu.Unlock()
	dp.quit = false
	dp.desk = desk
	dp.ecoin = dp.GetCoin()
}

// OnPlayerOverTime 玩家超时
func (dp *Player) OnPlayerOverTime() {
	dp.mu.Lock()
	defer dp.mu.Unlock()
	dp.overTime++

	if dp.overTime >= dp.maxOverTime && !dp.tuoguan {
		dp.tuoguan = true
		dp.notifyTuoguan(dp.PlayerID, true)
	}
}

// IsTuoguan 玩家是否在托管中
func (dp *Player) IsTuoguan() bool {
	dp.mu.RLock()
	defer dp.mu.RUnlock()
	return dp.tuoguan
}

// SetTuoguan 设置托管
func (dp *Player) SetTuoguan(tuoguan bool, notify bool) {
	dp.mu.Lock()
	defer dp.mu.Unlock()
	dp.tuoguan = tuoguan
	if notify {
		dp.notifyTuoguan(dp.PlayerID, tuoguan)
	}
}

func (p *Player) GetCoin() uint64 {
	return playerdata.GetPlayerCoin(p.PlayerID)
}

func (p *Player) SetCoin(coin uint64) {
	playerdata.SetPlayerCoin(p.PlayerID, coin)
}

// GetUserName() string

// 判断玩家是否在线
func (p *Player) IsOnline() bool {
	return playerdata.GetPlayerGateAddr(p.PlayerID) != ""
}

func (dp *Player) notifyTuoguan(playerID uint64, tuoguan bool) {
	facade.SendMessageToPlayer(playerID, msgid.MsgID_ROOM_TUOGUAN_NTF, &room.RoomTuoGuanNtf{
		Tuoguan: proto.Bool(tuoguan),
	})
}
