package core

import (
	"steve/gateway/connection"
	"steve/gateway/gateservice"
	"steve/gateway/register"
	"steve/gateway/watchdog"
	"steve/server_pb/gateway"
	"steve/structs"
	"steve/structs/proto/gate_rpc"
	"steve/structs/service"
)

type gatewayCore struct {
	e *structs.Exposer
}

// NewService 创建服务
func NewService() service.Service {
	return new(gatewayCore)
}

func (c *gatewayCore) Init(e *structs.Exposer, param ...string) error {
	c.e = e
	if err := c.registSender(); err != nil {
		return err
	}
	register.RegisteHandlers(e.Exchanger)
	return c.registerGateService()
}

func (c *gatewayCore) Start() error {
	return watchdog.StartWatchDog(c.e, &observer{}, connection.GetConnectionMgr())
	// return c.startWatchDog()
}

func (c *gatewayCore) registSender() error {
	return c.e.RPCServer.RegisterService(steve_proto_gaterpc.RegisterMessageSenderServer, &sender{})
}

func (c *gatewayCore) registerGateService() error {
	return c.e.RPCServer.RegisterService(gateway.RegisterGateServiceServer, gateservice.Default())
}
