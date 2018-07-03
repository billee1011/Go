package deskbase

import (
	"steve/client_pb/room"
	"steve/room/interfaces/facade"
	"steve/room/interfaces/global"
	"steve/structs/exchanger"
	"steve/structs/proto/gate_rpc"

	"github.com/Sirupsen/logrus"
)

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
	deskPlayer := facade.GetDeskPlayerByID(desk, playerID)
	deskPlayer.SetTuoguan(false, true)
	return
}
