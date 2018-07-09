package connection

import (
	"context"
	"fmt"
	"steve/common/data/player"
	"steve/gateway/config"
	"time"

	"github.com/Sirupsen/logrus"
)

const (
	// 多长时间没有检测到心跳断开连接
	heartBeatInterval time.Duration = time.Minute
	// 多长时间没有认证断开连接
	attachInterval time.Duration = time.Minute
)

// Connection 连接
type Connection struct {
	playerID       uint64
	clientID       uint64
	heartBeatTimer *time.Timer
	attachTimer    *time.Timer
	connMgr        *ConnMgr
}

func newConnection(clientID uint64, connMgr *ConnMgr) *Connection {
	return &Connection{
		clientID: clientID,
		connMgr:  connMgr,
	}
}

func (c *Connection) run(ctx context.Context, finish func()) {
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

func (c *Connection) detachPlayerConnect() {
	if c.playerID != 0 {
		player.SetPlayerGateAddr(c.playerID, "")
		c.connMgr.SetPlayerConnectionID(c.playerID, 0)
	}
}

func (c *Connection) kick(reason string, finish func()) {
	entry := logrus.WithFields(logrus.Fields{
		"player_id": c.playerID,
		"client_id": c.clientID,
		"reason":    reason,
	})
	entry.Infoln("踢出玩家")
	c.detachPlayerConnect()
	finish()
}

// GetPlayerID 获取绑定的玩家 ID
func (c *Connection) GetPlayerID() uint64 {
	return c.playerID
}

// AttachPlayer 绑定玩家 ID
func (c *Connection) AttachPlayer(playerID uint64) bool {
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
	c.connMgr.SetPlayerConnectionID(c.playerID, c.clientID)
	entry.Infoln("绑定成功")
	return true
}

// GetClientID 获取连接 ID
func (c *Connection) GetClientID() uint64 {
	return c.clientID
}

// HeartBeat 心跳
func (c *Connection) HeartBeat() {
	c.heartBeatTimer.Reset(heartBeatInterval)
}

func (c *Connection) getGatewayAddr() string {
	return fmt.Sprintf("%s:%d", config.GetRPCAddr(), config.GetRPCPort())
}
