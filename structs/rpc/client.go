package rpc

import "google.golang.org/grpc"

// Client RPC客户端
type Client interface {
	// 通过服务名称获取客户端连接。serviceloader为此作动态负载均衡
	GetConnectByServerName(serverName string) (*grpc.ClientConn, error)

	// 通过服务名和分组实现服务分组，比如实现room和match服务按照游戏ID分组。serviceloader为此作动态负载均衡
	GetConnectByGroupName(serverName string, groupName string) (*grpc.ClientConn, error)
	// 通过服务名和服务ID获取服务连接
	GetConnectByServerId(serverName string, serverId string) (*grpc.ClientConn, error)

	// 通过服务名和组名获取服务列表，并对列表进行一致性Hash
	GetConnectByGroupHashId(serverName string, groupName string, hashId uint64) (*grpc.ClientConn, error)
	// 通过服务名获取服务列表，并对列表进行一致性Hash
	GetConnectByServerHashId(serverName string,  hashId uint64) (*grpc.ClientConn, error)

	// 通过服务地址获取客户端连接。
	GetConnectByAddr(addr string) (*grpc.ClientConn, error)
}
