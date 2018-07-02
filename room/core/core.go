package core

import (
	"fmt"
	"steve/room/config"
	"steve/room/core/exchanger"
	"steve/room/interfaces/global"
	"steve/room/loader_balancer"
	"steve/room/peipai"
	"steve/room/registers"
	"steve/structs"
	"steve/structs/net"
	"steve/structs/service"

	"github.com/Sirupsen/logrus"
	"github.com/spf13/viper"

	_ "steve/room/autoevent" // 引入 autoevent 包，设置工厂
	_ "steve/room/desks/factory"
	_ "steve/room/desks/mgr"
	_ "steve/room/playermgr"
	_ "steve/room/req_event_translator"
	_ "steve/room/settle"
)

var flags struct {
	useGateway bool
}

func parseFlags() {
	flags.useGateway = !viper.GetBool("independent")

	logrus.WithField("flag", flags).Infoln("解析 flags ")
}

type roomCore struct {
	e   *structs.Exposer
	dog net.WatchDog
}

// NewService 创建服务
func NewService() service.Service {
	return new(roomCore)
}

func (c *roomCore) Init(e *structs.Exposer, param ...string) error {
	logrus.Info("room init")
	parseFlags()
	c.e = e
	if !flags.useGateway {
		e.Exchanger = exchanger.CreateLocalExchanger(&connectObserver{})
	}
	global.SetMessageSender(e.Exchanger)
	registers.RegisterHandlers(e.Exchanger)
	registerLbReporter(e)
	return nil
}

func (c *roomCore) Start() error {
	go startPeipai()
	if !flags.useGateway {
		return c.startLocalExchanger()
	}
	return nil
}

func (c *roomCore) startLocalExchanger() error {
	listenIP := viper.GetString(config.ListenClientAddr)
	listenPort := viper.GetInt(config.ListenClientPort)

	logEntry := logrus.WithFields(logrus.Fields{
		"listen_ip":   listenIP,
		"listen_port": listenPort,
	})
	logEntry.Info("准备监听")

	addr := fmt.Sprintf("%s:%d", listenIP, listenPort)
	return exchanger.StartLocalExchanger(c.e.Exchanger, addr, net.TCP)
}

func startPeipai() error {
	peipaiAddr := viper.GetString(config.ListenPeipaiAddr)
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "startPeipai",
		"addr":      peipaiAddr,
	})
	if peipaiAddr != "" {
		logEntry.Info("启动配牌服务")
		err := peipai.Run(peipaiAddr)
		if err != nil {
			logEntry.WithError(err).Panic("配牌服务启动失败")
		}
		return err
	}
	logEntry.Info("未配置配牌")
	return nil
}

func registerLbReporter(exposer *structs.Exposer) {
	if err := lb.RegisterLBReporter(exposer.RPCServer); err != nil {
		logrus.WithError(err).Panicln("注册负载上报服务失败")
	}
}
