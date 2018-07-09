package match

import (
	"steve/client_pb/room"
	"steve/common/data/player"

	"github.com/Sirupsen/logrus"
)

const (
	playersOneDesk int = 4 // 一个桌子需要的人数
)

type matchData struct {
	gameID      int
	queue       chan uint64
	playerCount int
}

// Manager 匹配管理器
type Manager struct {
	queues map[int]matchData
}

var defaultManager = NewManager()

func createMatchData(gameID room.GameId, playerCount int) matchData {
	return matchData{
		gameID:      int(gameID),
		queue:       make(chan uint64),
		playerCount: playerCount,
	}
}

// NewManager 创建 Manager
func NewManager() *Manager {
	// TODO: 从配置加载
	queues := make(map[int]matchData)
	queues[int(room.GameId_GAMEID_XUELIU)] = createMatchData(room.GameId_GAMEID_XUELIU, 4)
	queues[int(room.GameId_GAMEID_XUEZHAN)] = createMatchData(room.GameId_GAMEID_XUEZHAN, 4)
	queues[int(room.GameId_GAMEID_DOUDIZHU)] = createMatchData(room.GameId_GAMEID_DOUDIZHU, 3)

	m := &Manager{
		queues: queues,
	}
	m.runMatchs()
	return m
}

func (m *Manager) runMatchs() {
	for gameID, md := range m.queues {
		go m.match(gameID, md.playerCount, md.queue)
	}
}

func (m *Manager) addPlayer(playerID uint64, gameID int) {
	md, ok := m.queues[gameID]
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
	var players []uint64
	s := NewSender()

	for {
		select {
		case playerID := <-ch:
			online := player.GetPlayerGateAddr(playerID)
			if online == "" {
				logEntry.Debugln("player is not online, remove the player[%d]", playerID)
			}
			players = append(players, playerID)

			if len(players) == playerCount {
				if _, err := s.createDesk(players, gameID); err != nil {
					logEntry.Debug(err.Error())
				}
				players = players[len(players):]
			}
		}
	}

	logEntry.Debugln("Manager::Match() end ...")
}
