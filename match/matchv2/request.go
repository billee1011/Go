package matchv2

import (
	"context"
	"steve/server_pb/room_mgr"
	"steve/structs"

	"github.com/Sirupsen/logrus"
)

// requestCreateDesk 向 room 请求创建牌桌
func requestCreateDesk(desk *desk) {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "mgr::requestCreate",
		"desk":      *desk,
	})
	e := structs.GetGlobalExposer()

	rs, err := e.RPCClient.GetConnectByServerName("room")
	if err != nil || rs == nil {
		logEntry.WithError(err).Errorln("get 'room' service failed!!!")
		return
	}
	createPlayers := []*roommgr.DeskPlayer{}
	for _, player := range desk.players {
		createPlayers = append(createPlayers, &roommgr.DeskPlayer{
			PlayerId:   player.playerID,
			RobotLevel: int32(player.robotLv),
		})
	}

	roomMgrClient := roommgr.NewRoomMgrClient(rs)
	_, err = roomMgrClient.CreateDesk(context.Background(), &roommgr.CreateDeskRequest{
		Players: createPlayers,
		GameId:  uint32(desk.gameID),
	})

	if err != nil {
		logEntry.WithError(err).Errorln("create desk failed!!!")
		return
	}
	logEntry.Debugln("create desk success.")
	return
}
