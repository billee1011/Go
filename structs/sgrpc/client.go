package sgrpc

import "google.golang.org/grpc"

// RPCClient RPC客户端
type RPCClient interface {
	GetClientConnByServerName(serverName string) (*grpc.ClientConn, error)
	GetClientConnByServerID(serverID string) (*grpc.ClientConn, error)
}
