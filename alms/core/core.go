package core

/*
	功能： 服务控制逻辑中心，实现服务定义，Client消息分派初始化，和服务启动逻辑

*/
import (
	"steve/alms/data"
	"steve/structs"
	"steve/structs/service"

	"github.com/Sirupsen/logrus"
)

type AlmsCore struct {
}

// NewService 创建服务
func NewService() service.Service {
	return new(AlmsCore)
}

func (a *AlmsCore) Init(e *structs.Exposer, param ...string) error {
	entry := logrus.WithField("name", "AlmsCore.Init")
	// 注册客户端Client消息处理器
	if err := registerHandles(e.Exchanger); err != nil {
		entry.WithError(err).Error("注册客户端Client消息处理器失败")
		return err
	}
	// 获取救济金配置存入,存入redis，用于检验
	acd, err := data.GetDBAlmsConfigData()
	if err != nil {
		entry.WithError(err).Errorln("Init get alms config 失败")
		return err
	}
	// 存储到redis
	if err = data.SetAlmsConfigWatchs(data.AlmsConfigToMap(acd)); err != nil {
		entry.WithError(err).Errorln("Init set alms config redis 失败")
		return err
	}
	entry.Debugf("AlmsCoreserver init succeed ...")
	return nil
}

func (a *AlmsCore) Start() error {
	logrus.Debugf("AlmsCore server start succeed ...")
	return nil
}
