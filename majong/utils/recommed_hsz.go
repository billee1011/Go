package utils

import (
	"math/rand"
	majongpb "steve/server_pb/majong"
	"time"

	"github.com/Sirupsen/logrus"
)

//CardDis 牌距离
type CardDis struct {
	Card    *majongpb.Card
	DisHand int32
	DisLast int32
}

// GetRecommedHuanSanZhang 获取推荐换三张牌
func GetRecommedHuanSanZhang(handCards []*majongpb.Card) []*majongpb.Card {
	colorCardsMap := GetColorStatistics(handCards)
	colorCardsMap = CheckTeShu(colorCardsMap) // 特殊
	if flag, cards := IsCanGetRecommedCars(colorCardsMap); flag {
		return cards
	}
	// 获取每种比较的颜色的优先级
	sortPriority, colorPrioMap := GetColorPriorityInfo(colorCardsMap)
	minPriority := sortPriority[len(sortPriority)-1] // 获取最小优先级 1为最大，18为最小优先级
	colors := GetColorByPriority(colorPrioMap, minPriority)
	if len(colors) == 1 {
		return DingCard(colorCardsMap[colors[0]])
	}
	return CardTypeIsSame(colors, colorCardsMap)
}

// CardTypeIsSame 牌型一样
func CardTypeIsSame(colors []majongpb.CardColor, colorCardsMap map[majongpb.CardColor][]*majongpb.Card) []*majongpb.Card {
	// 牌数不一样，选择最小牌数
	if flag, minCards := IsCardNumEqualAndMinCards(colors, colorCardsMap); !flag {
		logrus.WithFields(logrus.Fields{"func_name": "GetRecommedHuanSanZhang",
			"minCards": minCards, "colors": colors}).Info("牌型一样，牌数不一样，选择最小牌数")
		return DingCard(minCards)
	}
	//牌数一样，随机
	rd := rand.New(rand.NewSource(time.Now().UnixNano())) // 随机出颜色
	towards := rd.Intn(len(colors))
	cards := colorCardsMap[colors[towards]]
	logrus.WithFields(logrus.Fields{"func_name": "GetRecommedHuanSanZhang",
		"cards": cards, "towards": towards, "colors": colors}).Info("牌型一样，牌数一样，随机")
	return DingCard(cards)
}

//IsCardNumEqualAndMinCards 判断牌数是否相等，并返回最小的牌数组
func IsCardNumEqualAndMinCards(colors []majongpb.CardColor, colorCardsMap map[majongpb.CardColor][]*majongpb.Card) (bool, []*majongpb.Card) {
	minCards, flag := colorCardsMap[colors[0]], true
	for _, color := range colors {
		cards := colorCardsMap[color]
		if len(minCards) != len(cards) {
			flag = false
			if len(minCards) > len(cards) {
				minCards = cards
			}
		}
	}
	return flag, minCards
}

//IsCanGetRecommedCars 是否能获取推荐牌
func IsCanGetRecommedCars(colorCardsMap map[majongpb.CardColor][]*majongpb.Card) (bool, []*majongpb.Card) {
	min, mid, max := ColorSort(colorCardsMap)
	//如果有两门牌小于3，只有1门牌数大于3张
	if mid < 3 && max >= 3 {
		// 定牌逻辑
		color := GetCardColorsByLen(colorCardsMap, max)[0]
		cards := colorCardsMap[color]
		return true, DingCard(cards)
	}
	//选牌判断
	return SelectCards(min, mid, max, colorCardsMap)
}

//SelectCards 选牌判断
func SelectCards(min, mid, max int, colorCardsMap map[majongpb.CardColor][]*majongpb.Card) (bool, []*majongpb.Card) {
	if min < 3 {
		// 最少的花色的牌数少于3，不比较
		colors := GetCardColorsByLen(colorCardsMap, min)
		delete(colorCardsMap, colors[0])
		//现在最小花色，与最大花色差>= 2，直接选定现在最小花色
		if max-mid >= 2 {
			color := GetCardColorsByLen(colorCardsMap, mid)[0]
			cards := colorCardsMap[color]
			return true, DingCard(cards)
		}
	} else {
		// 中间花色，与最小花色差>= 2，直接选定最小花色
		if mid-min >= 2 {
			color := GetCardColorsByLen(colorCardsMap, min)[0]
			cards := colorCardsMap[color]
			return true, DingCard(cards)
		}
		// 最大花色，与中间花色差>= 2 ，删除最大花色
		if max-min >= 2 {
			colors := GetCardColorsByLen(colorCardsMap, max)
			delete(colorCardsMap, colors[0])
		}
	}
	return false, []*majongpb.Card{}
}

