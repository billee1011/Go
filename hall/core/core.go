package core

import (
	"steve/structs"
	"steve/structs/net"
	"steve/structs/service"

	"github.com/Sirupsen/logrus"
)

type hallCore struct {
	e *structs.Exposer

	dog net.WatchDog
}

// NewService 创建服务
func NewService() service.Service {
	return new(hallCore)
}

func (c *hallCore) Init(e *structs.Exposer, param ...string) error {
	entry := logrus.WithField("name", "hallCore.Init")

	c.e = e
	if err := registerHandles(e.Exchanger); err != nil {
		entry.WithError(err).Error("注册消息处理器失败")
		return err
	}
	return nil
}

func (c *hallCore) Start() error {
	return nil
}
