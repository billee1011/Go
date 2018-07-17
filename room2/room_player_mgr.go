package room2

import (
	"sync"
)

type RoomPlayerMgr struct {
	playerMap sync.Map
}

var roomPlayerMgr RoomPlayerMgr

func init(){
	roomPlayerMgr = RoomPlayerMgr{}
}

func GetRoomPlayerMgr() RoomPlayerMgr{
	return roomPlayerMgr
}

func (pm *RoomPlayerMgr) GetPlayer(playerID uint64) *RoomPlayer {
	result,ok  := pm.playerMap.Load(playerID)
	if !ok {
		return nil
	}
	player := result.(RoomPlayer)
	return &player
}

//TODO 第一次进入房间服初始化
func (pm *RoomPlayerMgr) InitPlayer(playerID uint64) {
	player := RoomPlayer{
		PlayerID:playerID,
	}
	pm.playerMap.Store(playerID,player)
}

//TODO 离开房间服删除
func (pm *RoomPlayerMgr) RemovePlayer(playerID uint64) {
	pm.playerMap.Delete(playerID)
}
