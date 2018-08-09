package matchv3

import (
	"context"
	"math/rand"
	"steve/client_pb/match"
	"steve/client_pb/msgid"
	"steve/external/gateclient"
	"steve/external/hallclient"
	"steve/server_pb/room_mgr"
	"steve/structs"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
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

	// 通知玩家，匹配成功，创建桌子
	// matchPlayer转换为deskPlayerInfo
	deskPlayers := []*match.DeskPlayerInfo{}
	for i := 0; i < len(desk.players); i++ {
		pDeskPlayer := translateToDeskPlayer(&desk.players[i])
		if pDeskPlayer == nil {
			logEntry.Errorln("把matchPlayer转换为deskPlayerInfo失败，跳过")
			continue
		}
		deskPlayers = append(deskPlayers, pDeskPlayer)
	}

	// 通知消息体
	ntf := match.MatchSucCreateDeskNtf{
		GameId:  &desk.gameID,
		LevelId: &desk.levelID,
		Players: deskPlayers,
	}

	// 广播给桌子内的所有真实玩家
	for i := 0; i < len(desk.players); i++ {
		if desk.players[i].robotLv == 0 {
			gateclient.SendPackageByPlayerID(desk.players[i].playerID, uint32(msgid.MsgID_MATCH_SUC_CREATE_DESK_NTF), &ntf)
		}
	}

	// 该桌子所有的玩家信息
	createPlayers := []*roommgr.DeskPlayer{}
	for i := 0; i < len(desk.players); i++ {

		deskPlayer := &roommgr.DeskPlayer{
			PlayerId:   desk.players[i].playerID,
			RobotLevel: desk.players[i].robotLv,
			Seat:       uint32(desk.players[i].seat),
		}

		createPlayers = append(createPlayers, deskPlayer)
	}

	roomMgrClient := roommgr.NewRoomMgrClient(rs)

	// 调用room服的创建桌子
	rsp, err := roomMgrClient.CreateDesk(context.Background(), &roommgr.CreateDeskRequest{
		GameId:  desk.gameID,
		LevelId: desk.levelID,
		DeskId:  desk.deskID,
		Players: createPlayers,
	})

	// 不成功时，报错，应该重新调用或者重新匹配
	if err != nil || rsp.GetErrCode() != roommgr.RoomError_SUCCESS {
		logEntry.WithError(err).Errorln("room服创建桌子失败，桌子被丢弃!!!")
		return
	}

	// 成功时的处理

	// 记录匹配成功的真实玩家同桌信息
	for i := 0; i < len(desk.players); i++ {
		if desk.players[i].robotLv == 0 {
			globalInfo.sucPlayers[desk.players[i].playerID] = desk.deskID
		}
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

// 把 matchPlayer 转换为 match.DeskPlayerInfo
func translateToDeskPlayer(player *matchPlayer) *match.DeskPlayerInfo {

	// 从hall服获取玩家基本信息
	playerInfo, err := hallclient.GetPlayerInfo(player.playerID)
	if err != nil || playerInfo == nil {
		logrus.WithError(err).Errorln("从hall服获取玩家信息失败，玩家ID:%v", player.playerID)
		return nil
	}

	deskPlayer := match.DeskPlayerInfo{
		PlayerId: &player.playerID,
		Name:     proto.String(playerInfo.GetNickName()),
		Coin:     &player.gold,
		Seat:     &player.seat,
		Gender:   proto.Uint32(playerInfo.GetGender()),
		Avatar:   proto.String(playerInfo.GetAvatar()),
	}

	return &deskPlayer
}
