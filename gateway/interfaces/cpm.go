package interfaces

// ConnectPlayerMap 连接到玩家的映射
type ConnectPlayerMap interface {
	// GetConnectPlayer 根据连接 ID 获取玩家 ID
	GetConnectPlayer(clientID uint64) uint64

	// GetPlayerConnect 根据玩家 ID 获取连接 ID
	GetPlayerConnect(playerID uint64) uint64

	// SaveConnectPlayer 保存连接 ID 和玩家 ID 的映射
	SaveConnectPlayer(clientID uint64, playerID uint64)

	// RemoveConnectPlayer 移除连接
	RemoveConnect(clientID uint64)
}
