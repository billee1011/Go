package matchv2

import (
	"steve/client_pb/room"
	"steve/common/data/player"
	"steve/gutils"
	"steve/match/web"
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
	applyChannel chan applyPlayer   // 申请通道
	maxDeskID    uint64             // 最大牌桌 ID
	gameConfig   map[int]gameConfig // gameID -> gameConfig
	desks        map[uint64]*desk   // 当前匹配中的牌桌
	playerDesk   map[uint64]uint64  // 匹配中的玩家， playerID -> deskID
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

// addContinueDesk 添加续局牌桌
func (m *mgr) addContinueDesk(players []deskPlayer, gameID int, fixBanker bool, bankerSeat int) {
	m.maxDeskID++
	// 有玩家在匹配中，不创建
	for _, player := range players {
		if _, ok := m.playerDesk[player.playerID]; ok {
			logrus.WithField("player_id", player.playerID).Infoln("添加续局牌桌时玩家已经在匹配中了")
			return
		}
	}
	deskID := m.maxDeskID
	desk := createContinueDesk(gameID, deskID, players, fixBanker, bankerSeat)
	for _, player := range players {
		m.playerDesk[player.playerID] = deskID
	}
	m.desks[deskID] = desk
}

// dismissContinueDesk 解散续局牌桌
func (m *mgr) dismissContinueDesk(desk *desk) {
	logrus.WithFields(logrus.Fields{
		"func_name":    "mgr.dismissContinueDesk",
		"ready_player": desk.players,
	}).Debugln("解散续局牌桌")
	readyPlayers := desk.players
	for _, player := range desk.players {
		delete(m.playerDesk, player.playerID)
	}
	for playerID := range desk.continueWaitPlayers {
		delete(m.playerDesk, playerID)
	}
	delete(m.desks, desk.deskID)
	// 已准备的玩家重新入队
	for _, player := range readyPlayers {
		if player.robotLv != 0 { // 机器人不入队
			continue
		}
		m.acceptApplyPlayer(desk.gameID, player.playerID, false)
	}
}

// addPlayer 添加匹配玩家
func (m *mgr) addPlayer(playerID uint64, gameID int, isContinual bool) {
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
	continueTick := time.NewTicker(time.Second * 1)
	for {
		select {
		case ap := <-m.applyChannel:
			{
				m.acceptApplyPlayer(ap.gameID, ap.playerID, ap.isContinual)
			}
		case <-robotTick.C:
			{
				m.handleRobotTick()
			}
		case <-continueTick.C:
			{
				m.checkContinueDesks()
			}
		}
	}
}

// acceptContinuePlayer 接收续局匹配玩家， 返回是否接收成功
func (m *mgr) acceptContinuePlayer(gameID int, playerID uint64, robotLv int) bool {
	entry := logrus.WithFields(logrus.Fields{
		"func_name":   "mgr.acceptContinuePlayer",
		"player_id":   playerID,
		"robot_level": robotLv,
	})

	deskID, ok := m.playerDesk[playerID]
	if !ok {
		return false
	}
	entry = entry.WithField("desk_id", deskID)
	desk, ok := m.desks[deskID]
	if !ok {
		delete(m.playerDesk, playerID)
		entry.Errorln("牌桌不存在")
		return false
	}
	// 非续局牌桌
	if !desk.isContinue {
		return true
	}
	player, ok := desk.continueWaitPlayers[playerID]
	if !ok {
		return true
	}
	entry.Debugln("接收续局玩家")
	delete(desk.continueWaitPlayers, playerID)
	m.addDeskPlayer2Desk(&player, desk)
	return true
}

// acceptApplyPlayer 接收申请匹配玩家
func (m *mgr) acceptApplyPlayer(gameID int, playerID uint64, isContinue bool) {
	// 续局匹配
	if isContinue {
		if m.acceptContinuePlayer(gameID, playerID, 0) {
			return
		}
	}
	deskID, ok := m.playerDesk[playerID]
	logrus.WithFields(logrus.Fields{
		"func_name":   "mgr.acceptApplyPlayer",
		"player_id":   playerID,
		"game_id":     gameID,
		"is_continue": isContinue,
		"old_desk_id": deskID,
	}).Debugln("接收申请匹配玩家")
	if ok {
		// 等待续局中
		if desk, exist := m.desks[deskID]; exist && desk.isContinue {
			m.dismissContinueDesk(desk)
		} else {
			return // 匹配中
		}
	}
	// 加入到牌桌
	for _, desk := range m.desks {
		if desk.gameID != gameID || desk.isContinue {
			continue
		}
		m.addDeskPlayer2Desk(&deskPlayer{
			playerID: playerID,
		}, desk)
		return
	}
	m.maxDeskID++
	desk := createDesk(gameID, m.maxDeskID)
	m.desks[desk.deskID] = desk
	m.addDeskPlayer2Desk(&deskPlayer{
		playerID: playerID,
	}, desk)
}

// addDeskPlayer2Desk 将玩家添加到牌桌
func (m *mgr) addDeskPlayer2Desk(player *deskPlayer, desk *desk) {
	desk.players = append(desk.players, *player)
	m.playerDesk[player.playerID] = desk.deskID
	m.removeOfflines(desk)
	config := m.gameConfig[desk.gameID]
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
		m.addDeskPlayer2Desk(&deskPlayer{
			playerID: GetIdleRobot(1),
			robotLv:  1,
		}, desk)
	}
}

// removeOfflines 移除 desk 中的离线玩家
func (m *mgr) removeOfflines(desk *desk) {
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
		if !desk.isContinue && time.Now().Sub(desk.createTime) >= web.GetRobotJoinTime() {
			m.fillRobots(desk)
		}
	}
}

// checkContinueDesks 检查续局牌桌，超过 20s 解散
func (m *mgr) checkContinueDesks() {
	// 避免遍历时删除
	deskIDs := make([]uint64, 0, len(m.desks))
	for deskID := range m.desks {
		deskIDs = append(deskIDs, deskID)
	}
	for _, deskID := range deskIDs {
		desk := m.desks[deskID]
		// 非续局牌桌
		if !desk.isContinue {
			continue
		}
		interval := time.Now().Sub(desk.createTime)
		// 超过解散时间
		if interval >= web.GetContinueDismissTime() {
			m.dismissContinueDesk(desk)
			continue
		}
		// 超过机器人续局时间
		if interval >= web.GetContinueRobotTime() {
			m.robotContinue(desk)
		}
	}
}

// robotContinue 机器人作续局决策
func (m *mgr) robotContinue(desk *desk) {
	entry := logrus.WithFields(logrus.Fields{
		"func_name": "mgr.robotContinue",
		"desk":      desk.String(),
	})
	robots := make([]uint64, 0, len(desk.continueWaitPlayers))

	for playerID := range desk.continueWaitPlayers {
		robots = append(robots, playerID)
	}

	for _, playerID := range robots {
		player := desk.continueWaitPlayers[playerID]
		if player.robotLv == 0 {
			continue
		}
		rate := web.GetRobotContinueRate(player.winner)
		if !gutils.Probability(rate) {
			entry.Debugln("机器人不续局，解散")
			m.dismissContinueDesk(desk)
			return
		}
		m.acceptContinuePlayer(desk.gameID, playerID, player.robotLv)
	}
}
