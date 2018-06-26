package connect

import (
	"errors"
	"fmt"
	"steve/common/data/helper"
	"steve/common/data/redis"

	"github.com/Sirupsen/logrus"
)

var errRedisOperation = errors.New("redis 操作失败")

const (
	connectKey      string = "connect"
	gateAddrField   string = "gate_addr"
	maxConnectIDKey string = "max_connect_id"
)

func fmtConnectKey(clientID uint64) string {
	return fmt.Sprintf("%s:%v", connectKey, clientID)
}

// GetConnectGatewayAddr 获取客户端连接所在的网关服 RPC 地址
func GetConnectGatewayAddr(clientID uint64) (string, error) {
	// entry := logrus.WithFields(logrus.Fields{
	// 	"func_name": "GetConnectGatewayAddr",
	// })
	key := fmtConnectKey(clientID)
	redis := redis.GetRedisClient()
	cmd := redis.HGet(key, gateAddrField)
	if cmd.Err() != nil {
		return "", nil
	}
	return cmd.Val(), nil
}

// SetConnectGatewayAddr 设置客户端连接所在的网关服 RPC 地址
func SetConnectGatewayAddr(clientID uint64, addr string) error {
	entry := logrus.WithFields(logrus.Fields{
		"func_name": "SetConnectGatewayAddr",
	})
	key := fmtConnectKey(clientID)
	redis := redis.GetRedisClient()
	cmd := redis.HSet(key, gateAddrField, addr)
	if cmd.Err() != nil {
		entry.WithError(cmd.Err()).Errorln(errRedisOperation)
		return errRedisOperation
	}
	return nil
}

// RemoveConnect 移除连接
func RemoveConnect(clientID uint64) error {
	entry := logrus.WithFields(logrus.Fields{
		"func_name": "SetConnectGatewayAddr",
	})
	key := fmtConnectKey(clientID)
	redis := redis.GetRedisClient()
	cmd := redis.Del(key)
	if cmd.Err() != nil {
		entry.WithError(cmd.Err()).Errorln(errRedisOperation)
		return errRedisOperation
	}
	return nil
}

// AllocConnectID 分配客户端连接 ID
func AllocConnectID() (uint64, error) {
	return helper.AllocID(maxConnectIDKey)
}
