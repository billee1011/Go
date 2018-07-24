package ddz

import (
	"steve/client_pb/room"
	"steve/room/interfaces"
	"steve/room/interfaces/global"
	"steve/server_pb/ddz"

	"github.com/Sirupsen/logrus"
)

type grabStateAI struct {
}

// 注册 AI
func init() {
	g := global.GetDeskAutoEventGenerator()
	g.RegisterAI(int(room.GameId_GAMEID_DOUDIZHU), int32(ddz.StateID_state_grab), &grabStateAI{})
	g.RegisterAI(int(room.GameId_GAMEID_DOUDIZHU), int32(ddz.StateID_state_double), &doubleStateAI{})
	g.RegisterAI(int(room.GameId_GAMEID_DOUDIZHU), int32(ddz.StateID_state_playing), &playStateAI{})
}

// GenerateAIEvent 生成 出牌问询AI 事件
// 无论是超时、托管还是机器人，胡过了自动胡，没胡过的其他操作都默认弃， 并且产生相应的事件
func (h *grabStateAI) GenerateAIEvent(params interfaces.AIEventGenerateParams) (result interfaces.AIEventGenerateResult, err error) {
	result, err = interfaces.AIEventGenerateResult{
		Events: []interfaces.AIEvent{},
	}, nil

	context := params.DDZContext
	playerID := params.PlayerID
	// 没到自己抢庄
	if context.GetCurrentPlayerId() != playerID {
		return result, nil
	}

	request := ddz.GrabRequestEvent{
		Head: &ddz.RequestEventHead{
			PlayerId: playerID,
		}, Grab: false,
	}
	event := interfaces.AIEvent{
		ID:      int32(ddz.EventID_event_grab_request),
		Context: request,
	}
	result.Events = append(result.Events, event)
	logrus.WithField("player", playerID).WithField("result", result).Debug("grab timeout event")
	return
}
