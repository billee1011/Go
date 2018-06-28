package interfaces

import (
	"steve/structs/net"
)

// Connection 连接信息
type Connection interface {
	GetPlayerID() uint64
	// AttachPlayer 关联玩家，成功返回 true, 如果已经关联过玩家了返回 false
	// 如果以后再有其它错误，再修改此接口
	AttachPlayer(playerID uint64) bool
	GetClientID() uint64
	HeartBeat()
}

// ConnectionManager 连接管理器
type ConnectionManager interface {
	net.ConnectObserver
	GetConnection(clientID uint64) Connection
	SetKicker(kicker func(clientID uint64))
}
