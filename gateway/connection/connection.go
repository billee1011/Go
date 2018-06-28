package connection

import (
	"context"
	"fmt"
	"steve/common/data/player"
	"steve/gateway/config"
	"steve/gateway/global"
	"time"

	"github.com/Sirupsen/logrus"
)

const (
	// 多长时间没有检测到心跳断开连接
	heartBeatInterval time.Duration = time.Minute
	// 多长时间没有认证断开连接
	attachInterval time.Duration = time.Minute
)

type connection struct {
	playerID       uint64
	clientID       uint64
	heartBeatTimer *time.Timer
	attachTimer    *time.Timer
}

func newConnection(clientID uint64) *connection {
	return &connection{
		clientID: clientID,
	}
}

func (c *connection) run(ctx context.Context, finish func()) {
	c.heartBeatTimer = time.NewTimer(heartBeatInterval)
	c.attachTimer = time.NewTimer(attachInterval)

	go func() {
		defer c.heartBeatTimer.Stop()
		defer c.attachTimer.Stop()
		select {
		case <-ctx.Done():
			{
				c.detachPlayerConnect()
				return
			}
		case <-c.heartBeatTimer.C:
			{
				c.kick("无心跳", finish)
			}
		case <-c.attachTimer.C:
			{
				c.kick("未认证", finish)
			}
		}
	}()
}

func (c *connection) detachPlayerConnect() {
	if c.playerID != 0 {
		player.SetPlayerGateAddr(c.playerID, "")
		global.GetPlayerManager().SetPlayerConnectionID(c.playerID, 0)
	}
}

func (c *connection) kick(reason string, finish func()) {
	entry := logrus.WithFields(logrus.Fields{
		"player_id": c.playerID,
		"client_id": c.clientID,
		"reason":    reason,
	})
	entry.Infoln("踢出玩家")
	c.detachPlayerConnect()
	finish()
}

func (c *connection) GetPlayerID() uint64 {
	return c.playerID
}

func (c *connection) AttachPlayer(playerID uint64) bool {
	entry := logrus.WithFields(logrus.Fields{
		"func_name":        "connection.AttachPlayer",
		"player_id":        c.playerID,
		"client_id":        c.clientID,
		"attach_player_id": playerID,
	})
	if c.playerID != 0 {
		entry.Infoln("已绑定")
		return false
	}
	c.playerID = playerID
	c.attachTimer.Stop()
	player.SetPlayerGateAddr(playerID, c.getGatewayAddr())
	global.GetPlayerManager().SetPlayerConnectionID(c.playerID, c.clientID)
	entry.Infoln("绑定成功")
	return true
}

func (c *connection) GetClientID() uint64 {
	return c.clientID
}

func (c *connection) HeartBeat() {
	c.heartBeatTimer.Reset(heartBeatInterval)
}

func (c *connection) getGatewayAddr() string {
	return fmt.Sprintf("%s:%d", config.GetRPCAddr(), config.GetRPCPort())
}
