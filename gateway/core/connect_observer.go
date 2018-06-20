package core

import (
	"fmt"
	"steve/common/data/connect"
	"steve/gateway/global"
	"steve/gutils/topics"
	userpb "steve/server_pb/user"
	"steve/structs"
	"steve/structs/net"
	"steve/structs/pubsub"

	"github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
	"github.com/spf13/viper"
)

type connectObserver struct{}

var _ net.ConnectObserver = new(connectObserver)

func (co *connectObserver) OnClientConnect(clientID uint64) {
	logrus.WithField("client_id", clientID).Info("client connected")
	co.saveClientGateAddr(clientID)
}

func (co *connectObserver) saveClientGateAddr(clientID uint64) {
	entry := logrus.WithFields(logrus.Fields{
		"func_name": "connectObserver.saveClientGateAddr",
		"client_id": clientID,
	})
	if err := connect.SetConnectGatewayAddr(clientID, co.getRPCAddr()); err != nil {
		entry.WithError(err).Errorln("保存连接 ID 和网关 RPC 地址的映射关系失败")
	}
}

func (co *connectObserver) getRPCAddr() string {
	return fmt.Sprintf("%s:%d", viper.GetString("rpc_addr"), viper.GetInt("rpc_port"))
}

func (co *connectObserver) OnClientDisconnect(clientID uint64) {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "connectObserver.OnClientDisconnect",
		"client_id": clientID,
	})
	co.pubDisconnect(clientID)
	co.removeConnectPlayer(clientID)
	co.removeConnectGatewayAddr(clientID)
	logEntry.Info("client disconnected")
}

func (co *connectObserver) removeConnectPlayer(clientID uint64) {
	cpm := global.GetConnectPlayerMap()
	cpm.RemoveConnect(clientID)
}

func (co *connectObserver) pubDisconnect(clientID uint64) {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "connectObserver.pubDisconnect",
		"client_id": clientID,
	})
	pub := co.getPublisher()
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

func (co *connectObserver) getPublisher() pubsub.Publisher {
	exposer := structs.GetGlobalExposer()
	return exposer.Publisher
}

func (co *connectObserver) removeConnectGatewayAddr(clientID uint64) {
	if err := connect.RemoveConnect(clientID); err != nil {
		logrus.WithFields(logrus.Fields{
			"func_name": "removeConnectGatewayAddr",
			"client_id": clientID,
		}).WithError(err).Errorln("移除连接失败")
	}
}
