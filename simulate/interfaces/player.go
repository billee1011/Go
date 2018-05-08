package interfaces

// ClientPlayer 客户端玩家
type ClientPlayer interface {
	GetID() uint64
	GetCoin() uint64
	GetClient() Client
}
