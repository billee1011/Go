package ddz

import (
	"steve/entity/poker/ddz"
	"steve/room/interfaces"

	"github.com/Sirupsen/logrus"
)

type doubleStateAI struct {
}

// GenerateAIEvent 生成 出牌问询AI 事件
// 无论是超时、托管还是机器人，胡过了自动胡，没胡过的其他操作都默认弃， 并且产生相应的事件
func (h *doubleStateAI) GenerateAIEvent(params interfaces.AIEventGenerateParams) (result interfaces.AIEventGenerateResult, err error) {
	result, err = interfaces.AIEventGenerateResult{
		Events: []interfaces.AIEvent{},
	}, nil

	playerID := params.PlayerID

	context := params.DDZContext
	for _, doubledPlayer := range context.DoubledPlayers {
		if doubledPlayer == playerID { //此用户已加倍
			return
		}
	}

	request := &ddz.DoubleRequestEvent{Head: &ddz.RequestEventHead{
		PlayerId: playerID,
	}, IsDouble: false,
	}
	event := interfaces.AIEvent{
		ID:      int32(ddz.EventID_event_double_request),
		Context: request,
	}
	result.Events = append(result.Events, event)

	logrus.WithField("player", playerID).WithField("result", result).Debug("double timeout event")
	return
}
