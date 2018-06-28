package interfaces

// Player 玩家
type Player interface {
	GetID() uint64
	GetCoin() uint64
	// GetClientID() uint64
	// SetClientID(clientID uint64)
	SetCoin(coin uint64)
	// GetUserName() string

	// 判断玩家是否在线
	IsOnline() bool
}

// PlayerMgr 玩家管理器
type PlayerMgr interface {
	GetPlayer(ID uint64) Player
}
