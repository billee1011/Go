package ddz

import (
	"steve/entity/poker/ddz"
	"steve/room/ai"

	"github.com/Sirupsen/logrus"
)

type grabStateAI struct {
}

// GenerateAIEvent 生成 出牌问询AI 事件
// 无论是超时、托管还是机器人，胡过了自动胡，没胡过的其他操作都默认弃， 并且产生相应的事件
func (h *grabStateAI) GenerateAIEvent(params ai.AIEventGenerateParams) (result ai.AIEventGenerateResult, err error) {
	result, err = ai.AIEventGenerateResult{
		Events: []ai.AIEvent{},
	}, nil

	context := params.DDZContext
	playerID := params.PlayerID
	// 没到自己抢庄
	if context.GetCurrentPlayerId() != playerID {
		return result, nil
	}

	request := &ddz.GrabRequestEvent{
		Head: &ddz.RequestEventHead{
			PlayerId: playerID,
		}, Grab: false,
	}
	event := ai.AIEvent{
		ID:      int32(ddz.EventID_event_grab_request),
		Context: request,
	}
	result.Events = append(result.Events, event)
	logrus.WithField("player", playerID).WithField("result", result).Debug("grab timeout event")
	return
}
