package core

import (
	"context"
	msgid "steve/client_pb/msgid"
	"steve/client_pb/room"
	"steve/common/data/player"
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

	deskPlayers := make([]interfaces.DeskPlayer, 0, len(players))
	var seat uint32
	for _, pbPlayer := range players {
		robotLv := int(pbPlayer.GetRobotLevel())
		playerID := pbPlayer.GetPlayerId()
		deskPlayer := deskbase.CreateDeskPlayer(playerID, seat, player.GetPlayerCoin(playerID), 2, robotLv)
		deskPlayers = append(deskPlayers, deskPlayer)
		seat++
	}

	deskFactory := global.GetDeskFactory()
	deskMgr := global.GetDeskMgr()

	// 创建桌子
	result, err := deskFactory.CreateDesk(deskPlayers, int(req.GetGameId()), interfaces.CreateDeskOptions{})
	if err != nil {
		logEntry.WithFields(
			logrus.Fields{
				"players": deskPlayers,
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
