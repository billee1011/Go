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
func createExposer(opt option) *structs.Exposer {
	exposer := &structs.Exposer{}
	exposer.Configuration = createConfiguration()
	exposer.RPCServer = createRPCServer(opt.rpcKeyFile, opt.rpcCertiFile)
	exposer.RPCClient = createRPCClient(opt.rpcCAFile, opt.rpcCAServerName)
	exposer.WatchDogFactory = watchdog.NewFactory()
	exposer.Exchanger = createExchanger(exposer.RPCServer)
	exposer.RedisFactory = redisfactory.NewFactory(opt.redisAddr, opt.redisPasswd)
	exposer.Publisher = pubsub.CreatePublisher()
	exposer.Subscriber = pubsub.CreateSubscriber()
	structs.SetGlobalExposer(exposer)
	return exposer
}

// run 启动服务循环
func run(service service.Service, exposer *structs.Exposer, opt option) {
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
	wg.Wait()
}

// LoadService load service appointed by name
func LoadService(name string, options ...ServiceOption) {
	opt := loadOptions(options...)
	exposer := createExposer(opt)

	registerServer(&registerParams{
		serverName:   opt.rpcServerName,
		addr:         opt.rpcAddr,
		port:         opt.rpcPort,
		redisFactory: exposer.RedisFactory,
	})
	registerHealthServer(exposer.RPCServer)
	service := initService(name, exposer)
	run(service, exposer, opt)
}
