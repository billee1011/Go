package factory

import (
	"steve/client_pb/room"
	"steve/room/desks/ddzdesk"
	"steve/room/desks/mjdesk"
	"steve/room/interfaces"
	"steve/room/interfaces/global"
	//"steve/room/desks/majong"
	//"steve/room/desks/ddz"
	//"steve/client_pb/room"
)

type deskFactory struct{}

func (df *deskFactory) CreateDesk(deskPlayers []interfaces.DeskPlayer, gameID int, opt interfaces.CreateDeskOptions) (interfaces.CreateDeskResult, error) {
	switch room.GameId(gameID) {
	case room.GameId_GAMEID_DOUDIZHU:
		return ddzdesk.CreateDDZDesk(deskPlayers, gameID, opt, global.GetDeskIDAllocator())
	default:
		return mjdesk.CreateMajongDesk(deskPlayers, gameID, opt, global.GetDeskIDAllocator())

	}
}

func init() {
	global.SetDeskFactory(new(deskFactory))
}
