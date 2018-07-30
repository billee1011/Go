package data

import (
	"fmt"
	"steve/structs"
	"time"

	"github.com/go-redis/redis"
	"github.com/go-xorm/xorm"
)

func getRedisCli(redis string, db int) (*redis.Client, error) {
	exposer := structs.GetGlobalExposer()
	redisCli, err := exposer.RedisFactory.GetRedisClient(redis, db)
	if err != nil {
		return nil, fmt.Errorf("获取 redis client 失败: %v", err)
	}
	return redisCli, nil
}

func getMysqlEngine(mysqlName string) (*xorm.Engine, error) {
	exposer := structs.GetGlobalExposer()
	engine, err := exposer.MysqlEngineMgr.GetEngine(playerMysqlName)
	if err != nil {
		return nil, fmt.Errorf("获取 mysql 引擎失败：%v", err)
	}
	return engine, nil
}

// 单元测试通过这两个值修改 mysql 引擎获取和 redis cli 获取
var mysqlEngineGetter = getMysqlEngine
var redisCliGetter = getRedisCli

func getRedisUint64Val(redisName string, key string) (uint64, error) {
	redisCli, err := redisCliGetter(redisName, 0)
	if err != nil {
		return 0, err
	}
	redisCmd := redisCli.Get(key)
	if redisCmd.Err() == nil {
		val, err := redisCmd.Uint64()
		if err != nil {
			return 0, fmt.Errorf("获取 redis 数据失败")
		}
		return val, nil
	}
	return 0, fmt.Errorf("redis 命令执行失败: %v", redisCmd.Err())
}

func getRedisStringVal(redisName string, key string) (string, error) {
	redisCli, err := redisCliGetter(redisName, 0)
	if err != nil {
		return "", err
	}
	redisCmd := redisCli.Get(key)
	if redisCmd.Err() == nil {
		return redisCmd.String(), nil
	}
	return "", fmt.Errorf("redis 命令执行失败: %v", redisCmd.Err())
}

func hgetRedisUint64Val(redisName string, key string, field string) (uint64, error) {
	redisCli, err := redisCliGetter(redisName, 0)
	if err != nil {
		return 0, err
	}
	redisCmd := redisCli.HGet(key, field)
	if redisCmd.Err() == nil {
		playerID, err := redisCmd.Uint64()
		if err != nil {
			return 0, fmt.Errorf("获取 redis 数据失败")
		}
		return playerID, nil
	}
	return 0, fmt.Errorf("redis 命令执行失败: %v", redisCmd.Err())
}

func hmGetRedisFields(redisName string, key string, fields ...string) (map[string]interface{}, error) {
	redisCli, err := redisCliGetter(redisName, 0)
	if err != nil {
		return nil, err
	}
	vals, err := redisCli.HMGet(key, fields...).Result()
	if err != nil {
		return nil, fmt.Errorf("获取 redis 数据失败")
	}
	result := make(map[string]interface{}, len(fields))
	for i := 0; i < len(fields)-1; i++ {
		result[fields[i]] = vals[i]
	}
	return result, nil
}

func setRedisVal(redisName string, key string, val interface{}, duration time.Duration) error {
	redisCli, err := redisCliGetter(redisName, 0)
	if err != nil {
		return err
	}
	redisCmd := redisCli.Set(key, val, duration)
	if redisCmd.Err() != nil {
		return fmt.Errorf("redis 命令执行失败：%v", redisCmd.Err())
	}
	return nil
}

func hmSetRedisVal(redisName string, key string, fields map[string]interface{}) error {
	redisCli, err := redisCliGetter(redisName, 0)
	if err != nil {
		return err
	}
	err = redisCli.Watch(func(tx *redis.Tx) error {
		_, err := tx.Get(key).Uint64()
		if err != nil && err != redis.Nil {
			return fmt.Errorf("redis 查找key失败：%v", err)
		}

		_, err = tx.Pipelined(func(pipe redis.Pipeliner) error {
			pipe.HMSet(key, fields)
			return nil
		})
		return fmt.Errorf("redis 事务执行失败：%v", err)
	}, key)
	return err
}

func isKeyExists(redisName string, key string) bool {

	return false
}
