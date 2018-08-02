package hallclient

import (
	"context"
	"errors"
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
// return:  玩家状态，正在进行的游戏ID,错误信息
func GetPlayerState(uid uint64) (user.PlayerState, uint32, error) {

	// 得到服务连接
	con, err := getHallServer()
	if err != nil || con == nil {
		return 0, 0, errors.New("no connection")
	}

	// 新建Client
	client := user.NewPlayerDataClient(con)

	// 调用RPC方法
	rsp, err := client.GetPlayerState(context.Background(), &user.GetPlayerStateReq{
		PlayerId: uid,
	})

	// 检测返回值
	if err != nil {
		return 0, 0, err
	}

	if rsp.ErrCode != int32(user.ErrCode_EC_SUCCESS) {
		return 0, 0, errors.New("get player state failed")
	}
	return rsp.State, rsp.GameId, nil
}

// UpdatePlayerState 更新玩家状态
// param: uid 玩家ID, oldState 玩家当前状态， newState 要更新状态， serverType 服务类型，serverAddr 服务地址
// return: 更新结果，错误信息
func UpdatePlayerState(uid uint64, oldState, newState, serverType uint32, serverAddr string) (bool, error) {

	// 得到服务连接
	con, err := getHallServer()
	if err != nil || con == nil {
		return false, errors.New("no connection")
	}

	// 新建Client
	client := user.NewPlayerDataClient(con)

	// 调用RPC方法
	rsp, err := client.UpdatePlayerState(context.Background(), &user.UpdatePlayerStateReq{
		PlayerId:   uid,
		OldState:   user.PlayerState(oldState),
		NewState:   user.PlayerState(newState),
		ServerType: user.ServerType(serverType),
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

// GetPlayerInfo 获取玩家个人资料信息
// params: uid 玩家ID
// return: 玩家个人资料，错误信息
func GetPlayerInfo(uid uint64) (user.GetPlayerInfoRsp, error) {
	// 得到服务连接
	con, err := getHallServer()
	if err != nil || con == nil {
		return user.GetPlayerInfoRsp{}, errors.New("no connection")
	}

	// 新建Client
	client := user.NewPlayerDataClient(con)

	// 调用RPC方法
	rsp, err := client.GetPlayerInfo(context.Background(), &user.GetPlayerInfoReq{
		PlayerId: uid,
	})

	// 检测返回值
	if err != nil {
		return user.GetPlayerInfoRsp{}, err
	}

	if rsp.ErrCode != int32(user.ErrCode_EC_SUCCESS) {
		return user.GetPlayerInfoRsp{}, errors.New("get player info failed")
	}

	return *rsp, nil
}

// GetGameListInfo 获取游戏列表信息
// reutn: 游戏列表，错误信息
func GetGameListInfo() (user.GetGameListInfoRsp, error) {

	// 得到服务连接
	con, err := getHallServer()
	if err != nil || con == nil {
		return user.GetGameListInfoRsp{}, errors.New("no connection")
	}

	// 新建Client
	client := user.NewPlayerDataClient(con)

	// 调用RPC方法
	rsp, err := client.GetGameListInfo(context.Background(), &user.GetGameListInfoReq{})

	// 检测返回值
	if err != nil {
		return user.GetGameListInfoRsp{}, err
	}

	if rsp.ErrCode != int32(user.ErrCode_EC_SUCCESS) {
		return user.GetGameListInfoRsp{}, errors.New("get game list info failed")
	}

	return *rsp, nil
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
