package consul

import (
	"fmt"

	"github.com/go-redis/redis"
	consulapi "github.com/hashicorp/consul/api"
)

var gConsulClient *consulapi.Client

func initConsulClient() error {
	var err error
	gConsulClient, err = consulapi.NewClient(consulapi.DefaultConfig())
	if err != nil {
		return fmt.Errorf("new consul client failed: %v", err)
	}

	return nil
}

// Setup 启动服务发现， 包括启动服务收集以及注册服务
func Setup(serviceName string, addr string, port int, redisClient *redis.Client) error {
	if err := initConsulClient(); err != nil {
		return err
	}
	setupCollector()
	return registerService(serviceName, addr, port, redisClient)
}
