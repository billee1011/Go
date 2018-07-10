package core

// import (
// 	"context"
// 	"time"

// 	"github.com/Sirupsen/logrus"
// )

// const (
// 	kickInterval      = time.Minute     // 连接后多长时间自动断开连接
// 	checkKickInterval = time.Second * 5 // 断开连接检测间隔
// )

// // connectTime 连接时间信息
// type connectTime struct {
// 	clientID uint64
// 	kickTime time.Time // 要被踢出的时间
// }

// type connection struct {
// 	// connectTimes sorted by kick time
// 	connectTimes []connectTime
// 	kicker       func(clientID uint64)
// 	freshConnect chan connectTime
// }

// func newConnectionMgr() *connection {
// 	return &connection{
// 		connectTimes: make([]connectTime, 0, 200),
// 		freshConnect: make(chan connectTime),
// 	}
// }

// func (c *connection) OnClientConnect(clientID uint64) {
// 	entry := logrus.WithFields(logrus.Fields{
// 		"func_name": "connection.OnClientConnect",
// 		"client_id": clientID,
// 	})
// 	entry.Infoln("客户端连接")
// 	c.freshConnect <- connectTime{
// 		clientID: clientID,
// 		kickTime: time.Now().Add(kickInterval),
// 	}
// }

// func (c *connection) OnClientDisconnect(clientID uint64) {
// 	entry := logrus.WithFields(logrus.Fields{
// 		"func_name": "connection.OnClientConnect",
// 		"client_id": clientID,
// 	})
// 	entry.Infoln("客户端断开连接")
// }

// func (c *connection) run(ctx context.Context) {
// 	ticker := time.NewTicker(checkKickInterval)

// 	defer ticker.Stop()

// 	for {
// 		select {
// 		case <-ctx.Done():
// 			return
// 		case <-ticker.C:
// 			c.kickClients()
// 		case ct := <-c.freshConnect:
// 			c.addFreshConnect(&ct)
// 		}
// 	}
// }

// func (c *connection) addFreshConnect(ct *connectTime) {
// 	c.connectTimes = append(c.connectTimes, *ct)
// }

// // setKicker set kicker function used for kicking connect
// func (c *connection) setKicker(kicker func(uint64)) {
// 	c.kicker = kicker
// }

// // kickClients kick clients on overdue
// func (c *connection) kickClients() {
// 	now := time.Now()
// 	kickCount := 0
// 	for _, ct := range c.connectTimes {
// 		if now.After(ct.kickTime) || now.Equal(ct.kickTime) {
// 			c.KickClient(ct.clientID)
// 			kickCount++
// 			continue
// 		}
// 		break
// 	}
// 	c.connectTimes = c.connectTimes[kickCount:]
// }

// func (c *connection) KickClient(clientID uint64) {
// 	if c.kicker != nil {
// 		c.kicker(clientID)
// 	}
// }
