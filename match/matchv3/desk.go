package matchv3

import (
	"fmt"
	"steve/external/hallclient"
	"steve/server_pb/user"
	"time"

	"github.com/Sirupsen/logrus"
)

// deskPlayer 牌桌玩家
type deskPlayer struct {
	playerID uint64 // 玩家ID
	robotLv  int    // 机器人等级，为 0 时表示非机器人
	seat     int    // 座号
	winner   bool   // 上局是否为赢家，续局时有效
}

// deskPlayer转为字符串
func (dp *deskPlayer) String() string {
	return fmt.Sprintf("player_id: %d robot_level:%d", dp.playerID, dp.robotLv)
}

// matchPlayer 匹配中的玩家
type matchPlayer struct {
	playerID uint64 // 玩家ID
	robotLv  int32  // 机器人等级，为 0 时表示非机器人
	seat     int32  // 座号
	IP       uint32 // IP地址
	gold     int64  // 金币数
}

// matchPlayer转为字符串
func (pPlayer *matchPlayer) String() string {
	return fmt.Sprintf("player_id: %v, robot_level:%v, seat:%v, IP:%v", pPlayer.playerID, pPlayer.robotLv, pPlayer.seat, IPUInt32ToString(pPlayer.IP))
}

// matchDesk 匹配中的牌桌
type matchDesk struct {
	deskID          uint64        // 桌子唯一ID
	gameID          uint32        // 游戏ID
	levelID         uint32        // 场次ID
	aveGold         int64         // 桌子的平均金币
	needPlayerCount uint8         // 满桌需要的玩家数量
	players         []matchPlayer // 桌子中的所有玩家
	createTime      int64         // 桌子创建时间(单位：秒)
}

// 已成功的牌桌，用于计算玩家上局是否同桌
type sucDesk struct {
	gameID  uint32 // 游戏ID
	levelID uint32 // 场次ID
	sucTime int64  // 成功时间
}

// matchDesk转为字符串
func (pDesk *matchDesk) String() string {
	return fmt.Sprintf("gameID: %v, levelID: %v, gold: %v, needPlayerCount:%v, players:%v, createTime:%v",
		pDesk.gameID, pDesk.levelID, pDesk.aveGold, pDesk.needPlayerCount, pDesk.players, pDesk.createTime)
}

// createMatchDesk 创建一个新的匹配桌子
// deskID			: 桌子ID
// gameID 			: 游戏ID
// levelID 			: 级别ID
// needPlayerCount 	: 满桌需要的玩家数量
// gold				: 金币(第一个玩家的金币数)
func createMatchDesk(deskID uint64, gameID uint32, levelID uint32, needPlayerCount uint8, gold int64) *matchDesk {
	logrus.WithFields(logrus.Fields{
		"func_name":       "createMatchDesk",
		"deskID":          deskID,
		"gameID":          gameID,
		"levelID":         levelID,
		"needPlayerCount": needPlayerCount,
		"gold":            gold,
	}).Debugln("创建匹配牌桌")

	return &matchDesk{
		deskID:          deskID,
		gameID:          gameID,
		levelID:         levelID,
		aveGold:         gold,
		needPlayerCount: needPlayerCount,
		players:         make([]matchPlayer, 0, needPlayerCount),
		createTime:      time.Now().Unix(),
	}
}

// dealErrorDesk 处理出现错误的桌子
// 把桌子内的所有玩家更改为空闲状态
func dealErrorDesk(pDesk *matchDesk) bool {
	if pDesk == nil {
		logrus.Errorln("DealDeskError() 参数错误，pDesk == nil")
		return false
	}

	logEntry := logrus.WithFields(logrus.Fields{
		"pDesk": pDesk,
	})

	// 处理桌子内所有玩家
	for i := 0; i < len(pDesk.players); i++ {
		if !dealErrorPlayer(&pDesk.players[i]) {
			logEntry.Errorf("处理错误桌子时，处理错误玩家失败，玩家ID:%v", pDesk.players[i].playerID)
		}
	}

	return true
}

// dealErrorPlayer 处理出现错误的玩家
func dealErrorPlayer(pPlayer *matchPlayer) bool {
	if pPlayer == nil {
		logrus.Errorln("dealErrorPlayer() 参数错误，pPlayer == nil")
		return false
	}

	logEntry := logrus.WithFields(logrus.Fields{
		"pPlayer": pPlayer,
	})

	// 设置为匹配状态，后面匹配过程中出错删除时再标记为空闲状态，匹配成功时不需处理(room服会标记为游戏状态)
	bSuc, err := hallclient.UpdatePlayerState(pPlayer.playerID, user.PlayerState_PS_MATCHING, user.PlayerState_PS_IDIE, 0, 0)
	if err != nil || !bSuc {
		logEntry.WithError(err).Errorf("处理错误玩家时，通知hall服设置玩家状态为空闲状态时失败，玩家ID:%v", pPlayer.playerID)
		return false
	}

	// 更新玩家所在match服务器的地址
	bSuc, err = hallclient.UpdatePlayeServerAddr(pPlayer.playerID, user.ServerType_ST_MATCH, "")
	if err != nil || !bSuc {
		logEntry.WithError(err).Errorf("处理错误玩家时，通知hall服设置玩家的match服务器地址时失败，玩家ID:%v", pPlayer.playerID)
		return false
	}

	return true
}
