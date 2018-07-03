package tests

import (
	"steve/client_pb/room"
	"steve/room/desks/ddzdesk/flow/ddz/ddzmachine"
	"steve/room/desks/ddzdesk/flow/ddz/states"
	"steve/server_pb/ddz"
)

func createMachine(stateID ddz.StateID) *ddzmachine.DDZMachine {
	ddzContext := &ddz.DDZContext{
		GameId:    int32(room.GameId_GAMEID_DDZ),
		CurState:  stateID,
		Players:   []*ddz.Player{},
		WallCards: []uint32{},
	}
	statefactory := states.NewFactory()
	return ddzmachine.CreateDDZMachine(ddzContext, statefactory, nil)
}
