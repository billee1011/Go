package procedure

import (
	"steve/client_pb/room"
	"steve/server_pb/ddz"
)

// CreateInitDDZContext 创建初始斗地主现场
func CreateInitDDZContext(players []uint64) *ddz.DDZContext {
	return &ddz.DDZContext{
		GameId:            int32(room.GameId_GAMEID_DOUDIZHU),
		CurState:          ddz.StateID_state_init,
		Players:           createDDZPlayers(players),
		WallCards:         []uint32{},
		FirstGrabPlayerId: 0,
		GrabbedCount:      0,
		AllAbandonCount:   0,
		TotalGrab:         0,
		DoubledCount:      0,
		TotalDouble:       1,
		CurCardType:       ddz.CardType_CT_NONE,
		PassCount:         0,
		TotalBomb:         1,
		Spring:            true,
		AntiSpring:        true,
	}
}

func createDDZPlayers(players []uint64) []*ddz.Player {
	result := make([]*ddz.Player, 0, len(players))
	for _, playerID := range players {
		result = append(result, &ddz.Player{
			PlayerId:  playerID,
			HandCards: []uint32{},
		})
	}
	return result
}
