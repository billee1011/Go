package ddz

import (
	"github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
	"steve/room/interfaces"
	"steve/server_pb/ddz"
)

type doubleStateAI struct {
}

// GenerateAIEvent 生成 出牌问询AI 事件
// 无论是超时、托管还是机器人，胡过了自动胡，没胡过的其他操作都默认弃， 并且产生相应的事件
func (h *doubleStateAI) GenerateAIEvent(params interfaces.AIEventGenerateParams) (result interfaces.AIEventGenerateResult, err error) {
	result, err = interfaces.AIEventGenerateResult{
		Events: []interfaces.AIEvent{},
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
	event := interfaces.AIEvent{
		ID:      int32(ddz.EventID_event_double_request),
		Context: data,
	}
	result.Events = append(result.Events, event)

	logrus.WithField("player", playerId).WithField("result", result).Debug("double timeout event")
	return
}
