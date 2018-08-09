package player

import (
	"fmt"
	"steve/external/hallclient"
	user_pb "steve/server_pb/user"
	"sync"

	"github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
)

type PlayerMgr struct {
	playerMap sync.Map
}

var roomPlayerMgr *PlayerMgr

func init() {
	roomPlayerMgr = &PlayerMgr{}
}

func GetPlayerMgr() *PlayerMgr {
	return roomPlayerMgr
}

func (pm *PlayerMgr) GetPlayer(playerID uint64) *Player {
	result, ok := pm.playerMap.Load(playerID)
	if !ok {
		return nil
	}
	player := result.(*Player)
	return player
}

// InitDeskData init desk data
func (pm *PlayerMgr) InitDeskData(players []uint64, maxOverTime int, robotLv []int) {
	for seat, playerID := range players {
		player := pm.GetPlayer(playerID)
		if player == nil {
			pm.InitPlayer(playerID)
			player = pm.GetPlayer(playerID)
		}
		player.SetSeat(uint32(seat))
		player.SetEcoin(uint64(player.GetCoin()))
		player.SetMaxOverTime(maxOverTime)
		player.SetRobotLv(robotLv[seat])
		player.SetQuit(false)
	}
}

// 解除玩家的 room 服绑定
func (pm *PlayerMgr) UnbindPlayerRoomAddr(players []uint64) {
	entry := logrus.WithFields(logrus.Fields{
		"func_name": "deskFactory.unbindPlayerRoomAddr",
		"players":   players,
	})
	entry.Errorln("解除玩家的 room 服绑定-----------------------------------------")
	for _, playerID := range players {
		result, err := hallclient.UpdatePlayerState(playerID, user_pb.PlayerState_PS_GAMEING, user_pb.PlayerState_PS_IDIE, 0, 0)
		if !result && err != nil {
			entry.WithError(err).Errorln("设置玩家游戏状态失败")
		}
	}
}

// 绑定玩家所在 room 服
func (pm *PlayerMgr) BindPlayerRoomAddr(players []uint64, gameID int, levelID int) {
	entry := logrus.WithFields(logrus.Fields{
		"func_name": "deskFactory.bindPlayerRoomAddr",
		"players":   players,
		"gameID":    gameID,
		"levelID":   levelID,
	})
	roomIP := viper.GetString("rpc_addr")
	roomPort := viper.GetInt("rpc_port")
	roomAddr := fmt.Sprintf("%s:%d", roomIP, roomPort)
	for _, playerID := range players {

		// 更新玩家状态(从匹配状态改为游戏状态)
		result, err := hallclient.UpdatePlayerState(playerID, user_pb.PlayerState_PS_MATCHING, user_pb.PlayerState_PS_GAMEING, uint32(gameID), uint32(levelID))
		if !result || err != nil {
			entry.WithError(err).Errorln("设置玩家游戏状态失败")
		}

		// 更新玩家所在的服务器类型和地址
		result, err = hallclient.UpdatePlayeServerAddr(playerID, user_pb.ServerType_ST_ROOM, roomAddr)
		if !result || err != nil {
			entry.WithError(err).Errorln("设置玩家room服地址失败")
		}
	}
}

//TODO 第一次进入房间服初始化
func (pm *PlayerMgr) InitPlayer(playerID uint64) {
	player := &Player{
		PlayerID: playerID,
	}
	pm.playerMap.Store(playerID, player)
}

//TODO 离开房间服删除
func (pm *PlayerMgr) RemovePlayer(playerID uint64) {
	pm.playerMap.Delete(playerID)
}

func (pm *PlayerMgr) PlayerOverTime(player *Player) {
	player.OnPlayerOverTime()
}
