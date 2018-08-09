package gameconfig

import (
	"fmt"
	"steve/entity/db"
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
	engine, err := exposer.MysqlEngineMgr.GetEngine(mysqlName)
	if err != nil {
		return nil, fmt.Errorf("获取 mysql 引擎失败：%v", err)
	}
	return engine, nil
}

// 单元测试通过这两个值修改 mysql 引擎获取和 redis cli 获取
var mysqlEngineGetter = getMysqlEngine
var redisCliGetter = getRedisCli

func getRedisField(redisName string, key string, field ...string) ([]interface{}, error) {
	redisCli, err := redisCliGetter(redisName, 0)
	if err != nil {
		return nil, err
	}
	result, err := redisCli.HMGet(key, field...).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, err
		}
		return nil, fmt.Errorf("redis 命令执行失败: %v", err)
	}
	return result, nil
}

func setRedisFields(redisName string, key string, fields map[string]string, duration time.Duration) error {
	redisCli, err := redisCliGetter(redisName, 0)
	if err != nil {
		return err
	}
	kv := make(map[string]interface{}, len(fields))
	for k, field := range fields {
		kv[k] = field
	}
	status := redisCli.HMSet(key, kv)
	if status.Err() != nil {
		return fmt.Errorf("设置失败(%v)", status.Err())
	}
	redisCli.Expire(key, duration)
	return nil
}

func generateDbPlayer(playerID uint64, info map[string]string, fields ...string) (dbPlayer *db.TPlayer, err error) {
	dbPlayer, err = new(db.TPlayer), nil
	for _, field := range fields {
		v, ok := info[field]
		if !ok {
			return nil, fmt.Errorf("错误的数据类型。field=%s val=%v", field, info)
		}
		if err = setDBPlayerByField(dbPlayer, field, v); err != nil {
			return nil, err
		}
	}
	return
}

func generateDbPlayerGame(playerID uint64, gameID uint32, info map[string]string, fields ...string) (dbPlayerGame *db.TPlayerGame, err error) {
	dbPlayerGame, err = new(db.TPlayerGame), nil

	for _, field := range fields {
		v, ok := info[field]
		if !ok {
			return nil, fmt.Errorf("错误的数据类型。field=%s val=%v", field, info)
		}
		if err = setDBPlayerGameByField(dbPlayerGame, field, v); err != nil {
			return nil, err
		}
	}
	return
}
