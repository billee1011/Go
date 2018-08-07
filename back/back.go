package main

import (
	"steve/back/core"
	"steve/structs/service"
)

// GetService 获取服务接口，被 serviceloader 调用
func GetService() service.Service {
	return core.NewService()
}

func main() {}
