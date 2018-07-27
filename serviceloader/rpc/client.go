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

// 通过服务名和分组实现服务分组，比如实现room和match服务按照游戏ID分组。serviceloader为此作动态负载均衡
func (ccm *ClientConnMgr) GetConnectByGroupName(serverName string, groupName string) (*grpc.ClientConn, error) {
	addr, err := ccm.loadBalancer.getServerAddr(serverName + "," + groupName)
	if err != nil {
		return nil, err
	}
	return ccm.connectPool.getConnect(addr)
}


// 通过服务名和服务ID获取服务连接
func (ccm *ClientConnMgr) GetConnectByServerId(serverName string, serverId string) (*grpc.ClientConn, error) {
	addr, err := ccm.loadBalancer.getServerAddrByServerId(serverName, serverId)
	if err != nil {
		return nil, err
	}
	return ccm.connectPool.getConnect(addr)
}

// 通过服务名和组名获取服务列表，并对列表进行一致性Hash
func (ccm *ClientConnMgr) GetConnectByGroupHashId(serverName string, groupName string, hashId uint64) (*grpc.ClientConn, error) {
	addr, err := ccm.loadBalancer.getServerAddrByHashId(serverName + "," + groupName, hashId)
	if err != nil {
		return nil, err
	}
	return ccm.connectPool.getConnect(addr)
}

// 通过服务名获取服务列表，并对列表进行一致性Hash
func (ccm *ClientConnMgr) GetConnectByServerHashId(serverName string,  hashId uint64) (*grpc.ClientConn, error) {
	addr, err := ccm.loadBalancer.getServerAddrByHashId(serverName, hashId)
	if err != nil {
		return nil, err
	}
	return ccm.connectPool.getConnect(addr)
}


// GetConnectByAddr 根据地址获取连接
func (ccm *ClientConnMgr) GetConnectByAddr(addr string) (*grpc.ClientConn, error) {
	return ccm.connectPool.getConnect(addr)
}
