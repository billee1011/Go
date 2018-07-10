package scxlai

import (
	"steve/gutils"
	"steve/room/interfaces"
	"steve/server_pb/majong"

	"github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
)

type chupaiWenxunStateAI struct {
}

// GenerateAIEvent 生成 出牌问询AI 事件
// 无论是超时、托管还是机器人，胡过了自动胡，没胡过的其他操作都默认弃， 并且产生相应的事件
func (h *chupaiWenxunStateAI) GenerateAIEvent(params interfaces.AIEventGenerateParams) (result interfaces.AIEventGenerateResult, err error) {
	result, err = interfaces.AIEventGenerateResult{
		Events: []interfaces.AIEvent{},
	}, nil

	mjContext := params.MajongContext
	player := gutils.GetMajongPlayer(params.PlayerID, mjContext)
	if player.GetHasSelected() {
		return
	}
	if len(player.GetPossibleActions()) == 0 {
		return
	}
	if event := h.chupaiWenxun(player); event != nil {
		result.Events = append(result.Events, *event)
	}
	return
}

// getAction 获取问询动作
func (h *chupaiWenxunStateAI) getAction(player *majong.Player) majong.Action {
	action := majong.Action_action_qi
	if len(player.HuCards) != 0 {
		for _, possibleAction := range player.GetPossibleActions() {
			if possibleAction == majong.Action_action_hu {
				action = majong.Action_action_hu
			}
		}
	}
	return action
}

// chupaiWenxun 生成出牌问询请求事件
func (h *chupaiWenxunStateAI) chupaiWenxun(player *majong.Player) *interfaces.AIEvent {
	var (
		data    []byte
		err     error
		eventID majong.EventID
	)
	action := h.getAction(player)

	switch action {
	case majong.Action_action_hu:
		mjContext := majong.HuRequestEvent{
			Head: &majong.RequestEventHead{
				PlayerId: player.GetPalyerId(),
			},
		}
		eventID = majong.EventID_event_hu_request
		data, err = proto.Marshal(&mjContext)
	default:
		mjContext := majong.QiRequestEvent{
			Head: &majong.RequestEventHead{
				PlayerId: player.GetPalyerId(),
			},
		}
		eventID = majong.EventID_event_qi_request
		data, err = proto.Marshal(&mjContext)
	}
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"func_name": "chupaiWenxunStateAI.chupaiWenxun",
			"player_id": player.GetPalyerId(),
			"action":    action,
		}).Errorln("事件序列化失败")
		return nil
	}
	return &interfaces.AIEvent{
		ID:      int32(eventID),
		Context: data,
	}
}
