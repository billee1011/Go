package service

import "steve/structs"

// Service 是每个使用 serviceloader 的服务必须实现的接口，并且通过 GetService 函数导出
// serviceloader 在启动 service 时，会先调用 Init 函数作初始化， 然后调用 Start 函数，并且等待 Start 函数完成工作
type Service interface {

	// Init 作一些初始化操作。 参数 e 为服务提供了一些通用的接口。
	// 参数 param 是启动服务的命令行参数.
	// 如果该函数返回非 nil 错误， serviceloader 不会执行 Start 函数
	Init(e *structs.Exposer, param ...string) error

	// Start 启动服务，服务可以在该函数内部 block
	// serviceloader 在调用 Init 函数返回成功后，调用 Start 函数启动服务
	Start() error
}
