package desks

import (
	"steve/room/interfaces"
	"steve/room/interfaces/global"
)

type deskFactory struct{}

func (df *deskFactory) CreateDesk(players []uint64, gameID int, opt interfaces.CreateDeskOptions) (interfaces.CreateDeskResult, error) {
	return newDesk(players, gameID, opt)
}

func init() {
	global.SetDeskFactory(new(deskFactory))
}
