package core

import (
	"fmt"
	"steve/structs"

	"github.com/go-redis/redis"
)

// showUID 最大展示uid
var showUID = "max_show_uid"

var playerRedisName = "player"

// InitServer 初始化服务
func InitServer() error {
	// redis设置showUID
	// redisCli, err := getRedisCli(playerRedisName, 0)
	// if err != nil {
	// 	return fmt.Errorf("InitServer 获取 redis 客户端失败(%s)", err.Error())
	// }
	// redisCli.Set(showUID, 10000*10000*10, -1)
	return nil
}

func getRedisCli(redis string, db int) (*redis.Client, error) {
	exposer := structs.GetGlobalExposer()
	redisCli, err := exposer.RedisFactory.GetRedisClient(redis, db)
	if err != nil {
		return nil, fmt.Errorf("获取 redis client 失败: %v", err)
	}
	return redisCli, nil
}
