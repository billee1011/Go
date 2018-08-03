package core

import (
	"steve/robot/data"
	"steve/robot/robotservice"
	"steve/server_pb/robot"
	"steve/structs"
	"steve/structs/service"

	"github.com/Sirupsen/logrus"
)

//RobotCore 机器人
type RobotCore struct {
	e *structs.Exposer
}

// NewService 创建服务
func NewService() service.Service {
	return new(RobotCore)
}

func (r *RobotCore) Init(e *structs.Exposer, param ...string) error {
	entry := logrus.WithField("name", "RobotCore.Init")
	r.e = e
	// 注册当前模块RPC服务处理器
	if err := e.RPCServer.RegisterService(robot.RegisterRobotServiceServer, robotservice.DefaultRobot()); err != nil {
		entry.WithError(err).Error("注册RPC服务处理器失败")
		return err
	}
	data.InitRobotRedis() //从mysql获取到机器人,存到redis
	entry.Debugf("RobotCoreserver init succeed ...")
	return nil
}

func (r *RobotCore) Start() error {
	return nil
}
