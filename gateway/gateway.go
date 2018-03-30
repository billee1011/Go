package main

import (
	"steve/structs"
	"steve/structs/net"
	"steve/structs/service"

	"github.com/Sirupsen/logrus"
)

type gateway struct{}

var _ service.Service = new(gateway)

func (gate *gateway) Start(e *structs.Exposer, param ...string) error {

	logrus.Debug("启动服务")

	wd := e.WatchDogFactory.NewWatchDog(nil, nil, nil)

	if err := wd.Start("localhost:36001", net.TCP); err != nil {
		return err
	}
	return nil
}

// GetService 获取服务接口，被 serviceloader 调用
func GetService() service.Service {
	return new(gateway)
}

func main() {}
