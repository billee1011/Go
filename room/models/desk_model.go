package models

import (
	"steve/room/desk"
	"steve/structs/proto/gate_rpc"
)

// DeskModel 牌桌模型
type DeskModel interface {
	GetName() string
	Active() // 初始化
	Start()
	Stop()
	GetDesk() *desk.Desk
	SetDesk(desk *desk.Desk)
}

// DeskEventModel 事件模型
type DeskEventModel interface {
	DeskModel
	PushRequest(playerID uint64, head *steve_proto_gaterpc.Header, bodyData []byte)
	// 开始处理事件
	StartProcessEvents()
}
