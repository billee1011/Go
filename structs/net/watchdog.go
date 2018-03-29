package net

import (
	"steve/structs/proto/base"

	"github.com/golang/protobuf/proto"
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
// 当收到客户端发来的消息时， OnRecv 函数会被调用
type MessageObserver interface {
	OnRecv(clientID uint64, header *steve_proto_base.Header, body []byte)
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
	SendPackage(clientID uint64, header *steve_proto_base.Header, bodyMsg proto.Message) error
	BroadPackage(clientIDs []uint64, header *steve_proto_base.Header, bodyMsg proto.Message) error
	Disconnect(clientID uint64) error
}
