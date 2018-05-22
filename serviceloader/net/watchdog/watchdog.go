package watchdog

import (
	"fmt"
	"steve/structs/net"
	"steve/structs/proto/base"
	"sync"

	"github.com/Sirupsen/logrus"
)

type defaultIDAllocator struct {
	maxClientID uint64
	mutex       sync.Mutex
}

func (alloc *defaultIDAllocator) NewClientID() uint64 {
	alloc.mutex.Lock()
	alloc.maxClientID++
	ret := alloc.maxClientID
	alloc.mutex.Unlock()
	return ret
}

type watchDogImpl struct {
	alloc        net.IDAllocator
	msgObserver  net.MessageObserver
	connObserver net.ConnectObserver

	clientMap sync.Map
	callback  clientCallback

	serverMap      map[net.ServerType]server
	serverMapMutex sync.Mutex
}

type clientCallbackImpl struct {
	dog      *watchDogImpl
	clientID uint64
}

func (cc *clientCallbackImpl) onRecvPkg(header *steve_proto_base.Header, body []byte) {
	if cc.dog.msgObserver != nil {
		cc.dog.msgObserver.OnRecv(cc.clientID, header, body)
	}
}

func (cc *clientCallbackImpl) onError(err error) {
	logrus.WithField("client_id", cc.clientID).WithError(err).Debug("client error")
}

func (cc *clientCallbackImpl) onClientClose() {
	if cc.dog.connObserver != nil {
		cc.dog.connObserver.OnClientDisconnect(cc.clientID)
	}
}

func (dog *watchDogImpl) workOnExchanger(e exchanger) error {
	clientID := dog.alloc.NewClientID()
	client := newClientV2(e, &clientCallbackImpl{
		dog:      dog,
		clientID: clientID,
	})
	dog.clientMap.Store(clientID, client)

	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		// 由 run 函数自己处理异常
		err := client.run(func() {
			dog.clientMap.Delete(clientID)
		})
		logrus.WithField("client_id", clientID).WithError(err).Debug("client finished")
		wg.Done()
	}()

	logrus.WithField("client_id", clientID).Debug("client comming")

	if dog.connObserver != nil {
		dog.connObserver.OnClientConnect(clientID)
	}
	wg.Wait()
	return nil
}

func (dog *watchDogImpl) createServer(addr string, serverType net.ServerType) (server, error) {
	dog.serverMapMutex.Lock()
	defer dog.serverMapMutex.Unlock()

	if _, ok := dog.serverMap[serverType]; ok {
		return nil, fmt.Errorf("该类型的服务已经启动了")
	}

	server := newServer(addr, serverType)
	dog.serverMap[serverType] = server
	return server, nil
}

func (dog *watchDogImpl) Start(addr string, serverType net.ServerType) error {
	server, err := dog.createServer(addr, serverType)
	if err != nil {
		return err
	}

	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := server.Serve(addr, workerFunc(dog.workOnExchanger))
		if err != nil {
			panic(err)
		}
	}()
	wg.Wait()
	return nil
}

func (dog *watchDogImpl) Stop(serverType net.ServerType) error {
	dog.serverMapMutex.Lock()
	defer dog.serverMapMutex.Unlock()

	server, ok := dog.serverMap[serverType]
	if !ok {
		return fmt.Errorf("服务不存在")
	}
	delete(dog.serverMap, serverType)
	server.Close()
	return nil
}

func (dog *watchDogImpl) SendPackage(clientID uint64, header *steve_proto_base.Header, body []byte) error {
	return dog.pushClientMessage(clientID, header, body)
}

func (dog *watchDogImpl) pushClientMessage(clientID uint64, header *steve_proto_base.Header, body []byte) error {
	var c *clientV2
	// TODO : 此处有一个 BUG， 拿出来后在发送消息前客户端可能已被关闭
	if tmp, ok := dog.clientMap.Load(clientID); ok {
		c = tmp.(*clientV2)
	} else {
		return fmt.Errorf("clientID %v not exists", clientID)
	}
	return c.pushMessage(header, body)
}

func (dog *watchDogImpl) BroadPackage(clientIDs []uint64, header *steve_proto_base.Header, body []byte) error {
	for _, clientID := range clientIDs {
		dog.pushClientMessage(clientID, header, body)
	}
	return nil
}

func (dog *watchDogImpl) Disconnect(clientID uint64) error {
	var c *clientV2
	if tmp, ok := dog.clientMap.Load(clientID); ok {
		c = tmp.(*clientV2)
	} else {
		return fmt.Errorf("clientID %v not exists", clientID)
	}
	dog.clientMap.Delete(clientID)
	c.close()
	return nil
}

type factory struct{}

func (f *factory) NewWatchDog(alloc net.IDAllocator, msgObserver net.MessageObserver, connObserver net.ConnectObserver) net.WatchDog {
	if alloc == nil {
		alloc = &defaultIDAllocator{}
	}
	return &watchDogImpl{
		alloc:        alloc,
		msgObserver:  msgObserver,
		connObserver: connObserver,
		serverMap:    make(map[net.ServerType]server),
	}
}

// NewFactory 创建 WatchDogFactory
func NewFactory() net.WatchDogFactory {
	return new(factory)
}
