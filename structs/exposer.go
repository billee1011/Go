package structs

import (
	"steve/structs/configuration"
	"steve/structs/exchanger"
	"steve/structs/net"
	"steve/structs/redisfactory"
	"steve/structs/sgrpc"
)

// Exposer provide common interfaces for services
type Exposer struct {
	RPCServer       sgrpc.RPCServer
	RPCClient       sgrpc.RPCClient
	Configuration   configuration.Configuration
	WatchDogFactory net.WatchDogFactory
	Exchanger       exchanger.Exchanger
	RedisFactory    redisfactory.RedisFactory
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
