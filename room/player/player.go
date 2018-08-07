package player

import (
	"steve/client_pb/msgid"
	"steve/client_pb/room"
	"steve/external/goldclient"
	"steve/external/hallclient"
	"steve/room/desk"
	"steve/room/util"
	"steve/server_pb/gold"
	server_gold "steve/server_pb/gold"
	"sync"

	"github.com/golang/protobuf/proto"
)

type Player struct {
	PlayerID    uint64
	seat        uint32 // 座号
	ecoin       uint64 // 进牌桌金币数
	quit        bool   // 是否已经退出牌桌
	overTime    int    // 超时计数
	maxOverTime int    // 最大超时次数
	tuoguan     bool   // 是否在托管中
	robotLv     int    // 机器人等级
	desk        *desk.Desk

	mu sync.RWMutex
}

// GetDesk 获取玩家所在牌桌
func (dp *Player) GetDesk() *desk.Desk {
	return dp.desk
}

// SetDesk 设置玩家所在牌桌
func (dp *Player) SetDesk(deskObj *desk.Desk) {
	dp.mu.Lock()
	dp.desk = deskObj
	dp.mu.Unlock()
}

// SetQuit 设置玩家退出状态
func (dp *Player) SetQuit(quit bool) {
	dp.mu.Lock()
	dp.quit = quit
	dp.mu.Unlock()
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
func (dp *Player) GetEcoin() uint64 {
	return dp.ecoin
}

func (dp *Player) SetEcoin(coin uint64) {
	dp.ecoin = coin
}

func (p *Player) SetMaxOverTime(time int) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	p.overTime = time
}

func (p *Player) SetRobotLv(lv int) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	p.robotLv = lv
}

func (p *Player) GetRobotLv() int {
	return p.robotLv
}

// IsQuit 是否已经退出
func (dp *Player) IsQuit() bool {
	return dp.quit
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
	coin, err := goldclient.GetGold(p.PlayerID, int16(gold.GoldType_GOLD_COIN))
	if err != nil {
		return 0
	}
	return uint64(coin)
}

func (p *Player) SetCoin(coin uint64) {
	gold, err := goldclient.GetGold(p.PlayerID, int16(server_gold.GoldType_GOLD_COIN))
	if err != nil {
		return
	}
	goldclient.AddGold(p.PlayerID, int16(server_gold.GoldType_GOLD_COIN), int64(coin)-gold, 0, 0)
}

// IsOnline 判断玩家是否在线
func (p *Player) IsOnline() bool {
	online, _ := hallclient.GetGateAddr(p.PlayerID)
	return online != ""
}

func (dp *Player) notifyTuoguan(playerID uint64, tuoguan bool) {
	util.SendMessageToPlayer(playerID, msgid.MsgID_ROOM_TUOGUAN_NTF, &room.RoomTuoGuanNtf{
		Tuoguan: proto.Bool(tuoguan),
	})
}
