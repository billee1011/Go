package fantype

import (
	"sort"
	"steve/gutils"
	majongpb "steve/server_pb/majong"
)

type intSlice []int32

func (b intSlice) Len() int {
	return len(b)
}

func (b intSlice) Less(i, j int) bool {
	return b[i] < b[j]
}

func (b intSlice) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

// uint card sort 排序
func cardSort(cardInt []int32) {
	// 从小到大排序
	sort.Sort(intSlice(cardInt))
}

//sortRemoveDuplicate 排序加去重
func sortRemoveDuplicate(card []int32) (ret []int32) {
	cardSort(card)
	for i := 0; i < len(card); i++ {
		if (i > 0 && card[i-1] == card[i]) || card[i] == 0 {
			continue
		}
		ret = append(ret, card[i])
	}
	return
}

// 吃转牌
func chiToCards(chis []*majongpb.ChiCard) []*majongpb.Card {
	cards := make([]*majongpb.Card, 0)
	for _, chi := range chis {
		// min
		cardMin := chi.GetCard()
		cards = append(cards, cardMin)
		//mid
		cardmid := &majongpb.Card{
			Color: cardMin.GetColor(),
			Point: cardMin.GetPoint() + 1,
		}
		cards = append(cards, cardmid)
		// max
		cardMax := &majongpb.Card{
			Color: cardMin.GetColor(),
			Point: cardMin.GetPoint() + 2,
		}
		cards = append(cards, cardMax)
	}
	return cards
}

// 碰转牌
func pengToCards(pengs []*majongpb.PengCard) []*majongpb.Card {
	cards := make([]*majongpb.Card, 0)
	for _, peng := range pengs {
		cards = append(cards, peng.GetCard(), peng.GetCard(), peng.GetCard())
	}
	return cards
}

// 杠转牌
func gangToCards(gangs []*majongpb.GangCard) []*majongpb.Card {
	cards := make([]*majongpb.Card, 0)
	for _, gang := range gangs {
		cards = append(cards, gang.GetCard(), gang.GetCard(), gang.GetCard(), gang.GetCard())
	}
	return cards
}

// getPlayerCardAll 获取玩家所有牌,手，胡,碰，杠，吃牌
func getPlayerCardAll(tc *typeCalculator) []*majongpb.Card {
	// 所有牌
	cardAll := make([]*majongpb.Card, 0, len(tc.getHandCards()))
	// 手
	cardAll = append(cardAll, tc.getHandCards()...)
	// 胡牌
	cardAll = append(cardAll, tc.getHuCard().GetCard())
	// 吃
	cardAll = append(cardAll, chiToCards(tc.getChiCards())...)
	// 碰
	cardAll = append(cardAll, pengToCards(tc.getPengCards())...)
	// 杠
	cardAll = append(cardAll, gangToCards(tc.getGangCards())...)
	return cardAll
}

//getAssignCardMap 获取指定牌的映射数量,minCard 最小牌值，maxCard 最大牌值
func getAssignCardMap(cards []*majongpb.Card, minCard, maxCard uint32) map[uint32]int {
	cardCountMap := make(map[uint32]int)
	for _, card := range cards {
		cardValue := gutils.ServerCard2Number(card)
		if minCard <= cardValue && cardValue <= maxCard {
			cardCountMap[cardValue] = cardCountMap[cardValue] + 1
		}
	}
	return cardCountMap
}

// huJoinHandCard 将胡牌加入到手牌中
func huJoinHandCard(handCards []*majongpb.Card, huCard *majongpb.HuCard) []*majongpb.Card {
	checkCard := handCards
	if huCard != nil {
		checkCard = append(checkCard, huCard.GetCard())
	}
	return checkCard
}

// int 转 Card
func intToCard(cardInt int) *majongpb.Card {
	var color majongpb.CardColor
	switch cardInt / 10 {
	case 1:
		color = majongpb.CardColor_ColorWan
	case 2:
		color = majongpb.CardColor_ColorTiao
	case 3:
		color = majongpb.CardColor_ColorTong
	case 4:
		color = majongpb.CardColor_ColorFeng
	case 5:
		color = majongpb.CardColor_ColorHua
	}
	point := int32(cardInt % 10)
	card := &majongpb.Card{
		Color: color,
		Point: point,
	}
	return card
}

func intsToCards(cardInts []int) []*majongpb.Card {
	newCard := make([]*majongpb.Card, 0, len(cardInts))
	for _, card := range cardInts {
		newCard = append(newCard, intToCard(card))
	}
	return newCard
}

//IsXuShuCard 判断是否是序数牌（万，条，筒）
func IsXuShuCard(card *majongpb.Card) bool {
	currColor := card.GetColor()
	return currColor == majongpb.CardColor_ColorWan || currColor == majongpb.CardColor_ColorTiao || currColor == majongpb.CardColor_ColorTong
}

//getPlayerMaxAnKeNum 获取玩家最大暗刻子数
func getPlayerMaxAnKeNum(combines []Combine, num int) bool {
	for _, combine := range combines {
		keLen := len(combine.kes)
		if keLen >= num {
			return true
		}
	}
	return false
}
