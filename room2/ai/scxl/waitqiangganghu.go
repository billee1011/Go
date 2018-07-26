package scxlai

import (
	"fmt"
	"steve/gutils"

	"steve/entity/majong"
	"time"

	"steve/room2/ai"

	"github.com/Sirupsen/logrus"
<<<<<<< HEAD
	"github.com/golang/protobuf/proto"
=======
>>>>>>> majong
)

type waitQiangganghuStateAI struct {
	maxDingqueTime time.Duration // 最大定缺时间
}

// 注册 AI
// func init() {
// 	g := global.GetDeskAutoEventGenerator()
// 	g.RegisterAI(gGameID, majong.StateID_state_waitqiangganghu, &waitQiangganghuStateAI{})
// }

// GenerateAIEvent 生成 AI 事件
// 等待抢杠胡的状态下
// 首先判断请求的自动事件是否可以进行操作
// 可以的话处理
// 如果玩家开过胡,那么自动给胡
// 如果玩家没开过胡,那么选择过
func (h *waitQiangganghuStateAI) GenerateAIEvent(params ai.AIEventGenerateParams) (result ai.AIEventGenerateResult, err error) {
	result, err = ai.AIEventGenerateResult{
		Events: []ai.AIEvent{},
	}, nil
	var aiEvent ai.AIEvent
	mjContext := params.MajongContext
	player := gutils.GetMajongPlayer(params.PlayerID, mjContext)
	if player.GetPalyerId() == mjContext.GetLastGangPlayer() {
		return result, fmt.Errorf("玩家%v是补杠的玩家,不允许抢杠", player.GetPalyerId())

	}
	if len(player.GetHandCards())%3+1 != 2 {
		return result, fmt.Errorf("玩家%v手牌不符合查胡要求", player.GetPalyerId())
	}

	if gutils.CheckHasDingQueCard(mjContext, player) {
		return result, fmt.Errorf("")
	}
	canhu := false
	for _, act := range player.GetPossibleActions() {
		if act == majong.Action_action_hu {
			canhu = true
			break
		}
	}
	entry := logrus.WithFields(logrus.Fields{
		"playerID":   player.GetPalyerId(),
		"handCards":  gutils.FmtMajongpbCards(player.GetHandCards()),
		"bugangCard": gutils.FmtMajongpbCards([]*majong.Card{mjContext.GetGangCard()}),
		"canhu":      canhu,
	})
	if canhu {
		if len(player.GetHuCards()) > 0 {
			aiEvent = h.hu(player)
			entry.Info("生成抢杠胡的自动事件")
		} else {
			aiEvent = h.qi(player)
			entry.Info("生成弃的自动事件")
		}
	} else {
		return result, fmt.Errorf("")
	}
	result.Events = append(result.Events, aiEvent)
	return
}

func (h *waitQiangganghuStateAI) qi(player *majong.Player) ai.AIEvent {
	eventContext := &majong.QiRequestEvent{
		Head: &majong.RequestEventHead{
			PlayerId: player.GetPalyerId(),
		},
	}

	return ai.AIEvent{
		ID:      majong.EventID_event_qi_request,
		Context: eventContext,
	}
}

func (h *waitQiangganghuStateAI) hu(player *majong.Player) ai.AIEvent {
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
