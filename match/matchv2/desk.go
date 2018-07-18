package matchv2

import (
	"fmt"
	"time"
)

// deskPlayer 牌桌玩家
type deskPlayer struct {
	playerID uint64
	robotLv  int  // 机器人等级，为 0 时表示非机器人
	seat     int  // 座号，续局时有效
	robotWin bool // 为机器人时，上局是否为赢家，续局时有效
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
}

func (d *desk) String() string {
	return fmt.Sprintf("game_id: %d player:%v desk_id:%d ", d.gameID, d.players, d.deskID)
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
func createContinueDesk(gameID int, deskID uint64, players []deskPlayer) *desk {
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
	}
}
