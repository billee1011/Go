package scxlai

import (
	"fmt"
	"steve/common/mjoption"
	"steve/entity/majong"
	"steve/gutils"
	"time"

	"steve/room2/ai"
)

type zixunStateAI struct {
	maxDingqueTime time.Duration // 最大定缺时间
}

// // 注册 AI
// func init() {
// 	g := global.GetDeskAutoEventGenerator()
// 	g.RegisterAI(gGameID, majong.StateID_state_zixun, &zixunStateAI{})
// }

// GenerateAIEvent 生成 AI 事件
// 前端排序，定缺牌在最右侧，其他手牌按花色万条筒、以及点数大小从左到右排序
// 首先判断玩家是否时当前可以操作的玩家
// 是的话,判断当前玩家是否可以执行自动事件
// 可以的话,根据玩家状态生成不同的自动事件
// 1,玩家是碰自询:
//	 			之前胡过,自动事件:出最右的一张牌
//				之前没有胡过,自动事件:出最右的一张牌(如果有定缺牌，优先出定缺牌)
// 2,玩家是摸牌自询:
//	 			之前胡过,自动事件:
//								可胡,等待三秒,然后自动胡牌
//								不可胡,无需等待,直接出牌
//				之前没有胡过,自动事件:
// 								1,出摸到的那张牌
//								2,如果是庄家首次出牌,出最右侧的牌
func (h *zixunStateAI) GenerateAIEvent(params ai.AIEventGenerateParams) (result ai.AIEventGenerateResult, err error) {
	result, err = ai.AIEventGenerateResult{
		Events: []ai.AIEvent{},
	}, nil
	var aiEvent ai.AIEvent
	mjContext := params.MajongContext
	player := gutils.GetMajongPlayer(params.PlayerID, mjContext)
	handCards := player.GetHandCards()
	if gutils.GetZixunPlayer(mjContext) != params.PlayerID {
		return result, fmt.Errorf("当前玩家不允许进行自动操作")
	}
	if len(handCards) < 2 {
		return result, fmt.Errorf("手牌数量少于2")
	}
	switch mjContext.GetZixunType() {
	case majong.ZixunType_ZXT_PENG:
		{
			//有定缺牌，出最大的定缺牌
			hasChuPai := false
			for i := len(handCards) - 1; i >= 0; i-- {
				hc := handCards[i]
				if hc.GetColor() == player.GetDingqueColor() {
					aiEvent = h.chupai(player, hc)
					hasChuPai = true
					break
				}
			}
			if !hasChuPai {
				aiEvent = h.chupai(player, handCards[len(handCards)-1])
			}
		}
	case majong.ZixunType_ZXT_NORMAL:
		{
			zxRecord := player.GetZixunRecord()
			canHu := zxRecord.GetEnableZimo()
			if (gutils.IsTing(player) || gutils.IsHu(player)) && canHu && !gutils.CheckHasDingQueCard(mjContext, player) {
				aiEvent = h.hu(player)
			} else {
				//先判断是否有定缺牌，有的话，先出定缺牌
				//有定缺牌，出最大的定缺牌
				hasChuPai := false
				for i := len(handCards) - 1; i >= 0; i-- {
					hc := handCards[i]
					if mjoption.GetXingpaiOption(int(mjContext.GetXingpaiOptionId())).EnableDingque &&
						hc.GetColor() == player.GetDingqueColor() {
						aiEvent = h.chupai(player, hc)
						hasChuPai = true
						break
					}
				}
				if !hasChuPai {
					if player.GetMopaiCount() == 0 {
						aiEvent = h.chupai(player, handCards[len(handCards)-1])
					} else {
						aiEvent = h.chupai(player, mjContext.GetLastMopaiCard())
					}
				}
			}
		}
	default:
		return
	}
	result.Events = append(result.Events, aiEvent)
	return
}

func (h *zixunStateAI) chupai(player *majong.Player, card *majong.Card) ai.AIEvent {
	eventContext := &majong.ChupaiRequestEvent{
		Head: &majong.RequestEventHead{
			PlayerId: player.GetPalyerId(),
		},
		Cards: card,
	}

	return ai.AIEvent{
		ID:      majong.EventID_event_chupai_request,
		Context: eventContext,
	}
}

func (h *zixunStateAI) hu(player *majong.Player) ai.AIEvent {
	eventContext := &majong.HuRequestEvent{
		Head: &majong.RequestEventHead{
			PlayerId: player.GetPalyerId(),
		},
	}

	return ai.AIEvent{
		ID:      majong.EventID_event_hu_request,
		Context: eventContext,
	}
}
