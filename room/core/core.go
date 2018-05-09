package core

import (
	"fmt"
	"steve/room/interfaces/global"
	"steve/room/registers"
	"steve/structs"
	"steve/structs/net"
	"steve/structs/service"

	"github.com/Sirupsen/logrus"
	"github.com/spf13/viper"

	_ "steve/room/desks"
	_ "steve/room/playermgr"
	_ "steve/room/req_event_translator"
	_ "steve/room/settle"
)

type roomCore struct {
	e         *structs.Exposer
	exchanger exchangerImpl

	dog net.WatchDog
}

// NewService 创建服务
func NewService() service.Service {
	return new(roomCore)
}

func (c *roomCore) Init(e *structs.Exposer, param ...string) error {
	logrus.Info("room init")
	c.e = e
	e.Exchanger = &c.exchanger
	structs.SetGlobalExposer(c.e)
	global.SetMessageSender(&c.exchanger)

	registers.RegisterHandlers(&c.exchanger)
	return nil
}

func (c *roomCore) Start() error {
	return c.startWatchDog()
}

func (c *roomCore) startWatchDog() error {
	listenIP := viper.GetString(ListenClientAddr)
	listenPort := viper.GetInt(ListenClientPort)

	logEntry := logrus.WithFields(logrus.Fields{
		"listen_ip":   listenIP,
		"listen_port": listenPort,
	})

	mo := &messageObserver{
		core: c,
	}
	co := &connectObserver{}

	c.dog = c.e.WatchDogFactory.NewWatchDog(nil, mo, co)
	c.exchanger.watchDog = c.dog

	if c.dog == nil {
		logEntry.Error("创建 watchdog 失败")
		return fmt.Errorf("创建 watchdog 失败")
	}
	logEntry.Info("准备监听")

	addr := fmt.Sprintf("%s:%d", listenIP, listenPort)
	return c.dog.Start(addr, net.TCP)
}