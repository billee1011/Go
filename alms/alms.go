package main

/*
	 功能：
		1. 服务关联启动代码。(./core/*.go)
		2. 服务通过plugin编译成so，并且通过serviceloader加载。(./core/*.go)
		3. 服务支持定义RPC服务。(./server/*.go)
		4. 服务支持处理Client请求消息。(./msg/*.go)
		5. 服务DB和redis逻辑代码。(/data/*.go)
		6. 业务逻辑代码。 (/logic/*.go)
		7. 常量定义代码。 (define/*.go)
		8. 服务支持下发通知消息给Client。
		9. 服务支持调用其他RPC服务API。(./external/*.go)
*/
import (
	"steve/alms/core"
	"steve/structs/service"
)

// GetService 获取服务接口，被 serviceloader 调用
func GetService() service.Service {
	return core.NewService()
}

func main() {}
