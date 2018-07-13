package match

import (
	"steve/client_pb/room"
	"steve/common/data/player"
	"time"

	"github.com/Sirupsen/logrus"
)

const (
	playersOneDesk int = 4 // 一个桌子需要的人数
)

type matchData struct {
	gameID      int
	queue       chan uint64
	playerCount int
	desks       map[uint64]*Desk // 当前正在匹配中的牌桌
	maxDeskID   uint64
}

// Manager 匹配管理器
type Manager struct {
	matchDataMap map[int]matchData
}

var defaultManager = NewManager()

func createMatchData(gameID room.GameId, playerCount int) matchData {
	return matchData{
		gameID:      int(gameID),
		queue:       make(chan uint64),
		playerCount: playerCount,
		desks:       make(map[uint64]*Desk, 16),
	}
}

// NewManager 创建 Manager
func NewManager() *Manager {
	// TODO: 从配置加载
	queues := make(map[int]matchData)
	queues[int(room.GameId_GAMEID_XUELIU)] = createMatchData(room.GameId_GAMEID_XUELIU, 4)
	queues[int(room.GameId_GAMEID_XUEZHAN)] = createMatchData(room.GameId_GAMEID_XUEZHAN, 4)
	queues[int(room.GameId_GAMEID_DOUDIZHU)] = createMatchData(room.GameId_GAMEID_DOUDIZHU, 3)
	queues[int(room.GameId_GAMEID_ERRENMJ)] = createMatchData(room.GameId_GAMEID_ERRENMJ, 2)

	m := &Manager{
		matchDataMap: queues,
	}
	m.runMatchs()
	return m
}

func (m *Manager) runMatchs() {
	for gameID, md := range m.matchDataMap {
		go m.match(gameID, md.playerCount, md.queue)
	}
}

func (m *Manager) addPlayer(playerID uint64, gameID int) {
	md, ok := m.matchDataMap[gameID]
	if !ok {
		return
	}
	md.queue <- playerID
}

// 具体的匹配操作
func (m *Manager) match(gameID int, playerCount int, ch chan uint64) {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "Manager::match()",
	})

	logEntry.Debugln("Manager::Match() start ...")
	robotTicker := time.NewTicker(time.Second)
	for {
		select {
		case playerID := <-ch:
			online := player.GetPlayerGateAddr(playerID)
			if online == "" {
				logEntry.Debugln("player is not online, remove the player[%d]", playerID)
				continue
			}
			m.receivePlayer(gameID, playerID, 0)
		case <-robotTicker.C:
			{
				m.addRobots(gameID)
			}
		}

	}
}

func (m *Manager) onDeskMatchFinish(desk *Desk) {
	gameID := desk.GetGameID()
	md := m.matchDataMap[gameID]
	deskID := desk.GetID()

	delete(md.desks, deskID)

	players := desk.GetPlayers()
	sender := Sender{}
	sender.createDesk(players, gameID)
	return
}

func (m *Manager) receivePlayer(gameID int, playerID uint64, robotLv int) {
	md := m.matchDataMap[gameID]
	var desk *Desk
	if len(md.desks) == 0 {
		md.maxDeskID++
		deskID := md.maxDeskID
		desk = CreateMatchDesk(deskID, gameID, md.playerCount, m.onDeskMatchFinish)
		md.desks[deskID] = desk
	} else {
		// TODO: 选一个合适的 Desk
		for _, desk = range md.desks {
			break
		}
	}
	desk.AddPlayer(playerID, robotLv)
}

// addRobots 添加机器人
func (m *Manager) addRobots(gameID int) {
	md := m.matchDataMap[gameID]
	now := time.Now()
	for _, desk := range md.desks {
		lastAddTime := desk.GetLastAddPlayerTime()
		if now.Sub(lastAddTime) >= time.Second*5 {
			playerCount := len(desk.GetPlayers())
			addCount := md.playerCount - playerCount
			logrus.WithField("count", addCount).Debugln("添加机器人")
			for i := 0; i < addCount; i++ {
				desk.AddPlayer(GetIdleRobot(1), 1)
			}
		}
	}
}
