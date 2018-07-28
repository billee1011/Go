package data

import (
	"fmt"
	"steve/entity/cache"
	"steve/structs"
	"time"

	"github.com/go-redis/redis"
)

func getPlayerRedisCli() (*redis.Client, error) {
	return structs.GetGlobalExposer().RedisFactory.GetRedisClient("player", 0)
}

// GetPlayerToken 获取玩家认证 token
func GetPlayerToken(playerID uint64) (string, error) {
	redisCli, err := getPlayerRedisCli()
	if err != nil {
		return "", fmt.Errorf("获取 redis 连接失败：%v", err)
	}
	tokenKey := cache.FmtPlayerTokenKey(playerID)
	cmd := redisCli.Get(tokenKey)
	if cmd.Err() != nil {
		return "", fmt.Errorf("获取 redis 数据失败: %v", cmd.Err())
	}
	return cmd.String(), nil
}

// SetPlayerToken 设置玩家认证 token
func SetPlayerToken(playerID uint64, token string, duration time.Duration) error {
	redisCli, err := getPlayerRedisCli()
	if err != nil {
		return fmt.Errorf("获取 redis 连接失败：%v", err)
	}
	tokenKey := cache.FmtPlayerTokenKey(playerID)
	cmd := redisCli.Set(tokenKey, token, duration)
	if cmd.Err() != nil {
		return fmt.Errorf("redis 设置数据失败：%v", cmd.Err())
	}
	return nil
}