//CheckTeShu 特殊说明：牌型组合为只有四张并形成杠的牌型，此牌型为特殊牌型，不会作为换牌换出。四张的杠与八张的两个杠比较的话，也是选八张为选牌。
func CheckTeShu(colorCardsMap map[majongpb.CardColor][]*majongpb.Card) map[majongpb.CardColor][]*majongpb.Card {
	newMap := CopyColorCardMap(colorCardsMap)
	for color, cards := range colorCardsMap {
		if len(cards) == 4 {
			currCard, flag := cards[0], true
			for _, card := range cards {
				if !CardEqual(currCard, card) {
					flag = false
					break
				}
			}
			if flag {
				delete(newMap, color)
			}
		}
	}
	//以后可能会出现，三个都是杠的情况
	if len(newMap) == 0 {
		return colorCardsMap
	}
	return newMap
}

//DingCard 定牌
func DingCard(colorCards []*majongpb.Card) []*majongpb.Card {
	if len(colorCards) == 3 {
		return colorCards
	}
	SortCards(colorCards) // 升序
	hand := colorCards[0]
	last := colorCards[len(colorCards)-1]
	//获取要比较牌的距离信息
	readyCards, compareCardDis := GetCompareCardDisInfo(hand, last, colorCards)
	if len(readyCards) == 3 {
		return readyCards
	}
	//比较最远的牌
	randCards := CompareFarDistance(compareCardDis)
	randL := len(randCards)
	if randL > 1 { // 有相同的级别的定牌
		rd := rand.New(rand.NewSource(time.Now().UnixNano()))
		towards := rd.Intn(randL)
		readyCards = append(readyCards, randCards[towards])
	} else {
		readyCards = append(readyCards, randCards[0])
	}
	return readyCards
}

//CompareFarDistance 比较最远的牌
func CompareFarDistance(compareCardDis []CardDis) []*majongpb.Card {
	// 初始比较的牌
	compareCard := compareCardDis[0]
	minDis := compareCard.DisHand
	if compareCard.DisHand > compareCard.DisLast {
		minDis = compareCard.DisLast
	}
	// 初始最小级别的牌数组
	randCards := []*majongpb.Card{compareCard.Card}
	for i := 1; i < len(compareCardDis); i++ {
		// 最小的距离
		currMin := compareCardDis[i].DisHand
		if compareCardDis[i].DisHand > compareCardDis[i].DisLast {
			currMin = compareCardDis[i].DisLast
		}
		//比较最小距离
		if currMin == minDis { //有相同级别的
			randCards = append(randCards, compareCardDis[i].Card)
			break
		} else if currMin > minDis {
			minDis = currMin
			compareCard = compareCardDis[i]
			randCards = []*majongpb.Card{compareCard.Card}
		}
	}
	return randCards
}

//GetCompareCardDisInfo 获取要比较牌的距离信息
func GetCompareCardDisInfo(hand, last *majongpb.Card, colorCards []*majongpb.Card) ([]*majongpb.Card, []CardDis) {
	// 获取中间的牌离，头和尾的距离
	sum := (last.GetPoint() + hand.GetPoint())
	l := len(colorCards)
	readyCards := append([]*majongpb.Card{}, hand, last)
	compareCardDis := make([]CardDis, 0)
	for _, card := range NotDuplicatesCards(colorCards[1 : l-1]) {
		currPoint := card.GetPoint()
		if sum%2 == 0 && currPoint == sum/2 {
			return append(readyCards, card), compareCardDis
		}
		disHand := currPoint - hand.GetPoint()
		disLast := last.GetPoint() - currPoint
		compareCardDis = append(compareCardDis, CardDis{
			Card:    card,
			DisHand: disHand,
			DisLast: disLast,
		})
	}
	return readyCards, compareCardDis
}
