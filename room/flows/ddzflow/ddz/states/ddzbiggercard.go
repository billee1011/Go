/*purpose：给定拥有的牌，获取指定牌，指定牌型的压制牌
author   : 李全林
date    ：2018-08-03
*/

package states

import (
	"steve/entity/poker"
)

// GetBoom 若有炸弹，返回炸弹;没有则返回false
func GetBoom(handCards []Poker) (bool, []Poker) {
	bomb := FindSecondaryCards(handCards, 4, 1)
	return bomb != nil, bomb
}

// GetKingBoom 若有炸弹，返回炸弹;没有则返回false
func GetKingBoom(handCards []Poker) (bool, []Poker) {
	has := Contains(handCards, BlackJoker) && Contains(handCards, RedJoker)
	return has, append([]Poker{BlackJoker}, RedJoker)
}

func GetMinBiggerCards(handCards []Poker, outCards []Poker) (bool, []Poker) {
	cardType, pivot := GetCardType(outCards)
	if cardType == poker.CardType_CT_KINGBOMB {
		return false, nil
	} else if cardType == poker.CardType_CT_BOMB {
		bomb := FindMinBiggerCards(handCards, 4, 1, pivot)
		return bomb != nil, bomb
	} else if cardType == poker.CardType_CT_4SAND2S {
		bomb := FindMinBiggerCards(handCards, 4, 1, pivot)
		remain := RemoveAll(handCards, bomb)
		pairs := FindSecondaryCards(remain, 2, 2)
		return bomb != nil && pairs != nil, append(bomb, pairs...)
	} else if cardType == poker.CardType_CT_4SAND1S {
		bomb := FindMinBiggerCards(handCards, 4, 1, pivot)
		remain := RemoveAll(handCards, bomb)
		singles := FindSecondaryCards(remain, 1, 2)
		return bomb != nil && singles != nil, append(bomb, singles...)
	} else if cardType == poker.CardType_CT_TRIPLES {
		triple := FindMinBiggerCards(handCards, 3, len(outCards)/3, pivot)
		return triple != nil, triple
	} else if cardType == poker.CardType_CT_3SAND2S {
		triples := FindMinBiggerCards(handCards, 3, len(outCards)/5, pivot)
		remain := RemoveAll(handCards, triples)
		pairs := FindSecondaryCards(remain, 2, len(outCards)/5)
		return triples != nil && pairs != nil, append(triples, pairs...)
	} else if cardType == poker.CardType_CT_3SAND1S {
		triples := FindMinBiggerCards(handCards, 3, len(outCards)/4, pivot)
		remain := RemoveAll(handCards, triples)
		singles := FindSecondaryCards(remain, 1, len(outCards)/4)
		return triples != nil && singles != nil, append(triples, singles...)
	} else if cardType == poker.CardType_CT_PAIRS {
		pairs := FindMinBiggerCards(handCards, 2, len(outCards)/2, pivot)
		return pairs != nil, pairs
	} else if cardType == poker.CardType_CT_SHUNZI {
		shunzi := FindMinBiggerCards(handCards, 1, len(outCards), pivot)
		return shunzi != nil, shunzi
	} else if cardType == poker.CardType_CT_TRIPLE {
		triple := FindMinBiggerCards(handCards, 3, 1, pivot)
		return triple != nil, triple
	} else if cardType == poker.CardType_CT_3AND2 {
		triple := FindMinBiggerCards(handCards, 3, 1, pivot)
		remain := RemoveAll(handCards, triple)
		pair := FindSecondaryCards(remain, 2, 1)
		return triple != nil && pair != nil, append(triple, pair...)
	} else if cardType == poker.CardType_CT_3AND1 {
		triple := FindMinBiggerCards(handCards, 3, 1, pivot)
		remain := RemoveAll(handCards, triple)
		single := FindSecondaryCards(remain, 1, 1)
		return triple != nil && single != nil, append(triple, single...)
	} else if cardType == poker.CardType_CT_PAIR {
		pair := FindMinBiggerCards(handCards, 2, 1, pivot)
		return pair != nil, pair
	} else if cardType == poker.CardType_CT_SINGLE {
		single := FindMinBiggerCards(handCards, 1, 1, pivot)
		return single != nil, single
	}
	return false, nil
}

