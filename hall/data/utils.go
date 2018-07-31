package data

import (
	"fmt"
	"steve/entity/cache"
	"steve/server_pb/user"
	"steve/structs"
	"strconv"
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

func getRedisField(redisName string, key string, field string) (uint64, error) {
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
func setRedisWatch(redisName string, key string, val interface{}, duration time.Duration) error {
	redisCli, err := redisCliGetter(redisName, 0)
	err = redisCli.Watch(func(tx *redis.Tx) error {
		_, err := tx.Get(key).Result()
		if err != nil && err != redis.Nil {
			return err
		}

		_, err = tx.Pipelined(func(pipe redis.Pipeliner) error {
			pipe.Set(key, val, duration)
			return nil
		})
		return err
	}, key)
	return err
}

func trans2hallPlayer(cp *cache.HallPlayer, info map[string]string) {
	cp.NickName = info[cache.NickNameField]
	cp.Avatar = info[cache.AvatarField]
	gender, _ := strconv.ParseInt(info[cache.GenderField], 10, 16)
	cp.Gender = uint64(gender)
	cp.Name = info[cache.NameField]
	cp.Phone = info[cache.PhoneField]
}

func transToGameInfo(configs []gameConfigDetail) (gameConfigs []*user.GameConfigInfo) {
	gameConfigs = make([]*user.GameConfigInfo, 0)
	for _, config := range configs {
		gameConfigs = append(gameConfigs, &user.GameConfigInfo{
			GameId:     uint32(config.TGameConfig.Gameid),
			GameName:   config.TGameConfig.Name,
			GameType:   uint32(config.TGameConfig.Type),
			LevelId:    uint32(config.TGameLevelConfig.Levelid),
			BaseScores: uint32(config.TGameLevelConfig.Basescores),
			LowScores:  uint32(config.TGameLevelConfig.Lowscores),
			HighScores: uint32(config.TGameLevelConfig.Highscores),
			MinPeople:  uint32(config.TGameLevelConfig.Minpeople),
			MaxPeople:  uint32(config.TGameLevelConfig.Maxpeople),
		})
	}
	return gameConfigs
}
