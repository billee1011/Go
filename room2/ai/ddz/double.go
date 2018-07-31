package ddz

import (
	"steve/room2/ai"
	"steve/server_pb/ddz"

	"github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
)

type doubleStateAI struct {
}

// GenerateAIEvent 生成 出牌问询AI 事件
// 无论是超时、托管还是机器人，胡过了自动胡，没胡过的其他操作都默认弃， 并且产生相应的事件
func (h *doubleStateAI) GenerateAIEvent(params ai.AIEventGenerateParams) (result ai.AIEventGenerateResult, err error) {
	result, err = ai.AIEventGenerateResult{
		Events: []ai.AIEvent{},
	}, nil

	playerId := params.PlayerID

	context := params.DDZContext
	for _, doubledPlayer := range context.DoubledPlayers {
		if doubledPlayer == playerId { //此用户已加倍
			return
		}
	}

	request := ddz.DoubleRequestEvent{Head: &ddz.RequestEventHead{
		PlayerId: playerId,
	}, IsDouble: false,
	}
	data, _ := proto.Marshal(&request)
	event := ai.AIEvent{
		ID:      int32(ddz.EventID_event_double_request),
		Context: data,
	}
	result.Events = append(result.Events, event)

	logrus.WithField("player", playerId).WithField("result", result).Debug("double timeout event")
	return
}
