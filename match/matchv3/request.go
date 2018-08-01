package matchv3

import (
	"context"
	"fmt"
	"math/rand"
	"steve/server_pb/gold"
	"steve/server_pb/room_mgr"
	"steve/structs"

	"github.com/Sirupsen/logrus"
)

// randSeat 给desk的所有玩家分配座位号
func randSeat(desk *matchDesk) {
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

// requestPlayerGold 向gold服请求指定玩家，指定货币类型的数量
// playerID : 玩家playerID
// goldType : 货币类型，对应枚举 gold.GoldType
func requestPlayerGold(playerID uint64, goldType gold.GoldType) (int64, error) {
	logEntry := logrus.WithFields(logrus.Fields{
		"playerID": playerID,
		"goldType": goldType,
	})

	logEntry.Debugln("进入函数")

	exposer := structs.GetGlobalExposer()

	// 获取gold的Connection
	goldConnection, err := exposer.RPCClient.GetConnectByServerName("gold")
	if err != nil || goldConnection == nil {
		logEntry.WithError(err).Errorln("获得gold服的gRPC失败!!!")
		return 0, fmt.Errorf("获得gold服的gRPC失败!!!")
	}

	// 请求结构体
	reqGold := gold.GetGoldReq{
		Item: &gold.GetItem{
			Uid:      playerID,
			GoldType: int32(goldType),
			Value:    0,
		},
	}

	goldClient := gold.NewGoldClient(goldConnection)

	// 调用room服的创建桌子
	rspGold, err := goldClient.GetGold(context.Background(), &reqGold)

	// 不成功时，报错
	if err != nil || rspGold == nil {
		logEntry.WithError(err).Errorln("从gold服获取玩家金钱数据失败!!!")
		return 0, fmt.Errorf("从gold服获取玩家金钱数据失败!!!")
	}

	// 其他出错
	if rspGold.GetErrCode() != gold.ResultStat_SUCCEED || rspGold.GetErrDesc() != "" {
		logEntry.Errorf("从gold服获取玩家金钱数据出错，errCode:%v，errDesc:%v \n", rspGold.GetErrCode(), rspGold.GetErrDesc())
		return 0, fmt.Errorf("从gold服获取玩家金钱数据出错，errCode:%v，errDesc:%v \n", rspGold.GetErrCode(), rspGold.GetErrDesc())
	}

	logEntry.Debugln("create desk success.")
	return rspGold.GetItem().GetValue(), nil
}

// requestPlayerWinRate 向hall服请求指定玩家，指定游戏的胜率
// playerID : 玩家playerID
// gameID 	: 游戏ID
func requestPlayerWinRate(playerID uint64, gameID uint32) (float32, error) {
	logEntry := logrus.WithFields(logrus.Fields{
		"playerID": playerID,
		"gameID":   gameID,
	})

	logEntry.Debugln("进入函数")

	exposer := structs.GetGlobalExposer()

	// 获取robot服的Connection
	robotConnection, err := exposer.RPCClient.GetConnectByServerName("robot")
	if err != nil || robotConnection == nil {
		logEntry.WithError(err).Errorln("获得robot服的gRPC失败!!!")
		return 0, fmt.Errorf("获得robot服的gRPC失败!!!")
	}

	// 请求结构体
	/*reqGold := gold.GetGoldReq{
		Item: &gold.GetItem{
			Uid:      playerID,
			GoldType: int32(0),
			Value:    0,
		},
	}

	 	goldClient := gold.NewGoldClient(goldConnection)

	   	// 调用room服的创建桌子
	   	rspGold, err := goldClient.GetGold(context.Background(), &reqGold)

	   	// 不成功时，报错
	   	if err != nil || rspGold == nil {
	   		logEntry.WithError(err).Errorln("从gold服获取玩家金钱数据失败!!!")
	   		return 0, fmt.Errorf("从gold服获取玩家金钱数据失败!!!")
	   	}

	   	// 其他出错
	   	if rspGold.GetErrCode() != gold.ResultStat_SUCCEED || rspGold.GetErrDesc() != "" {
	   		logEntry.Errorf("从gold服获取玩家金钱数据出错，errCode:%v，errDesc:%v \n", rspGold.GetErrCode(), rspGold.GetErrDesc())
	   		return 0, fmt.Errorf("从gold服获取玩家金钱数据出错，errCode:%v，errDesc:%v \n", rspGold.GetErrCode(), rspGold.GetErrDesc())
	   	}

	   	logEntry.Debugln("create desk success.") */
	return 0.5, nil
}
