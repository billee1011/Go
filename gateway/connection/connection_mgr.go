package connection

import (
	"context"
	"fmt"
	"steve/gateway/interfaces"
	"steve/gutils/topics"
	"steve/structs"
	"steve/structs/pubsub"
	"sync"

	userpb "steve/server_pb/user"

	"github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
	"github.com/spf13/viper"
)

type connectionWithCancelFunc struct {
	connection *connection
	cancel     context.CancelFunc
}

type connectionMgr struct {
	kicker      func(clientID uint64)
	connections sync.Map // clientID: *connectionWithContext
}

// NewConnectionMgr 创建 ConnectionManager
func NewConnectionMgr() interfaces.ConnectionManager {
	return &connectionMgr{}
}

func (cm *connectionMgr) SetKicker(kicker func(clientID uint64)) {
	cm.kicker = kicker
}

func (cm *connectionMgr) OnClientConnect(clientID uint64) {
	logrus.WithField("client_id", clientID).Info("client connected")
	connection := newConnection(clientID)
	ctx, cancel := context.WithCancel(context.Background())
	cm.connections.Store(clientID, &connectionWithCancelFunc{
		connection: connection,
		cancel:     cancel,
	})
	connection.run(ctx, func() {
		cm.removeConnection(clientID)
		cm.kicker(clientID)
	})
}

func (cm *connectionMgr) GetConnection(clientID uint64) interfaces.Connection {
	_connection, ok := cm.connections.Load(clientID)
	if !ok || _connection == nil {
		return nil
	}
	connection := _connection.(*connectionWithCancelFunc)
	return connection.connection
}

func (cm *connectionMgr) OnClientDisconnect(clientID uint64) {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "connectionMgr.OnClientDisconnect",
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

func (cm *connectionMgr) pubDisconnect(clientID uint64) {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "connectionMgr.pubDisconnect",
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

func (cm *connectionMgr) getPublisher() pubsub.Publisher {
	exposer := structs.GetGlobalExposer()
	return exposer.Publisher
}

func (cm *connectionMgr) getRPCAddr() string {
	return fmt.Sprintf("%s:%d", viper.GetString("rpc_addr"), viper.GetInt("rpc_port"))
}

func (cm *connectionMgr) removeConnection(clientID uint64) {
	cm.connections.Delete(clientID)
}
