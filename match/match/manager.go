package match

import (
	"steve/client_pb/room"
	"steve/common/data/player"

	"github.com/Sirupsen/logrus"
)

const (
	playersOneDesk int = 4 // 一个桌子需要的人数
)

type Manager struct {
	queue  chan uint64
	queues map[int](chan uint64)
}

var defaultManager = NewManager()

// makeQueues 创建所有游戏的队列
func makeQueues() map[int](chan uint64) {
	// TODO: 从配置加载
	queues := make(map[int](chan uint64))
	queues[int(room.GameId_GAMEID_XUELIU)] = make(chan uint64)
	queues[int(room.GameId_GAMEID_XUEZHAN)] = make(chan uint64)
	return queues
}

// NewManager 创建 Manager
func NewManager() *Manager {

	m := &Manager{
		queue:  make(chan uint64),
		queues: makeQueues(),
	}
	m.runMatchs()
	return m
}

func (m *Manager) runMatchs() {
	for gameID, ch := range m.queues {
		go m.match(gameID, ch)
	}
}

func (m *Manager) addPlayer(playerID uint64, gameID int) {
	if m.queues[gameID] == nil {
		return
	}
	m.queues[gameID] <- playerID
}

// 具体的匹配操作
func (m *Manager) match(gameID int, ch chan uint64) {
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

			if len(players) == playersOneDesk {
				if _, err := s.createDesk(players, gameID); err != nil {
					logEntry.Debug(err.Error())
				}
				players = players[len(players):]
			}
		}
	}

	logEntry.Debugln("Manager::Match() end ...")
}
