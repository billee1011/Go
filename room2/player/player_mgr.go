package player

import (
	"sync"
	"errors"
)

type PlayerMgr struct {
	playerMap sync.Map
}

var roomPlayerMgr *PlayerMgr

func init(){
	roomPlayerMgr = &PlayerMgr{}
}

func GetPlayerMgr() *PlayerMgr {
	return roomPlayerMgr
}

func (pm *PlayerMgr) SetPlayerGold(playerID uint64,gold uint64) error{
	player := pm.GetPlayer(playerID)
	if player == nil {
		return errors.New("player not find")
	}
	player.SetCoin(gold)
	return nil
}

func (pm *PlayerMgr) GetPlayer(playerID uint64) *Player {
	result,ok  := pm.playerMap.Load(playerID)
	if !ok {
		return nil
	}
	player := result.(Player)
	return &player
}

func (pm *PlayerMgr) InitDeskData(players []uint64,maxOverTime int,robotLv []int){
	for seat,playerId := range players{
		player := pm.GetPlayer(playerId)
		player.SetSeat(uint32(seat))
		player.SetEcoin(int(player.GetCoin()))
		player.SetMaxOverTime(maxOverTime)
		player.SetRobotLv(robotLv[seat])
	}
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


func (pm *PlayerMgr) PlayerOverTime(player *Player){
	player.OnPlayerOverTime()
}