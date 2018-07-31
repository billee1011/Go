package tests

import (
	"steve/client_pb/room"
	"steve/room2/flows/ddzflow/ddz/ddzmachine"
	"steve/room2/flows/ddzflow/ddz/states"
	"steve/server_pb/ddz"
)

func createMachine(stateID ddz.StateID) *ddzmachine.DDZMachine {
	ddzContext := &ddz.DDZContext{
		GameId:   int32(room.GameId_GAMEID_DOUDIZHU),
		CurState: stateID,
		Players:  []*ddz.Player{{PlayerId: 1}, {PlayerId: 2}, {PlayerId: 3}},
	}
	statefactory := states.NewFactory()
	return ddzmachine.CreateDDZMachine(ddzContext, statefactory, nil)
}
