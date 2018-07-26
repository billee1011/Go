package matchv2

import (
	"fmt"
	"time"
)

// deskPlayer 牌桌玩家
type deskPlayer struct {
	playerID uint64
	robotLv  int  // 机器人等级，为 0 时表示非机器人
	seat     int  // 座号
	winner   bool // 上局是否为赢家，续局时有效
}

func (dp *deskPlayer) String() string {
	return fmt.Sprintf("player_id: %d robot_level:%d", dp.playerID, dp.robotLv)
}

// desk 匹配中的牌桌
type desk struct {
	gameID              int
	players             []deskPlayer
	deskID              uint64
	createTime          time.Time
	isContinue          bool                  // 是否为续局牌桌
	continueWaitPlayers map[uint64]deskPlayer // 续局等待玩家
	fixBanker           bool                  // 是否固定庄家位置
	bankerSeat          int                   // 庄家位置
}

func (d *desk) String() string {
	return fmt.Sprintf("game_id: %d player:%v desk_id:%d continue:%v fixBanker:%v bankerSeat:%v", d.gameID, d.players, d.deskID, d.isContinue, d.fixBanker, d.bankerSeat)
}

// removePlayer 移除玩家
func (d *desk) removePlayer(playerID uint64) {
	newPlayers := make([]deskPlayer, 0, 4)
	for _, player := range d.players {
		if playerID == player.playerID {
			continue
		}
		newPlayers = append(newPlayers, player)
	}
	d.players = newPlayers
}

// createDesk 创建牌桌
func createDesk(gameID int, deskID uint64) *desk {
	// logrus.WithFields(logrus.Fields{
	// 	"func_name": "createDesk",
	// 	"game_id":   gameID,
	// 	"desk_id":   deskID,
	// }).Debugln("创建牌桌")
	return &desk{
		gameID:     gameID,
		players:    make([]deskPlayer, 0, 4),
		deskID:     deskID,
		createTime: time.Now(),
	}
}

// createContinueDesk 创建续局牌桌
func createContinueDesk(gameID int, deskID uint64, players []deskPlayer, fixBanker bool, bankerSeat int) *desk {
	waitPlayers := make(map[uint64]deskPlayer, len(players))
	for _, player := range players {
		waitPlayers[player.playerID] = player
	}
	return &desk{
		gameID:              gameID,
		players:             make([]deskPlayer, 0, len(players)),
		deskID:              deskID,
		createTime:          time.Now(),
		isContinue:          true,
		continueWaitPlayers: waitPlayers,
		fixBanker:           fixBanker,
		bankerSeat:          bankerSeat,
	}
}
