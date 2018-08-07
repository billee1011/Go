package robotclient

import (
	"context"
	"errors"
	"steve/server_pb/robot"
	"steve/server_pb/user"
	"steve/structs"

	"google.golang.org/grpc"
)

// GetRobot 获取玩家游戏信息
// gameID:   游戏ID
// levelID:	 场次ID
// beginRate:   胜率起值
// endRate:	 胜率
// gameID:   游戏ID
// levelID:	 场次ID
// return:   获取到的机器人信息
func GetRobot(gameID uint32, levelID uint32, beginRate int8, endRate int8, beginGold int64, endGold int64) (*robot.GetRobotPlayerIDRsp, error) {

	// 得到服务连接
	con, err := getRobotServer()
	if err != nil || con == nil {
		return nil, errors.New("no robot connection")
	}

	// 新建Client
	client := robot.NewRobotServiceClient(con)

	// 调用RPC方法
	rsp, err := client.GetRobotPlayerIDByInfo(context.Background(), &robot.GetRobotPlayerIDReq{
		Game:         &robot.GameConfig{GameId: gameID, LevelId: levelID},
		WinRateRange: &robot.WinRateRange{Low: int32(beginRate), High: int32(endRate)},
		CoinsRange:   &robot.CoinsRange{Low: int64(beginGold), High: int64(endGold)},
	})

	// 检测返回值
	if err != nil || rsp == nil {
		return nil, err
	}

	if rsp.ErrCode != int32(user.ErrCode_EC_SUCCESS) {
		return nil, errors.New("GetRobot()获取机器人成功 ，但rsp.ErrCode显示失败")
	}

	return rsp, nil
}

func getRobotServer() (*grpc.ClientConn, error) {
	e := structs.GetGlobalExposer()

	// 对uid进行一致性hash路由策略.
	con, err := e.RPCClient.GetConnectByServerName("robot")
	if err != nil || con == nil {
		return nil, errors.New("获取robot connection 失败")
	}

	return con, nil
}
