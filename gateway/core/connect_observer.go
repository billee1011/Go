package core

import (
	"steve/gutils/topics"
	userpb "steve/server_pb/user"
	"steve/structs"
	"steve/structs/net"
	"steve/structs/pubsub"

	"github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
)

type connectObserver struct{}

var _ net.ConnectObserver = new(connectObserver)

func (co *connectObserver) OnClientConnect(clientID uint64) {
	logrus.WithField("client_id", clientID).Info("client connected")
}

func (co *connectObserver) OnClientDisconnect(clientID uint64) {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "connectObserver.OnClientDisconnect",
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
	logEntry.Info("client disconnected")
}

func (co *connectObserver) getPublisher() pubsub.Publisher {
	exposer := structs.GetGlobalExposer()
	return exposer.Publisher
}
