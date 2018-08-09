package matchv2

import (
	"context"
	"math/rand"
	"steve/server_pb/room_mgr"
	"steve/structs"

	"github.com/Sirupsen/logrus"
)

// randSeat 分配座号
func randSeat(desk *desk) {
	// 续局牌桌不重新分配
	if desk.isContinue {
		return
	}
	seat := 0
	for i := range desk.players {
		desk.players[i].seat = seat
		seat++
	}
	rand.Shuffle(len(desk.players), func(i, j int) {
		desk.players[i].seat, desk.players[j].seat = desk.players[j].seat, desk.players[i].seat
	})
}

// requestCreateDesk 向 room 请求创建牌桌
func requestCreateDesk(desk *desk) {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "mgr::requestCreate",
		"desk":      desk.String(),
	})
	e := structs.GetGlobalExposer()

	rs, err := e.RPCClient.GetConnectByServerName("room")
	if err != nil || rs == nil {
		logEntry.WithError(err).Errorln("get 'room' service failed!!!")
		return
	}
	randSeat(desk)

	createPlayers := []*roommgr.DeskPlayer{}
	for _, player := range desk.players {
		createPlayers = append(createPlayers, &roommgr.DeskPlayer{
			PlayerId:   player.playerID,
			RobotLevel: int32(player.robotLv),
			Seat:       uint32(player.seat),
		})
	}
	logEntry = logEntry.WithField("create_players", createPlayers)

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
