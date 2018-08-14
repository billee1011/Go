package core

import (
	"context"
	"fmt"
	"steve/hall/logic"
	"steve/structs"
	"time"

	"github.com/go-redis/redis"
)

// showUID 最大展示uid
var showUID = "max_show_uid"

var playerRedisName = "player"

// InitServer 初始化服务
func InitServer() error {
	// redis设置showUID
	redisCli, err := getRedisCli(playerRedisName, 0)
	if err != nil {
		return fmt.Errorf("InitServer 获取 redis 客户端失败(%s)", err.Error())
	}
	_, err = redisCli.Get(showUID).Result()
	if err != nil {
		redisCli.Set(showUID, 10000*10000*10, -1)
	}

	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)

	// 初始化游戏场次配置
	go logic.InitGameConfig(ctx)

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
