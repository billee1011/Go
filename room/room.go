package main

import (
	"steve/room/core"
	_ "steve/room/desks"
	_ "steve/room/playermgr"
	_ "steve/room/req_event_translator"
	_ "steve/room/settle"
	"steve/structs/service"
)

// GetService 获取服务接口，被 serviceloader 调用
func GetService() service.Service {
	return core.NewService()
}

func main() {}
