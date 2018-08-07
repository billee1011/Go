package scxlai

import (
	"steve/common/mjoption"
	"steve/entity/majong"
	"steve/gutils"
	"steve/room/ai"
)

func (h *zixunStateAI) getMiddleAIEvent(player *majong.Player, mjContext *majong.MajongContext) (aiEvent ai.AIEvent) {
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
	colors := divideByColor(handCards)
	for _, colorCards := range colors {
		shunZis1 := getShunZi(colorCards)
		remain1 := RemoveSplits(colorCards, shunZis1)
		keZis1 := getKeZi(remain1)
		remain1 = RemoveSplits(remain1, keZis1)

		keZis2 := getKeZi(colorCards)
		remain2 := RemoveSplits(colorCards, keZis2)
		shunZis2 := getShunZi(remain2)
		remain2 = RemoveSplits(remain2, shunZis2)

		if len(shunZis1)+len(keZis1) == len(shunZis2)+len(keZis2) {
			noShunZi1 := RemoveSplits(colorCards, shunZis1)
			gang1 := getGang(noShunZi1)
			if len(gang1) > 0 {
				h.gang(player, &gang1[0].cards[0])
				return
			}

			noShunZi2 := RemoveSplits(colorCards, shunZis2)
			gang2 := getGang(noShunZi2)
			if len(gang2) > 0 {
				h.gang(player, &gang2[0].cards[0])
				return
			}
		} else if len(shunZis1)+len(keZis1) > len(shunZis2)+len(keZis2) {
			noShunZi1 := RemoveSplits(colorCards, shunZis1)
			gang1 := getGang(noShunZi1)
			if len(gang1) > 0 {
				h.gang(player, &gang1[0].cards[0])
				return
			}
		} else if len(shunZis1)+len(keZis1) > len(shunZis2)+len(keZis2) {
			noShunZi2 := RemoveSplits(colorCards, shunZis2)
			gang2 := getGang(noShunZi2)
			if len(gang2) > 0 {
				h.gang(player, &gang2[0].cards[0])
				return
			}
		}

	}

	//if player.GetMopaiCount() == 0 {
	//	aiEvent = h.chupai(player, handCards[len(handCards)-1])
	//} else {
	//	aiEvent = h.chupai(player, mjContext.GetLastMopaiCard())
	//}
	return
}

func SplitCards(cards []majong.Card, shunZiFirst bool) (shunZis []Split, keZis []Split, pairs []Split, doubleChas []Split, singleChas []Split, singles []Split) {
	remain := cards
	if shunZiFirst {
		shunZis = getShunZi(remain)
		remain = RemoveSplits(remain, shunZis)
		keZis = getKeZi(remain)
		remain = RemoveSplits(remain, keZis)
	} else {
		keZis = getKeZi(remain)
		remain := RemoveSplits(remain, keZis)
		shunZis = getShunZi(remain)
		remain = RemoveSplits(remain, shunZis)
	}
	pairs = getPair(remain)
	remain = RemoveSplits(remain, pairs)
	doubleChas = getDoubleCha(remain)
	remain = RemoveSplits(remain, doubleChas)
	singleChas = getSingleCha(remain)
	remain = RemoveSplits(remain, singleChas)
	singles = getSingle(remain)
	return
}

