package data

import (
	"fmt"
	"steve/entity/db"
	"steve/server_pb/user"
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

func getRedisByteVal(redisName string, key string) ([]byte, error) {
	redisCli, err := redisCliGetter(redisName, 0)
	if err != nil {
		return []byte{}, err
	}
	data, err := redisCli.Get(key).Bytes()
	if err == nil {
		return data, nil
	}
	return []byte{}, fmt.Errorf("redis 命令执行失败: %v", err)
}

func getRedisField(redisName string, key string, field ...string) ([]interface{}, error) {
	redisCli, err := redisCliGetter(redisName, 0)
	if err != nil {
		return nil, err
	}
	result, err := redisCli.HMGet(key, field...).Result()
	if err == nil {
		return result, nil
	}
	return nil, fmt.Errorf("redis 命令执行失败: %v", err)
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

// setRedisWatch 事务
func setRedisWatch(redisName string, key string, fields map[string]string, duration time.Duration) error {
	redisCli, err := redisCliGetter(redisName, 0)

	list := make(map[string]interface{}, len(fields))
	for k, v := range fields {
		list[k] = v
	}

	err = redisCli.Watch(func(tx *redis.Tx) error {
		err := tx.HKeys(key).Err()
		if err != nil && err != redis.Nil {
			return err
		}
		cmd := tx.HMSet(key, list)
		if cmd.Err() != nil {
			return fmt.Errorf("set redis watch err: %v ", cmd.Err())
		}
		redisCli.Expire(key, duration)
		return nil
	}, key)
	return err
}

func dbGameConfig2serverGameConfig(dbGameConfigs []db.TGameConfig) (gameInfos []*user.GameConfig) {
	gameInfos = make([]*user.GameConfig, 0)
	for _, dbGameConfig := range dbGameConfigs {
		gameInfo := &user.GameConfig{
			GameId:   uint32(dbGameConfig.Gameid),
			GameName: dbGameConfig.Name,
			GameType: uint32(dbGameConfig.Type),
		}

		gameInfos = append(gameInfos, gameInfo)
	}
	return
}

func dbGamelevelConfig2serverGameConfig(dbGameConfigs []db.TGameLevelConfig) (gamelevelConfigs []*user.GameLevelConfig) {
	gamelevelConfigs = make([]*user.GameLevelConfig, 0)
	for _, dbGameConfig := range dbGameConfigs {
		gamelevelConfig := &user.GameLevelConfig{
			GameId:     uint32(dbGameConfig.Gameid),
			LevelId:    uint32(dbGameConfig.Levelid),
			LevelName:  dbGameConfig.Name,
			BaseScores: uint32(dbGameConfig.Basescores),
			LowScores:  uint32(dbGameConfig.Lowscores),
			HighScores: uint32(dbGameConfig.Highscores),
			MinPeople:  uint32(dbGameConfig.Minpeople),
			MaxPeople:  uint32(dbGameConfig.Maxpeople),
		}

		gamelevelConfigs = append(gamelevelConfigs, gamelevelConfig)
	}
	return
}
