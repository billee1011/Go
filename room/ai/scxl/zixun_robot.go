package scxlai

import (
	"sort"
	"steve/common/mjoption"
	"steve/entity/majong"
	"steve/gutils"
	"steve/room/ai"
)

func (h *zixunStateAI) getRobotAIEvent(player *majong.Player, mjContext *majong.MajongContext) (aiEvent ai.AIEvent) {
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
	if player.GetMopaiCount() == 0 {
		aiEvent = h.chupai(player, handCards[len(handCards)-1])
	} else {
		aiEvent = h.chupai(player, mjContext.GetLastMopaiCard())
	}
	return
}

func getOutCard(cards []*majong.Card) {
	colors := splitCard(cards)
	for _, colorCards := range colors {
		shunzis1 := getShunZi(colorCards)
		remain1 := RemoveAll(colorCards, shunzis1)
		kezis1 := getKeZi(remain1)
		remain1 = RemoveAll(remain1, kezis1)

		kezis2 := getKeZi(colorCards)
		remain2 := RemoveAll(colorCards, kezis2)
		shunzis2 := getShunZi(remain2)
		remain2 = RemoveAll(remain2, shunzis2)

	}
}

func RemoveAll(cards []majong.Card, splits []Split) []majong.Card {
	var result []majong.Card
	for _, card := range cards {
		if !Contains(splits, card) {
			result = append(result, card)
		}
	}
	return result
}

func Contains(splits []Split, inCard majong.Card) bool {
	for _, split := range splits {
		for _, card := range split.cards {
			if card == inCard {
				return true
			}
		}
	}
	return false
}

// 按万、条、筒、字拆分手牌
func splitCard(cards []*majong.Card) map[majong.CardColor][]majong.Card {
	colors := make(map[majong.CardColor][]majong.Card)
	for _, card := range cards {
		colors[card.Color] = append(colors[card.Color], *card)
	}
	return colors
}

type MJCardSlice []majong.Card

func (cs MJCardSlice) Len() int      { return len(cs) }
func (cs MJCardSlice) Swap(i, j int) { cs[i], cs[j] = cs[j], cs[i] }
func (cs MJCardSlice) Less(i, j int) bool {
	return gutils.ServerCard2Number(&cs[i]) < gutils.ServerCard2Number(&cs[j])
}

func MJCardSort(cards []majong.Card) {
	cs := MJCardSlice(cards)
	sort.Sort(cs)
}

type SplitType int

const (
	GANG       SplitType = iota //杠，四张相同的牌，已成牌
	KEZI                        //刻子，三张相同的牌，已成牌
	SHUNZI                      //顺子，如567，已成牌
	PAIR                        //对子，如55，一步成刻
	DOUBLE_CHA                  //双茬，如56，一步成顺
	SINGLE_CHA                  //单茬，如57，89，一步成顺
	SINGLE                      //单牌，如5，两步成牌
)

type Split struct {
	t     SplitType
	cards []majong.Card
}

// 相同Color的牌获取顺子
func getShunZi(cards []majong.Card) (result []Split) {
	return getTriple(cards, SHUNZI)
}

// 相同Color的牌获取刻子
func getKeZi(cards []majong.Card) (result []Split) {
	return getTriple(cards, KEZI)
}

//两边夹击获取顺子
func getTriple(cards []majong.Card, splitType SplitType) (result []Split) {
	if len(cards) < 3 {
		return
	}
	MJCardSort(cards)

	var diff int32
	if splitType == SHUNZI {
		diff = 1
	} else {
		diff = 0
	}
	i := 0
	j := len(cards) - 1
	for {
		if i+2 < len(cards)-1 && cards[i+2].Point-cards[i+1].Point == diff && cards[i+1].Point-cards[i].Point == diff {
			result = append(result, Split{splitType, []majong.Card{cards[i], cards[i+1], cards[i+2]}})
			i += 3
		} else {
			i++
		}

		if j-2 >= 0 && i+2 < j-2 && cards[j].Point-cards[j-1].Point == diff && cards[j-1].Point-cards[j-2].Point == diff {
			result = append(result, Split{splitType, []majong.Card{cards[j-2], cards[j-1], cards[j]}})
			j -= 3
		} else {
			j--
		}
		if i >= j {
			break
		}
	}
	return
}

func getDoubleCha() {

}

func getDouble(cards []majong.Card, splitType SplitType) {

}
