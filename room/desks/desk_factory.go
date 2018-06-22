package desks

import (
	"steve/client_pb/room"
	"steve/room/interfaces"
	"steve/room/interfaces/global"
)

type deskFactory struct{}

func (df *deskFactory) CreateDesk(players []uint64, gameID int, opt interfaces.CreateDeskOptions, infos map[uint64][]*room.GeographicalLocation) (interfaces.CreateDeskResult, error) {
	return newDesk(players, gameID, opt, infos)
}

func init() {
	global.SetDeskFactory(new(deskFactory))
}
