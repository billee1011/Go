package data

import (
	"errors"
	"fmt"
	"steve/structs"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/go-redis/redis"
)

const (
	almsConfigKey = "alms_config_key" // 救济金配置redis key
	almsPlayerKey = "almsPlayerID:%v" // 救济玩家对应已经领取的数量
	// AlmsGetNorm 救济线
	AlmsGetNorm = "getNorm"
	// AlmsGetTimes 最多领取次数
	AlmsGetTimes = "getTimes"
	// AlmsGetNumber 领取数量
	AlmsGetNumber = "getNumber"
	// AlmsCountDonw   救济倒计时，时间是秒
	AlmsCountDonw = "almsCountDown"
	// DepositCountDonw 快充倒计时，时间是秒
	DepositCountDonw = "depositCountDown"
	// GameLeveConfigs 游戏场次配置
	GameLeveConfigs = "gameLeveConfigs"
	// AlmsVersion 救济金配置表版本号,初始1
	AlmsVersion = "version"
	// AlmsLowScores 下限金币
	AlmsLowScores = "lowscores"
)

//AlmsConfig redis 救济金配置
type AlmsConfig struct {
	GetNorm          int64             // 救济线
	GetTimes         int               // 最多领取次数
	GetNumber        int64             // 领取数量
	AlmsCountDonw    int               // 救济倒计时，时间是秒
	DepositCountDonw int               // 快充倒计时，时间是秒
	GameLeveConfigs  []*GameLeveConfig // 游戏场次是否开启救济金
	PlayerGotTimes   int               // 玩家已领取数量
	Version          int               // 救济金配置表版本号,初始1
}

//GameLeveConfig redis 游戏场次是否有救济金
type GameLeveConfig struct {
	GameID    int32 // 游戏ID
	LevelID   int32 // 场次ID
	LowScores int64 // 下限金币
	IsOpen    int   // 是否为救济金场，0：关闭，1：开启
}

var redisClifunc = getAlmsRedis //获取redisClien
var errRobotRedisGain = errors.New("robot_redis 获取失败")
var errRobotRedisOpertaion = errors.New("robot_redis 操作失败")

// RedisTimeOut 过期时间 1小时
var RedisTimeOut = time.Hour

// getAlmsRedis 获取redis
func getAlmsRedis() *redis.Client {
	e := structs.GetGlobalExposer()
	redis, err := e.RedisFactory.NewClient()
	if err != nil {
		logrus.WithError(err).Errorln(errRobotRedisGain)
		return nil
	}
	return redis
}

//GetAlmsPlayerGotTimes 获取玩家已经领取的数量
func GetAlmsPlayerGotTimes(playerID uint64) (int, error) {
	entry := logrus.WithFields(logrus.Fields{
		"func_name": "GetAlmsPlayerGotTimes",
		"playerID":  playerID,
	})
	client := redisClifunc()
	key := fmt.Sprintf(almsPlayerKey, playerID)
	data, err := client.Get(key).Int64()
	if err != nil {
		entry.WithError(err).Errorln("redis 命令执行失败")
		return 0, err
	}
	return int(data), nil
}

//UpdateAlmsPlayerGotTimes 修改玩家救济金已领取数量
func UpdateAlmsPlayerGotTimes(playerID uint64, val int, date time.Duration) error {
	entry := logrus.WithFields(logrus.Fields{
		"func_name": "UpdateAlmsPlayerGotTimes",
		"playerID":  playerID,
		"val":       val,
	})
	redisCli := redisClifunc()
	key := fmt.Sprintf(almsPlayerKey, playerID)
	err := redisCli.Watch(func(tx *redis.Tx) error {
		err := tx.Get(key).Err()
		if err != nil && err != redis.Nil {
			return err
		}
		_, err = tx.Pipelined(func(pipe redis.Pipeliner) error {
			pipe.Set(key, val, date)
			return nil
		})
		return err
	}, key)
	if err == redis.TxFailedErr {
		entry.WithError(err).Errorln("重试")
		return UpdateAlmsPlayerGotTimes(playerID, val, date)
	}
	return nil
}

// GetAlmsConfigFiled 获取救济金配置
func GetAlmsConfigFiled(fieldName string) (string, error) {
	entry := logrus.WithFields(logrus.Fields{
		"func_name": "GetAlmsConfigFiled",
		"fieldName": fieldName,
	})
	client := redisClifunc()
	key := almsConfigKey
	val, err := client.HGet(key, fieldName).Result()
	if err != nil {
		entry.WithError(err).Errorln(errRobotRedisOpertaion)
		return "", err
	}
	return val, nil
}

// GetAlmsConfigFileds 获取救济金配置多个属性
func GetAlmsConfigFileds(fields ...string) (map[string]interface{}, error) {
	entry := logrus.WithFields(logrus.Fields{
		"func_name": "GetAlmsConfigFileds",
		"fields":    fields,
	})
	client := redisClifunc()
	key := almsConfigKey
	vals, err := client.HMGet(key, fields...).Result()
	if err != nil {
		entry.WithError(err).Errorln(errRobotRedisOpertaion)
		return nil, err
	}
	result := make(map[string]interface{}, len(fields))
	for i := 0; i < len(fields); i++ {
		if vals[i] == nil {
			continue
		}
		result[fields[i]] = vals[i]
	}
	return result, nil
}

// SetAlmsConfigWatch 设置救济金配置
func SetAlmsConfigWatch(fieldName string, val interface{}, duration time.Duration) error {
	entry := logrus.WithFields(logrus.Fields{
		"func_name": "SetAlmsConfigWatch",
		"fieldName": fieldName,
		"val":       val,
	})
	redisCli := redisClifunc()
	key := almsConfigKey
	err := redisCli.Watch(func(tx *redis.Tx) error {
		err := tx.HKeys(key).Err()
		if err != nil && err != redis.Nil {
			return err
		}
		_, err = tx.Pipelined(func(pipe redis.Pipeliner) error {
			pipe.HSet(key, fieldName, val)
			return nil
		})
		redisCli.Expire(key, duration)
		return err
	}, key)
	if err == redis.TxFailedErr {
		entry.WithError(err).Errorln("重试")
		return SetAlmsConfigWatch(fieldName, val, duration)
	}
	return err
}

// SetAlmsConfigWatchs 设置救济金配置多个属性
func SetAlmsConfigWatchs(fields map[string]interface{}, duration time.Duration) error {
	entry := logrus.WithFields(logrus.Fields{
		"func_name": "SetAlmsConfigWatchs",
		"fields":    fields,
	})
	redisCli := redisClifunc()
	key := almsConfigKey
	err := redisCli.Watch(func(tx *redis.Tx) error {
		err := tx.HKeys(key).Err()
		if err != nil && err != redis.Nil {
			return err
		}
		_, err = tx.Pipelined(func(pipe redis.Pipeliner) error {
			pipe.HMSet(key, fields)
			return nil
		})
		redisCli.Expire(key, duration)
		return err
	}, key)
	if err == redis.TxFailedErr {
		entry.WithError(err).Errorln("重试")
		return SetAlmsConfigWatchs(fields, duration)
	}
	return err
}
