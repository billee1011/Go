package core

import (
	"net/http"
	"steve/match/register"
	mservice "steve/match/service"
	"steve/server_pb/match"
	"steve/structs"
	"steve/structs/service"

	"github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
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

	e.RPCServer.RegisterService(match.RegisterMatchServer, &mservice.MatchService{})

	return nil
}

func (c *matchCore) Start() error {
	httpAddr := viper.GetString("http_addr")
	if httpAddr != "" {
		go http.ListenAndServe(httpAddr, nil)
		logrus.Infoln("启动 http 服务")
	}
	return nil
}
