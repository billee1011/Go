package matchv3

import (
	"context"
	"math/rand"
	"steve/server_pb/room_mgr"
	"steve/structs"

	"github.com/Sirupsen/logrus"
)

// randSeat 给desk的所有玩家分配座位号
func randSeat(desk *desk) {
	// 续局牌桌不重新分配
	if desk.isContinue {
		return
	}

	// 先分配为：0,1,2,3......
	seat := 0
	for i := range desk.players {
		desk.players[i].seat = seat
		seat++
	}

	// 随机交换
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

	exposer := structs.GetGlobalExposer()

	// 获取room的service
	rs, err := exposer.RPCClient.GetConnectByServerName("room")
	if err != nil || rs == nil {
		logEntry.WithError(err).Errorln("获得room服的gRPC失败!!!")
		return
	}

	// 给desk的所有玩家分配座位号
	randSeat(desk)

	// 该桌子所有的玩家信息
	createPlayers := []*roommgr.DeskPlayer{}
	for _, player := range desk.players {

		deskPlayer := &roommgr.DeskPlayer{
			PlayerId:   player.playerID,
			RobotLevel: int32(player.robotLv),
			Seat:       uint32(player.seat),
		}

		createPlayers = append(createPlayers, deskPlayer)
	}

	logEntry = logEntry.WithField("create_players", createPlayers)

	roomMgrClient := roommgr.NewRoomMgrClient(rs)

	// 调用room服的创建桌子
	_, err = roomMgrClient.CreateDesk(context.Background(), &roommgr.CreateDeskRequest{
		Players:    createPlayers,
		GameId:     uint32(desk.gameID),
		FixBanker:  desk.fixBanker,
		BankerSeat: uint32(desk.bankerSeat),
	})

	// 不成功时，报错，应该重新调用或者重新匹配
	if err != nil {
		logEntry.WithError(err).Errorln("create desk failed!!!")
		return
	}

	logEntry.Debugln("create desk success.")
	return
}

// sendCreateDesk 向room服请求创建牌桌，创建失败时则重新请求
func sendCreateDesk(pDesk *matchDesk) {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "sendCreateDesk",
		"desk":      pDesk.String(),
	})

	exposer := structs.GetGlobalExposer()

	// 获取room的service
	rs, err := exposer.RPCClient.GetConnectByServerName("room")
	if err != nil || rs == nil {
		logEntry.WithError(err).Errorln("获得room服的gRPC失败!!!")
		return
	}

	// 给desk的所有玩家分配座位号
	//randSeat(desk)

	// 该桌子所有的玩家信息
	createPlayers := []*roommgr.DeskPlayer{}
	for _, player := range pDesk.players {

		deskPlayer := &roommgr.DeskPlayer{
			PlayerId:   player.playerID,
			RobotLevel: int32(player.robotLv),
			Seat:       uint32(player.seat),
		}

		createPlayers = append(createPlayers, deskPlayer)
	}

	logEntry = logEntry.WithField("create_players", createPlayers)

	roomMgrClient := roommgr.NewRoomMgrClient(rs)

	// 调用room服的创建桌子
	_, err = roomMgrClient.CreateDesk(context.Background(), &roommgr.CreateDeskRequest{
		Players:    createPlayers,
		GameId:     uint32(pDesk.gameID),
		FixBanker:  true,
		BankerSeat: uint32(0),
	})

	// 不成功时，报错，应该重新调用或者重新匹配
	if err != nil {
		logEntry.WithError(err).Errorln("create desk failed!!!")
		return
	}

	logEntry.Debugln("create desk success.")
	return
}
