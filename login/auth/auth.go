package auth

import (
	"fmt"
	"steve/common/data/player"

	"github.com/Sirupsen/logrus"
)

// HandleLoginRequest 处理登录请求
func HandleLoginRequest(accountID uint64) uint64 {
	playerID := player.GetAccountPlayerID(accountID)
	if playerID == 0 {
		playerID = newPlayer(accountID)
	}
	return playerID
}

// newPlayer 创建玩家
func newPlayer(accountID uint64) uint64 {
	entry := logrus.WithFields(logrus.Fields{
		"func_name":  "newPlayer",
		"account_id": accountID,
	})
	playerID, err := player.AllocPlayerID()
	if err != nil {
		entry.WithError(err).Errorln("分配玩家 ID 失败")
		return 0
	}
	if err := player.NewPlayer(accountID, playerID); err != nil {
		entry.WithError(err).Errorln("创建玩家失败")
		return 0
	}
	initPlayerData(playerID)
	return playerID
}

func initPlayerData(playerID uint64) {
	player.SetPlayerCoin(playerID, 10000)
	player.SetPlayerNickName(playerID, fmt.Sprintf("玩家%v", playerID))
}
