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
	robotRedisKey = "Robot:%v"

	//玩家ID字段名
	robotPlayerIDField string = "playerID"
	//玩家金币数字段名
	robotPlayerCoinField string = "coin"
	//玩家状态字段名
	robotPlayerStateField string = "state"
	//玩家游戏 ID 字段名
	robotPlayerGameIDField string = "game_id"
	//玩家昵称字段
	robotPlayerNickNameField string = "nick_name"
	//玩家头像字段
	robotPlayerHeadImageField string = "head_image"
)

//RedisClifunc 获取redisClien
var RedisClifunc = getRobotRedis
var errRobotRedisGain = errors.New("robot_redis 获取失败")
var errRobotRedisOpertaion = errors.New("robot_redis 操作失败")

// RedisTimeOut 过期时间
var RedisTimeOut = time.Hour * 24 * 30

// getRobotRedis 获取大厅服redis
func getRobotRedis() *redis.Client {
	e := structs.GetGlobalExposer()
	redis, err := e.RedisFactory.NewClient()
	if err != nil {
		logrus.WithError(err).Errorln(errRobotRedisGain)
		return nil
	}
	return redis
}

// SetRobotWatch 设置机器人属性
func SetRobotWatch(playerID uint64, fieldName string, val interface{}, duration time.Duration) error {
	entry := logrus.WithFields(logrus.Fields{
		"func_name": "SetRobotFiled",
		"playerID":  playerID,
		"fieldName": fieldName,
		"val":       val,
	})
	redisCli := RedisClifunc()
	key := fmt.Sprintf(robotRedisKey, playerID)
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
		return SetRobotWatch(playerID, fieldName, val, duration)
	}
	return err
}

// GetRobotFields 获取机器人多个属性
func GetRobotFields(playerID uint64, fields ...string) (map[string]interface{}, error) {
	entry := logrus.WithFields(logrus.Fields{
		"func_name": "GetRobotFields",
		"playerID":  playerID,
		"fields":    fields,
	})
	client := RedisClifunc()
	key := fmt.Sprintf(robotRedisKey, playerID)
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

// GetRobotStringFiled 获取机器人属性
func GetRobotStringFiled(playerID uint64, fieldName string) (string, error) {
	entry := logrus.WithFields(logrus.Fields{
		"func_name": "GetRobotStringFiled",
		"playerID":  playerID,
		"fieldName": fieldName,
	})
	client := RedisClifunc()
	key := fmt.Sprintf(robotRedisKey, playerID)
	val, err := client.HGet(key, fieldName).Result()
	if err != nil {
		entry.WithError(err).Errorln(errRobotRedisOpertaion)
		return "", err
	}
	return val, nil
}

// SetRobotPlayerWatchs 设置机器人玩家多个属性
func SetRobotPlayerWatchs(playerID uint64, fields map[string]interface{}, duration time.Duration) error {
	entry := logrus.WithFields(logrus.Fields{
		"func_name": "SetRobotPlayerFields",
		"playerID":  playerID,
		"fields":    fields,
	})
	redisCli := RedisClifunc()
	key := fmt.Sprintf(robotRedisKey, playerID)

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
		return SetRobotPlayerWatchs(playerID, fields, duration)
	}
	return err
}
