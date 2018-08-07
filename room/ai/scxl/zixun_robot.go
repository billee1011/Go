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
		shunZis1 := SplitShunZi(colorCards)
		remain1 := RemoveSplits(colorCards, shunZis1)
		keZis1 := SplitKeZi(remain1)
		remain1 = RemoveSplits(remain1, keZis1)

		keZis2 := SplitKeZi(colorCards)
		remain2 := RemoveSplits(colorCards, keZis2)
		shunZis2 := SplitShunZi(remain2)
		remain2 = RemoveSplits(remain2, shunZis2)

		if len(shunZis1)+len(keZis1) == len(shunZis2)+len(keZis2) {
			noShunZi1 := RemoveSplits(colorCards, shunZis1)
			gang1 := SplitGang(noShunZi1)
			if len(gang1) > 0 {
				h.gang(player, &gang1[0].cards[0])
				return
			}

			noShunZi2 := RemoveSplits(colorCards, shunZis2)
			gang2 := SplitGang(noShunZi2)
			if len(gang2) > 0 {
				h.gang(player, &gang2[0].cards[0])
				return
			}
		} else if len(shunZis1)+len(keZis1) > len(shunZis2)+len(keZis2) {
			noShunZi1 := RemoveSplits(colorCards, shunZis1)
			gang1 := SplitGang(noShunZi1)
			if len(gang1) > 0 {
				h.gang(player, &gang1[0].cards[0])
				return
			}
		} else if len(shunZis1)+len(keZis1) > len(shunZis2)+len(keZis2) {
			noShunZi2 := RemoveSplits(colorCards, shunZis2)
			gang2 := SplitGang(noShunZi2)
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
		shunZis = SplitShunZi(remain)
		remain = RemoveSplits(remain, shunZis)
		keZis = SplitKeZi(remain)
		remain = RemoveSplits(remain, keZis)
	} else {
		keZis = SplitKeZi(remain)
		remain = RemoveSplits(remain, keZis)
		shunZis = SplitShunZi(remain)
		remain = RemoveSplits(remain, shunZis)
	}
	pairs = SplitPair(remain)
	remain = RemoveSplits(remain, pairs)
	doubleChas = SplitDoubleCha(remain)
	remain = RemoveSplits(remain, doubleChas)
	singleChas = SplitSingleCha(remain)
	remain = RemoveSplits(remain, singleChas)
	singles = SplitSingle(remain)
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

// 拆出所有顺子
func SplitShunZi(handCards []majong.Card) (result []Split) {
	shunZis := FindAllCards(handCards, 1, 3)
	for _, shunZi := range shunZis {
		result = append(result, Split{SHUNZI, shunZi})
	}
	return
}

// 拆出所有刻子
func SplitKeZi(handCards []majong.Card) (result []Split) {
	keZis := FindAllCards(handCards, 3, 1)
	for _, keZi := range keZis {
		result = append(result, Split{KEZI, keZi})
	}
	return
}

func SplitDoubleCha(cards []majong.Card) []Split {
	doubleCha, _ := getNearShunZi(cards)
	return doubleCha
}

func SplitSingleCha(cards []majong.Card) []Split {
	_, singleCha := getNearShunZi(cards)
	remain := RemoveSplits(cards, singleCha)
	singleCha = append(singleCha, getSpaceShunZi(remain)...)
	return singleCha
}

func SplitSingle(cards []majong.Card) []Split {
	var result []Split
	for _, card := range cards {
		result = append(result, Split{t: SINGLE, cards: []majong.Card{card}})
	}
	return result
}

func getNearShunZi(handCards []majong.Card) (doubleCha []Split, singleCha []Split) {
	result := FindAllCards(handCards, 1, 2)

	for _, split := range result {
		if ContainsEdge(split) {
			singleCha = append(singleCha, Split{SINGLE_CHA, split})
		} else {
			doubleCha = append(doubleCha, Split{DOUBLE_CHA, split})
		}
	}
	return
}

func SplitPair(handCards []majong.Card) (result []Split) {
	pairs := FindAllCards(handCards, 2, 1)
	for _, pair := range pairs {
		result = append(result, Split{PAIR, pair})
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

func SplitGang(handCards []majong.Card) (result []Split) {
	gangs := FindAllCards(handCards, 4, 1)
	for _, gang := range gangs {
		result = append(result, Split{GANG, gang})
	}
	return
}

/**
FindAllCards 双向夹击，找出手牌中所有顺子长度为shunZiLen，重复次数为duplicateCount的牌
*/
func FindAllCards(handCards []majong.Card, duplicateCount int, shunZiLen int) (result [][]majong.Card) {
	countMap := CountCard(handCards)
	var matchCards []majong.Card
	for card, count := range countMap {
		if count >= duplicateCount {
			matchCards = append(matchCards, card)
		}
	}
	MJCardSort(matchCards)

	gap := shunZiLen - 1

	i := 0
	j := len(matchCards) - 1
	for {
		if i+gap <= len(matchCards)-1 && matchCards[i+gap].Color == matchCards[i].Color && matchCards[i+gap].Point-matchCards[i].Point == int32(gap) { //从1向9取
			shunZi := matchCards[i : i+gap+1]
			inflated := InflateAll(shunZi, duplicateCount)
			result = append(result, inflated)
			decreaseAll(countMap, shunZi, duplicateCount)
			if existAll(countMap, matchCards, i, i+gap, duplicateCount) {
				continue //重复取
			} else {
				i += shunZiLen
			}
		} else {
			i++
		}

		if j-gap >= 0 && i+gap < j-gap && matchCards[j-gap].Color == matchCards[j].Color && matchCards[j].Point-matchCards[j-gap].Point == int32(gap) { //从9向1取
			shunZi := matchCards[j-gap : j+1]
			inflated := InflateAll(shunZi, duplicateCount)
			result = append(result, inflated)
			decreaseAll(countMap, shunZi, duplicateCount)
			if existAll(countMap, matchCards, j-gap, j, duplicateCount) {
				continue //重复取
			} else {
				j -= shunZiLen
			}
		} else {
			j--
		}
		if i > j {
			break
		}
	}
	return
}

func decreaseAll(countMap map[majong.Card]int, shunZi []majong.Card, duplicateCount int) {
	for _, card := range shunZi {
		countMap[card] -= duplicateCount
	}
}

func existAll(countMap map[majong.Card]int, matchCards []majong.Card, start int, end int, duplicateCount int) bool {
	for i := start; i <= end; i++ {
		card := matchCards[i]
		if countMap[card] < duplicateCount {
			return false
		}
	}
	return true
}

func InflateAll(cards []majong.Card, duplicateCount int) (result []majong.Card) {
	for _, card := range cards {
		for i := 0; i < duplicateCount; i++ {
			result = append(result, card)
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
