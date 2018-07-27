package structs

import (
	"steve/structs/configuration"
	"steve/structs/exchanger"
	"steve/structs/net"
	"steve/structs/pubsub"
	"steve/structs/redisfactory"
	"steve/structs/rpc"
	"steve/structs/sgrpc"
	"steve/structs/arg"
)

// Exposer provide common interfaces for services
type Exposer struct {
	RPCServer       sgrpc.RPCServer
	RPCClient       rpc.Client
	Configuration   configuration.Configuration
	WatchDogFactory net.WatchDogFactory
	Exchanger       exchanger.Exchanger
	RedisFactory    redisfactory.RedisFactory
	Publisher       pubsub.Publisher
	Subscriber      pubsub.Subscriber
	Option          arg.Option
}

var gExposer *Exposer

// GetGlobalExposer 获取全局 exposer 对象
// 全局对象在 servieloader 调用 Init 函数之前设置
func GetGlobalExposer() *Exposer {
	return gExposer
}

// SetGlobalExposer 设置全局 exposer
func SetGlobalExposer(e *Exposer) {
	gExposer = e
}
