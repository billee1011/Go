package interfaces

// PlayerManager 玩家管理器
type PlayerManager interface {
	GetPlayerConnectionID(playerID uint64) uint64
	SetPlayerConnectionID(playerID uint64, connectionID uint64)
}
