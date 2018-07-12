package launcher

import (
	"github.com/spf13/viper"
	gatewaycore "steve/gateway/core"
	matchcore "steve/match/core"
	roomcore "steve/room/core"
	"steve/serviceloader/loader"
	"steve/structs/service"
)

func Init(args []string) {
	LoadService(args[0],
		loader.WithRPCParams(viper.GetString("rpc_certi_file"), viper.GetString("rpc_key_file"), viper.GetString("rpc_addr"), viper.GetInt("rpc_port"),
			viper.GetString("rpc_server_name")),
		loader.WithClientRPCCA(viper.GetString("rpc_ca_file"), viper.GetString("certi_server_name")),
		loader.WithRedisOption(viper.GetString("redis_addr"), viper.GetString("redis_passwd")),
		loader.WithConsulAddr(viper.GetString("consul_addr")),
		loader.WithParams(args[1:]))
}

// LoadService load service appointed by name
func LoadService(name string, options ...loader.ServiceOption) {
	opt := loader.LoadOptions(options...)
	exposer := loader.CreateExposer(opt)

	loader.RegisterServer2(&opt)
	loader.RegisterHealthServer(exposer.RPCServer)
	// service := initService(name, exposer)
	var service service.Service
	switch name {
	case "match":
		service = matchcore.NewService()
	case "room":
		service = roomcore.NewService()
	case "gateway":
		service = gatewaycore.NewService()
	}
	loader.Run(service, exposer, opt)
}
