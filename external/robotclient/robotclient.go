package robotclient

import (
	"context"
	"errors"
	"steve/server_pb/robot"
	"steve/structs"
	"steve/structs/common"

	"github.com/Sirupsen/logrus"
	"google.golang.org/grpc"
)

/*
	功能：机器人服的Client API封装,实现调用
	作者： wuhongwei
	日期： 2018-8-03
*/

//LeisureRobotReqInfo 空闲机器人请求信息
type LeisureRobotReqInfo struct {
	CoinHigh    int64
	CoinLow     int64
	WinRateHigh int32
	WinRateLow  int32
	GameID      uint32
	LevelID     uint32
}

// GetLeisureRobotInfoByInfo 获取空闲机器人
// param:   LeisureRobotReqInfo
// 返回:
// uint64:	机器人玩家ID
// int64:	机器人金豆数
// int32:	机器人胜率
// error:	错误信息
func GetLeisureRobotInfoByInfo(leisureRobotReqInfo LeisureRobotReqInfo) (uint64, int64, float64, error) {
	// 得到服务连接
	con, err := getRobotServer()
	if err != nil || con == nil {
		return 0, 0, 0, errors.New("no connection")
	}
	// 新建Client
	client := robot.NewRobotServiceClient(con)
	coinsR := &robot.CoinsRange{
		High: leisureRobotReqInfo.CoinHigh,
		Low:  leisureRobotReqInfo.CoinLow,
	}
	winR := &robot.WinRateRange{
		High: leisureRobotReqInfo.WinRateHigh,
		Low:  leisureRobotReqInfo.WinRateLow,
	}
	game := &robot.GameConfig{
		GameId:  leisureRobotReqInfo.GameID,
		LevelId: leisureRobotReqInfo.LevelID,
	}
	// 调用RPC方法
	rsp, err := client.GetLeisureRobotInfoByInfo(context.Background(), &robot.GetLeisureRobotInfoReq{
		Game:         game,
		CoinsRange:   coinsR,
		WinRateRange: winR,
		NewState:     robot.RobotPlayerState_RPS_MATCHING,
	})
	// 检测返回值
	if err != nil {
		return 0, 0, 0, err
	}
	if rsp.ErrCode != int32(robot.ErrCode_EC_SUCCESS) {
		return 0, 0, 0, errors.New(" get leisure robot failed")
	}

	return rsp.GetRobotPlayerId(), rsp.GetCoin(), rsp.GetWinRate(), nil
}

// SetRobotPlayerState 更新机器人状态
// param:  playerID,oldState 玩家当前状态， newState 要更新状态， serverType 服务类型，serverAddr 服务地址
// return: 更新结果，错误信息
func SetRobotPlayerState(playerID uint64, oldState, newState, serverType uint32, serverAddr string) (bool, error) {
	// 得到服务连接
	con, err := getRobotServer()
	if err != nil || con == nil {
		return false, errors.New("no connection")
	}

	// 新建Client
	client := robot.NewRobotServiceClient(con)
	// 调用RPC方法
	rsp, err := client.SetRobotPlayerState(context.Background(), &robot.SetRobotPlayerStateReq{
		RobotPlayerId: playerID,
		NewState:      robot.RobotPlayerState(newState),
		OldState:      robot.RobotPlayerState(oldState),
		ServerType:    robot.ServerType(serverType),
		ServerAddr:    serverAddr,
	})
	// 检测返回值
	if err != nil {
		return false, err
	}
	if rsp.ErrCode != int32(robot.ErrCode_EC_NOTROBOT) && rsp.ErrCode != int32(robot.ErrCode_EC_SUCCESS) {
		return false, errors.New("update player state failed")
	}
	return rsp.Result, nil
}

// UpdataRobotPlayerWinRate 更新机器人胜率
// param:  playerID,oldWinRate 玩家当前胜率， newWinRate 要更新胜率，
// return: 更新结果，错误信息
func UpdataRobotPlayerWinRate(playerID uint64, gameID int32, oldWinRate, newWinRate float64) (bool, error) {
	// 得到服务连接
	con, err := getRobotServer()
	if err != nil || con == nil {
		return false, errors.New("no connection")
	}
	// 新建Client
	client := robot.NewRobotServiceClient(con)
	// 调用RPC方法
	rsp, err := client.UpdataRobotGameWinRate(context.Background(), &robot.UpdataRobotGameWinRateReq{
		RobotPlayerId: playerID,
		GameId:        gameID,
		OldWinRate:    oldWinRate,
		NewWinRate:    newWinRate,
	})
	// 检测返回值
	if err != nil {
		return false, err
	}
	if rsp.ErrCode != int32(robot.ErrCode_EC_NOTROBOT) && rsp.ErrCode != int32(robot.ErrCode_EC_SUCCESS) {
		return false, errors.New("update player winRate failed")
	}
	return rsp.Result, nil
}

// IsRobotPlayer 判断是否时机器人
// param:  playerID,
// return: 判断结果，错误信息
func IsRobotPlayer(playerID uint64) (bool, error) {
	// 得到服务连接
	con, err := getRobotServer()
	if err != nil || con == nil {
		return false, errors.New("no connection")
	}

	// 新建Client
	client := robot.NewRobotServiceClient(con)
	// 调用RPC方法
	rsp, err := client.IsRobotPlayer(context.Background(), &robot.IsRobotPlayerReq{
		RobotPlayerId: playerID,
	})
	// 检测返回值
	if err != nil {
		logrus.WithError(err).Infoln("判断是否时机器人")
		return false, err
	}
	return rsp.GetResult(), nil
}

func getRobotServer() (*grpc.ClientConn, error) {
	e := structs.GetGlobalExposer()
	con, err := e.RPCClient.GetConnectByServerName(common.RobotServiceName)
	if err != nil || con == nil {
		return nil, errors.New("no connection")
	}
	return con, nil
}
