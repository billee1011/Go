package scxlai

import (
	"fmt"
	"steve/common/mjoption"
	"steve/entity/majong"
	"steve/gutils"
	"steve/room/ai"
	"time"
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
	mjContext := params.MajongContext
	player := gutils.GetMajongPlayer(params.PlayerID, mjContext)
	if h.checkAIEvent(player, mjContext, params) != nil {
		return
	}
	var aiEvent ai.AIEvent
	switch mjContext.GetZixunType() {
	case majong.ZixunType_ZXT_PENG, majong.ZixunType_ZXT_CHI:
		aiEvent = h.handleOtherZixun(player, mjContext)
	case majong.ZixunType_ZXT_NORMAL:
		aiEvent = h.handleNormalZixun(player, mjContext, params)
	default:
		return
	}

	result, err = ai.AIEventGenerateResult{
		Events: []ai.AIEvent{},
	}, nil

	result.Events = append(result.Events, aiEvent)
	return
}

func (h *zixunStateAI) handleNormalZixun(player *majong.Player, mjContext *majong.MajongContext, params ai.AIEventGenerateParams) (aiEvent ai.AIEvent) {
	switch params.AIType {
	case ai.OverTimeAI, ai.SpecialOverTimeAI:
		aiEvent = h.generateOverTime(player, mjContext)
	case ai.TuoGuangAI:
		aiEvent = h.getMiddleAIEvent(player, mjContext)
	case ai.RobotAI:
		aiEvent = h.getMiddleAIEvent(player, mjContext)
	case ai.TingAI:
		aiEvent = h.generateTing(player, mjContext)
	case ai.HuAI:
		aiEvent = h.generateHu(player, mjContext)
	}
	return aiEvent
}

func (h *zixunStateAI) handleOtherZixun(player *majong.Player, mjContext *majong.MajongContext) (aiEvent ai.AIEvent) {
	//有定缺牌，出最大的定缺牌
	handCards := player.GetHandCards()
	//优先出定缺牌
	for i := len(handCards) - 1; i >= 0; i-- {
		hc := handCards[i]
		if mjoption.GetXingpaiOption(int(mjContext.GetXingpaiOptionId())).EnableDingque &&
			hc.GetColor() == player.GetDingqueColor() {
			aiEvent = h.chupai(player, hc)
			return
		}
	}

	aiEvent = h.chupai(player, handCards[len(handCards)-1])
	return
}

func (h *zixunStateAI) generateOverTime(player *majong.Player, mjContext *majong.MajongContext) (aiEvent ai.AIEvent) {
	// 生成超时AI事件
	// 超时时，可以胡就是胡，否则是出牌
	return h.getNormalZiXunAIEvent(player, mjContext)
}

func (h *zixunStateAI) generateTuoGuang(player *majong.Player, mjContext *majong.MajongContext) (aiEvent ai.AIEvent) {
	// 生成托管AI事件
	// 无状态下才会托管，因此托管就是出牌
	return h.getNormalZiXunAIEvent(player, mjContext)
}

func (h *zixunStateAI) generateRobot(player *majong.Player, mjContext *majong.MajongContext) (aiEvent ai.AIEvent) {
	// 生成机器人AI事件
	// 一直无状态，所以只会出摸到的牌
	return h.getNormalZiXunAIEvent(player, mjContext)
}

func (h *zixunStateAI) generateTing(player *majong.Player, mjContext *majong.MajongContext) (aiEvent ai.AIEvent) {
	// 生成听AI事件
	// 听状态下，能胡不做操作等玩家自行选择或者等超时事件，不能胡就打出摸到的牌
	return h.getNormalZiXunTingStateAIEvent(player, mjContext)
}

func (h *zixunStateAI) generateHu(player *majong.Player, mjContext *majong.MajongContext) (aiEvent ai.AIEvent) {
	// 生成胡AI事件
	// 胡状态下，能胡直接让胡，不能胡就打出摸到的牌
	return h.getNormalZiXunHuStateAIEvent(player, mjContext)
}

func (h *zixunStateAI) checkAIEvent(player *majong.Player, mjContext *majong.MajongContext, params ai.AIEventGenerateParams) error {

	if mjContext.GetCurState() != majong.StateID_state_zixun {
		return fmt.Errorf("当前不是自询状态")
	}
	if gutils.GetZixunPlayer(mjContext) != params.PlayerID {
		return fmt.Errorf("当前玩家不允许进行自动操作")
	}
	if len(player.GetHandCards()) < 2 {
		return fmt.Errorf("手牌数量少于2")
	}
	record := player.GetZixunRecord()
	if params.AIType == ai.TingAI {
		if record.GetEnableZimo() || len(record.GetEnableAngangCards()) > 0 ||
			len(record.GetEnableBugangCards()) > 0 {
			return fmt.Errorf("听的类型下，玩家有特殊操作的时候不处理")
		}
	}
	if params.AIType == ai.HuAI {
		if len(record.GetEnableAngangCards()) > 0 ||
			len(record.GetEnableBugangCards()) > 0 {
			return fmt.Errorf("胡的类型下，玩家有除了胡之外的特殊操作不处理")
		}
	}
	return nil
}
