package core

import "github.com/Sirupsen/logrus"

type connection struct{}

func (c *connection) OnClientConnect(clientID uint64) {
	entry := logrus.WithFields(logrus.Fields{
		"func_name": "connection.OnClientConnect",
		"client_id": clientID,
	})
	entry.Infoln("客户端连接")
}

func (c *connection) OnClientDisconnect(clientID uint64) {
	entry := logrus.WithFields(logrus.Fields{
		"func_name": "connection.OnClientConnect",
		"client_id": clientID,
	})
	entry.Infoln("客户端断开连接")
}
