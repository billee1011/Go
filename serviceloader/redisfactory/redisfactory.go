package redisfactory

import (
	"fmt"
	ifac "steve/structs/redisfactory"

	"github.com/go-redis/redis"
)

type factory struct {
	address string
	passwd  string
	client  *redis.Client
}

// NewFactory 创建 RedisFactory
func NewFactory(addr, passwd string) ifac.RedisFactory {
	return &factory{
		address: addr,
		passwd:  passwd,
	}
}

var _ ifac.RedisFactory = new(factory)

func (f *factory) NewClient() (*redis.Client, error) {
	if f.client != nil {
		return f.client, nil
	}
	c := redis.NewClient(&redis.Options{
		Addr:     f.address,
		Password: f.passwd,
		DB:       0,
	})
	result := c.Ping()
	if result == nil {
		return nil, fmt.Errorf("连接 redis 服务失败")
	}
	if result.Err() != nil {
		return nil, fmt.Errorf("连接 redis 服务失败: %v", result.Err())
	}
	f.client = c
	return c, nil
}
