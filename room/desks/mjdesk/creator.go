package mjdesk

import (
	"steve/room/desks/deskbase"
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
	deskPlayerMgr := deskbase.CreateDeskPlayerMgr()
	deskPlayerMgr.SetPlayers(players)
	return interfaces.CreateDeskResult{
		Desk: &desk{
			DeskPlayerMgr: deskPlayerMgr,
			deskUID:       id,
			gameID:        gameID,
			settler:       global.GetDeskSettleFactory().CreateDeskSettler(gameID),
			event:         make(chan deskEvent, 16),
		},
	}, nil
}
