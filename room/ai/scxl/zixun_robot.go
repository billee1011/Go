package scxlai

import (
	"github.com/Sirupsen/logrus"
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

	var cards []majong.Card
	for _, handCard := range handCards {
		cards = append(cards, *handCard)
	}

	// 拆牌
	var shunZis, keZis, pairs, doubleChas, singleChas, singles []Split
	shunZis1, keZis1, pairs1, doubleChas1, singleChas1, singles1 := SplitCards(cards, true)
	shunZis2, keZis2, pairs2, doubleChas2, singleChas2, singles2 := SplitCards(cards, false)
	if len(shunZis1)+len(keZis1) > len(shunZis2)+len(keZis2) {
		remain := RemoveSplits(cards, shunZis1)
		gangs := SplitGang(remain)
		if len(gangs) > 0 {
			h.gang(player, &gangs[0].cards[0])
			return
		}
		goto assign1
	} else if len(shunZis1)+len(keZis1) == len(shunZis2)+len(keZis2) {
		remain1 := RemoveSplits(cards, shunZis1)
		gangs := SplitGang(remain1)
		if len(gangs) > 0 {
			h.gang(player, &gangs[0].cards[0])
			return
		}

		remain2 := RemoveSplits(cards, shunZis1)
		gangs = SplitGang(remain2)
		if len(gangs) > 0 {
			h.gang(player, &gangs[0].cards[0])
			return
		}
		if len(pairs1) > len(pairs2) {
			goto assign1
		} else if len(pairs1) == len(pairs2) {
			if len(doubleChas1) > len(doubleChas2) {
				goto assign1
			} else if len(doubleChas1) == len(doubleChas2) {
				if len(singleChas1) > len(singleChas2) {
					goto assign1
				} else if len(singleChas1) == len(singleChas2) {
					goto assign2
				} else {
					goto assign2
				}
			} else {
				goto assign2
			}
		} else {
			goto assign2
		}
	} else {
		remain := RemoveSplits(cards, shunZis2)
		gangs := SplitGang(remain)
		if len(gangs) > 0 {
			h.gang(player, &gangs[0].cards[0])
			return
		}
		goto assign2
	}
assign1:
	shunZis = shunZis1
	keZis = keZis1
	pairs = pairs1
	doubleChas = doubleChas1
	singleChas = singleChas1
	singles = singles1
	goto analysis
assign2:
	shunZis = shunZis2
	keZis = keZis2
	pairs = pairs2
	doubleChas = doubleChas2
	singleChas = singleChas2
	singles = singles2
	goto analysis
analysis:
	logrus.WithFields(logrus.Fields{"手牌": cards, "顺子": shunZis, "刻子": keZis, "对子": pairs, "双茬": doubleChas, "单茬": singleChas, "单牌": singles}).Infoln("中级AI拆牌结果")

	if len(singles) == 1 {
		h.chupai(player, &singles[0].cards[0])
		return
	}

	//var wallCards []majong.Card
	//for _, wallCard := range mjContext.WallCards {
	//	wallCards = append(wallCards, *wallCard)
	//}
	//
	//remainCards := CountCard(wallCards)

	var visibleCards []*majong.Card
	visibleCards = append(visibleCards, player.HandCards...)
	for _, player := range mjContext.Players {
		visibleCards = append(visibleCards, player.OutCards...)
	}

	countMap := make(map[majong.Card]int)
	for _, visuableCard := range visibleCards {
		countMap[*visuableCard]++
	}

	var remainCards map[majong.Card]int
	for k, v := range countMap {
		remainCards[k] = 4 - v
	}

	if len(singles) > 1 {
		min := 99
		var outCard majong.Card
		if len(pairs) >= 1 { //有将，比较成茬机会数
			for _, single := range singles {
				validCards := getValidCard(single)
				chance := countValidCard(remainCards, validCards)
				if chance < min {
					min = chance
					outCard = single.cards[0]
				}
			}

		} else { //无将，比较成将机会数
			for _, single := range singles {
				chance := remainCards[single.cards[0]]
				if chance < min {
					min = chance
					outCard = single.cards[0]
				}
			}
		}
		h.chupai(player, &outCard)
		return
	}

	var twoCards []Split
	for _, singleCha := range singleChas {
		twoCards = append(twoCards, singleCha)
	}
	for _, doubleCha := range doubleChas {
		twoCards = append(twoCards, doubleCha)
	}
	if len(pairs) > 1 { //一个对子留作将，多于一个对子看成牌机会
		for _, pair := range pairs {
			twoCards = append(twoCards, pair)
		}
	}

	var chances map[Split]int
	for _, twoCard := range twoCards {
		chances[twoCard] = countValidCard(remainCards, getValidCard(twoCard))
	}

	var needChai Split
	minChance := 99
	for split, chance := range chances {
		if chance < minChance {
			minChance = chance
			needChai = split
		} else if chance == minChance && split.t > needChai.t {
			needChai = split
		}
	}

	h.chupai(player, &needChai.cards[0])
	return
}

func countValidCard(remainCards map[majong.Card]int, validCards []majong.Card) int {
	total := 0
	for _, validCard := range validCards {
		total += remainCards[validCard]
	}
	return total
}

