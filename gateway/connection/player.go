package connection

// import "sync"

// // PlayerMgr 玩家连接管理
// type PlayerMgr struct {
// 	playerConnects sync.Map // playerID -> connectionID
// }

// // GetPlayerConnectionID 获取玩家的连接
// func (pm *PlayerMgr) GetPlayerConnectionID(playerID uint64) uint64 {
// 	_connectID, ok := pm.playerConnects.Load(playerID)
// 	if !ok {
// 		return 0
// 	}
// 	return _connectID.(uint64)
// }

// // SetPlayerConnectionID 设置玩家的连接
// func (pm *PlayerMgr) setPlayerConnectionID(playerID uint64, connectionID uint64) {
// 	pm.playerConnects.Store(playerID, connectionID)
// }
