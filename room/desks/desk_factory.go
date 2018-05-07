package desks

import "steve/room/interfaces"

type deskFactory struct{}

func (df *deskFactory) CreateDesk(players []uint64, gameID int, opt interfaces.CreateDeskOptions) (interfaces.CreateDeskResult, error) {
	// TODO
	return interfaces.CreateDeskResult{}, nil
}
