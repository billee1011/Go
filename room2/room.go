package core

import (
	"steve/structs/service"
	"steve/structs"
	"steve/structs/net"
	"steve/server_pb/room_mgr"
	"github.com/spf13/viper"
	"github.com/Sirupsen/logrus"
	"steve/room2/util"
	"steve/room2/register"
	"steve/room2/common"
	"steve/room2/models"
	"steve/room2/fixed"
	_"steve/room2/contexts"
)

type roomCore struct {
	e   *structs.Exposer
	dog net.WatchDog
}

// GetService 获取服务接口，被 serviceloader 调用
func GetService() service.Service {
	return new(roomCore)
}

// NewService 创建服务
func NewService() service.Service {
	return new(roomCore)
}

func (c *roomCore) Init(e *structs.Exposer, param ...string) error {
	logrus.Info("room init")
	c.e = e
	util.SetMessageSender(e.Exchanger)
	registers.RegisterHandlers(e.Exchanger)
	//registerLbReporter(e)

	rpcServer := e.RPCServer
	deskRpc := models.GetDeskMgr()
	err := rpcServer.RegisterService(roommgr.RegisterRoomMgrServer, deskRpc)
	if err != nil {
		return err
	}

	return nil
}

/*func registerLbReporter(exposer *structs.Exposer) {
	if err := lb.RegisterLBReporter(exposer.RPCServer); err != nil {
		logrus.WithError(err).Panicln("注册负载上报服务失败")
	}
}*/

func (c *roomCore) Start() error {
	go startPeipai()
	return nil
}

func startPeipai() error {
	peipaiAddr := viper.GetString(fixed.ListenPeipaiAddr)
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "startPeipai",
		"addr":      peipaiAddr,
	})
	if peipaiAddr != "" {
		logEntry.Info("启动配牌服务")
		err := common.RunPeiPai(peipaiAddr)
		if err != nil {
			logEntry.WithError(err).Panic("配牌服务启动失败")
		}
		return err
	}
	logEntry.Info("未配置配牌")
	return nil
}

func main() {}
