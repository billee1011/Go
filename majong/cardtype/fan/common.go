package fan

import (
	"steve/majong/interfaces"
	"steve/majong/utils"
	majongpb "steve/server_pb/majong"
)

// Fan 番型
type Fan struct {
	name      majongpb.CardType
	value     uint32
	Condition condition
}

type condition func(cardCalcParams interfaces.CardCalcParams) bool

// GetFanName 获取番型名字
func (f *Fan) GetFanName() majongpb.CardType {
	return f.name
}

// GetFanValue 获取番型倍数
func (f *Fan) GetFanValue() uint32 {
	return f.value
}

// checkPingHu 平胡-不包含其他番型
func checkPingHu(cardCalcParams interfaces.CardCalcParams) bool {
	if checkQingYiSe(cardCalcParams) {
		return false
	} else if checkQiDui(cardCalcParams) {
		return false
	} else if checkPengPengHu(cardCalcParams) {
		return false
	} else if checkJingGouDiao(cardCalcParams) {
		return false
	} else if checkShiBaLuoHan(cardCalcParams) {
		return false
	}
	return true
}

// checkQingYiSe 清一色-所有牌同一花色
func checkQingYiSe(cardCalcParams interfaces.CardCalcParams) bool {
	checkCards := getCheckCards(cardCalcParams.HandCard, cardCalcParams.HuCard)
	for _, pengCard := range cardCalcParams.PengCard {
		checkCards = append(checkCards, pengCard)
	}
	for _, gangCard := range cardCalcParams.GangCard {
		checkCards = append(checkCards, gangCard)
	}
	color := majongpb.CardColor(-1)
	for _, card := range checkCards {
		if color == -1 {
			color = card.Color
		} else if color != card.Color {
			return false
		}
	}
	return true
}

// checkQiDui 七对-由七个对子组成--不能有碰杠
func checkQiDui(cardCalcParams interfaces.CardCalcParams) bool {
	checkCards := getCheckCards(cardCalcParams.HandCard, cardCalcParams.HuCard)

	if len(checkCards) != 14 {
		return false
	}
	cardCount := make(map[int32]int)
	for _, card := range checkCards {
		cardValue, _ := utils.CardToInt(*card)
		cardCount[*cardValue] = cardCount[*cardValue] + 1
	}
	for _, v := range cardCount {
		if v%2 != 0 {
			return false
		}
	}
	return true
}

// checkQingQiDui 清七对-清一色+七对
func checkQingQiDui(cardCalcParams interfaces.CardCalcParams) bool {
	if checkQiDui(cardCalcParams) && checkQingYiSe(cardCalcParams) {
		return true
	}
	return false
}

// checkLongQiDui 龙七对-至少有一个根的七对
func checkLongQiDui(cardCalcParams interfaces.CardCalcParams) bool {
	if checkQiDui(cardCalcParams) && GetGenCount(cardCalcParams) > 0 {
		return true
	}
	return false
}

// checkQingLongQiDui 清龙七对-清一色+龙七对
func checkQingLongQiDui(cardCalcParams interfaces.CardCalcParams) bool {
	if checkLongQiDui(cardCalcParams) && checkQingYiSe(cardCalcParams) {
		return true
	}
	return false
}

// checkPengPengHu 对对(碰碰)胡-刻子或碰或杠，加將牌组成就是没有顺子
func checkPengPengHu(cardCalcParams interfaces.CardCalcParams) bool {
	checkCards := getCheckCards(cardCalcParams.HandCard, cardCalcParams.HuCard)
	//开牌，即碰杠这些
	openCardSum := len(cardCalcParams.PengCard) + len(cardCalcParams.GangCard)
	if openCardSum >= 4 {
		return true
	}
	// 手牌中重复3个的个数
	handCardSum := 0
	cardCount := make(map[int32]int)
	for _, card := range checkCards {
		cardValue, _ := utils.CardToInt(*card)
		cardCount[*cardValue] = cardCount[*cardValue] + 1
	}
	cards := []int32{}
	for cardPoint, v := range cardCount {
		if v == 4 {
			return false
		} else if v == 3 {
			handCardSum++
		} else if v == 1 {
			cards = append(cards, cardPoint)
		}
	}
	if openCardSum+handCardSum >= 4 && len(cards) == 0 {
		return true
	}
	return false
}

// checkQingPeng 清碰碰胡-清一色+碰碰胡
func checkQingPeng(cardCalcParams interfaces.CardCalcParams) bool {
	if checkPengPengHu(cardCalcParams) && checkQingYiSe(cardCalcParams) {
		return true
	}
	return false
}

// checkJingGouDiao 金钩钓-胡牌时手里只剩一张，并且单钓一这张，其他的牌都被杠或碰了,不计碰碰胡。
func checkJingGouDiao(cardCalcParams interfaces.CardCalcParams) bool {
	checkCards := getCheckCards(cardCalcParams.HandCard, cardCalcParams.HuCard)
	if len(checkCards) == 2 {
		if utils.CardEqual(checkCards[0], checkCards[1]) {
			if len(cardCalcParams.PengCard) != 0 {
				return true
			}
		}
	}
	return false
}

// checkQingJingGouDiao 清金钩钓-清一色+金钩钓
func checkQingJingGouDiao(cardCalcParams interfaces.CardCalcParams) bool {
	if checkJingGouDiao(cardCalcParams) && checkQingYiSe(cardCalcParams) {
		return true
	}
	return false
}

// checkShiBaLuoHan 十八罗汉-胡牌时手上只剩一张牌单吊，其他手牌形成四个杠，此时不计四根和碰碰胡。
func checkShiBaLuoHan(cardCalcParams interfaces.CardCalcParams) bool {
	l := len(cardCalcParams.GangCard)
	if l == 4 {
		return true
	}
	return false
}

//checkQingShiBaLuoHan 清十八罗汉-清一色+十八罗汉
func checkQingShiBaLuoHan(cardCalcParams interfaces.CardCalcParams) bool {
	if checkShiBaLuoHan(cardCalcParams) && checkQingYiSe(cardCalcParams) {
		return true
	}
	return false
}

// GetGenCount 获取玩家牌型根的数目
func GetGenCount(cardCalcParams interfaces.CardCalcParams) uint32 {
	var gCount uint32
	gangCards := cardCalcParams.GangCard
	pengCards := cardCalcParams.PengCard

	checkCards := getCheckCards(cardCalcParams.HandCard, cardCalcParams.HuCard)

	cardCount := make(map[int32]int)
	for _, card := range checkCards {
		cardValue, _ := utils.CardToInt(*card)
		cardCount[*cardValue] = cardCount[*cardValue] + 1
	}
	for card, sum := range cardCount {
		if sum >= 4 {
			gCount++
		} else if sum == 1 {
			for _, pengCard := range pengCards {
				cardValue, _ := utils.CardToInt(*pengCard)
				if *cardValue == card {
					gCount++
				}
			}
		}
	}
	gCount = gCount + uint32(len(gangCards))
	return gCount
}

// getCheckCards 获取校验的牌组
func getCheckCards(handCards []*majongpb.Card, huCard *majongpb.Card) []*majongpb.Card {
	checkCard := handCards
	if huCard != nil {
		checkCard = append(checkCard, huCard)
	}
	return checkCard
}
