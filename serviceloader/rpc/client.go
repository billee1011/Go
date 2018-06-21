package rpc

import (
	"steve/structs/rpc"

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
	addr, err := ccm.loadBalancer.getServerAddr(serverName)
	if err != nil {
		return nil, err
	}
	return ccm.connectPool.getConnect(addr)
}

// GetConnectByAddr 根据地址获取连接
func (ccm *ClientConnMgr) GetConnectByAddr(addr string) (*grpc.ClientConn, error) {
	return ccm.connectPool.getConnect(addr)
}