func getValidCard(split Split) (result []majong.Card) {
	if split.t == SINGLE { //单牌成茬有效牌
		single := split.cards[0]
		if single.Color == majong.CardColor_ColorHua || single.Color == majong.CardColor_ColorZi {
			return
		}
		for _, addend := range []int32{-2, -1, 1, 2} {
			if single.Point+addend >= 1 && single.Point+addend <= 9 {
				result = append(result, majong.Card{Color: single.Color, Point: single.Point + addend})
			}
		}
	}
	if split.t == PAIR { //对子成刻有效牌
		result = append(result, split.cards[0])
	}
	if split.t == DOUBLE_CHA {
		small := split.cards[0]
		result = append(result, majong.Card{Color: small.Color, Point: small.Point - 1})
		result = append(result, majong.Card{Color: small.Color, Point: small.Point + 2})
	}
	if split.t == SINGLE_CHA {
		small := split.cards[0]
		if ContainsEdge(split.cards) { // 12 89
			if small.Point == 1 {
				result = append(result, majong.Card{Color: small.Color, Point: 3})
			} else {
				result = append(result, majong.Card{Color: small.Color, Point: 7})
			}
		} else { // 13 24 35 ... 79
			result = append(result, majong.Card{Color: small.Color, Point: small.Point + 1})
		}
	}
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

// 拆出所有杠
func SplitGang(handCards []majong.Card) (result []Split) {
	gangs := FindAllShunZi(handCards, 4, 1)
	for _, gang := range gangs {
		result = append(result, Split{GANG, gang})
	}
	return
}

// 拆出所有刻子
func SplitKeZi(handCards []majong.Card) (result []Split) {
	keZis := FindAllShunZi(handCards, 3, 1)
	for _, keZi := range keZis {
		result = append(result, Split{KEZI, keZi})
	}
	return
}

// 拆出所有顺子
func SplitShunZi(handCards []majong.Card) (result []Split) {
	shunZis := FindAllShunZi(handCards, 1, 3)
	for _, shunZi := range shunZis {
		result = append(result, Split{SHUNZI, shunZi})
	}
	return
}

// 拆出所有对子
func SplitPair(handCards []majong.Card) (result []Split) {
	pairs := FindAllShunZi(handCards, 2, 1)
	for _, pair := range pairs {
		result = append(result, Split{PAIR, pair})
	}
	return
}

// 拆出所有双茬
func SplitDoubleCha(cards []majong.Card) []Split {
	doubleCha, _ := getNearShunZi(cards)
	return doubleCha
}

func getNearShunZi(handCards []majong.Card) (doubleCha []Split, singleCha []Split) {
	result := FindAllShunZi(handCards, 1, 2)

	for _, split := range result {
		if ContainsEdge(split) {
			singleCha = append(singleCha, Split{SINGLE_CHA, split})
		} else {
			doubleCha = append(doubleCha, Split{DOUBLE_CHA, split})
		}
	}
	return
}

// 拆出所有单茬
func SplitSingleCha(cards []majong.Card) []Split {
	_, singleCha := getNearShunZi(cards)
	remain := RemoveSplits(cards, singleCha)
	singleCha = append(singleCha, getSpaceShunZi(remain)...)
	return singleCha
}

func getSpaceShunZi(handCards []majong.Card) (result []Split) {
	spaceShunZis := FindAllCommonShunZi(handCards, 1, 2, 2)
	for _, spaceShunZi := range spaceShunZis {
		result = append(result, Split{SINGLE_CHA, spaceShunZi})
	}
	return
}

// 拆成单牌
func SplitSingle(cards []majong.Card) []Split {
	var result []Split
	for _, card := range cards {
		result = append(result, Split{t: SINGLE, cards: []majong.Card{card}})
	}
	return result
}

/**
FindAllShunZi 双向夹击，找出手牌中所有顺子长度为shunZiLen，重复次数为duplicateCount的牌
*/
func FindAllShunZi(handCards []majong.Card, duplicateCount int, shunZiLen int) (result [][]majong.Card) {
	return FindAllCommonShunZi(handCards, duplicateCount, shunZiLen, 1)
}
func FindAllCommonShunZi(handCards []majong.Card, duplicateCount int, shunZiLen int, shunZiGap int) (result [][]majong.Card) {
	countMap := CountCard(handCards)
	var matchCards []majong.Card
	for card, count := range countMap {
		if count >= duplicateCount {
			matchCards = append(matchCards, card)
		}
	}
	MJCardSort(matchCards)

	gap := shunZiLen - 1

	colorCards := divideByColor(matchCards)
	for color, cards := range colorCards {
		if color == majong.CardColor_ColorHua || color == majong.CardColor_ColorZi && shunZiLen != 1 {
			continue //花牌都按单牌处理，字牌没有顺子
		}
		i := 0
		j := len(cards) - 1
		for {
			if i+gap <= len(cards)-1 && cards[i+gap].Color == cards[i].Color && cards[i+gap].Point-cards[i].Point == int32(gap*shunZiGap) { //从1向9取
				shunZi := cards[i : i+gap+1]
				inflated := InflateAll(shunZi, duplicateCount)
				result = append(result, inflated)
				decreaseAll(countMap, shunZi, duplicateCount)
				if existAll(countMap, cards, i, i+gap, duplicateCount) {
					continue //重复取
				} else {
					i += shunZiLen
				}
			} else {
				i++
			}

			if j-gap >= 0 && i+gap <= j-gap && cards[j-gap].Color == cards[j].Color && cards[j].Point-cards[j-gap].Point == int32(gap*shunZiGap) { //从9向1取
				shunZi := cards[j-gap : j+1]
				inflated := InflateAll(shunZi, duplicateCount)
				result = append(result, inflated)
				decreaseAll(countMap, shunZi, duplicateCount)
				if existAll(countMap, cards, j-gap, j, duplicateCount) {
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
	}
	return
}

// 按万、条、筒、字拆分手牌
func divideByColor(cards []majong.Card) map[majong.CardColor][]majong.Card {
	colors := make(map[majong.CardColor][]majong.Card)
	for _, card := range cards {
		colors[card.Color] = append(colors[card.Color], card)
	}
	return colors
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
