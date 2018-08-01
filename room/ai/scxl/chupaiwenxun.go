package scxlai

import (
	"fmt"
	"steve/gutils"
	"steve/server_pb/majong"

	"steve/room/ai"

	"github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
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
	if h.checkAIEvent(player, mjContext, params) != nil {
		return
	}
	// if len(player.GetPossibleActions()) == 0 {
	// 	return
	// }
	switch params.AIType {
	case ai.HuAI:
		{
			if gutils.IsTing(player) {
				return
			}
			if h.containAction(player, majong.Action_action_gang) {
				return
			}
			if gutils.IsHu(player) && h.containAction(player, majong.Action_action_hu) {
				//执行胡操作
				if event := h.chupaiWenxun(player); event != nil {
					result.Events = append(result.Events, *event)
				}
			}
		}
	case ai.TingAI:
		return
	case ai.RobotAI, ai.OverTimeAI, ai.TuoGuangAI:
		{
			if event := h.chupaiWenxun(player); event != nil {
				result.Events = append(result.Events, *event)
			}
		}
	}

	return
}

func (h *chupaiWenxunStateAI) containAction(player *majong.Player, action majong.Action) bool {
	for _, possibleAction := range player.GetPossibleActions() {
		if possibleAction == action {
			return true
		}
	}
	return false
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
	return &ai.AIEvent{
		ID:      int32(eventID),
		Context: data,
	}
}

func (h *chupaiWenxunStateAI) checkAIEvent(player *majong.Player, mjContext *majong.MajongContext, params ai.AIEventGenerateParams) error {
	err := fmt.Errorf("不生成自动事件")
	if player.GetHasSelected() {
		return err
	}
	if len(player.GetPossibleActions()) == 0 {
		return err
	}
	if mjContext.GetCurState() != majong.StateID_state_chupaiwenxun {
		return err
	}

	return nil
}