// 按万、条、筒、字拆分手牌
func divideByColor(cards []*majong.Card) map[majong.CardColor][]majong.Card {
	colors := make(map[majong.CardColor][]majong.Card)
	for _, card := range cards {
		colors[card.Color] = append(colors[card.Color], *card)
	}
	return colors
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
		if i+2 <= len(cards)-1 && cards[i+2].Point-cards[i+1].Point == diff && cards[i+1].Point-cards[i].Point == diff {
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

func getDoubleCha(cards []majong.Card) []Split {
	doubleCha, _ := getNearShunZi(cards)
	return doubleCha
}

func getSingleCha(cards []majong.Card) []Split {
	_, singleCha := getNearShunZi(cards)
	remain := RemoveSplits(cards, singleCha)
	singleCha = append(singleCha, getSpaceShunZi(remain)...)
	return singleCha
}

func getSingle(cards []majong.Card) []Split {
	var result []Split
	for _, card := range cards {
		result = append(result, Split{t: SINGLE, cards: []majong.Card{card}})
	}
	return result
}

func getNearShunZi(cards []majong.Card) (doubleCha []Split, singleCha []Split) {
	result := getDouble(cards, DOUBLE_CHA)
	for _, split := range result {
		if ContainsEdge(split) {
			split.t = SINGLE_CHA
		}
	}
	for _, split := range result {
		if split.t == DOUBLE_CHA {
			doubleCha = append(doubleCha, split)
		}
		if split.t == SINGLE_CHA {
			singleCha = append(singleCha, split)
		}
	}
	return
}

func getPair(cards []majong.Card) (result []Split) {
	return getDouble(cards, PAIR)
}

func getDouble(cards []majong.Card, splitType SplitType) (result []Split) {
	if len(cards) < 2 {
		return
	}
	MJCardSort(cards)

	var diff int32
	if splitType == DOUBLE_CHA {
		diff = 1
	} else {
		diff = 0
	}
	i := 0
	j := len(cards) - 1
	for {
		if i+1 <= len(cards)-1 && cards[i+1].Point-cards[i].Point == diff {
			result = append(result, Split{splitType, []majong.Card{cards[i], cards[i+1]}})
			i += 2
		} else {
			i++
		}

		if j-1 >= 0 && i+1 < j-1 && cards[j].Point-cards[j-1].Point == diff {
			result = append(result, Split{splitType, []majong.Card{cards[j-1], cards[j]}})
			j -= 2
		} else {
			j--
		}
		if i >= j {
			break
		}
	}
	return
}

func getSpaceShunZi(cards []majong.Card) (result []Split) {
	if len(cards) < 2 {
		return
	}
	MJCardSort(cards)

	i := 0
	j := len(cards) - 1
	for {
		if i+1 < len(cards)-1 && cards[i+1].Point-cards[i].Point == 2 {
			result = append(result, Split{SINGLE_CHA, []majong.Card{cards[i], cards[i+1]}})
			i += 2
		} else {
			i++
		}

		if j-1 >= 0 && i+1 < j-1 && cards[j].Point-cards[j-1].Point == 2 {
			result = append(result, Split{SINGLE_CHA, []majong.Card{cards[j-1], cards[j]}})
			j -= 2
		} else {
			j--
		}
		if i >= j {
			break
		}
	}
	return
}

func getGang(cards []majong.Card) (result []Split) {
	countMap := make(map[majong.Card]int)
	for _, card := range cards {
		countMap[card]++
	}

	for card, count := range countMap {
		if count >= 4 {
			result = append(result, Split{GANG, []majong.Card{card, card, card, card}})
			break
		}
	}
	return
}

func CountCard(cards []majong.Card) map[majong.Card]int {
	countMap := make(map[majong.Card]int)
	for _, card := range cards {
		countMap[card]++
	}
	return countMap
}

func FindCards(cards []majong.Card, duplicateCount int, shunZiLen int) []majong.Card {
	countMap := CountCard(cards)
	var matchCards []majong.Card
	for card, count := range countMap {
		if count >= duplicateCount {
			matchCards = append(matchCards, card)
		}
	}
	MJCardSort(matchCards)

	gap := shunZiLen - 1
	for i, _ := range matchCards {
		if i >= gap && matchCards[i].Color == matchCards[i-gap].Color && matchCards[i].Point-matchCards[i-gap].Point == int32(gap) {
			shunZi := matchCards[i-gap : i+1]
			var result []majong.Card
			for _, card := range shunZi {
				inflated := Inflate(card, duplicateCount)
				result = append(result, inflated...)
			}
			return result
		}
	}
	return nil
}

func Inflate(card majong.Card, duplicateCount int) (result []majong.Card) {
	for i := 0; i < duplicateCount; i++ {
		result = append(result, card)
	}
	return
}