/* FindMinBiggerCards 从手牌中找重复数量为duplicateCount，连续长度为shunZiLen，最大牌比maxPivot大一的顺子
理论基础：顺子等价于在去重的情况下，最大牌减最小牌等于集合数量减1
handCards:手牌
duplicateCount:重复数量，单牌为1，对子为2，飞机为3，炸弹为4
shunZiLen:顺子长度，符合这个长度就返回
maxPivot:上手牌的最大牌，顺子传最大
*/
func FindMinBiggerCards(handCards []Poker, duplicateCount int, shunZiLen int, maxPivot *Poker) []Poker {
	countMap := CountSamePointPoker(handCards)
	var matchCards []Poker
	for card, count := range countMap {
		if count >= duplicateCount {
			matchCards = append(matchCards, card)
		}
	}
	DDZPokerSort(matchCards)

	gap := shunZiLen - 1
	for i, card := range matchCards {
		if maxPivot != nil && card.PointWeight <= maxPivot.PointWeight {
			continue
		}
		if i >= gap && matchCards[i].PointWeight-matchCards[i-gap].PointWeight == uint32(gap) {
			shunZi := matchCards[i-gap : i+1]
			var result []Poker
			for _, card := range shunZi {
				inflated := Inflate(handCards, card.PointWeight, duplicateCount)
				result = append(result, inflated...)
			}
			return result
		}
	}
	return nil
}

// Inflate 根据点数还原牌
func Inflate(handCards []Poker, pointWeight uint32, duplicateCount int) (result []Poker) {
	countMap := make(map[uint32]int)
	DDZPokerSort(handCards)
	for _, card := range handCards {
		if card.PointWeight == pointWeight && countMap[pointWeight] < duplicateCount {
			result = append(result, card)
			countMap[card.PointWeight]++
		}
	}
	return
}

/* FindSecondaryCards寻找副牌，如三带一对中的对子
handCards:手牌
duplicateCount:重复数量，单牌为1，对子为2
num:副牌数量
*/
func FindSecondaryCards(handCards []Poker, duplicatCount int, num int) []Poker {
	remain := handCards
	var result []Poker
	for i := 0; i < num; i++ {
		temp := SearchDuplicateCards(remain, duplicatCount, false) //非破坏性找牌
		if temp == nil {
			temp = SearchDuplicateCards(remain, duplicatCount, true) //破坏性找牌，三张会被拆成对子
			if temp == nil {
				return nil
			}
		}
		remain = RemoveAll(remain, temp)
		result = append(result, temp...)
	}
	return result
}

// SearchDuplicateCards，寻找飞机带对子的对子，或者单牌
func SearchDuplicateCards(handCards []Poker, duplicateCount int, chai bool) []Poker {
	countMap := CountSamePointPoker(handCards)
	var matchCards []Poker
	for card, count := range countMap {
		if (chai && count > duplicateCount) || (!chai && count == duplicateCount) {
			matchCards = append(matchCards, card)
		}
	}

	if len(matchCards) > 0 {
		DDZPokerSort(matchCards)
		return Inflate(handCards, matchCards[0].PointWeight, duplicateCount)
	} else {
		return nil
	}
}

// CountSamePointPoker 四个6返回 map<黑桃6, 4>
func CountSamePointPoker(cards []Poker) map[Poker]int {
	countMap := make(map[uint32]int)
	for _, card := range cards {
		countMap[card.PointWeight]++
	}
	cardMap := make(map[Poker]int)
	for _, card := range cards {
		count := countMap[card.PointWeight]
		if count != -1 {
			cardMap[card] = count
			countMap[card.PointWeight] = -1
		}
	}
	return cardMap
}
