package manager

import (
	"github.com/Sirupsen/logrus"
	"steve/common/data/player"
	"steve/match/core"
)

const (
	playersOneDesk int = 4 // 一个桌子需要的人数
)


type Manager struct {
	queue chan uint64
}

func New() *Manager{
	return &Manager{
		queue: make(chan uint64),
	}
}

// 加入一个新的匹配玩家
func (m *Manager) AddPlayer(playerID uint64) {
	m.queue<-playerID
}

// 具体的匹配操作
func (m *Manager) Match(sender *core.Sender) {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "Manager::Match start ...",
	})

	logEntry.Debugln("Manager::Match()")
	var players []uint64

	select {
	case playerID := <- m.queue:
		online := player.GetPlayerGateAddr(playerID)
		if online == "" {
			logEntry.Debugln("玩家不在线，剔除玩家[%d]", playerID)
		}
		players = append(players, playerID)

		if len(players) == playersOneDesk {
			if err := sender.CreateDesk(players); err != nil {
				logEntry.Debug(err.Error())
			}

			players = players[len(players):]
		}
	}

	logEntry.Debugln("Manager::Match() end ...")
}
