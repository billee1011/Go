package sgrpc

import "google.golang.org/grpc"

// ServiceInfo ...
type ServiceInfo struct {
	ClientConn *grpc.ClientConn
}

// BindType used for BindService
type BindType string

const (
	// UserBind used for bind users to someserver
	UserBind = BindType("user")
	// DeskBind used for bind desks to someserver
	DeskBind = BindType("desk")
)

// RPCClient RPC客户端
type RPCClient interface {
	// GetServiceInfo query single service info.
	GetServiceInfo(serverName string, addr string) (*ServiceInfo, error)
	BindService(serverName string, bindType BindType, bindData string, addr string) error
	GetBindAddr(serverName string, bindType BindType, bindData string) string
	InvalidateBind(serverName string, bindType BindType, bindData string) error
}
