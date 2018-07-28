package router

import (
	"fmt"
	"steve/common/data/player"
	"steve/structs"

	"google.golang.org/grpc"
)

// Strategy 路由策略
type Strategy interface {
	// GetServerAddr 获取服务地址
	GetServerAddr(serverName string, playerID uint64, router uint32) (string, error)
}

// -----------------------------------------------------------------------------------------------------

// defaultStrategy 默认路由策略
type defaultStrategy struct {
	canRouter bool // 是否可以指定路由节点
}

func (ds *defaultStrategy) GetServerAddr(serverName string, playerID uint64, router uint32) (string, error) {
	rpcClient := structs.GetGlobalExposer().RPCClient

	var tags []string
	if router != 0 {
		tags = append(tags, fmt.Sprintf("node_%d", router))
	}
	addr, err := rpcClient.GetServerAddr(serverName, tags)
	if err != nil {
		return "", fmt.Errorf("获取服务地址失败，服务名：%s, tags:%v，错误：%v", serverName, tags, err)
	}
	return addr, nil
}

// -----------------------------------------------------------------------------------------------------

// roomStrategy room 服路由策略
type roomStrategy struct{}

func (rs *roomStrategy) GetServerAddr(serverName string, playerID uint64, router uint32) (string, error) {
	addr := player.GetPlayerRoomAddr(playerID)
	if addr == "" {
		s := defaultStrategy{}
		return s.GetServerAddr(serverName, playerID, router)
	}
	return addr, nil
}

// -----------------------------------------------------------------------------------------------------

var strategys = map[string]Strategy{
	"room":  &roomStrategy{},
	"match": &defaultStrategy{canRouter: true},
}

// -----------------------------------------------------------------------------------------------------

// GetConnection 获取连接
func GetConnection(serverName string, playerID uint64, router uint32) (*grpc.ClientConn, error) {
	strategy, exist := strategys[serverName]
	if !exist {
		strategy = &defaultStrategy{canRouter: false}
	}
	addr, err := strategy.GetServerAddr(serverName, playerID, router)
	if err != nil || addr == "" {
		return nil, fmt.Errorf("获取服务地址失败：%v", err)
	}

	exposer := structs.GetGlobalExposer()
	cc, err := exposer.RPCClient.GetConnectByAddr(addr)
	if err != nil {
		return nil, fmt.Errorf("获取服务连接失败，地址：%s, 错误：%v", addr, err)
	}
	return cc, nil
}
