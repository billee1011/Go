package hallclient

import (
	"context"
	"errors"
	"steve/server_pb/user"
	"steve/structs"

	"google.golang.org/grpc"
)

/*
	功能：玩家服的Client API封装,实现调用
	作者： Zhengxuzhang
	日期： 2018-8-02
*/

// GetPlayerState 获取玩家状态
// param:   uid:玩家ID
// return:  玩家状态，正在进行的游戏ID,正在进行的场次ID,错误信息
func GetPlayerState(uid uint64) (user.PlayerState, uint32, uint32, error) {

	// 得到服务连接
	con, err := getHallServer(uid)
	if err != nil || con == nil {
		return 0, 0, 0, errors.New("no hall connection")
	}

	// 新建Client
	client := user.NewPlayerDataClient(con)

	// 调用RPC方法
	rsp, err := client.GetPlayerState(context.Background(), &user.GetPlayerStateReq{
		PlayerId: uid,
	})

	// 检测返回值
	if err != nil {
		return 0, 0, 0, err
	}

	if rsp.ErrCode != int32(user.ErrCode_EC_SUCCESS) {
		return 0, 0, 0, errors.New("get player state from hall failed")
	}
	return rsp.GetState(), rsp.GetGameId(), rsp.GetLevelId(), nil
}

// UpdatePlayerState 更新玩家状态
// param: uid:玩家ID, oldState 玩家当前状态， newState 要更新状态， serverType 服务类型，serverAddr 服务地址
// return: 更新结果，错误信息
func UpdatePlayerState(uid uint64, oldState user.PlayerState, newState user.PlayerState, gameID uint32, levelID uint32, serverType user.ServerType, serverAddr string) (bool, error) {

	// 得到服务连接
	con, err := getHallServer(uid)
	if err != nil || con == nil {
		return false, errors.New("no connection")
	}

	// 新建Client
	client := user.NewPlayerDataClient(con)

	// 调用RPC方法
	rsp, err := client.UpdatePlayerState(context.Background(), &user.UpdatePlayerStateReq{
		PlayerId:   uid,
		OldState:   oldState,
		NewState:   newState,
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
	con, err := getHallServer(uid)
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
	con, err := getHallServer(uid)
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

func getHallServer(uid uint64) (*grpc.ClientConn, error) {
	e := structs.GetGlobalExposer()
	// 对uid进行一致性hash路由策略.
	con, err := e.RPCClient.GetConnectByServerHashId("hall", uid)
	if err != nil || con == nil {
		return nil, errors.New("no connection")
	}

	return con, nil
}
