package main

import (
	"steve/structs"
	"steve/structs/service"
)

type gateway struct{}

var _ service.Service = new(gateway)

func (gate *gateway) Start(e *structs.Exposer, param ...string) error {
	return nil
}

// GetService 获取服务接口，被 serviceloader 调用
func GetService() service.Service {
	return new(gateway)
}

func main() {}
