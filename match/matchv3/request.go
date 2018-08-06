package matchv3

import (
	"context"
	"math/rand"
	"steve/server_pb/room_mgr"
	"steve/structs"
	"time"

	"github.com/Sirupsen/logrus"
)

// randSeat 给desk的所有玩家分配座位号
func randSeat(desk *matchDesk) {
	// 先分配为：0,1,2,3......
	var seat int32 = 0
	for i := range desk.players {
		desk.players[i].seat = seat
		seat++
	}

	// 随机交换
	rand.Shuffle(len(desk.players), func(i, j int) {
		desk.players[i].seat, desk.players[j].seat = desk.players[j].seat, desk.players[i].seat
	})
}

// sendCreateDesk 向room服请求创建牌桌，创建失败时则重新请求
func sendCreateDesk(desk matchDesk, globalInfo *levelGlobalInfo) {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "sendCreateDesk",
		"desk":      desk,
	})

	logEntry.Debugln("进入函数，准备向room服请求创建桌子")

	exposer := structs.GetGlobalExposer()

	// 获取room的service
	rs, err := exposer.RPCClient.GetConnectByServerName("room")
	if err != nil || rs == nil {
		logEntry.WithError(err).Errorln("获得room服的gRPC失败，桌子被丢弃!!!")
		return
	}

	// 给desk的所有玩家分配座位号
	randSeat(&desk)

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
	rsp, err := roomMgrClient.CreateDesk(context.Background(), &roommgr.CreateDeskRequest{
		GameId:  desk.gameID,
		LevelId: desk.levelID,
		Players: createPlayers,
	})

	// 不成功时，报错，应该重新调用或者重新匹配
	if err != nil || rsp.GetErrCode() != roommgr.RoomError_SUCCESS {
		logEntry.WithError(err).Errorln("room服创建桌子失败，桌子被丢弃!!!")
		return
	}

	// 成功时的处理

	// 通知桌子的玩家
	// todo

	// 记录匹配成功的玩家同桌信息
	for i := 0; i < len(desk.players); i++ {
		globalInfo.sucPlayers[desk.players[i].playerID] = desk.deskID
	}

	// 记录匹配成功的桌子信息
	newSucDesk := sucDesk{
		gameID:  desk.gameID,
		levelID: desk.levelID,
		sucTime: time.Now().Unix(),
	}
	globalInfo.sucDesks[desk.deskID] = &newSucDesk

	logEntry.Debugln("离开函数，room服创建桌子成功")

	return
}
