package structs

import (
	"steve/structs/configuration"
	"steve/structs/net"
	"steve/structs/sgrpc"
)

// Exposer provide common interfaces for services
type Exposer struct {
	RPCServer       sgrpc.RPCServer
	RPCClient       sgrpc.RPCClient
	Configuration   configuration.Configuration
	WatchDogFactory net.WatchDogFactory
}
