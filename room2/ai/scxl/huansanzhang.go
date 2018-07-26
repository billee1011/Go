package scxlai

import (
	"steve/entity/majong"
	"steve/gutils"
<<<<<<< HEAD

	"steve/room2/ai"

	"github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
=======

	"steve/room2/ai"
>>>>>>> majong
)

type huansanzhangStateAI struct {
}

// GenerateAIEvent 生成 换三张AI 事件
// 无论是超时、托管还是机器人，若已存在换三张的牌，则直接换该三张牌，否则取花色最少的三张手牌换三张， 并且产生相应的事件
func (h *huansanzhangStateAI) GenerateAIEvent(params ai.AIEventGenerateParams) (result ai.AIEventGenerateResult, err error) {
	result, err = ai.AIEventGenerateResult{
		Events: []ai.AIEvent{},
	}, nil

	mjContext := params.MajongContext
	player := gutils.GetMajongPlayer(params.PlayerID, mjContext)
	if player.GetHuansanzhangSure() {
		return
	}
	if event := h.huansanzhang(player); event != nil {
		result.Events = append(result.Events, *event)
	}
	return
}

// getHszCards 获取换三张的牌
func (h *huansanzhangStateAI) getHszCards(player *majong.Player) (hszCards []*majong.Card) {
	if len(player.HuansanzhangCards) == 3 {
		hszCards = player.HuansanzhangCards
		return hszCards
	}
	// 随即获取最小花色的三张牌
	cards := player.GetHandCards()
	colorMap := map[majong.CardColor]int{}
	for _, card := range cards {
		colorMap[card.GetColor()] = colorMap[card.GetColor()] + 1
	}

	leastColor := majong.CardColor(-1)
	leastCount := 0
	hszCards = make([]*majong.Card, 0, 3)
	for _, color := range new(dingqueStateAI).allColor() {
		if colorMap[color] >= 3 {
			if leastColor == -1 || colorMap[color] < leastCount {
				leastColor = color
				leastCount = colorMap[color]
			}
		}
	}
	for _, card := range cards {
		if len(hszCards) == 3 {
			break
		}
		if card.GetColor() == leastColor {
			hszCards = append(hszCards, card)
		}
	}
	return hszCards
}

// huansanzhang 生成换三张请求事件
func (h *huansanzhangStateAI) huansanzhang(player *majong.Player) *ai.AIEvent {
	hszCards := h.getHszCards(player)

	eventContext := &majong.HuansanzhangRequestEvent{
		Head: &majong.RequestEventHead{
			PlayerId: player.GetPalyerId(),
		},
		Cards: hszCards,
		Sure:  true,
	}

	return &ai.AIEvent{
		ID:      majong.EventID_event_huansanzhang_request,
		Context: eventContext,
	}
}
