package helper

import (
	"errors"
	"steve/common/data/redis"

	"github.com/Sirupsen/logrus"
)

var errRedisOperation = errors.New("redis 操作失败")

// AllocID 从 redis 中分配 ID
func AllocID(key string) (uint64, error) {
	entry := logrus.WithFields(logrus.Fields{
		"func_name": "AllocID",
	})
	redis := redis.GetRedisClient()
	cmd := redis.Incr(key)
	if cmd.Err() != nil {
		entry.WithError(cmd.Err()).Errorln(errRedisOperation)
		return 0, errRedisOperation
	}
	ID, err := cmd.Result()
	if err != nil {
		entry.WithError(err).Errorln(errRedisOperation)
		return 0, errRedisOperation
	}
	return uint64(ID), nil
}
