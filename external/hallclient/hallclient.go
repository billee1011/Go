package userclient

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
// return:  玩家状态，正在进行的游戏ID,错误信息
func GetPlayerState(uid uint64) (uint32, uint32, error) {

	// 得到服务连接
	con, err := getHallServer(uid)
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
	return uint32(rsp.State), rsp.GameId, nil
}

// UpdatePlayerState 更新玩家状态
// param: uid:玩家ID, oldState 玩家当前状态， newState 要更新状态， serverType 服务类型，serverAddr 服务地址
// return: 更新结果，错误信息
func UpdatePlayerState(uid uint64, oldState, newState, serverType uint32, serverAddr string) (bool, error) {

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

func getHallServer(uid uint64) (*grpc.ClientConn, error) {
	e := structs.GetGlobalExposer()
	// 对uid进行一致性hash路由策略.
	con, err := e.RPCClient.GetConnectByServerHashId("hall", uid)
	if err != nil || con == nil {
		return nil, errors.New("no connection")
	}

	return con, nil
}
