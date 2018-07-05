package interfaces

import (
	msgid "steve/client_pb/msgId"
	room "steve/client_pb/room"
	"steve/structs/proto/gate_rpc"
)

// TuoGuanMgr 牌桌托管管理器
type TuoGuanMgr interface {
	// GetTuoGuanPlayers 获取托管玩家
	GetTuoGuanPlayers() []uint64
	// OnPlayerTimeOut 玩家超时
	OnPlayerTimeOut(playerID uint64)
	// SetTuoGuan 设置玩家托管
	SetTuoGuan(playerID uint64, set bool, notify bool)

	// IsTuoGuan 是否托管中
	IsTuoGuan(playerID uint64) bool
}

// DeskPlayer 牌桌玩家
type DeskPlayer interface {
	// GetPlayerID 获取玩家 ID
	GetPlayerID() uint64
	// GetSeat 获取座号
	GetSeat() int
	// GetEcoin 获取进入时金币数
	GetEcoin() int
	// IsQuit 是否已经退出
	IsQuit() bool
}

// Desk 牌桌
type Desk interface {
	// GetUID 获取牌桌 UID
	GetUID() uint64

	// GetGameID 获取游戏 ID
	GetGameID() int

	// GetPlayers 获取牌桌玩家数据 (将会被废弃，不要使用， 改为 GetDeskPlayers 代替)
	GetPlayers() []*room.RoomPlayerInfo

	// GetDeskPlayers 获取牌桌玩家
	GetDeskPlayers() []DeskPlayer

	// Start 启动牌桌逻辑
	// finish : 当牌桌逻辑完成时调用
	Start(finish func()) error

	// Stop 停止牌桌
	Stop() error

	// PushRequest 压入玩家请求
	PushRequest(playerID uint64, head *steve_proto_gaterpc.Header, bodyData []byte)

	// PushEvent 压入事件
	PushEvent(event Event)

	// PlayerQuit 玩家退出
	PlayerQuit(playerID uint64)

	// PlayerEnter 玩家进入
	PlayerEnter(playerID uint64)

	// ChangePlayer 换对手
	ChangePlayer(playerID uint64) error

	// GetTuoGuanMgr 获取托管管理器
	GetTuoGuanMgr() TuoGuanMgr

	// BroadcastMessage 广播消息给牌桌玩家
	// playerIDs ： 目标玩家，如果为 nil 或者长度为0，则针对牌桌所有玩家
	// exceptQuit ： 已经退出的玩家是否排除
	BroadcastMessage(playerIDs []uint64, msgID msgid.MsgID, body []byte, exceptQuit bool)
}

// DeskMgr 牌桌管理器
type DeskMgr interface {
	// RunDesk 运转牌桌
	RunDesk(desk Desk) error

	// HandlePlayerRequest 处理玩家请求
	HandlePlayerRequest(playerID uint64, head *steve_proto_gaterpc.Header, bodyData []byte)

	// GetRunDeskByPlayerID 获取该玩家所在牌桌
	GetRunDeskByPlayerID(playerID uint64) (Desk, error)

	// RemoveDeskPlayerByPlayerID 移除某个在桌子上的玩家
	RemoveDeskPlayerByPlayerID(playerID uint64)
}

// CreateDeskOptions 创建牌桌选项
type CreateDeskOptions struct{}

// CreateDeskResult 创建房间结果
type CreateDeskResult struct {
	Desk Desk
}

// DeskFactory 牌桌工厂
type DeskFactory interface {
	// CreateDesk 创建牌桌
	CreateDesk(players []uint64, gameID int, opt CreateDeskOptions) (CreateDeskResult, error)
}

// DeskIDAllocator 牌桌 ID 分配器
type DeskIDAllocator interface {
	AllocDeskID() (uint64, error)
}
