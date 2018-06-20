package player

import (
	"errors"
	"fmt"
	"steve/common/data/helper"
	"steve/common/data/redis"

	"github.com/Sirupsen/logrus"
)

var errRedisOperation = errors.New("redis 操作失败")

// fmtAccountPlayerKey 账号 ID 到玩家 ID 映射的 key
func fmtAccountPlayerKey(accountID uint64) string {
	return fmt.Sprintf("account:player:%v", accountID)
}

// GetAccountPlayerID 根据账号 ID 获取玩家 ID
func GetAccountPlayerID(accountID uint64) uint64 {
	redis := redis.GetRedisClient()
	key := fmtAccountPlayerKey(accountID)
	cmd := redis.Get(key)
	playerID, _ := cmd.Uint64()
	return playerID
}

// NewPlayer 创建玩家
func NewPlayer(accountID uint64, playerID uint64) error {
	entry := logrus.WithFields(logrus.Fields{
		"func_name": "AllocPlayerID",
	})
	redis := redis.GetRedisClient()
	key := fmtAccountPlayerKey(accountID)
	cmd := redis.SetNX(key, playerID, 0)
	if cmd.Err() != nil {
		entry.WithError(cmd.Err()).Errorln(errRedisOperation)
		return errRedisOperation
	}
	return nil
}

// AllocPlayerID 分配玩家 ID
func AllocPlayerID() (uint64, error) {
	return helper.AllocID("max_player_id")
}
