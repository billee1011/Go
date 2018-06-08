package interfaces

// Player 玩家
type Player interface {
	GetID() uint64
	GetCoin() uint64
	GetClientID() uint64
	SetClientID(clientID uint64)
	SetCoin(coin uint64)
	GetUserName() string
}

// PlayerMgr 玩家管理器
type PlayerMgr interface {
	AddPlayer(p Player)
	GetPlayer(ID uint64) Player
	GetPlayerByClientID(clientID uint64) Player
	OnClientDisconnect(clientID uint64)
	GetPlayerByUserName(userName string) Player
	UpdatePlayerClientID(playerID uint64, clientID uint64)
}
