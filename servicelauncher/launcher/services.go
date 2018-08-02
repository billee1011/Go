package launcher

import (
	gatewaycore "steve/gateway/core"
	goldcore "steve/gold/core"
	hallcore "steve/hall/core"
	logincore "steve/login/core"
	matchcore "steve/match/core"
	"steve/serviceloader/loader"
	"github.com/spf13/viper"
	"steve/structs/service"
	"github.com/Sirupsen/logrus"
	testcore "steve/testserver/core"

)

func Init(args []string, flagList map[string]*string) {
	// 处理命令行
	for k, v := range flagList {
		loader.SetArg(k, *v)
	}
	LoadService(args[0],
		loader.WithRPCParams(viper.GetString("rpc_certi_file"), viper.GetString("rpc_key_file"), viper.GetString("rpc_addr"), viper.GetInt("rpc_port"),
			viper.GetString("rpc_server_name")),
		loader.WithClientRPCCA(viper.GetString("rpc_ca_file"), viper.GetString("certi_server_name")),
		loader.WithRedisOption(viper.GetString("redis_addr"), viper.GetString("redis_passwd")),
		loader.WithConsulAddr(viper.GetString("consul_addr")),
		loader.WithPProf(viper.GetString("pprofExposeType"), viper.GetInt("pprofHttpPort")),
		loader.WithHealthPort(viper.GetInt("health_port")),
		loader.WithGroupName(viper.GetString("group_name")),

		loader.WithParams(args[1:]))

}

// LoadService load service appointed by name
func LoadService(name string, options ...loader.ServiceOption) {
	opt := loader.LoadOptions(options...)
	exposer := loader.CreateExposer(&opt)

	loader.RegisterServer2(&opt)
	loader.RegisterHealthServer(exposer.RPCServer)
	// service := initService(name, exposer)
	var svr service.Service
	switch name {
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
	case "gateway":
		svr = gatewaycore.NewService()
	case "gold":
		svr = goldcore.NewService()
	}
	if svr != nil {
		svr.Init(exposer)
		loader.Run(svr, exposer, opt)
	} else {
		logrus.Errorln("no service found service name : ", svr)
		panic("no service found")
	}
}
