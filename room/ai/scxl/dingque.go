package scxlai

import (
	"steve/gutils"
	"steve/server_pb/majong"
	"time"

	"steve/room/ai"

	"github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
)

type dingqueStateAI struct {
	maxDingqueTime time.Duration // 最大定缺时间
}

// GenerateAIEvent 生成 AI 事件
// 无论是超时、托管还是机器人，都选最少的牌作为定缺牌， 并且产生相应的事件
func (h *dingqueStateAI) GenerateAIEvent(params ai.AIEventGenerateParams) (result ai.AIEventGenerateResult, err error) {
	result, err = ai.AIEventGenerateResult{
		Events: []ai.AIEvent{},
	}, nil

	mjContext := params.MajongContext
	player := gutils.GetMajongPlayer(params.PlayerID, mjContext)
	if player.GetHasDingque() {
		return
	}
	if event := h.dingque(player); event != nil {
		result.Events = append(result.Events, *event)
	}
	return
}

// allColor 所有的麻将花色
func (h *dingqueStateAI) allColor() []majong.CardColor {
	return []majong.CardColor{majong.CardColor_ColorWan, majong.CardColor_ColorTiao, majong.CardColor_ColorTong}
}

// getColor 获取定缺花色
func (h *dingqueStateAI) getColor(player *majong.Player) majong.CardColor {
	return player.GetDingqueColor() // 在进入定缺状态时，会设置推荐定缺颜色
}

// dingque 生成定缺请求事件
func (h *dingqueStateAI) dingque(player *majong.Player) *ai.AIEvent {
	color := h.getColor(player)
	mjContext := majong.DingqueRequestEvent{
		Head: &majong.RequestEventHead{
			PlayerId: player.GetPalyerId(),
		},
		Color: color,
	}

	data, err := proto.Marshal(&mjContext)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"func_name": "dingqueStateAI.dingque",
			"player_id": player.GetPalyerId(),
			"color":     color,
		}).Errorln("事件序列化失败")
		return nil
	}
	return &ai.AIEvent{
		ID:      int32(majong.EventID_event_dingque_request),
		Context: data,
	}
}
