package matchv3

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
)

// UInt32ToIP 整形IP地址转为字符串型IP
func UInt32ToIP(intIP uint32) net.IP {
	var bytes [4]byte
	bytes[0] = byte(intIP & 0xFF)
	bytes[1] = byte((intIP >> 8) & 0xFF)
	bytes[2] = byte((intIP >> 16) & 0xFF)
	bytes[3] = byte((intIP >> 24) & 0xFF)

	return net.IPv4(bytes[3], bytes[2], bytes[1], bytes[0])
}

// IPToUInt32 字符串型IP转为uint32型
func IPToUInt32(ipnr net.IP) uint32 {
	bits := strings.Split(ipnr.String(), ".")

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
	robotLv  int    // 机器人等级，为 0 时表示非机器人
	seat     int    // 座号
	IP       uint32 // IP地址
}

// matchPlayer转为字符串
func (pPlayer *matchPlayer) String() string {
	return fmt.Sprintf("player_id: %v, robot_level:%v, seat:%v, IP:%v", pPlayer.playerID, pPlayer.robotLv, pPlayer.seat, UInt32ToIP(pPlayer.IP))
}

// desk 匹配中的牌桌
type desk struct {
	gameID              int                   // 游戏ID
	deskID              uint64                // 桌子唯一ID
	players             []deskPlayer          // 桌子中的所有玩家
	createTime          time.Time             // 桌子创建时间
	isContinue          bool                  // 是否为续局牌桌，默认为false
	continueWaitPlayers map[uint64]deskPlayer // 续局牌桌等待的玩家，key:玩家ID,value:deskPlayer
	fixBanker           bool                  // 是否固定庄家位置
	bankerSeat          int                   // 庄家位置
	winRate             uint8                 // 创建时的胜率
}

// desk转为字符串
func (d *desk) String() string {
	return fmt.Sprintf("game_id: %d player:%v desk_id:%d continue:%v fixBanker:%v bankerSeat:%v",
		d.gameID, d.players, d.deskID, d.isContinue, d.fixBanker, d.bankerSeat)
}

// matchDesk 匹配中的牌桌
type matchDesk struct {
	gameID          int32         // 游戏ID
	levelID         int32         // 场次ID
	gold            uint64        // 桌子的金币（也是第一个玩家的金币数）
	needPlayerCount uint8         // 满桌需要的玩家数量
	players         []matchPlayer // 桌子中的所有玩家
	createTime      time.Time     // 桌子创建时间
}

// matchDesk转为字符串
func (pDesk *matchDesk) String() string {
	return fmt.Sprintf("gameID: %v, levelID: %v, gold: %v, needPlayerCount:%v, players:%v, createTime:%v",
		pDesk.gameID, pDesk.levelID, pDesk.gold, pDesk.needPlayerCount, pDesk.players, pDesk.createTime)
}

// createDesk 创建一个新牌桌
// gameID 	:	游戏ID
// deskID	:	桌子唯一ID
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

// createMatchDesk 创建一个新的匹配桌子
// gameID 			: 游戏ID
// levelID 			: 级别ID
// needPlayerCount 	: 满桌需要的玩家数量
// gold				: 金币(第一个玩家的金币数)
func createMatchDesk(gameID int32, levelID int32, needPlayerCount uint8, gold uint64) *matchDesk {
	logrus.WithFields(logrus.Fields{
		"func_name":       "createMatchDesk",
		"gameID":          gameID,
		"levelID":         levelID,
		"needPlayerCount": needPlayerCount,
		"gold":            gold,
	}).Debugln("创建匹配牌桌")

	return &matchDesk{
		gold:            gold,
		needPlayerCount: needPlayerCount,
		players:         make([]matchPlayer, 0, needPlayerCount),
		createTime:      time.Now(),
	}
}

// createContinueDesk 创建续局牌桌
// gameID		:	游戏ID
// deskID		:	桌子唯一ID
// players		:	等待的所有玩家
// fixBanker	:	是否固定庄家位置
// bankerSeat	:	庄家座位号
func createContinueDesk(gameID int, deskID uint64, players []deskPlayer, fixBanker bool, bankerSeat int) *desk {
	waitPlayers := make(map[uint64]deskPlayer, len(players))

	// 等待的玩家信息
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
