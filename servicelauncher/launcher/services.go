package launcher

import (
	gatewaycore "steve/gateway/core"
	goldcore "steve/gold/core"
	hallcore "steve/hall/core"
	logincore "steve/login/core"
	matchcore "steve/match/core"
	testcore "steve/testserver/core"
	"steve/serviceloader/loader"
	"steve/structs/service"

	"github.com/Sirupsen/logrus"
	"steve/structs"
	"steve/servicelauncher/cmd"
)


// LoadService load service appointed by name
func LoadService() {
	var svr service.Service
	switch cmd.ServiceName {
	case "hall":
		svr = hallcore.NewService()
	case "login":
		svr = logincore.NewService()
	case "match":
		svr = matchcore.NewService()
	// case "room":
	// 	svr = roomcore.NewService()
	case "testserver":
		svr = testcore.NewService()
	case "msgserver":
		svr = msgcore.NewService()
	case "gateway":
		svr = gatewaycore.NewService()
	case "gold":
		svr = goldcore.NewService()
	}
	if svr != nil {
		exposer := structs.GetGlobalExposer()
		svr.Init(exposer)
		loader.Run(svr, exposer, cmd.Option)
	} else {
		logrus.Errorln("no service found service name : ", svr)
		panic("no service found")
	}
}
