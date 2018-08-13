package loader

import (
	"runtime/debug"
	"steve/serviceloader/mysql"
	"steve/serviceloader/net/watchdog"
	"steve/serviceloader/pubsub"
	"steve/serviceloader/redisfactory"
	"steve/serviceloader/structs/configuration"
	"steve/structs"
	"steve/structs/service"
	"sync"

	"github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
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

// 用启动命令行参数，取代文件配置项
func ArgReplaceOption(opt *Option) {
	// 如果命令行启动参数定义了服务ID，启用启动参数定义的服务ID
	port , ok := IntArg("port")
	if ok && port > 100 {
		opt.rpcPort = int(port)
	}
	hport , ok := IntArg("hport")
	if ok && hport > 100 {
		opt.healthPort = int(hport)
	}
	// 配置文件中的分组名称+启动参数中的分组ID，一起合成最后的分组ID
	groupArg, ok := StringArg("gid")
	if ok &&  len(groupArg) > 0 {
		if len(opt.groupName) > 0 {
			opt.groupName += ","
		}
		opt.groupName += groupArg
	}


}

// createExposer 创建 exposer 对象
func CreateExposer(opt *Option) *structs.Exposer {
	ArgReplaceOption(opt)
	exposer := &structs.Exposer{}
	exposer.Configuration = createConfiguration()
	exposer.RPCServer = createRPCServer(opt.rpcKeyFile, opt.rpcCertiFile)
	exposer.RPCClient = createRPCClient(opt.rpcCAFile, opt.rpcCAServerName, opt.consulAddr)
	exposer.WatchDogFactory = watchdog.NewFactory()
	exposer.Exchanger = createExchanger(exposer.RPCServer)
	exposer.RedisFactory = redisfactory.NewFactory(opt.redisAddr, opt.redisPasswd)
	exposer.MysqlEngineMgr = mysql.CreateMysqlEngineMgr()
	exposer.Publisher = pubsub.CreatePublisher()
	exposer.Subscriber = pubsub.CreateSubscriber()
	exposer.Option = opt
	exposer.ConsulReq = &ConsulRequestImp{}

	structs.SetGlobalExposer(exposer)
	// 开启通用的负载报告服务
	RegisterLBReporter(exposer.RPCServer)
	// Hash路由方式，将server Id 设置为负载值
	if viper.GetString("rpc_lb") == "hash" {
		ridArg, ok := IntArg("rid")
		if ok && ridArg >= 0 && ridArg < 10000{
			exposer.RPCServer.SetScore(int64(ridArg))
		}
	}

	return exposer
}

// Run 启动服务循环
func Run(service service.Service, exposer *structs.Exposer, opt Option) {
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		defer recoverPanic()
		// 从consul删除服务节点
		//defer DeleteMyConsulAgent()
		runRPCServer(exposer.RPCServer, opt.rpcAddr, opt.rpcPort)
	}()

	go func() {
		defer wg.Done()
		defer recoverPanic()
		runService(service)
	}()
	//time.Sleep(time.Second * 15)
	//exposer.RPCClient.GetConnectByServerHashId("gold", 1000)
	wg.Wait()
	// 从consul删除服务节点
	//DeleteMyConsulAgent()
}

func runService(s service.Service) {
	if err := s.Start(); err != nil {
		logrus.WithError(err).Fatalln("服务启动失败")
	}
}
