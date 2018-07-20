package sgrpc

// RPCServer RPC服务
type RPCServer interface {
	// RegisterService 使用protoc生成的函数注册服务
	// Example: RegisterService(RegisterRsctlServiceServer, someimpl)
	RegisterService(f interface{}, service interface{}) error

	// 设置退休
	EnableRetire(bEnable bool)
	// 是否退休
	IsRetire() bool

	// 设置负载值
	SetScore(score int64)
	// 获取负载值
	GetScore() int64
}

/*
// RegisterLoadReporterService 注册负载上报服务
// 服务可以使用 grpclb/reporter 下面的通用 Reporter 或者使用自己实现的 Reporter 来注册
func RegisterLoadReporterService(lps grpclb.LoadReportServer, s RPCServer) error {
	return s.RegisterService(grpclb.RegisterLoadReportServer, lps)
}
*/