package data

import (
	"fmt"
	"steve/structs"

	"github.com/go-redis/redis"
)

const redisName = "steve"

func getRedisCli(redis string, db int) (*redis.Client, error) {
	exposer := structs.GetGlobalExposer()
	redisCli, err := exposer.RedisFactory.GetRedisClient(redis, db)
	if err != nil {
		return nil, fmt.Errorf("获取 redis client 失败: %v", err)
	}
	return redisCli, nil
}

// RedisCliGetter 单元测试通过这两个值修改 mysql 引擎获取和 redis cli 获取
var RedisCliGetter = getRedisCli

// SetPlayerMaxwinningstream 储存最大连胜
func SetPlayerMaxwinningstream(key string, maxStream int) error {
	redisCli, err := RedisCliGetter(redisName, 0)
	if err != nil {
		return err
	}
	err = redisCli.Set(key, maxStream, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

// GetPlayerMaxwinningstream 获取最大连胜
func GetPlayerMaxwinningstream(key string) (int, error) {
	redisCli, err := RedisCliGetter(redisName, 0)
	if err != nil {
		return 0, err
	}
	streamCmd := redisCli.Get(key)
	MaxStream, err := streamCmd.Int64()
	if err != nil {
		return 0, err
	}
	return int(MaxStream), nil
}
