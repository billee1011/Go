package core

import (
	"steve/room/config"
	"steve/room/interfaces/global"
	"steve/room/peipai"
	"steve/room/registers"
	"steve/server_pb/room_mgr"
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
	c.e = e
	global.SetMessageSender(e.Exchanger)
	registers.RegisterHandlers(e.Exchanger)
	// 使用seviceloader的通用的负载报告模块，
	//registerLbReporter(e)

	rpcServer := e.RPCServer
	err := rpcServer.RegisterService(roommgr.RegisterRoomMgrServer, &RoomService{})
	if err != nil {
		return err
	}

	return nil
}

func (c *roomCore) Start() error {
	go startPeipai()
	return nil
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

/*
func registerLbReporter(exposer *structs.Exposer) {
	if err := lb.RegisterLBReporter(exposer.RPCServer); err != nil {
		logrus.WithError(err).Panicln("注册负载上报服务失败")
	}
}
*/