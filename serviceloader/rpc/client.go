package rpc

import (
	"fmt"
	"steve/structs/rpc"
	"strings"

	"google.golang.org/grpc"
)

// ClientConnMgr 客户端连接管理器
type ClientConnMgr struct {
	loadBalancer *loadBalancer
	connectPool  *connectPool
}

// NewClient 创建对象
func NewClient(caFile string, tlsServerName string, consulAddr string) rpc.Client {
	return &ClientConnMgr{
		loadBalancer: newLoadBalancer(consulAddr),
		connectPool:  newConnectPool(caFile, tlsServerName),
	}
}

// GetConnectByServerName 根据服务名返回连接
func (ccm *ClientConnMgr) GetConnectByServerName(serverName string) (*grpc.ClientConn, error) {
	return ccm.getConnectByServerNameAndTags(serverName, nil)
}

// GetConnectByAddr 根据地址获取连接
func (ccm *ClientConnMgr) GetConnectByAddr(addr string) (*grpc.ClientConn, error) {
	return ccm.connectPool.getConnect(addr)
}

// GetServerAddr 根据服务名称和 tags 获取服务地址，如果有多个服务满足要求，serviceloader 为此作负载均衡
func (ccm *ClientConnMgr) GetServerAddr(serverName string, tags []string) (string, error) {
	tagstr := strings.Join(tags, ",")
	addr, err := ccm.loadBalancer.getServerAddr(serverName + "," + tagstr)
	return addr, err
}

// getConnectByServerNameAndTags 根据服务名称和 tags 获取连接，如果有多个服务满足要求，serviceloader 为此作负载均衡
func (ccm *ClientConnMgr) getConnectByServerNameAndTags(serverName string, tags []string) (*grpc.ClientConn, error) {
	addr, err := ccm.GetServerAddr(serverName, tags)
	if err != nil {
		return nil, fmt.Errorf("获取服务失败:%v", err)
	}
	return ccm.connectPool.getConnect(addr)
}
