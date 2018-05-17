package loader

import (
	"fmt"
	"reflect"
	"steve/room/core"
	"steve/serviceloader/exchanger"
	"steve/serviceloader/net/watchdog"
	"steve/serviceloader/redisfactory"
	"steve/serviceloader/structs/configuration"
	"steve/serviceloader/structs/sgrpc"
	"steve/structs"
	iexchanger "steve/structs/exchanger"
	"steve/structs/proto/gate_rpc"
	"sync"

	"github.com/Sirupsen/logrus"
	"github.com/go-redis/redis"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type option struct {
	rpcCertiFile    string
	rpcKeyFile      string
	rpcAddr         string // RPC服务监听地址
	rpcPort         int    // RPC端口号
	rpcServerName   string // 服务器名称
	params          []string
	rpcCAFile       string // RPC客户端的CA文件
	rpcCAServerName string // 证书中的服务器名称
	redisAddr       string // redis 服务地址
	redisPasswd     string // redis 密码
}

var defaultOption = option{
	redisAddr:   "127.0.0.1:6379",
	redisPasswd: "",
}

// ServiceOption ...
type ServiceOption func(opt *option)

// WithParams 参数选项， 参数将透传给 plugin
func WithParams(params []string) ServiceOption {
	return func(opt *option) {
		opt.params = params
	}
}

// WithRedisOption 设置 redis 选项
func WithRedisOption(addr, passwd string) ServiceOption {
	return func(opt *option) {
		opt.redisAddr = addr
		opt.redisPasswd = passwd
	}
}

// WithRPCParams RPC 选项， certiFile 为证书文件， keyFile 为私钥文件， addr 为 RPC 服务监听地址， port 为 RPC 服务监听端口
// serverName 为 RPC 服务名字
func WithRPCParams(certiFile string, keyFile string, addr string, port int, serverName string) ServiceOption {
	return func(opt *option) {
		opt.rpcCertiFile = certiFile
		opt.rpcKeyFile = keyFile
		opt.rpcAddr = addr
		opt.rpcPort = port
		opt.rpcServerName = serverName
	}
}

// WithClientRPCCA 客户端 RPC CA 证书选项， caFile 为 CA 证书文件， serverName 为服务的证书域名字段
func WithClientRPCCA(caFile, serverName string) ServiceOption {
	return func(opt *option) {
		opt.rpcCAFile = caFile
		opt.rpcCAServerName = serverName
	}
}

func createConfiguration() *configuration.ConfigurationImpl {
	c, err := configuration.NewConfiguration()
	if err != nil {
		panic(err)
	}
	return c
}

func createRPCServer(opt *option, redisClient *redis.Client) (*sgrpc.RPCServerImpl, error) {
	if err := sgrpc.Setup(&sgrpc.Options{
		ServerName:  opt.rpcServerName,
		Addr:        opt.rpcAddr,
		Port:        opt.rpcPort,
		RedisClient: redisClient,
	}); err != nil {
		return nil, fmt.Errorf("启动 rpc 服务失败： %v", err)
	}
	rpcOption := []grpc.ServerOption{}
	if opt.rpcKeyFile != "" {
		cred, err := credentials.NewServerTLSFromFile(opt.rpcCertiFile, opt.rpcKeyFile)
		if err != nil {
			return nil, fmt.Errorf("create server tls failed. certifile=%s, keyfile=%s, err=%v", opt.rpcCertiFile, opt.rpcKeyFile, err)
		}
		rpcOption = append(rpcOption, grpc.Creds(cred))
	}
	return sgrpc.NewRPCServer(rpcOption...), nil
}

func createRPCClient(opt *option) *sgrpc.ClientImpl {
	return sgrpc.NewClientImpl(opt.rpcCAFile, opt.rpcCAServerName)
}

func createExchanger(rpcServer *sgrpc.RPCServerImpl, opt *option) (iexchanger.Exchanger, error) {
	h, e := exchanger.NewMessageHandlerServer()
	if err := rpcServer.RegisterService(steve_proto_gaterpc.RegisterMessageHandlerServer, h); err != nil {
		return nil, fmt.Errorf("注册服务失败： %v", err)
	}
	return e, nil
}

// func getPluginService(name string) (service.Service, error) {
// 	// return nil, nil // 调试
// 	if !strings.HasSuffix(name, ".so") {
// 		name += ".so"
// 	}
// 	p, err := plugin.Open(name)
// 	if err != nil {
// 		return nil, err
// 	}
// 	f, err := p.Lookup("GetService")
// 	if err != nil {
// 		return nil, err
// 	}
// 	getter := f.(func() service.Service)
// 	service := getter()
// 	return service, nil
// }

func infoOption(opt option) {

	fields := make(logrus.Fields)
	t := reflect.TypeOf(opt)
	v := reflect.ValueOf(opt)
	for i := 0; i < t.NumField(); i++ {
		ft := t.Field(i)
		fName := ft.Name
		fv := v.Field(i)
		fields[fName] = fv
	}
	logrus.WithFields(fields).Info("服务选项列表")
}

// LoadService load service appointed by name
func LoadService(name string, options ...ServiceOption) {
	opt := defaultOption
	for _, option := range options {
		option(&opt)
	}
	infoOption(opt)

	// service, err := getPluginService(name)
	// if err != nil {
	// 	panic(err)
	// }
	// 调试用
	service := core.NewService()

	redisFacotry := redisfactory.NewFactory(opt.redisAddr, opt.redisPasswd)
	redisClient, err := redisFacotry.NewClient()
	if err != nil {
		logrus.WithError(err).Panic("创建 redis 客户端失败")
	}

	rpcServer, err := createRPCServer(&opt, redisClient)
	if err != nil {
		panic(err)
	}
	if err = registerHealthServer(rpcServer); err != nil {
		panic(err)
	}

	exchanger, err := createExchanger(rpcServer, &opt)
	if err != nil {
		logrus.WithError(err).Panic("创建交互器失败")
	}

	exposer := structs.Exposer{
		Configuration:   createConfiguration(),
		RPCServer:       rpcServer,
		RPCClient:       createRPCClient(&opt),
		WatchDogFactory: watchdog.NewFactory(),
		Exchanger:       exchanger,
		RedisFactory:    redisFacotry,
	}
	structs.SetGlobalExposer(&exposer)

	err = service.Init(&exposer, opt.params...)
	if err != nil {
		panic(err)
	}

	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		if opt.rpcAddr != "" && opt.rpcPort != 0 {
			if err := rpcServer.Work(opt.rpcAddr, opt.rpcPort); err != nil {
				logrus.WithFields(logrus.Fields{
					"addr": opt.rpcAddr,
					"port": opt.rpcPort,
				}).WithError(err).Fatalln("RPC 服务启动失败")
			}
		} else {
			logrus.Info("未配置 RPC 地址或者端口，不启动 RPC 服务")
		}
	}()

	go func() {
		defer wg.Done()
		if err := service.Start(); err != nil {
			logrus.WithError(err).Fatalln("服务启动失败")
		}
	}()

	wg.Wait()
}
