package structs

import (
	"steve/structs/configuration"
	"steve/structs/logger"
	"steve/structs/sgrpc"
)

// Exposer provide common interfaces for services
type Exposer struct {
	RPCServer     sgrpc.RPCServer
	Configuration configuration.Configuration
	Logger        logger.Logger
}
