package net

import (
	"steve/structs/proto/base"
)

// ServerType 服务类型
type ServerType int

const (
	// RPC 使用 gRPC 作为基础通信框架
	RPC ServerType = iota
	// TCP 使用 TCP 协议作为基础通信框架
	TCP
)

// MessageObserver 消息观察者
type MessageObserver interface {
	// OnRecv 收到消息回调
	OnRecv(clientID uint64, header *base.Header, body []byte)
	// AfterSend 消息发送完成之后的回调
	// err 表示消息发送错误信息
	AfterSend(clientID uint64, header *base.Header, body []byte, err error)
}

// ConnectObserver 连接观察者
// 当客户端连接或者断开连接时会触发回调
type ConnectObserver interface {
	OnClientConnect(clientID uint64)
	OnClientDisconnect(clientID uint64)
}

// IDAllocator ID 分配器，用来分配客户端连接 ID
type IDAllocator interface {
	NewClientID() uint64
}

// WatchDog 看门狗，用来管理基础网络连接
type WatchDog interface {
	Start(addr string, serverType ServerType) error
	Stop(serverType ServerType) error
	SendPackage(clientID uint64, header *base.Header, body []byte) error
	BroadPackage(clientIDs []uint64, header *base.Header, body []byte) error
	Disconnect(clientID uint64) error
}

// WatchDogFactory 用来创建 WatchDog
type WatchDogFactory interface {

	// 创建 WatchDog 对象。alloc 为客户端 ID 分配器， 为空时，将使用默认的 ID 分配器
	// msgObserver 用来观察客户端消息事件
	// connObserver 用来观察客户端连接事件
	NewWatchDog(alloc IDAllocator, msgObserver MessageObserver, connObserver ConnectObserver) WatchDog
}
