package core

import (
	localuser "steve/hall/user"
	"steve/server_pb/user"
	"steve/structs"
	"steve/structs/service"

	"github.com/Sirupsen/logrus"
)

type hallCore struct {
	e *structs.Exposer
}

// NewService 创建服务
func NewService() service.Service {
	return new(hallCore)
}

func (c *hallCore) Init(e *structs.Exposer, param ...string) error {
	entry := logrus.WithField("name", "hallCore.Init")

	// 注册当前模块RPC服务处理器
	if err := e.RPCServer.RegisterService(user.RegisterPlayerDataServer, localuser.PlayerDataService{}); err != nil {
		entry.WithError(err).Error("注册RPC服务处理器失败")
		return err
	}

	entry.Debugf("server init succeed ...")

	return nil
}

func (c *hallCore) Start() error {
	return nil
}
