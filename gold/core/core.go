package core
/*
	功能： 服务控制逻辑中心，实现服务定义，Client消息分派初始化，和服务启动逻辑

 */
import (
	"steve/structs"
	"steve/structs/service"
	"github.com/Sirupsen/logrus"
	"steve/structs/exchanger"
	"steve/client_pb/msgid"
)


// 全局控制总线
var gExposer *structs.Exposer
func GetExposer() *structs.Exposer {
	return gExposer
}

type goldCore struct {
}

// NewService 创建服务
func NewService() service.Service {
	return new(goldCore)
}

// 3.注册客户端Client Msg 消息分派
func (c *goldCore) dispatchClientMsg(e exchanger.Exchanger) error {

	if len(mapMsg) == 0 {
		return nil
	}
	regmsg := func(msgID msgid.MsgID, h interface{}) {
		if err := e.RegisterHandle(uint32(msgID), h); err != nil {
			logrus.WithField("msg_id", msgID).Panic(err)
		}
	}

	for k, v := range  mapMsg {
		regmsg(k, v)
	}

	return nil
}

// 服务初始化
func (c *goldCore) Init(e *structs.Exposer, param ...string) error {
	entry := logrus.WithField("name", "goldCore.Init")

	gExposer = e

	// 1.[RPC API]注册当前模块RPC服务处理器
	if pbService != nil {
		if err := e.RPCServer.RegisterService(pbService, pbServerImp); err != nil {
			entry.WithError(err).Error("注册RPC服务处理器失败")
			return err
		}
	}

	// 2.[C/S消息]分派客户端消息(Client Msg),进行MsgID -->Func()
	if err := c.dispatchClientMsg(e.Exchanger); err != nil {
		entry.WithError(err).Error("注册客户端Client消息处理器失败")
		return err
	}

	return nil
}

// 服务启动逻辑
func (c *goldCore) Start() error {
	return nil
}

