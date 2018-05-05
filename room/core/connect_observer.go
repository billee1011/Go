package core

import (
	"steve/structs/net"

	"github.com/Sirupsen/logrus"
)

type connectObserver struct{}

var _ net.ConnectObserver = new(connectObserver)

func (co *connectObserver) OnClientConnect(clientID uint64) {
	logrus.WithField("client_id", clientID).Info("client connected")
}

func (co *connectObserver) OnClientDisconnect(clientID uint64) {
	logrus.WithField("client_id", clientID).Info("client disconnected")
}
