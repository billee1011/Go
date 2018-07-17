package matchv2

import (
	"steve/client_pb/room"
	"steve/common/data/player"
	"time"

	"github.com/Sirupsen/logrus"
)

// applyPlayer 申请匹配的玩家
type applyPlayer struct {
	playerID    uint64 // 玩家 ID
	gameID      int    // 游戏 ID
	isContinual bool   // 是否为续局匹配
}

// gameConfig 游戏数据
type gameConfig struct {
	needPlayerCount int // 所需玩家数量
}

// mgr 匹配管理
type mgr struct {
	applyChannel chan applyPlayer
	maxDeskID    uint64 // 最大牌桌 ID

	gameConfig map[int]gameConfig // gameID->gameConfig
	desks      map[uint64]*desk   // 当前匹配中的牌桌
	playerDesk map[uint64]uint64  // playerID -> deskID
}

// defaultMgr 默认匹配管理
var defaultMgr = &mgr{
	applyChannel: make(chan applyPlayer, 128),
	maxDeskID:    0,
	gameConfig: map[int]gameConfig{
		int(room.GameId_GAMEID_XUELIU):   gameConfig{needPlayerCount: 4},
		int(room.GameId_GAMEID_XUEZHAN):  gameConfig{needPlayerCount: 4},
		int(room.GameId_GAMEID_DOUDIZHU): gameConfig{needPlayerCount: 3},
		int(room.GameId_GAMEID_ERRENMJ):  gameConfig{needPlayerCount: 2},
	},
	desks:      make(map[uint64]*desk, 128),
	playerDesk: make(map[uint64]uint64, 1024),
}

func init() {
	go defaultMgr.run()
}

// addPlayer 添加匹配玩家
func (m *mgr) addPlayer(playerID uint64, gameID int, isContinual bool) {
	// logrus.WithFields(logrus.Fields{
	// 	"func_name":   "mgr.addPlayer",
	// 	"player_id":   playerID,
	// 	"is_continue": isContinual,
	// }).Debugln("添加匹配玩家")

	m.applyChannel <- applyPlayer{
		playerID:    playerID,
		gameID:      gameID,
		isContinual: isContinual,
	}
	return
}

// run 执行匹配流程
func (m *mgr) run() {
	robotTick := time.NewTicker(time.Second * 1)
	for {
		select {
		case ap := <-m.applyChannel:
			{
				m.acceptApplyPlayer(&ap)
			}
		case <-robotTick.C:
			{
				m.handleRobotTick()
			}
		}
	}
}

// acceptApplyPlayer 接收申请匹配玩家
func (m *mgr) acceptApplyPlayer(ap *applyPlayer) {
	logrus.WithFields(logrus.Fields{
		"func_name":   "mgr.acceptApplyPlayer",
		"player_id":   ap.playerID,
		"game_id":     ap.gameID,
		"is_continue": ap.isContinual,
	}).Debugln("接收申请匹配玩家")

	_, ok := m.playerDesk[ap.playerID]

	// 匹配中
	if ok {
		return
	}

	// 加入到牌桌
	for _, desk := range m.desks {
		if desk.gameID != ap.gameID {
			continue
		}
		m.add2Desk(ap.gameID, ap.playerID, 0, desk)
		return
	}
	m.maxDeskID++
	desk := createDesk(ap.gameID, m.maxDeskID)
	m.desks[desk.deskID] = desk
	m.add2Desk(ap.gameID, ap.playerID, 0, desk)
}

// add2Desk 添加到牌桌
func (m *mgr) add2Desk(gameID int, playerID uint64, robotLv int, desk *desk) {
	// logrus.WithFields(logrus.Fields{
	// 	"func_name": "mgr.add2Desk",
	// 	"player_id": playerID,
	// 	"game_id":   gameID,
	// 	"desk":      fmt.Sprintf("%#v", desk),
	// }).Debugln("添加到牌桌")
	desk.players = append(desk.players, deskPlayer{
		playerID: playerID,
		robotLv:  robotLv,
	})
	m.playerDesk[playerID] = desk.deskID
	m.removeOfflines(desk)
	config := m.gameConfig[gameID]
	if len(desk.players) >= config.needPlayerCount {
		m.onDeskFinish(desk)
	}
}

// fillRobots 填充机器人
func (m *mgr) fillRobots(desk *desk) {
	config := m.gameConfig[desk.gameID]
	logrus.WithFields(logrus.Fields{
		"func_name":  "mgr.fillRobots",
		"desk":       desk,
		"need_count": config.needPlayerCount,
	}).Debugln("加入机器人")
	curPlayerCount := len(desk.players)
	for i := curPlayerCount; i < config.needPlayerCount; i++ {
		m.add2Desk(desk.gameID, GetIdleRobot(1), 1, desk)
	}
}

// removeOfflines 移除 desk 中的离线玩家
func (m *mgr) removeOfflines(desk *desk) {
	// entry := logrus.WithFields(logrus.Fields{
	// 	"desk": fmt.Sprintf("%#v", desk),
	// })
	newPlayers := make([]deskPlayer, 0, len(desk.players))
	for _, deskPlayer := range desk.players {
		// 机器人不移除
		if deskPlayer.robotLv != 0 {
			newPlayers = append(newPlayers, deskPlayer)
			continue
		}
		online := (player.GetPlayerGateAddr(deskPlayer.playerID) != "")
		if online {
			newPlayers = append(newPlayers, deskPlayer)
		} else {
			delete(m.playerDesk, deskPlayer.playerID)
		}
	}
	desk.players = newPlayers
	// entry.WithField("new_players", desk.players).Debugln("移除离线玩家")
}

// onDeskFinish 牌桌匹配完成
func (m *mgr) onDeskFinish(desk *desk) {
	requestCreateDesk(desk)
	players := desk.players
	// 解除关联
	for _, player := range players {
		delete(m.playerDesk, player.playerID)
	}
	// 移除 desk
	delete(m.desks, desk.deskID)
}

// handleRobotTick 处理机器人 tick
func (m *mgr) handleRobotTick() {
	// 避免遍历时删除
	deskIDs := make([]uint64, 0, len(m.desks))
	for deskID := range m.desks {
		deskIDs = append(deskIDs, deskID)
	}
	for _, deskID := range deskIDs {
		desk := m.desks[deskID]
		if time.Now().Sub(desk.createTime) >= time.Second*5 {
			m.fillRobots(desk)
		}
	}
}
