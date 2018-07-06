package core

import (
	"steve/structs"
	"steve/structs/service"
	"github.com/Sirupsen/logrus"
	"steve/match/register"
)


type matchCore struct {
	e *structs.Exposer
}

// NewService 创建服务
func NewService() service.Service {
	return new(matchCore)
}

func (c *matchCore) Init(e *structs.Exposer, param ...string) error {
	entry := logrus.WithField("name", "matchCore.Init")

	c.e = e

	if err := register.RegisterHandles(e.Exchanger); err != nil {
		entry.WithError(err).Error("注册消息处理器失败")
		return err
	}

	return nil
}

func (c *matchCore) Start() error {
	return nil
}

