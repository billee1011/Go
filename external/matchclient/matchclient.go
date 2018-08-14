package matchclient

import (
	"context"
	"steve/server_pb/match"
	"steve/structs"
	"steve/structs/common"

	"github.com/Sirupsen/logrus"
	"google.golang.org/grpc"
)

// ClearAllMatch 清空所有的匹配
func ClearAllMatch() {

	// 获取match服的grpc连接
	con, err := getMatchServer()
	if err != nil || con == nil {
		logrus.WithError(err).Errorln("获取match服的grpc失败！！！")
		return
	}

	// 新建一个matchClient
	matchClient := match.NewMatchClient(con)

	// 调用RPC方法
	rsp, err := matchClient.ClearAllMatch(context.Background(), &match.ClearAllMatchReq{})

	// 检测返回值
	if err != nil || rsp == nil {
		logrus.WithError(err).Errorln("调用match服的ClearAllMatch()失败！！！")
		return
	}

	return
}

// 获取match服的grpc连接
func getMatchServer() (*grpc.ClientConn, error) {
	e := structs.GetGlobalExposer()
	con, err := e.RPCClient.GetConnectByServerName(common.MatchServiceName)
	if err != nil || con == nil {
		return nil, err
	}

	return con, nil
}
