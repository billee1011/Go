package hallclient

import (
	"context"
	"errors"
	"steve/external/robotclient"
	"steve/server_pb/user"
	"steve/structs"
	"steve/structs/common"

	"google.golang.org/grpc"
)

/*
	功能：玩家服的Client API封装,实现调用
	作者： Zhengxuzhang
	日期： 2018-8-02
*/

// GetPlayerByAccount 根据账号获取玩家 ID
// param:   account:账号ID
// return:  玩家ID,错误信息
func GetPlayerByAccount(accountID uint64) (uint64, error) {

	// 得到服务连接
	con, err := getHallServer()
	if err != nil || con == nil {
		return 0, errors.New("no connection")
	}

	// 新建Client
	client := user.NewPlayerDataClient(con)

	// 调用RPC方法
	rsp, err := client.GetPlayerByAccount(context.Background(), &user.GetPlayerByAccountReq{
		AccountId: accountID,
	})

	// 检测返回值
	if err != nil {
		return 0, err
	}

	if rsp.ErrCode != int32(user.ErrCode_EC_SUCCESS) {
		return 0, errors.New("get player by account failed")
	}
	return rsp.PlayerId, nil
}

// GetPlayerState 获取玩家状态
// param:   uid:玩家ID
// return:  玩家状态，正在进行的游戏ID,正在进行的场次ID,错误信息
func GetPlayerState(uid uint64) (*user.GetPlayerStateRsp, error) {

	// 得到服务连接
	con, err := getHallServer()
	if err != nil || con == nil {
		return nil, errors.New("no hall connection")
	}

	// 新建Client
	client := user.NewPlayerDataClient(con)

	// 调用RPC方法
	rsp, err := client.GetPlayerState(context.Background(), &user.GetPlayerStateReq{
		PlayerId: uid,
	})

	// 检测返回值
	if err != nil {
		return nil, err
	}

	return rsp, nil
}

// GetGateAddr 获取玩家网关服地址
// param:   uid:玩家ID
// return:  网关服地址,错误信息
func GetGateAddr(uid uint64) (string, error) {

	// 得到服务连接
	con, err := getHallServer()
	if err != nil || con == nil {
		return "", errors.New("no hall connection")
	}

	// 新建Client
	client := user.NewPlayerDataClient(con)

	// 调用RPC方法
	rsp, err := client.GetPlayerState(context.Background(), &user.GetPlayerStateReq{
		PlayerId: uid,
	})

	// 检测返回值
	if err != nil {
		return "", err
	}

	return rsp.GetGateAddr(), nil
}

// GetRoomAddr 获取玩家房间服地址
// param:   uid:玩家ID
// return:  房间服地址,错误信息
func GetRoomAddr(uid uint64) (string, error) {

	// 得到服务连接
	con, err := getHallServer()
	if err != nil || con == nil {
		return "", errors.New("no hall connection")
	}

	// 新建Client
	client := user.NewPlayerDataClient(con)

	// 调用RPC方法
	rsp, err := client.GetPlayerState(context.Background(), &user.GetPlayerStateReq{
		PlayerId: uid,
	})

	// 检测返回值
	if err != nil {
		return "", err
	}

	return rsp.GetRoomAddr(), nil
}

// UpdatePlayerState 更新玩家状态
// param: uid 玩家ID, oldState 玩家当前状态， newState 要更新状态， gameID 游戏ID，levelID 场次ID
// return: 更新结果，错误信息
func UpdatePlayerState(uid uint64, oldState user.PlayerState, newState user.PlayerState, gameID uint32, levelID uint32) (bool, error) {

	// 得到服务连接
	con, err := getHallServer()
	if err != nil || con == nil {
		return false, errors.New("no connection")
	}

	// 新建Client
	client := user.NewPlayerDataClient(con)

	// 调用RPC方法
	rsp, err := client.UpdatePlayerState(context.Background(), &user.UpdatePlayerStateReq{
		PlayerId: uid,
		OldState: oldState,
		NewState: newState,
		GameId:   gameID,
		LevelId:  levelID,
	})

	// 检测返回值
	if err != nil {
		return false, err
	}

	if rsp.ErrCode != int32(user.ErrCode_EC_SUCCESS) {
		return false, errors.New("update player state failed")
	}
	// 更改机器人为未使用
	if newState == user.PlayerState_PS_IDIE { // 空闲
		flag, _ := robotclient.IsRobotPlayer(uid)
		if flag { // 是机器人
			robotclient.SetRobotPlayerState(uid, false)
		}
	}
	return true, nil
}

