package redisfactory

import (
	"fmt"
	ifac "steve/structs/redisfactory"
	"sync"

	"github.com/Sirupsen/logrus"
	"github.com/spf13/cast"

	"github.com/spf13/viper"

	"github.com/go-redis/redis"
)

// redisConf redis 配置信息
type redisConf struct {
	address string
	passwd  string
}

type factory struct {
	defaultConf redisConf
	client      *redis.Client
	mu          sync.Mutex

	configs map[string]redisConf
	clients sync.Map
}

// NewFactory 创建 RedisFactory
func NewFactory(addr, passwd string) ifac.RedisFactory {
	f := &factory{
		defaultConf: redisConf{
			address: addr,
			passwd:  passwd,
		},
		configs: make(map[string]redisConf, 16),
	}
	f.init()
	return f
}

var _ ifac.RedisFactory = new(factory)

// NewClient 创建默认客户端
func (f *factory) NewClient() (*redis.Client, error) {
	if f.client != nil {
		return f.client, nil
	}

	f.mu.Lock()
	if f.client != nil {
		f.mu.Unlock()
		return f.client, nil
	}
	client, err := f.createClient(&f.defaultConf, 0)
	if err != nil {
		f.mu.Unlock()
		return nil, err
	}
	f.client = client
	f.mu.Unlock()
	return f.client, nil
}

// createClient 创建 redis.Client
func (f *factory) createClient(conf *redisConf, db int) (*redis.Client, error) {
	c := redis.NewClient(&redis.Options{
		Addr:     conf.address,
		Password: conf.passwd,
		DB:       db,
	})
	ping := c.Ping()
	if ping.Err() != nil {
		return nil, fmt.Errorf("连接 redis 服务器失败:%v", ping.Err())
	}
	return c, nil
}

// GetRedisClient 根据名字获取 redis 客户端
func (f *factory) GetRedisClient(name string, db int) (*redis.Client, error) {
	if _client, ok := f.clients.Load(name); ok {
		return _client.(*redis.Client), nil
	}
	conf, ok := f.configs[name]
	if !ok {
		return nil, fmt.Errorf("没有对应的配置")
	}
	client, err := f.createClient(&conf, db)
	if err != nil {
		return nil, err
	}
	actual, loaded := f.clients.LoadOrStore(name, client)
	if loaded {
		client.Close()
	}
	return actual.(*redis.Client), nil
}

func (f *factory) init() {
	configs := viper.GetStringMap("redis_list")

	for name, _conf := range configs {
		conf := cast.ToStringMapString(_conf)
		addr := cast.ToString(conf["addr"])
		passwd := cast.ToString(conf["passwd"])

		f.configs[name] = redisConf{
			address: addr,
			passwd:  passwd,
		}
	}
	logrus.WithField("configs", f.configs).Infoln("redis 配置初始化")
}
