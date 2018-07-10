package ddz

import (
	"steve/room/interfaces"
	"github.com/golang/protobuf/proto"
	"steve/room/interfaces/global"
	"steve/client_pb/room"
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
}

// GenerateAIEvent 生成 出牌问询AI 事件
// 无论是超时、托管还是机器人，胡过了自动胡，没胡过的其他操作都默认弃， 并且产生相应的事件
func (h *grabStateAI) GenerateAIEvent(params interfaces.AIEventGenerateParams) (result interfaces.AIEventGenerateResult, err error) {
	result, err = interfaces.AIEventGenerateResult{
		Events: []interfaces.AIEvent{},
	}, nil

	context := params.DDZContext
	for _, playerId := range context.CountDownPlayers {
		request := ddz.GrabRequestEvent{Head:
			&ddz.RequestEventHead{
				PlayerId: playerId,
			}, Grab:false,
		}
		data, _ := proto.Marshal(&request)
		event := interfaces.AIEvent{
			ID:      int32(ddz.EventID_event_grab_request),
			Context: data,
		}
		result.Events = append(result.Events, event)
	}
	logrus.WithField("players", context.CountDownPlayers).WithField("result", result).Debug("grab timeout event")

	context.Duration = 0//清除倒计时
	return
}
