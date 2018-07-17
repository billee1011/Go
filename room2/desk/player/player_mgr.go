package player

import (
	"sync"
)

type PlayerMgr struct {
	playerMap sync.Map
}

var roomPlayerMgr PlayerMgr

func init(){
	roomPlayerMgr = PlayerMgr{}
}

func GetRoomPlayerMgr() PlayerMgr {
	return roomPlayerMgr
}

func (pm *PlayerMgr) GetPlayer(playerID uint64) *Player {
	result,ok  := pm.playerMap.Load(playerID)
	if !ok {
		return nil
	}
	player := result.(Player)
	return &player
}

//TODO 第一次进入房间服初始化
func (pm *PlayerMgr) InitPlayer(playerID uint64) {
	player := Player{
		PlayerID:playerID,
	}
	pm.playerMap.Store(playerID,player)
}

//TODO 离开房间服删除
func (pm *PlayerMgr) RemovePlayer(playerID uint64) {
	pm.playerMap.Delete(playerID)
}
