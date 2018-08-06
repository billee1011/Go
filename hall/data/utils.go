package data

import (
	"fmt"
	"steve/entity/cache"
	"steve/entity/db"
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

func generateDbPlayer(playerID uint64, info map[string]string) *db.TPlayer {
	gender, _ := strconv.Atoi(info[cache.Gender])
	channelID, _ := strconv.Atoi(info[cache.ChannelID])
	provinceID, _ := strconv.Atoi(info[cache.ProvinceID])
	cityID, _ := strconv.Atoi(info[cache.CityID])

	return &db.TPlayer{
		Playerid:   int64(playerID),
		Gender:     int(gender),
		Nickname:   info[cache.NickName],
		Avatar:     info[cache.Avatar],
		Channelid:  int(channelID),
		Provinceid: int(provinceID),
		Cityid:     int(cityID),
	}
}

func generateDbPlayerGame(playerID uint64, gameID uint32, info map[string]string) *db.TPlayerGame {
	winningRate, _ := strconv.Atoi(info[cache.WinningRate])

	return &db.TPlayerGame{
		Playerid:    int64(playerID),
		Gameid:      int(gameID),
		Winningrate: int(winningRate),
	}
}
