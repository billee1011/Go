package sgrpc

import "google.golang.org/grpc"

// RPCClient RPC客户端
type RPCClient interface {
	// 通过服务名称获取客户端连接。serviceloader为此作动态负载均衡
	GetClientConnByServerName(serverName string) (*grpc.ClientConn, error)
	// 通过服务ID获取客户端连接。
	GetClientConnByServerID(serverID string) (*grpc.ClientConn, error)
}
