package player

import (
	"errors"
	"fmt"
	"steve/common/data/helper"
	"steve/common/data/redis"

	"github.com/Sirupsen/logrus"
)

const (
	// 玩家金币数字段名
	playerCoinField string = "coin"
	// playerGatewayAddrField 玩家网关地址字段名
	playerGatewayAddrField string = "gate_addr"
	// playerRoomAddrField 玩家所在 room 地址字段名
	playerRoomAddrField string = "room_addr"

	// playerNickNameField 玩家昵称字段
	playerNickNameField string = "nick_name"
)

var errRedisOperation = errors.New("redis 操作失败")

// fmtAccountPlayerKey 账号 ID 到玩家 ID 映射的 key
func fmtAccountPlayerKey(accountID uint64) string {
	return fmt.Sprintf("account:player:%v", accountID)
}

// fmtPlayerKey 返回玩家的 key
func fmtPlayerKey(playerID uint64) string {
	return fmt.Sprintf("player:%v", playerID)
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

// getPlayerUint64Field 获取玩家 uint64 字段值
func getPlayerUint64Field(playerID uint64, fieldName string) uint64 {
	entry := logrus.WithFields(logrus.Fields{
		"func_name":  "getPlayerUint64Field",
		"player_id":  playerID,
		"field_name": fieldName,
	})
	redis := redis.GetRedisClient()
	key := fmtPlayerKey(playerID)
	cmd := redis.HGet(key, fieldName)
	if cmd.Err() != nil {
		entry.WithError(cmd.Err()).Errorln(errRedisOperation)
		return 0
	}
	val, err := cmd.Uint64()
	if err != nil {
		entry.WithError(err).Errorln("数据格式错误")
		return 0
	}
	return val
}

func setPlayerUint64Field(playerID uint64, fieldName string, val uint64) error {
	entry := logrus.WithFields(logrus.Fields{
		"func_name":  "setPlayerUint64Field",
		"player_id":  playerID,
		"feild_name": fieldName,
		"val":        val,
	})
	redis := redis.GetRedisClient()
	key := fmtPlayerKey(playerID)
	cmd := redis.HSet(key, fieldName, val)
	if cmd.Err() != nil {
		entry.WithError(cmd.Err()).Errorln(errRedisOperation)
		return errRedisOperation
	}
	return nil
}

// getPlayerStringField 获取玩家 string 字段值
func getPlayerStringField(playerID uint64, fieldName string) string {
	entry := logrus.WithFields(logrus.Fields{
		"func_name":  "getPlayerStringField",
		"player_id":  playerID,
		"field_name": fieldName,
	})
	redis := redis.GetRedisClient()
	key := fmtPlayerKey(playerID)
	cmd := redis.HGet(key, fieldName)
	if cmd.Err() != nil {
		entry.WithError(cmd.Err()).Errorln(errRedisOperation)
		return ""
	}
	return cmd.Val()
}

// setPlayerStringField 设置玩家 string 字段值
func setPlayerStringField(playerID uint64, fieldName string, value string) error {
	redis := redis.GetRedisClient()
	key := fmtPlayerKey(playerID)
	cmd := redis.HSet(key, fieldName, value)
	return cmd.Err()
}

// GetPlayerCoin 获取玩家的金币数
func GetPlayerCoin(playerID uint64) uint64 {
	return getPlayerUint64Field(playerID, playerCoinField)
}

// SetPlayerCoin 设置玩家金币数
func SetPlayerCoin(playerID uint64, coin uint64) error {
	return setPlayerUint64Field(playerID, playerCoinField, coin)
}

// SetPlayerNickName 设置玩家昵称
func SetPlayerNickName(playerID uint64, nickName string) {
	setPlayerStringField(playerID, playerNickNameField, nickName)
}

// GetPlayerNickName 获取玩家昵称
func GetPlayerNickName(playerID uint64) string {
	return getPlayerStringField(playerID, playerNickNameField)
}

// GetPlayerGateAddr 获取玩家所在的网关地址
func GetPlayerGateAddr(playerID uint64) string {
	return getPlayerStringField(playerID, playerGatewayAddrField)
}

// SetPlayerGateAddr 设置玩家所在网关地址
func SetPlayerGateAddr(playerID uint64, addr string) error {
	return setPlayerStringField(playerID, playerGatewayAddrField, addr)
}

// GetPlayerRoomAddr 获取玩家所在 room 地址
func GetPlayerRoomAddr(playerID uint64) string {
	return getPlayerStringField(playerID, playerRoomAddrField)
}

// SetPlayerRoomAddr 设置玩家所在 room 地址
func SetPlayerRoomAddr(playerID uint64, addr string) error {
	return setPlayerStringField(playerID, playerRoomAddrField, addr)
}
