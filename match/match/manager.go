package match

import (
	"github.com/Sirupsen/logrus"
	"steve/common/data/player"
)

const (
	playersOneDesk int = 4 // 一个桌子需要的人数
)

type Manager struct {
	queue chan uint64
}

var defaultManager = NewManager()

func NewManager() *Manager{
	m := &Manager {
		queue: make(chan uint64),
	}
	go m.match()
	return m
}


// 加入一个新的匹配玩家
func (m *Manager)addPlayer(playerID uint64) {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "Manager::addPlayer()",
	})

	logEntry.Debugln("add player = ", playerID)
	m.queue <- playerID
}

// 具体的匹配操作
func (m *Manager)match() {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "Manager::match()",
	})

	logEntry.Debugln("Manager::Match() start ...")
	var players []uint64
	s := NewSender()

	for {
		select {
		case playerID := <-m.queue:
			online := player.GetPlayerGateAddr(playerID)
			if online == "" {
				logEntry.Debugln("player is not online, remove the player[%d]", playerID)
			}
			players = append(players, playerID)

			if len(players) == playersOneDesk {
				if _, err := s.createDesk(players); err != nil {
					logEntry.Debug(err.Error())
				}
				players = players[len(players):]
			}
		}
	}

	logEntry.Debugln("Manager::Match() end ...")
}
