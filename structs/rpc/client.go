package rpc

import "google.golang.org/grpc"

// Client RPC客户端
type Client interface {
	// 通过服务名称获取客户端连接。serviceloader为此作动态负载均衡
	GetConnectByServerName(serverName string) (*grpc.ClientConn, error)
	// 通过服务地址获取客户端连接。
	GetConnectByAddr(addr string) (*grpc.ClientConn, error)
}
