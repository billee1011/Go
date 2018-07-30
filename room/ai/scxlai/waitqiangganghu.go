package scxlai

import (
	"fmt"
	"steve/gutils"

	"steve/room/interfaces"
	"steve/server_pb/majong"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
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
func (h *waitQiangganghuStateAI) GenerateAIEvent(params interfaces.AIEventGenerateParams) (result interfaces.AIEventGenerateResult, err error) {
	result, err = interfaces.AIEventGenerateResult{
		Events: []interfaces.AIEvent{},
	}, nil
	var aiEvent interfaces.AIEvent
	mjContext := params.MajongContext
	player := gutils.GetMajongPlayer(params.PlayerID, mjContext)
	if h.checkAIEvent(player, mjContext, params) != nil {
		return
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
		if gutils.IsHu(player) || gutils.IsTing(player) {
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

func (h *waitQiangganghuStateAI) qi(player *majong.Player) interfaces.AIEvent {
	eventContext := majong.QiRequestEvent{
		Head: &majong.RequestEventHead{
			PlayerId: player.GetPalyerId(),
		},
	}
	data, err := proto.Marshal(&eventContext)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"func_name": "waitQiangganghuStateAI.qi",
			"player_id": player.GetPalyerId(),
		}).Errorln("事件序列化失败")
	}
	return interfaces.AIEvent{
		ID:      int32(majong.EventID_event_qi_request),
		Context: data,
	}
}

func (h *waitQiangganghuStateAI) hu(player *majong.Player) interfaces.AIEvent {
	eventContext := majong.HuRequestEvent{
		Head: &majong.RequestEventHead{
			PlayerId: player.GetPalyerId(),
		},
	}
	data, err := proto.Marshal(&eventContext)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"func_name": "waitQiangganghuStateAI.hu",
			"player_id": player.GetPalyerId(),
		}).Errorln("事件序列化失败")
	}
	return interfaces.AIEvent{
		ID:      int32(majong.EventID_event_hu_request),
		Context: data,
	}
}

func (h *waitQiangganghuStateAI) checkAIEvent(player *majong.Player, mjContext *majong.MajongContext, params interfaces.AIEventGenerateParams) error {
	err := fmt.Errorf("不生成自动事件")
	if mjContext.GetCurState() != majong.StateID_state_waitqiangganghu ||
		player.GetPalyerId() == mjContext.GetLastGangPlayer() ||
		len(player.GetHandCards())%3+1 != 2 ||
		gutils.CheckHasDingQueCard(mjContext, player) ||
		len(player.GetPossibleActions()) == 0 ||
		params.AIType == interfaces.TingAI {
		return err
	}
	return nil
}
