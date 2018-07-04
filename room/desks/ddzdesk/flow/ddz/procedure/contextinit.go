package procedure

import (
	"steve/client_pb/room"
	"steve/server_pb/ddz"
)

// CreateInitDDZContext 创建初始斗地主现场
func CreateInitDDZContext(players []uint64) *ddz.DDZContext {
	return &ddz.DDZContext{
		GameId:    int32(room.GameId_GAMEID_DDZ),
		CurState:  ddz.StateID_state_init,
		Players:   createDDZPlayers(players),
		WallCards: []uint32{},
	}
}

func createDDZPlayers(players []uint64) []*ddz.Player {
	result := make([]*ddz.Player, 0, len(players))
	for _, playerID := range players {
		result = append(result, &ddz.Player{
			PalyerId: playerID,
			HandCards:    []uint32{},
		})
	}
	return result
}
