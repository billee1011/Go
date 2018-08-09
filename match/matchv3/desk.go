package matchv3

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
)

// IPUInt32ToString 整形IP地址转为字符串型IP
func IPUInt32ToString(intIP uint32) string {
	var bytes [4]byte
	bytes[0] = byte(intIP & 0xFF)
	bytes[1] = byte((intIP >> 8) & 0xFF)
	bytes[2] = byte((intIP >> 16) & 0xFF)
	bytes[3] = byte((intIP >> 24) & 0xFF)

	return net.IPv4(bytes[3], bytes[2], bytes[1], bytes[0]).String()
}

// IPStringToUInt32 字符串型IP转为uint32型
func IPStringToUInt32(ipStr string) uint32 {
	bits := strings.Split(ipStr, ".")

	if len(bits) != 4 {
		logrus.Errorln("IPStringToUInt32() 参数错误，ipStr = ", ipStr)
		return 0
	}

	b0, _ := strconv.Atoi(bits[0])
	b1, _ := strconv.Atoi(bits[1])
	b2, _ := strconv.Atoi(bits[2])
	b3, _ := strconv.Atoi(bits[3])

	var sum uint32

	sum += uint32(b0) << 24
	sum += uint32(b1) << 16
	sum += uint32(b2) << 8
	sum += uint32(b3)

	return sum
}

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
