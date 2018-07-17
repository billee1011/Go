package matchv2

import (
	"fmt"
	"time"
)

// deskPlayer 牌桌玩家
type deskPlayer struct {
	playerID uint64
	robotLv  int // 机器人等级，为 0 时表示非机器人
}

func (dp *deskPlayer) String() string {
	return fmt.Sprintf("player_id: %d robot_level:%d", dp.playerID, dp.robotLv)
}

// desk 匹配中的牌桌
type desk struct {
	gameID     int
	players    []deskPlayer
	deskID     uint64
	createTime time.Time
}

func (d *desk) String() string {
	return fmt.Sprintf("game_id: %d player:%v desk_id:%d ", d.gameID, d.players, d.deskID)
}

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
