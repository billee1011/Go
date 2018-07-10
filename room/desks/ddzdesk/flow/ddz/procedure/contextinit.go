package procedure

import (
	"steve/client_pb/room"
	"steve/server_pb/ddz"
	"steve/room/peipai/handle"
)

// CreateInitDDZContext 创建初始斗地主现场
func CreateInitDDZContext(players []uint64) *ddz.DDZContext {
	return &ddz.DDZContext{
		GameId:    int32(room.GameId_GAMEID_DOUDIZHU),
		CurState:  ddz.StateID_state_init,
		Players:   createDDZPlayers(players),
		WallCards: []uint32{},
		GrabbedCount: 0,
		AllAbandonCount: 0,
		TotalGrab: 1,
		DoubledCount: 0,
		TotalDouble: 1,
		CurCardType: ddz.CardType_CT_NONE,
		PassCount: 0,
		TotalBomb: 1,
		Spring: true,
		AntiSpring: true,
		Peipai:    handle.GetPeiPai(int(room.GameId_GAMEID_DOUDIZHU)),
	}
}

func createDDZPlayers(players []uint64) []*ddz.Player {
	result := make([]*ddz.Player, 0, len(players))
	for _, playerID := range players {
		result = append(result, &ddz.Player{
			PalyerId: playerID,
			HandCards: []uint32{},
		})
	}
	return result
}
