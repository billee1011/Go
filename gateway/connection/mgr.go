package connection

import (
	"context"
	"fmt"
	"steve/gateway/watchdog"
	"steve/gutils/topics"
	"steve/structs"
	"steve/structs/net"
	"steve/structs/pubsub"
	"sync"

	userpb "steve/server_pb/user"

	"github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
	"github.com/spf13/viper"
)

type connectionWithCancelFunc struct {
	connection *Connection
	cancel     context.CancelFunc
}

// ConnMgr 连接管理
type ConnMgr struct {
	connections         sync.Map // clientID: *connectionWithContext
	playerConnectionMap sync.Map // playerID: clientID
}

var defaultConnectionMgr = &ConnMgr{}
var _ net.ConnectObserver = defaultConnectionMgr

// GetConnectionMgr 获取连接管理
func GetConnectionMgr() *ConnMgr {
	return defaultConnectionMgr
}

func (cm *ConnMgr) kickClient(clientID uint64) {
	dog := watchdog.Get()
	dog.Disconnect(clientID)
}

// OnClientConnect 客户端断开连接
func (cm *ConnMgr) OnClientConnect(clientID uint64) {
	logrus.WithField("client_id", clientID).Info("client connected")
	connection := newConnection(clientID, cm)
	ctx, cancel := context.WithCancel(context.Background())
	cm.connections.Store(clientID, &connectionWithCancelFunc{
		connection: connection,
		cancel:     cancel,
	})
	connection.run(ctx, func() {
		cm.removeConnection(clientID)
		cm.kickClient(clientID)
	})
}

// GetPlayerConnection 获取玩家的连接对象
func (cm *ConnMgr) GetPlayerConnection(playerID uint64) *Connection {
	_clientID, ok := cm.playerConnectionMap.Load(playerID)
	if !ok {
		return nil
	}
	clientID := _clientID.(uint64)
	return cm.GetConnection(clientID)
}

// SetPlayerConnectionID 设置玩家的连接
func (cm *ConnMgr) setPlayerConnectionID(playerID uint64, connectionID uint64) {
	cm.playerConnectionMap.Store(playerID, connectionID)
}

// GetConnection 获取连接
func (cm *ConnMgr) GetConnection(clientID uint64) *Connection {
	_connection, ok := cm.connections.Load(clientID)
	if !ok || _connection == nil {
		return nil
	}
	connection := _connection.(*connectionWithCancelFunc)
	return connection.connection
}

// OnClientDisconnect 连接断开
func (cm *ConnMgr) OnClientDisconnect(clientID uint64) {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "ConnMgr.OnClientDisconnect",
		"client_id": clientID,
	})
	logEntry.Info("client disconnected")
	_connection, ok := cm.connections.Load(clientID)
	if !ok || _connection == nil {
		return
	}
	connection := _connection.(*connectionWithCancelFunc)
	connection.cancel()
	cm.removeConnection(clientID)

	cm.pubDisconnect(clientID)
}

func (cm *ConnMgr) pubDisconnect(clientID uint64) {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "ConnMgr.pubDisconnect",
		"client_id": clientID,
	})
	pub := cm.getPublisher()
	message := userpb.ClientDisconnect{
		ClientId: clientID,
	}
	var data []byte
	var err error
	if data, err = proto.Marshal(&message); err != nil {
		logEntry.WithError(err).Errorln("序列化消息失败")
		return
	}
	if err := pub.Publish(topics.ClientDisconnect, data); err != nil {
		logEntry.WithError(err).Errorln("发布消息失败")
	}
}

func (cm *ConnMgr) getPublisher() pubsub.Publisher {
	exposer := structs.GetGlobalExposer()
	return exposer.Publisher
}

func (cm *ConnMgr) getRPCAddr() string {
	return fmt.Sprintf("%s:%d", viper.GetString("rpc_addr"), viper.GetInt("rpc_port"))
}

func (cm *ConnMgr) removeConnection(clientID uint64) {
	cm.connections.Delete(clientID)
}
