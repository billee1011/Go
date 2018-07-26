package player

import (
	"sync"
	"errors"
	"github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
	"fmt"
	"steve/common/data/player"
	"steve/client_pb/common"
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
	player := result.(*Player)
	return player
}

func (pm *PlayerMgr) InitDeskData(players []uint64,maxOverTime int,robotLv []int){
	for seat,playerId := range players{
		player := pm.GetPlayer(playerId)
		if player == nil{
			pm.InitPlayer(playerId)
			player = pm.GetPlayer(playerId)
		}
		player.SetSeat(uint32(seat))
		player.SetEcoin(int(player.GetCoin()))
		player.SetMaxOverTime(maxOverTime)
		player.SetRobotLv(robotLv[seat])
	}
}

// 解除玩家的 room 服绑定
func (pm *PlayerMgr) UnbindPlayerRoomAddr(players []uint64) {
	entry := logrus.WithFields(logrus.Fields{
		"func_name": "deskFactory.unbindPlayerRoomAddr",
		"players":   players,
	})
	for _, playerID := range players {
		if err := player.SetPlayerPlayState(playerID, int(common.PlayerState_PS_IDLE)); err != nil {
			entry.WithError(err).Errorln("设置玩家游戏状态失败")
		}
	}
}

// 绑定玩家所在 room 服
func (pm *PlayerMgr) BindPlayerRoomAddr(players []uint64, gameID int) {
	entry := logrus.WithFields(logrus.Fields{
		"func_name": "deskFactory.bindPlayerRoomAddr",
		"players":   players,
	})
	roomIP := viper.GetString("rpc_addr")
	roomPort := viper.GetInt("rpc_port")
	roomAddr := fmt.Sprintf("%s:%d", roomIP, roomPort)
	for _, playerID := range players {
		if err := player.SetPlayerPlayStates(playerID, player.PlayStates{
			GameID:   gameID,
			State:    int(common.PlayerState_PS_GAMEING),
			RoomAddr: roomAddr,
		}); err != nil {
			entry.WithError(err).Errorln("设置玩家游戏状态失败")
		}
	}
}

//TODO 第一次进入房间服初始化
func (pm *PlayerMgr) InitPlayer(playerID uint64) {
	player := &Player{
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