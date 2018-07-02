package mjdesk

import (
	"steve/room/desks/deskplayer"
	"steve/room/desks/tuoguan"
	"steve/room/interfaces"
	"steve/room/interfaces/global"

	"github.com/Sirupsen/logrus"
)

// CreateMajongDesk 创建麻将房间
func CreateMajongDesk(players []uint64, gameID int, opt interfaces.CreateDeskOptions, alloc interfaces.DeskIDAllocator) (result interfaces.CreateDeskResult, err error) {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "newDesk",
		"game_id":   gameID,
		"players":   players,
	})
	id, err := alloc.AllocDeskID()
	if err != nil {
		logEntry.Errorln(errAllocDeskIDFailed)
		err = errAllocDeskIDFailed
		return
	}
	logEntry = logEntry.WithField("desk_uid", id)
	deskPlayers, err := makeDeskPlayers(logEntry, players)
	if err != nil {
		return
	}
	return interfaces.CreateDeskResult{
		Desk: &desk{
			deskUID:      id,
			gameID:       gameID,
			createOption: opt,
			settler:      global.GetDeskSettleFactory().CreateDeskSettler(gameID),
			players:      deskPlayers,
			event:        make(chan deskEvent, 16),
			tuoGuanMgr:   tuoguan.CreateTuoguanManager(),
			enterQuits:   make(chan enterQuitInfo),
		},
	}, nil
}

func makeDeskPlayers(logEntry *logrus.Entry, players []uint64) (map[uint32]interfaces.DeskPlayer, error) {
	playerMgr := global.GetPlayerMgr()
	deskPlayers := make(map[uint32]interfaces.DeskPlayer, 4)
	seat := uint32(0)
	for _, playerID := range players {
		player := playerMgr.GetPlayer(playerID)
		if player == nil {
			logEntry.WithField("player_id", playerID).Errorln(errPlayerNotExist)
			return nil, errPlayerNotExist
		}
		deskPlayers[seat] = deskplayer.CreateDeskPlayer(playerID, seat)
		seat++
	}
	return deskPlayers, nil
}
