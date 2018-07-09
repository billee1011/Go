package watchdog

import (
	"fmt"
	"steve/gateway/config"
	"steve/structs"
	"steve/structs/net"

	"github.com/Sirupsen/logrus"
)

var gWatchDog net.WatchDog

// Get 获取 watch dog
func Get() net.WatchDog {
	return gWatchDog
}

// StartWatchDog 启动 Watch dog
func StartWatchDog(e *structs.Exposer, messageObserver net.MessageObserver, connectObserver net.ConnectObserver) error {
	listenIP := config.GetListenClientAddr()
	listenPort := config.GetListenClientPort()

	logEntry := logrus.WithFields(logrus.Fields{
		"listen_ip":   listenIP,
		"listen_port": listenPort,
	})

	gWatchDog = e.WatchDogFactory.NewWatchDog(&idAllocator{}, messageObserver, connectObserver)

	logEntry.Info("准备监听")
	addr := fmt.Sprintf("%s:%d", listenIP, listenPort)
	return gWatchDog.Start(addr, net.TCP)
}
