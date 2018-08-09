package matchv2

import (
	"math/rand"
	"steve/external/hallclient"
	"time"

	"github.com/Sirupsen/logrus"
)

// GetIdleRobot 取一个空闲的机器人
func GetIdleRobot(level int) uint64 {
	entry := logrus.WithFields(logrus.Fields{
		"func_name": "GetIdleRobot",
		"level":     level,
	})
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	// 先直接分配一个
	playerID, err := hallclient.GetPlayerByAccount(uint64(r.Intn(1000000)))
	if err == nil {
		entry.Debugf("GetIdleRobot playerId:%d", playerID)
		return playerID
	}
	return 0
}
