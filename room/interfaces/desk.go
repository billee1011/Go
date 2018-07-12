package interfaces

import (
	msgid "steve/client_pb/msgid"
	"steve/structs/proto/gate_rpc"
)

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
	// QuitDesk 退出房间
	QuitDesk()
	// EnterDesk 进入房间
	EnterDesk()
	// OnPlayerOverTime 玩家超时
	OnPlayerOverTime()
	// IsTuoguan 玩家是否在托管中
	IsTuoguan() bool
	// SetTuoguan 设置托管
	SetTuoguan(tuoguan bool, notify bool)
	// 获取机器人等级
	GetRobotLv() int
}

// PlayerEnterQuitInfo 玩家退出进入信息
type PlayerEnterQuitInfo struct {
	PlayerID      uint64
	Quit          bool          // true 为退出， false 为进入
	FinishChannel chan struct{} // 完成通道
}

// DeskPlayerMgr 牌桌玩家管理器
type DeskPlayerMgr interface {

	// GetDeskPlayers 获取牌桌玩家
	GetDeskPlayers() []DeskPlayer

	// PlayerQuit 玩家退出
	PlayerQuit(playerID uint64) chan struct{}

	// PlayerEnter 玩家进入
	PlayerEnter(playerID uint64) chan struct{}

	// BroadcastMessage 广播消息给牌桌玩家
	// playerIDs ： 目标玩家，如果为 nil 或者长度为0，则针对牌桌所有玩家
	// exceptQuit ： 已经退出的玩家是否排除
	BroadcastMessage(playerIDs []uint64, msgID msgid.MsgID, body []byte, exceptQuit bool)

	// PlayerEnterQuitChannel 获取玩家进入退出通道
	PlayerEnterQuitChannel() <-chan PlayerEnterQuitInfo
}

// Desk 牌桌
type Desk interface {
	DeskPlayerMgr

	// GetUID 获取牌桌 UID
	GetUID() uint64

	// GetGameID 获取游戏 ID
	GetGameID() int

	// Start 启动牌桌逻辑
	// finish : 当牌桌逻辑完成时调用
	Start(finish func()) error

	// Stop 停止牌桌
	Stop() error

	// PushRequest 压入玩家请求
	PushRequest(playerID uint64, head *steve_proto_gaterpc.Header, bodyData []byte)

	// PushEvent 压入事件
	PushEvent(event Event)

	ChangePlayer(playerID uint64) error
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
	// GetDeskCount 获取牌桌数量
	GetDeskCount() int
}

// CreateDeskOptions 创建牌桌选项
type CreateDeskOptions struct {
}

// CreateDeskResult 创建房间结果
type CreateDeskResult struct {
	Desk Desk
}

// DeskFactory 牌桌工厂
type DeskFactory interface {
	// CreateDesk 创建牌桌
	CreateDesk(deskPlayers []DeskPlayer, gameID int, opt CreateDeskOptions) (CreateDeskResult, error)
}

// DeskIDAllocator 牌桌 ID 分配器
type DeskIDAllocator interface {
	AllocDeskID() (uint64, error)
}
