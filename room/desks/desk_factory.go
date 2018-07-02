package desks

import (
	"steve/room/interfaces"
	"steve/room/interfaces/global"
	//"steve/room/desks/majong"
	//"steve/room/desks/ddz"
	//"steve/client_pb/room"
)

type deskFactory struct{}

func (df *deskFactory) CreateDesk(players []uint64, gameID int, opt interfaces.CreateDeskOptions) (interfaces.CreateDeskResult, error) {
	//if gameID == int(room.GameId_GAMEID_DDZ) {
	//	return ddz.NewDesk(players, gameID, opt)
	//}
	//return majong.NewDesk(players, gameID, opt)
	return newDesk(players, gameID, opt)
}

func init() {
	global.SetDeskFactory(new(deskFactory))
}
