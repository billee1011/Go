package sgrpc

// RPCServer RPC服务
type RPCServer interface {
	// RegisterService 使用protoc生成的函数注册服务
	// Example: RegisterService(RegisterRsctlServiceServer, someimpl)
	RegisterService(f interface{}, service interface{}) error
}
