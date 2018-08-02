package matchv3

import (
	"context"
	"fmt"
	"math/rand"
	"steve/match/web"
	"steve/server_pb/gold"
	"steve/server_pb/room_mgr"
	"steve/server_pb/user"
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

	// 获取hall服的Connection
	hallConnection, err := exposer.RPCClient.GetConnectByServerName("hall")
	if err != nil || hallConnection == nil {
		logEntry.WithError(err).Errorln("获得hall服的gRPC失败!!!")
		return 0, fmt.Errorf("获得hall服的gRPC失败!!!")
	}

	hallClient := user.NewPlayerDataClient(hallConnection)

	// 向hall服请求游戏信息
	rsp, err := hallClient.GetPlayerGameInfo(context.Background(), &user.GetPlayerGameInfoReq{PlayerId: playerID, GameId: gameID})

	// 不成功时，报错
	if err != nil || rsp == nil {
		logrus.WithError(err).Errorln("从hall服获取玩家的胜率失败!!!")
		return 0, fmt.Errorf("从hall服获取玩家的胜率失败!!!")
	}

	// 返回的不是成功，报错
	if rsp.GetErrCode() != int32(user.ErrCode_EC_SUCCESS) {
		logrus.WithError(err).Errorln("从hall服获取玩家胜率成功，但errCode显示失败!!!")
		return 0, fmt.Errorf("从hall服获取玩家胜率成功，但errCode显示失败!!!")
	}

	// 计算胜率
	var winRate float64 = 0.0

	if rsp.GetTotalBurea() < web.GetMinGameTimes() {
		winRate = 0.5
		logEntry.Debugf("玩家总局数为：%v，少于规定的%v场，所以胜率定为50%", rsp.GetTotalBurea(), web.GetMinGameTimes())
	} else {
		winRate = float64(rsp.GetMaxWinningStream()) / float64(rsp.GetTotalBurea())
		logEntry.Debugf("玩家总局数为：%v，胜利局数为：%v，计算得到胜率为50%", rsp.GetTotalBurea(), web.GetMinGameTimes())
	}

	result := float32(winRate) * 100

	return result, nil
}

// requestPlayerIP 向hall服请求指定玩家的信息
// playerID : 玩家playerID
// gameID 	: 游戏ID
func requestPlayerIP(playerID uint64) (string, error) {
	logEntry := logrus.WithFields(logrus.Fields{
		"playerID": playerID,
	})

	logEntry.Debugln("进入函数")

	return "127.0.0.1", nil

	/* 	exposer := structs.GetGlobalExposer()

	   	// 获取hall服的Connection
	   	hallConnection, err := exposer.RPCClient.GetConnectByServerName("hall")
	   	if err != nil || hallConnection == nil {
	   		logEntry.WithError(err).Errorln("获得hall服的gRPC失败!!!")
	   		return 0, fmt.Errorf("获得hall服的gRPC失败!!!")
	   	}

	   	hallClient := user.NewPlayerDataClient(hallConnection)

	   	// 向hall服请求游戏信息
	   	rsp, err := hallClient.GetPlayerGameInfo(context.Background(), &user.GetPlayerGameInfoReq{PlayerId: playerID, GameId: gameID})

	   	// 不成功时，报错
	   	if err != nil || rsp == nil {
	   		logrus.WithError(err).Errorln("从hall服获取玩家的胜率失败!!!")
	   		return 0, fmt.Errorf("从hall服获取玩家的胜率失败!!!")
	   	}

	   	// 返回的不是成功，报错
	   	if rsp.GetErrCode() != int32(user.ErrCode_EC_SUCCESS) {
	   		logrus.WithError(err).Errorln("从hall服获取玩家胜率成功，但errCode显示失败!!!")
	   		return 0, fmt.Errorf("从hall服获取玩家胜率成功，但errCode显示失败!!!")
	   	}

	   	// 计算胜率
	   	var winRate float64 = 0.0

	   	if rsp.GetTotalBurea() < web.GetMinGameTimes() {
	   		winRate = 0.5
	   		logEntry.Debugf("玩家总局数为：%v，少于规定的%v场，所以胜率定为50%", rsp.GetTotalBurea(), web.GetMinGameTimes())
	   	} else {
	   		winRate = float64(rsp.GetMaxWinningStream()) / float64(rsp.GetTotalBurea())
	   		logEntry.Debugf("玩家总局数为：%v，胜利局数为：%v，计算得到胜率为50%", rsp.GetTotalBurea(), web.GetMinGameTimes())
	   	}

	   	result := float32(winRate) * 100

	   	return result, nil */
}
