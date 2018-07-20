package loader

import (
	"runtime/debug"
	"steve/serviceloader/net/watchdog"
	"steve/serviceloader/pubsub"
	"steve/serviceloader/redisfactory"
	"steve/serviceloader/structs/configuration"
	"steve/structs"
	"steve/structs/service"
	"sync"

	"github.com/Sirupsen/logrus"

)

func createConfiguration() *configuration.ConfigurationImpl {
	c, err := configuration.NewConfiguration()
	if err != nil {
		panic(err)
	}
	return c
}

// recoverPanic 异常恢复
func recoverPanic() {
	if x := recover(); x != nil {
		stack := debug.Stack()
		logrus.Errorln(string(stack))
	}
}

// createExposer 创建 exposer 对象
func CreateExposer(opt Option) *structs.Exposer {
	exposer := &structs.Exposer{}
	exposer.Configuration = createConfiguration()
	exposer.RPCServer = createRPCServer(opt.rpcKeyFile, opt.rpcCertiFile)
	exposer.RPCClient = createRPCClient(opt.rpcCAFile, opt.rpcCAServerName, opt.consulAddr)
	exposer.WatchDogFactory = watchdog.NewFactory()
	exposer.Exchanger = createExchanger(exposer.RPCServer)
	exposer.RedisFactory = redisfactory.NewFactory(opt.redisAddr, opt.redisPasswd)
	exposer.Publisher = pubsub.CreatePublisher()
	exposer.Subscriber = pubsub.CreateSubscriber()
	structs.SetGlobalExposer(exposer)
	// 开启通用的负载报告服务
	RegisterLBReporter(exposer.RPCServer)


	return exposer
}

// run 启动服务循环
func Run(service service.Service, exposer *structs.Exposer, opt Option) {
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		defer recoverPanic()
		runRPCServer(exposer.RPCServer, opt.rpcAddr, opt.rpcPort)
	}()

	go func() {
		defer wg.Done()
		defer recoverPanic()
		runService(service)
	}()
	//exposer.RPCClient.GetConnectByServerName("match")
	wg.Wait()
	// 从consul删除服务节点
	DeleteMyConsulAgent()
}

func runService(s service.Service) {
	if err := s.Start(); err != nil {
		logrus.WithError(err).Fatalln("服务启动失败")
	}
}