package playermgr

import (
	"steve/room/interfaces"
	"steve/room/interfaces/global"
	"sync"
)

type playerMgr struct {
	playerMap sync.Map // playerID -> player
	clientMap sync.Map // clientID -> playerID

	mu sync.RWMutex
}

func (pm *playerMgr) AddPlayer(p interfaces.Player) {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	pm.playerMap.Store(p.GetID(), p)
	pm.clientMap.Store(p.GetClientID(), p.GetID())
}

func (pm *playerMgr) GetPlayer(playerID uint64) interfaces.Player {
	return pm.getPlayer(playerID)
}

func (pm *playerMgr) getPlayer(playerID uint64) interfaces.Player {
	v, ok := pm.playerMap.Load(playerID)
	if !ok {
		return nil
	}
	return v.(interfaces.Player)
}

func (pm *playerMgr) GetPlayerByClientID(clientID uint64) interfaces.Player {
	pm.mu.RLock()
	defer pm.mu.RLock()
	v, ok := pm.clientMap.Load(clientID)
	if !ok {
		return nil
	}
	playerID := v.(uint64)
	return pm.getPlayer(playerID)
}

var setupOnce = sync.Once{}

// SetupPlayerMgr 初始化玩家管理器
func SetupPlayerMgr() {
	setupOnce.Do(func() {
		pm := &playerMgr{}
		global.SetPlayerMgr(pm)
	})
}
