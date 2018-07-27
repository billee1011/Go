package deskbase

import (
	"steve/client_pb/room"
	"steve/room/interfaces/facade"
	"steve/room/interfaces/global"
	"steve/structs/exchanger"
	"steve/structs/proto/gate_rpc"

	"github.com/Sirupsen/logrus"
)

// HandleCancelTuoGuanReq 处理取消托管请求 @Deprecated
func HandleCancelTuoGuanReq(playerID uint64, header *steve_proto_gaterpc.Header, req room.RoomCancelTuoGuanReq) (ret []exchanger.ResponseMsg) {
	ret = []exchanger.ResponseMsg{}

	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "HandleCancelTuoGuanReq",
		"player_id": playerID,
	})
	playerMgr := global.GetPlayerMgr()
	player := playerMgr.GetPlayer(playerID)
	if player == nil {
		logEntry.Debugln("获取玩家失败")
		return
	}
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

// HandleTuoGuanReq 处理取消托管请求
func HandleTuoGuanReq(playerID uint64, header *steve_proto_gaterpc.Header, req room.RoomTuoGuanReq) (ret []exchanger.ResponseMsg) {
	ret = []exchanger.ResponseMsg{}

	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "HandleTuoGuanReq",
		"player_id": playerID,
	})
	playerMgr := global.GetPlayerMgr()
	player := playerMgr.GetPlayer(playerID)
	if player == nil {
		logEntry.Debugln("获取玩家失败")
		return
	}
	deskMgr := global.GetDeskMgr()
	desk, _ := deskMgr.GetRunDeskByPlayerID(playerID)
	if desk == nil {
		logEntry.Debugln("玩家不在房间中")
		return
	}
	deskPlayer := facade.GetDeskPlayerByID(desk, playerID)
	deskPlayer.SetTuoguan(req.GetTuoguan(), true)
	return
}
