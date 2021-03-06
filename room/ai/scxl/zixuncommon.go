package scxlai

import (
	"steve/common/mjoption"
	"steve/entity/majong"
	"steve/gutils"
	"steve/room/ai"
)

func (h *zixunStateAI) getNormalZiXunAIEvent(player *majong.Player, mjContext *majong.MajongContext) (aiEvent ai.AIEvent) {
	zxRecord := player.GetZixunRecord()
	handCards := player.GetHandCards()
	canHu := zxRecord.GetEnableZimo()
	if (gutils.IsHu(player) || gutils.IsTing(player)) && canHu {
		aiEvent = h.hu(player)
		return
	}
	// 优先出定缺牌
	if gutils.CheckHasDingQueCard(mjContext, player) {
		for i := len(handCards) - 1; i >= 0; i-- {
			hc := handCards[i]
			if mjoption.GetXingpaiOption(int(mjContext.GetXingpaiOptionId())).EnableDingque &&
				hc.GetColor() == player.GetDingqueColor() {
				aiEvent = h.chupai(player, hc)
				return
			}
		}
	}

	// 正常出牌
	if player.GetMopaiCount() == 0 || mjContext.GetZixunType() == majong.ZixunType_ZXT_CHI || mjContext.GetZixunType() == majong.ZixunType_ZXT_PENG {
		aiEvent = h.chupai(player, handCards[len(handCards)-1])
	} else {
		aiEvent = h.chupai(player, mjContext.GetLastMopaiCard())
	}
	return
}

func (h *zixunStateAI) getNormalZiXunTingStateAIEvent(player *majong.Player, mjContext *majong.MajongContext) (aiEvent ai.AIEvent) {
	// 生成听AI事件
	// 听状态下，能胡不做操作等玩家自行选择或者等超时事件，不能胡就打出摸到的牌
	zxRecord := player.GetZixunRecord()
	if gutils.IsTing(player) {
		canHu := zxRecord.GetEnableZimo()
		if !canHu {
			aiEvent = h.chupai(player, mjContext.GetLastMopaiCard())
		}
	}
	return
}

func (h *zixunStateAI) getNormalZiXunHuStateAIEvent(player *majong.Player, mjContext *majong.MajongContext) (aiEvent ai.AIEvent) {
	// 生成胡AI事件
	// 胡状态下，能胡直接让胡，不能胡就打出摸到的牌
	zxRecord := player.GetZixunRecord()
	if gutils.IsHu(player) {
		canHu := zxRecord.GetEnableZimo()
		if canHu {
			aiEvent = h.hu(player)
		} else {
			aiEvent = h.chupai(player, mjContext.GetLastMopaiCard())
		}
	}
	return
}

func (h *zixunStateAI) chupai(player *majong.Player, card *majong.Card) ai.AIEvent {
	eventContext := majong.ChupaiRequestEvent{
		Head: &majong.RequestEventHead{
			PlayerId: player.GetPlayerId(),
		},
		Cards: card,
	}
	return ai.AIEvent{
		ID:      int32(majong.EventID_event_chupai_request),
		Context: &eventContext,
	}
}

func (h *zixunStateAI) gang(player *majong.Player, card *majong.Card) ai.AIEvent {
	eventContext := majong.GangRequestEvent{
		Head: &majong.RequestEventHead{
			PlayerId: player.GetPlayerId(),
		},
		Card: card,
	}
	return ai.AIEvent{
		ID:      int32(majong.EventID_event_gang_request),
		Context: &eventContext,
	}
}

func (h *zixunStateAI) hu(player *majong.Player) ai.AIEvent {
	eventContext := majong.HuRequestEvent{
		Head: &majong.RequestEventHead{
			PlayerId: player.GetPlayerId(),
		},
	}
	return ai.AIEvent{
		ID:      int32(majong.EventID_event_hu_request),
		Context: &eventContext,
	}
}
