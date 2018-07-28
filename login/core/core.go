package core

import (
	"steve/login/loginservice"
	"steve/server_pb/login"
	"steve/structs"
	"steve/structs/service"
)

type loginService struct {
}

// NewService 创建服务
func NewService() service.Service {
	return &loginService{}
}

func (s *loginService) Init(e *structs.Exposer, param ...string) error {
	e.RPCServer.RegisterService(login.RegisterLoginServiceServer, loginservice.Default())
	return nil
}

func (s *loginService) Start() error {
	return nil
	// return s.startWatchDog()
}

// func (s *loginService) startWatchDog() error {
// 	listenIP := viper.GetString(config.ListenClientAddr)
// 	listenPort := viper.GetInt(config.ListenClientPort)
// 	logEntry := logrus.WithFields(logrus.Fields{
// 		"listen_ip":   listenIP,
// 		"listen_port": listenPort,
// 	})
// 	exposer := structs.GetGlobalExposer()

// 	mo := NewReceiver()
// 	co := newConnectionMgr()
// 	ctx, cancelFunc := context.WithCancel(context.Background())
// 	defer cancelFunc()

// 	dog := exposer.WatchDogFactory.NewWatchDog(nil, mo, co)
// 	if dog == nil {
// 		logEntry.Error("创建 watchdog 失败")
// 		return fmt.Errorf("创建 watchdog 失败")
// 	}
// 	co.setKicker(func(clientID uint64) {
// 		dog.Disconnect(clientID)
// 	})
// 	go co.run(ctx)

// 	global.SetMessageSender(&sender{
// 		watchDog: dog,
// 	})
// 	logEntry.Info("准备监听")

// 	addr := fmt.Sprintf("%s:%d", listenIP, listenPort)
// 	return dog.Start(addr, net.TCP)
// }
