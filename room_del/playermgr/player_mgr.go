package playermgr

import (
	"steve/room/interfaces"
	"steve/room/interfaces/global"
	"sync"
)

type playerMgr struct {
	playerMap   sync.Map // playerID -> player
	clientMap   sync.Map // clientID -> playerID
	userNameMap sync.Map // userName-> playerID

	mu sync.RWMutex
}

func (pm *playerMgr) GetPlayer(playerID uint64) interfaces.Player {
	return &player{playerID: playerID}
}

func init() {
	global.SetPlayerMgr(&playerMgr{})
}
