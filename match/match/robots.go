package match

import (
	"fmt"
	"steve/common/data/player"

	"github.com/Sirupsen/logrus"
)

// GetIdleRobot 取一个空闲的机器人
func GetIdleRobot(level int) uint64 {
	entry := logrus.WithFields(logrus.Fields{
		"func_name": "GetIdleRobot",
		"level":     level,
	})
	// 先直接分配一个
	playerID, err := player.AllocPlayerID()
	if err != nil {
		entry.WithError(err).Errorln("分配玩家 ID 失败")
		return 0
	}
	player.SetPlayerCoin(playerID, 10000)
	player.SetPlayerNickName(playerID, fmt.Sprintf("Robot%v", playerID))
	return playerID
}
