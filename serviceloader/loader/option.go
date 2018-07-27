package loader

import (
	"reflect"

	"github.com/Sirupsen/logrus"
)

// Option 选项
type Option struct {
	node            int // 节点编号，同一种服务的节点编号不能相同
	rpcCertiFile    string
	rpcKeyFile      string
	rpcAddr         string // RPC服务监听地址
	rpcPort         int    // RPC端口号
	rpcServerName   string // 服务器名称
	params          []string
	rpcCAFile       string   // RPC客户端的CA文件
	rpcCAServerName string   // 证书中的服务器名称
	redisAddr       string   // redis 服务地址
	redisPasswd     string   // redis 密码
	consulAddr      string   // consul api 地址
	tags            []string // tags，用来注册 consul
	healthPort      int      // server health http port
	pprofExposeType string   // pprof 输出类型，空不输出
	pprofHttpPort   int      // pprof http输出端口
}

var defaultOption = Option{
	redisAddr:   "127.0.0.1:6379",
	redisPasswd: "",
	consulAddr:  "127.0.0.1:8500",
	tags:        []string{},
}

// ServiceOption ...
type ServiceOption func(opt *Option)

// WithTags with tags for register
func WithTags(tags []string) ServiceOption {
	return func(opt *Option) {
		opt.tags = tags
	}
}

// WithNode with node
func WithNode(node int) ServiceOption {
	return func(opt *Option) {
		opt.node = node
	}
}

// WithConsulAddr  with cosnul address
func WithConsulAddr(consulAddr string) ServiceOption {
	return func(opt *Option) {
		opt.consulAddr = consulAddr
	}
}

// WithHealthPort  with health port
func WithHealthPort(port int) ServiceOption {
	return func(opt *Option) {
		opt.healthPort = port
	}
}

// WithParams 参数选项， 参数将透传给 plugin
func WithParams(params []string) ServiceOption {
	return func(opt *Option) {
		opt.params = params
	}
}

// WithRedisOption 设置 redis 选项
func WithRedisOption(addr, passwd string) ServiceOption {
	return func(opt *Option) {
		opt.redisAddr = addr
		opt.redisPasswd = passwd
	}
}

// WithRPCParams RPC 选项， certiFile 为证书文件， keyFile 为私钥文件， addr 为 RPC 服务监听地址， port 为 RPC 服务监听端口
// serverName 为 RPC 服务名字
func WithRPCParams(certiFile string, keyFile string, addr string, port int, serverName string) ServiceOption {
	return func(opt *Option) {
		opt.rpcCertiFile = certiFile
		opt.rpcKeyFile = keyFile
		opt.rpcAddr = addr
		opt.rpcPort = port
		opt.rpcServerName = serverName
	}
}

// WithClientRPCCA 客户端 RPC CA 证书选项， caFile 为 CA 证书文件， serverName 为服务的证书域名字段
func WithClientRPCCA(caFile, serverName string) ServiceOption {
	return func(opt *Option) {
		opt.rpcCAFile = caFile
		opt.rpcCAServerName = serverName
	}
}

// WithPProf pprof配置
func WithPProf(exposeType string, httpPort int) ServiceOption {
	return func(opt *Option) {
		opt.pprofExposeType = exposeType
		opt.pprofHttpPort = httpPort
	}
}

// infoOption 输出选项信息
func infoOption(opt Option) {
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

// loadOptions 加载服务选项
func LoadOptions(options ...ServiceOption) Option {
	op := defaultOption

	for _, f := range options {
		f(&op)
	}
	infoOption(op)
	return op
}
