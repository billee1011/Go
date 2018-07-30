package matchv3

import (
	"fmt"
	"steve/common/data/player"

	"github.com/Sirupsen/logrus"
)

// GetIdleRobot 取一个空闲的机器人
// level 	:	机器人级别
func GetIdleRobot(level int) uint64 {
	entry := logrus.WithFields(logrus.Fields{
		"func_name": "GetIdleRobot",
		"level":     level,
	})

	// 先直接分配一个playerID
	playerID, err := player.AllocPlayerID()
	if err != nil {
		entry.WithError(err).Errorln("分配机器人的 playerID 失败")
		return 0
	}

	// redis中设置该playerID的金币数
	player.SetPlayerCoin(playerID, 10000)

	// redis中设置该playerID的昵称
	player.SetPlayerNickName(playerID, fmt.Sprintf("Robot%v", playerID))

	return playerID
}
