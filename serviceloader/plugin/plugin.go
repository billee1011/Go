package plugin

import (
	"steve/structs/service"
	"strings"

	"github.com/Sirupsen/logrus"
	"plugin"
	"steve/serviceloader/loader"
	"steve/structs"
)

// LoadService load service appointed by name
func LoadService(name string, options ...loader.ServiceOption) {
	opt := loader.LoadOptions(options...)
	exposer := loader.CreateExposer(&opt)

	loader.RegisterServer2(&opt)
	loader.RegisterHealthServer(exposer.RPCServer)
	svr := initService(name, exposer)
	loader.Run(svr, exposer, opt)
}

func initService(name string, ep *structs.Exposer) service.Service {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "initService",
		"name":      name,
	})
	svr, err := getPluginService(name)
	if err != nil {
		logEntry.Panicln(err)
	}
	if err := svr.Init(ep); err != nil {
		logEntry.Panicln(err)
	}
	logEntry.Infoln("初始化服务完成")
	return svr
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
	svr := getter()
	return svr, err
}
