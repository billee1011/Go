package loader

import (
	"plugin"
	"steve/structs"
	"steve/structs/service"
	"strings"

	"github.com/Sirupsen/logrus"
)

func initService(name string, ep *structs.Exposer) service.Service {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "initService",
		"name":      name,
	})
	service, err := getPluginService(name)
	if err != nil {
		logEntry.Panicln(err)
	}
	if err := service.Init(ep); err != nil {
		logEntry.Panicln(err)
	}
	logEntry.Infoln("初始化服务完成")
	return service
}

func runService(s service.Service) {
	if err := s.Start(); err != nil {
		logrus.WithError(err).Fatalln("服务启动失败")
	}
}

func getPluginService(name string) (service.Service, error) {
	if !strings.HasSuffix(name, ".so") {
		name += ".so"
	}
	p, err := plugin.Open(name)
	if err != nil {
		return nil, err
	}
	f, err := p.Lookup("GetService")
	if err != nil {
		return nil, err
	}
	getter := f.(func() service.Service)
	service := getter()
	return service, nil
}
