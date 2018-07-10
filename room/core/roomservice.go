package core

import (
	"context"
	msgid "steve/client_pb/msgId"
	"steve/client_pb/room"
	"steve/room/desks/deskbase"
	"steve/room/interfaces"
	"steve/room/interfaces/facade"
	"steve/room/interfaces/global"
	"steve/server_pb/room_mgr"

	"github.com/Sirupsen/logrus"
)

// RoomService room房间RPC服务
type RoomService struct {
}

// notifyDeskCreate 通知房间创建
func notifyDeskCreate(desk interfaces.Desk) {
	players := []*room.RoomPlayerInfo{}
	deskPlayers := desk.GetDeskPlayers()
	for _, player := range deskPlayers {
		roomPlayer := deskbase.TranslateToRoomPlayer(player)
		players = append(players, &roomPlayer)
	}
	ntf := room.RoomDeskCreatedNtf{
		Players: players,
	}
	facade.BroadCastDeskMessage(desk, nil, msgid.MsgID_ROOM_DESK_CREATED_NTF, &ntf, true)
}

// CreateDesk 创建牌桌
func (hws *RoomService) CreateDesk(ctx context.Context, req *roommgr.CreateDeskRequest) (rsp *roommgr.CreateDeskResponse, err error) {
	players := req.GetPlayers()
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "RoomService::CreateDesk",
		"players":   players,
	})
	// 回复match服的消息
	rsp = &roommgr.CreateDeskResponse{
		ErrCode: roommgr.RoomError_FAILED, // 默认是失败的
	}

	playerIDs := []uint64{}
	for _, player := range players {
		playerIDs = append(playerIDs, player.GetPlayerId())
	}

	deskFactory := global.GetDeskFactory()
	deskMgr := global.GetDeskMgr()

	// 创建桌子
	result, err := deskFactory.CreateDesk(playerIDs, int(req.GetGameId()), interfaces.CreateDeskOptions{})
	if err != nil {
		logEntry.WithFields(
			logrus.Fields{
				"players": playerIDs,
				"result":  result,
			},
		).WithError(err).Errorln("创建桌子失败")
		return
	}
	logEntry.Debugln("创建桌子成功")

	rsp.ErrCode = roommgr.RoomError_SUCCESS
	notifyDeskCreate(result.Desk)
	deskMgr.RunDesk(result.Desk)
	return
}
