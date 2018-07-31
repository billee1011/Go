package ddz

import (
	"steve/client_pb/room"
	"steve/room2/ai"
	"steve/server_pb/ddz"
)

// 注册 AI
func init() {
	ai.GetAtEvent().RegisterAI(int(room.GameId_GAMEID_DOUDIZHU), int32(ddz.StateID_state_grab), &grabStateAI{})
	ai.GetAtEvent().RegisterAI(int(room.GameId_GAMEID_DOUDIZHU), int32(ddz.StateID_state_double), &doubleStateAI{})
	ai.GetAtEvent().RegisterAI(int(room.GameId_GAMEID_DOUDIZHU), int32(ddz.StateID_state_playing), &playStateAI{})
}
