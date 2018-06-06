package desks

import (
	"steve/client_pb/msgId"
	"steve/client_pb/room"
	"steve/room/interfaces"
	"steve/room/interfaces/facade"
	"steve/room/interfaces/global"
	"steve/structs/exchanger"
	"steve/structs/proto/gate_rpc"
	"sync"

	"github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
)

type tuoGuanPlayer struct {
	overTimerCount int  // 超时计数
	tuoGuaning     bool // 是否在托管中
}

type tuoGuanMgr struct {
	players      map[uint64]*tuoGuanPlayer
	maxOverTimer int //	最大超时次数，超过此次数则进入托管状态
	mu           sync.RWMutex
}

// newTuoGuanMgr 创建托管管理器
func newTuoGuanMgr() interfaces.TuoGuanMgr {
	return &tuoGuanMgr{
		players:      make(map[uint64]*tuoGuanPlayer),
		maxOverTimer: 2,
	}
}

// GetTuoGuanPlayers 获取托管玩家
func (tg *tuoGuanMgr) GetTuoGuanPlayers() []uint64 {
	tg.mu.RLock()
	defer tg.mu.RUnlock()
	result := []uint64{}
	for playerID, player := range tg.players {
		if player.tuoGuaning {
			result = append(result, playerID)
		}
	}
	return result
}

// SetTuoGuan 设置玩家托管
func (tg *tuoGuanMgr) SetTuoGuan(playerID uint64, set bool, notify bool) {
	tg.mu.Lock()
	defer tg.mu.Unlock()

	player, exist := tg.players[playerID]
	if !exist {
		player = &tuoGuanPlayer{}
		tg.players[playerID] = player
	}
	player.tuoGuaning = set
	tg.notifyTuoguan(playerID, notify)
}

// OnPlayerTimeOut 处理完成超时事件
func (tg *tuoGuanMgr) OnPlayerTimeOut(playerID uint64) {
	tg.mu.Lock()
	defer tg.mu.Unlock()

	player, exist := tg.players[playerID]
	if !exist {
		player = &tuoGuanPlayer{}
		tg.players[playerID] = player
	}
	player.overTimerCount++
	if player.overTimerCount >= tg.maxOverTimer && !player.tuoGuaning {
		player.tuoGuaning = true
		tg.notifyTuoguan(playerID, true)
	}
}

// notifyTuoguan 通知玩家托管状态
func (tg *tuoGuanMgr) notifyTuoguan(playerID uint64, tuoguan bool) {
	facade.SendMessageToPlayer(playerID, msgid.MsgID_ROOM_TUOGUAN_NTF, &room.RoomTuoGuanNtf{
		Tuoguan: proto.Bool(tuoguan),
	})
}

// HandleCancelTuoGuanReq 处理取消托管请求
func HandleCancelTuoGuanReq(clientID uint64, header *steve_proto_gaterpc.Header, req room.RoomCancelTuoGuanReq) (ret []exchanger.ResponseMsg) {
	ret = []exchanger.ResponseMsg{}

	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "HandleCancelTuoGuanReq",
		"client_id": clientID,
	})
	playerMgr := global.GetPlayerMgr()
	player := playerMgr.GetPlayerByClientID(clientID)
	if player == nil {
		logEntry.Debugln("未登录的客户端")
		return
	}
	playerID := player.GetID()
	logEntry = logEntry.WithField("player_id", playerID)

	deskMgr := global.GetDeskMgr()
	desk, _ := deskMgr.GetRunDeskByPlayerID(playerID)
	if desk == nil {
		logEntry.Debugln("玩家不在房间中")
		return
	}

	tuoGuanMgr := desk.GetTuoGuanMgr()
	tuoGuanMgr.SetTuoGuan(playerID, false, true)
	logEntry.Debugln("玩家取消托管")
	return
}
