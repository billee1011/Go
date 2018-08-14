package scxlai

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
	"math/rand"
	"sort"
	"steve/entity/majong"
	"steve/gutils"
	"steve/room/ai"
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
	case ai.RobotAI:
		{
			if event := h.askMiddleAI(player, *mjContext.LastOutCard); event != nil {
				result.Events = append(result.Events, *event)
			}
		}
	case ai.OverTimeAI, ai.SpecialOverTimeAI, ai.TuoGuangAI:
		{
			if viper.GetBool("ai.test") {
				if event := h.askMiddleAI(player, *mjContext.LastOutCard); event != nil {
					result.Events = append(result.Events, *event)
				}
			} else {
				if event := h.chupaiWenxun(player); event != nil {
					result.Events = append(result.Events, *event)
				}
			}
		}
	}

	return
}

func (h *chupaiWenxunStateAI) askMiddleAI(player *majong.Player, lastOutCard majong.Card) *ai.AIEvent {
	logEntry := logrus.WithField("playerId", player.PlayerId)
	var (
		event ai.AIEvent
	)
	actions := player.GetPossibleActions()
	sort.Sort(sort.Reverse(majong.ActionSlice(actions))) //按优先级从高到低排列

	for _, action := range actions {
		switch action {
		case majong.Action_action_hu:
			event.Context = &majong.HuRequestEvent{
				Head: &majong.RequestEventHead{
					PlayerId: player.GetPlayerId(),
				},
			}
			event.ID = int32(majong.EventID_event_hu_request)
			logEntry.WithField("点炮牌", lastOutCard).Infoln("中级AI点炮胡牌")
			return &event
		case majong.Action_action_gang:
			_, keZis, _, _, _, _, _ := SplitBestCards(NonPointer(player.HandCards))
			if len(keZis) > 0 && Contains(keZis, lastOutCard) {
				event.Context = &majong.GangRequestEvent{
					Head: &majong.RequestEventHead{
						PlayerId: player.GetPlayerId(),
					},
					Card: &lastOutCard,
				}
				event.ID = int32(majong.EventID_event_gang_request)
				logEntry.WithField("明杠牌", lastOutCard).Infoln("中级AI明杠")
				return &event
			}
		case majong.Action_action_peng:
			_, _, pairs, _, _, _, _ := SplitBestCards(NonPointer(player.HandCards))
			if len(pairs) > 0 && Contains(pairs, lastOutCard) {
				r := rand.Intn(100)
				if len(pairs) >= 2 && r < 90 || len(pairs) == 1 && r < 10 { //多于1对时，碰牌概率90%；等于1对时，碰牌概率10%
					event.Context = &majong.PengRequestEvent{
						Head: &majong.RequestEventHead{
							PlayerId: player.GetPlayerId(),
						},
					}
					event.ID = int32(majong.EventID_event_peng_request)
					logEntry.WithField("碰牌", lastOutCard).Infoln("中级AI碰牌")
					return &event
				}
			}
		case majong.Action_action_chi:
			_, _, _, doubleChas, singleChas, _, _ := SplitBestCards(NonPointer(player.HandCards))
			if len(singleChas)+len(doubleChas) > 0 {
				for _, cha := range append(singleChas, doubleChas...) { //优先处理单茬
					validCards := getValidCard(cha)
					if ContainsCard(validCards, lastOutCard) {
						event.Context = &majong.ChiRequestEvent{
							Head: &majong.RequestEventHead{
								PlayerId: player.GetPlayerId(),
							},
							Cards: []*majong.Card{&cha.cards[0], &cha.cards[1], &lastOutCard},
						}
						event.ID = int32(majong.EventID_event_chi_request)
						logEntry.WithField("吃牌", lastOutCard).Infoln("中级AI吃牌")
						return &event
					}
				}
			}
		default:
			event.Context = &majong.QiRequestEvent{
				Head: &majong.RequestEventHead{
					PlayerId: player.GetPlayerId(),
				},
			}
			event.ID = int32(majong.EventID_event_qi_request)
		}
	}

	return &event
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
		event ai.AIEvent
	)
	action := h.getAction(player)

	switch action {
	case majong.Action_action_hu:
		event.Context = &majong.HuRequestEvent{
			Head: &majong.RequestEventHead{
				PlayerId: player.GetPlayerId(),
			},
		}
		event.ID = int32(majong.EventID_event_hu_request)
	default:
		event.Context = &majong.QiRequestEvent{
			Head: &majong.RequestEventHead{
				PlayerId: player.GetPlayerId(),
			},
		}
		event.ID = int32(majong.EventID_event_qi_request)
	}
	return &event
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
