package player

import (
	"steve/gateway/global"
	"sync"
)

type playerMgr struct {
	playerConnects sync.Map // playerID -> connectionID
}

func (pm *playerMgr) GetPlayerConnectionID(playerID uint64) uint64 {
	_connectID, ok := pm.playerConnects.Load(playerID)
	if !ok {
		return 0
	}
	return _connectID.(uint64)
}

func (pm *playerMgr) SetPlayerConnectionID(playerID uint64, connectionID uint64) {
	pm.playerConnects.Store(playerID, connectionID)
}

func init() {
	global.SetPlayerManager(&playerMgr{})
}
