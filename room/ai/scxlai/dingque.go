package scxlai

import (
	"steve/gutils"
	"steve/room/interfaces"
	"steve/room/interfaces/global"
	"steve/server_pb/majong"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
)

type dingqueStateAI struct {
	maxDingqueTime time.Duration // 最大定缺时间
}

// 注册 AI
func init() {
	g := global.GetDeskAutoEventGenerator()
	// 血流
	g.RegisterAI(scxlGameID, int32(majong.StateID_state_dingque), &dingqueStateAI{})
	g.RegisterAI(scxlGameID, int32(majong.StateID_state_huansanzhang), &huansanzhangStateAI{})
	g.RegisterAI(scxlGameID, int32(majong.StateID_state_chupaiwenxun), &chupaiWenxunStateAI{})
	g.RegisterAI(scxlGameID, int32(majong.StateID_state_zixun), &zixunStateAI{})
	g.RegisterAI(scxlGameID, int32(majong.StateID_state_waitqiangganghu), &waitQiangganghuStateAI{})
	// 血战
	g.RegisterAI(scxzGameID, int32(majong.StateID_state_dingque), &dingqueStateAI{})
	g.RegisterAI(scxzGameID, int32(majong.StateID_state_huansanzhang), &huansanzhangStateAI{})
	g.RegisterAI(scxzGameID, int32(majong.StateID_state_chupaiwenxun), &chupaiWenxunStateAI{})
	g.RegisterAI(scxzGameID, int32(majong.StateID_state_zixun), &zixunStateAI{})
	g.RegisterAI(scxzGameID, int32(majong.StateID_state_waitqiangganghu), &waitQiangganghuStateAI{})

}

// GenerateAIEvent 生成 AI 事件
// 无论是超时、托管还是机器人，都选最少的牌作为定缺牌， 并且产生相应的事件
func (h *dingqueStateAI) GenerateAIEvent(params interfaces.AIEventGenerateParams) (result interfaces.AIEventGenerateResult, err error) {
	result, err = interfaces.AIEventGenerateResult{
		Events: []interfaces.AIEvent{},
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
	cards := player.GetHandCards()
	colorMap := map[majong.CardColor]int{}

	for _, card := range cards {
		colorMap[card.GetColor()] = colorMap[card.GetColor()] + 1
	}
	colors := h.allColor()
	leastColor := colors[0]
	leastCount := colorMap[leastColor]
	for _, color := range h.allColor() {
		if colorMap[color] < leastCount {
			leastCount = colorMap[color]
			leastColor = color
		}
	}
	return leastColor
}

// dingque 生成定缺请求事件
func (h *dingqueStateAI) dingque(player *majong.Player) *interfaces.AIEvent {
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
	return &interfaces.AIEvent{
		ID:      int32(majong.EventID_event_dingque_request),
		Context: data,
	}
}
