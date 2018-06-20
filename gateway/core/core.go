package core

import (
	"fmt"
	"steve/gateway/config"
	"steve/gateway/gateservice"
	"steve/gateway/register"
	"steve/server_pb/gateway"
	"steve/structs"
	"steve/structs/net"
	"steve/structs/proto/gate_rpc"
	"steve/structs/service"

	"github.com/Sirupsen/logrus"
	"github.com/spf13/viper"

	_ "steve/gateway/cpm" // 初始化 connectPlayerMap
)

type gatewayCore struct {
	e *structs.Exposer

	dog net.WatchDog
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
	return c.startWatchDog()
}

func (c *gatewayCore) registSender() error {
	return c.e.RPCServer.RegisterService(steve_proto_gaterpc.RegisterMessageSenderServer, &sender{
		core: c,
	})
}

func (c *gatewayCore) registerGateService() error {
	return c.e.RPCServer.RegisterService(gateway.RegisterGateServiceServer, gateservice.New())
}

func (c *gatewayCore) startWatchDog() error {
	listenIP := viper.GetString(config.ListenClientAddr)
	listenPort := viper.GetInt(config.ListenClientPort)

	logEntry := logrus.WithFields(logrus.Fields{
		"listen_ip":   listenIP,
		"listen_port": listenPort,
	})

	mo := &receiver{
		core: c,
	}
	co := &connectObserver{}

	c.dog = c.e.WatchDogFactory.NewWatchDog(&idAllocator{}, mo, co)
	if c.dog == nil {
		logEntry.Error("创建 watchdog 失败")
		return fmt.Errorf("创建 watchdog 失败")
	}
	logEntry.Info("准备监听")

	addr := fmt.Sprintf("%s:%d", listenIP, listenPort)
	return c.dog.Start(addr, net.TCP)
}
