package sgrpc

import (
	"fmt"
	"steve/serviceloader/structs/sgrpc/consul"

	"github.com/go-redis/redis"
)

// Options 启动参数
type Options struct {
	ServerName  string
	Addr        string
	Port        int
	RedisClient *redis.Client
}

// Setup 启动 grpc
// TODO : 修改调用方
func Setup(opt *Options) error {
	if err := consul.Setup(opt.ServerName, opt.Addr, opt.Port, opt.RedisClient); err != nil {
		return fmt.Errorf("init consul failed:%v", err)
	}
	return nil
}
