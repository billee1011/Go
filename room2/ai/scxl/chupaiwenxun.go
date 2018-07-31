package scxlai

import (
	"steve/entity/majong"
	"steve/gutils"
	"steve/room2/ai"
)

type chupaiWenxunStateAI struct {
}

// GenerateAIEvent 生成 出牌问询AI 事件
// 无论是超时、托管还是机器人，胡过了自动胡，没胡过的其他操作都默认弃， 并且产生相应的事件
func (h *chupaiWenxunStateAI) GenerateAIEvent(params ai.AIEventGenerateParams) (result ai.AIEventGenerateResult, err error) {
	result, err = ai.AIEventGenerateResult{
		Events: []ai.AIEvent{},
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
	if gutils.IsTing(player) || gutils.IsHu(player) {
		for _, possibleAction := range player.GetPossibleActions() {
			if possibleAction == majong.Action_action_hu {
				action = majong.Action_action_hu
			}
		}
	}
	return action
}

// chupaiWenxun 生成出牌问询请求事件
func (h *chupaiWenxunStateAI) chupaiWenxun(player *majong.Player) *ai.AIEvent {
	var (
		event ai.AIEvent
	)
	action := h.getAction(player)

	switch action {
	case majong.Action_action_hu:
		event.Context = &majong.HuRequestEvent{
			Head: &majong.RequestEventHead{
				PlayerId: player.GetPalyerId(),
			},
		}
		event.ID = majong.EventID_event_hu_request
	default:
		event.Context = &majong.QiRequestEvent{
			Head: &majong.RequestEventHead{
				PlayerId: player.GetPalyerId(),
			},
		}
		event.ID = majong.EventID_event_qi_request
	}

	return &event
}
