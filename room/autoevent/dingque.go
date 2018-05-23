package autoevent

import (
	"steve/room/interfaces"
	"steve/server_pb/majong"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
)

type dingqueHandler struct {
	maxDingqueTime time.Duration // 最大定缺时间
}

func newDingqueHandler() *dingqueHandler {
	return &dingqueHandler{
		maxDingqueTime: 10 * time.Second,
	}
}

func (h *dingqueHandler) Generate(mjContext *majong.MajongContext, stateTime time.Time) (result []interfaces.Event) {
	duration := time.Now().Sub(stateTime)
	if duration < h.maxDingqueTime {
		return
	}
	players := mjContext.GetPlayers()
	for _, player := range players {
		if player.GetHasDingque() {
			continue
		}
		if event := h.dingque(player); event != nil {
			result = append(result, *event)
		}
	}
	return
}

// allColor 所有的麻将花色
func (h *dingqueHandler) allColor() []majong.CardColor {
	return []majong.CardColor{majong.CardColor_ColorWan, majong.CardColor_ColorTiao, majong.CardColor_ColorTong}
}

// getColor 获取定缺花色
func (h *dingqueHandler) getColor(player *majong.Player) majong.CardColor {
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
func (h *dingqueHandler) dingque(player *majong.Player) *interfaces.Event {
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
			"func_name": "dingqueHandler.dingque",
			"player_id": player.GetPalyerId(),
			"color":     color,
		}).Errorln("事件序列化失败")
		return nil
	}
	return &interfaces.Event{
		ID:      majong.EventID_event_dingque_request,
		Context: data,
	}
}
