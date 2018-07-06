package redis

import (
	"sync"

	"github.com/Sirupsen/logrus"
	"github.com/go-redis/redis"
	"github.com/spf13/viper"
)

var gRedisClient *redis.Client
var redisOnce sync.Once

// GetRedisClient 获取 redis 客户端
func GetRedisClient() *redis.Client {
	redisOnce.Do(initRedisClient)
	return gRedisClient
}

func initRedisClient() {
	redisAddr := viper.GetString("redis_addr")
	redisPsw := viper.GetString("redis_passwd")

	entry := logrus.WithFields(logrus.Fields{
		"func_name":      "initRedisClient",
		"redis_addr":     redisAddr,
		"redis_password": redisPsw,
	})

	gRedisClient = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPsw,
		DB:       0,
	})
	cmd := gRedisClient.Ping()
	if cmd.Err() != nil {
		entry.Panicln("连接 redis 失败")
	}
}