// UpdatePlayerGateInfo 更新玩家网关信息
// param: uid 玩家ID, ipAddr 客户端地址， gateAddr 网关服地址
// return: 更新结果，错误信息
func UpdatePlayerGateInfo(uid uint64, ipAddr, gateAddr string) (bool, error) {

	// 得到服务连接
	con, err := getHallServer()
	if err != nil || con == nil {
		return false, errors.New("no connection")
	}

	// 新建Client
	client := user.NewPlayerDataClient(con)

	// 调用RPC方法
	rsp, err := client.UpdatePlayerGateInfo(context.Background(), &user.UpdatePlayerGateInfoReq{
		PlayerId: uid,
		IpAddr:   ipAddr,
		GateAddr: gateAddr,
	})

	// 检测返回值
	if err != nil {
		return false, err
	}

	if rsp.ErrCode != int32(user.ErrCode_EC_SUCCESS) {
		return false, errors.New("update player state failed")
	}

	return true, nil
}

// UpdatePlayeServerAddr 更新玩家服务端地址
// param: uid 玩家ID, serverType 服务端类型 serverAddr 服务端地址
// return: 更新结果，错误信息
func UpdatePlayeServerAddr(uid uint64, serverType user.ServerType, serverAddr string) (bool, error) {

	// 得到服务连接
	con, err := getHallServer()
	if err != nil || con == nil {
		return false, errors.New("no connection")
	}

	// 新建Client
	client := user.NewPlayerDataClient(con)

	// 调用RPC方法
	rsp, err := client.UpdatePlayerServerAddr(context.Background(), &user.UpdatePlayerServerAddrReq{
		PlayerId:   uid,
		ServerType: serverType,
		ServerAddr: serverAddr,
	})

	// 检测返回值
	if err != nil {
		return false, err
	}

	if rsp.ErrCode != int32(user.ErrCode_EC_SUCCESS) {
		return false, errors.New("update player state failed")
	}

	return true, nil
}

// GetPlayerInfo 获取玩家信息
// param:   uid:玩家ID
// return:  玩家状态,错误信息
func GetPlayerInfo(uid uint64) (*user.GetPlayerInfoRsp, error) {

	// 得到服务连接
	con, err := getHallServer()
	if err != nil || con == nil {
		return nil, errors.New("no hall connection")
	}

	// 新建Client
	client := user.NewPlayerDataClient(con)

	// 调用RPC方法
	rsp, err := client.GetPlayerInfo(context.Background(), &user.GetPlayerInfoReq{
		PlayerId: uid,
	})

	// 检测返回值
	if err != nil || rsp == nil {
		return nil, err
	}

	if rsp.ErrCode != int32(user.ErrCode_EC_SUCCESS) {
		return nil, errors.New("GetPlayerInfo()成功，但rsp.ErrCode显示失败")
	}

	return rsp, nil
}

// GetPlayerGameInfo 获取玩家游戏信息
// param:   uid:玩家ID
// return:  玩家状态,错误信息
func GetPlayerGameInfo(uid uint64, gameID uint32) (*user.GetPlayerGameInfoRsp, error) {

	// 得到服务连接
	con, err := getHallServer()
	if err != nil || con == nil {
		return nil, errors.New("no hall connection")
	}

	// 新建Client
	client := user.NewPlayerDataClient(con)

	// 调用RPC方法
	rsp, err := client.GetPlayerGameInfo(context.Background(), &user.GetPlayerGameInfoReq{
		PlayerId: uid,
		GameId:   gameID,
	})

	// 检测返回值
	if err != nil || rsp == nil {
		return nil, err
	}

	if rsp.ErrCode != int32(user.ErrCode_EC_SUCCESS) {
		return nil, errors.New("GetPlayerGameInfo() ，但rsp.ErrCode显示失败")
	}

	return rsp, nil
}

// InitRobotPlayerState 初始化机器人玩家状态
func InitRobotPlayerState(robotPids []uint64) (*user.InitRobotPlayerStateRsp, error) {
	// 得到服务连接
	con, err := getHallServer()
	if err != nil || con == nil {
		return nil, errors.New("no hall connection")
	}

	// 新建Client
	client := user.NewPlayerDataClient(con)

	// 调用RPC方法
	rsp, err := client.InitRobotPlayerState(context.Background(), &user.InitRobotPlayerStateReq{
		RobotIds: robotPids,
	})

	// 检测返回值
	if err != nil || rsp == nil {
		return nil, err
	}

	if rsp.ErrCode != int32(user.ErrCode_EC_SUCCESS) {
		return nil, errors.New("InitRobotPlayerState() ，但rsp.ErrCode显示失败")
	}

	return rsp, nil

}

func getHallServer() (*grpc.ClientConn, error) {
	e := structs.GetGlobalExposer()
	// 对uid进行一致性hash路由策略.
	con, err := e.RPCClient.GetConnectByServerName(common.HallServiceName)
	if err != nil || con == nil {
		return nil, errors.New("no connection")
	}

	return con, nil
}
