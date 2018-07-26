package data

import (
	"errors"
	"steve/structs"

	"github.com/Sirupsen/logrus"
	"github.com/go-redis/redis"
)

var errHallRedisOpertaion = errors.New("hall_redis 操作失败")
var errHallRedisGain = errors.New("hall_redis 获取失败")

// GetHallRedis 获取大厅服redis
func GetHallRedis() *redis.Client {
	exposer := structs.GetGlobalExposer()
	redis, err := exposer.RedisFactory.NewClient()
	if err != nil {
		logrus.WithError(err).Errorln(errHallRedisGain)
		return nil
	}
	return redis
}

// GetAccountPlayerID 获取账号PlayerId
func GetAccountPlayerID(accountID uint64) (uint64, error) {
	entry := logrus.WithFields(logrus.Fields{
		"func_name": "GetAccountPlayerID",
		"accountID": accountID,
	})
	client := GetHallRedis()
	key := fmtAccountPlayerKey(accountID)
	cmd := client.Get(key)
	playerID, err := cmd.Uint64()
	if err != nil {
		entry.WithError(err).Errorln(errHallRedisOpertaion)
		return 0, err
	}
	return playerID, nil
}

// AllocIDIncr 获取递增的PlayerId
func AllocIDIncr(key string) (uint64, error) {
	entry := logrus.WithFields(logrus.Fields{
		"func_name": "AllocIdIncr",
		"key":       key,
	})
	client := GetHallRedis()
	cmd := client.Incr(key)
	if cmd.Err() != nil {
		entry.WithError(cmd.Err()).Errorln(errHallRedisOpertaion)
		return 0, errHallRedisOpertaion
	}
	ID, err := cmd.Result()
	if err != nil {
		entry.WithError(err).Errorln(errHallRedisOpertaion)
		return 0, errHallRedisOpertaion
	}
	return uint64(ID), nil
}

// NewPlayer 创建玩家
func NewPlayer(accountID uint64, playerID uint64) error {
	entry := logrus.WithFields(logrus.Fields{
		"func_name": "NewPlayer",
	})
	client := GetHallRedis()
	key := fmtAccountPlayerKey(accountID)
	cmd := client.SetNX(key, playerID, 0)
	if cmd.Err() != nil {
		entry.WithError(cmd.Err()).Errorln(errHallRedisOpertaion)
		return errHallRedisOpertaion
	}
	return nil
}

// SetPlayerFiled 设置玩家属性
func SetPlayerFiled(playerID uint64, fieldName string, val interface{}) error {
	entry := logrus.WithFields(logrus.Fields{
		"func_name": "SetPlayerFiled",
		"playerID":  playerID,
		"fieldName": fieldName,
		"val":       val,
	})
	client := GetHallRedis()
	key := fmtPlayerKey(playerID)

	err := client.Watch(func(tx *redis.Tx) error {
		_, err := tx.Get(key).Uint64()
		if err != nil && err != redis.Nil {
			entry.WithError(err).Errorln("key不存在")
			return err
		}

		_, err = tx.Pipelined(func(pipe redis.Pipeliner) error {
			pipe.HSet(key, fieldName, val)
			return nil
		})
		entry.WithError(err).Errorln("数据格式错误")
		return err
	}, key)
	if err == redis.TxFailedErr {
		entry.WithError(err).Errorln("重试")
		return SetPlayerFiled(playerID, fieldName, val)
	}
	return err
}

// SetPlayerFields 设置玩家多个属性
func SetPlayerFields(playerID uint64, fields map[string]interface{}) error {
	entry := logrus.WithFields(logrus.Fields{
		"func_name": "SetPlayerFileds",
		"playerID":  playerID,
		"fields":    fields,
	})
	client := GetHallRedis()
	key := fmtPlayerKey(playerID)

	err := client.Watch(func(tx *redis.Tx) error {
		_, err := tx.Get(key).Uint64()
		if err != nil && err != redis.Nil {
			entry.WithError(err).Errorln("key不存在")
			return err
		}

		_, err = tx.Pipelined(func(pipe redis.Pipeliner) error {
			pipe.HMSet(key, fields)
			return nil
		})
		entry.WithError(err).Errorln("数据格式错误")
		return err
	}, key)
	if err == redis.TxFailedErr {
		entry.WithError(err).Errorln("重试")
		return SetPlayerFields(playerID, fields)
	}
	return err
}

// GetPlayerFields 获取玩家属性
func GetPlayerFields(playerID uint64, fields ...string) (map[string]interface{}, error) {
	entry := logrus.WithFields(logrus.Fields{
		"func_name": "GetPlayerFields",
		"playerID":  playerID,
		"fields":    fields,
	})
	client := GetHallRedis()
	key := fmtPlayerKey(playerID)
	vals, err := client.HMGet(key, fields...).Result()
	if err != nil {
		entry.WithError(err).Errorln(errHallRedisOpertaion)
		return nil, err
	}
	result := make(map[string]interface{}, len(fields))
	for i := 0; i <= len(fields); i++ {
		result[fields[i]] = vals[i]
	}
	return result, nil
}

// GetPlayerStringFiled 获取玩家属性
func GetPlayerStringFiled(playerID uint64, fieldName string) (string, error) {
	entry := logrus.WithFields(logrus.Fields{
		"func_name": "GetPlayerFiled",
		"playerID":  playerID,
		"fieldName": fieldName,
	})
	client := GetHallRedis()
	key := fmtPlayerKey(playerID)
	val, err := client.HGet(key, fieldName).Result()
	if err != nil {
		entry.WithError(err).Errorln(errHallRedisOpertaion)
		return "", err
	}
	return val, nil
}

// GetPlayerUint64Filed 获取玩家属性
func GetPlayerUint64Filed(playerID uint64, fieldName string) (uint64, error) {
	entry := logrus.WithFields(logrus.Fields{
		"func_name": "GetPlayerUint64Filed",
		"playerID":  playerID,
		"fieldName": fieldName,
	})
	client := GetHallRedis()
	key := fmtPlayerKey(playerID)
	cmd := client.HGet(key, fieldName)
	if cmd.Err() != nil {
		entry.WithError(cmd.Err()).Errorln(errHallRedisOpertaion)
		return 0, cmd.Err()
	}
	val, err := cmd.Uint64()
	if err != nil {
		entry.WithError(err).Errorln("数据格式错误")
		return 0, err
	}
	return val, nil
}

// IsPlayerExists 玩家是否存在
func IsPlayerExists(playerID uint64) bool {
	entry := logrus.WithFields(logrus.Fields{
		"func_name": "IsPlayerExists",
		"playerID":  playerID,
	})
	client := GetHallRedis()
	key := fmtPlayerKey(playerID)
	cmd := client.Exists(key)
	if cmd.Err() != nil {
		entry.WithError(cmd.Err()).Errorln(errHallRedisOpertaion)
		return false
	}
	return cmd.Val() == 1
}
