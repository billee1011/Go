package states

import (
	"steve/server_pb/ddz"
)

func CanBiggerThan(mine ddz.CardType, other ddz.CardType) bool {
	if mine == ddz.CardType_CT_NONE {
		return false
	}
	if mine == ddz.CardType_CT_KINGBOMB {
		return true
	} else if mine == ddz.CardType_CT_BOMB && other != ddz.CardType_CT_KINGBOMB {
		return true
	} else {
		return mine == other
	}
}

func GetCardType(cards []Poker) (ddz.CardType, *Poker) {
	if yes, pivot := IsKingBomb(cards); yes {
		return ddz.CardType_CT_KINGBOMB, pivot
	} else if yes, pivot := IsBomb(cards); yes {
		return ddz.CardType_CT_BOMB, pivot
	} else if yes, pivot := IsBombAndPairs(cards); yes {
		return ddz.CardType_CT_4SAND2S, pivot
	} else if yes, pivot := IsBombAndSingles(cards); yes {
		return ddz.CardType_CT_4SAND1S, pivot
	} else if yes, pivot := IsTriples(cards); yes {
		return ddz.CardType_CT_TRIPLES, pivot
	} else if yes, pivot := IsTriplesAndPairs(cards); yes {
		return ddz.CardType_CT_3SAND2S, pivot
	} else if yes, pivot := IsTriplesAndSingles(cards); yes {
		return ddz.CardType_CT_3SAND1S, pivot
	} else if yes, pivot := IsPairs(cards); yes {
		return ddz.CardType_CT_PAIRS, pivot
	} else if yes, pivot := IsShunZi(cards); yes {
		return ddz.CardType_CT_SHUNZI, pivot
	} else if yes, pivot := IsTriple(cards); yes {
		return ddz.CardType_CT_TRIPLE, pivot
	} else if yes, pivot := IsTripleAndPair(cards); yes {
		return ddz.CardType_CT_3AND2, pivot
	} else if yes, pivot := IsTripleAndSingle(cards); yes {
		return ddz.CardType_CT_3AND1, pivot
	} else if yes, pivot := IsPair(cards); yes {
		return ddz.CardType_CT_PAIR, pivot
	} else if yes, pivot := IsSingle(cards); yes {
		return ddz.CardType_CT_SINGLE, pivot
	}

	return ddz.CardType_CT_NONE, nil
}

// 火箭
func IsKingBomb(cards []Poker) (bool, *Poker) {
	if len(cards) != 2 {
		return false, nil
	}
	return Contains(cards, BlackJoker) && Contains(cards, RedJoker), &RedJoker
}

// 炸弹
func IsBomb(cards []Poker) (bool, *Poker) {
	if len(cards) != 4 {
		return false, nil
	}
	return IsAllSamePoint(cards), GetMaxCard(cards)
}

// 四带两对
func IsBombAndPairs(cards []Poker) (bool, *Poker) {
	if len(cards) != 8 {
		return false, nil
	}

	bomb := GetMaxSamePointCards(cards)
	if len(bomb) != 4 {
		return false, nil //没有炸弹
	}

	remain := RemoveAll(cards, bomb)
	firstPair := GetMaxSamePointCards(remain)
	if len(firstPair) < 2 { //44445555视为四带两对
		return false, nil
	}

	remain = RemoveAll(remain, firstPair[0:2])
	if remain[0].Point != remain[1].Point {
		return false, nil
	}

	return true, GetMaxCard(bomb)
}

// 四带二单张
func IsBombAndSingles(cards []Poker) (bool, *Poker) {
	if len(cards) != 6 {
		return false, nil
	}

	bomb := GetMaxSamePointCards(cards)
	if len(bomb) != 4 {
		return false, nil
	}

	return true, GetMaxCard(bomb)
}

// 飞机
func IsTriples(cards []Poker) (bool, *Poker) {
	planeCount := len(cards) / 3
	if len(cards)%3 != 0 || planeCount < 2 {
		return false, nil
	}

	shunZi := make([]Poker, 0, planeCount)
	remain := cards
	for i := 0; i < planeCount; i++ {
		triple := GetMaxSamePointCards(remain)
		if len(triple) != 3 {
			return false, nil
		}
		shunZi = append(shunZi, triple[0])
		remain = RemoveAll(remain, triple)
	}
	return isMinShunZi(shunZi, 2)
}

// 飞机带对子
func IsTriplesAndPairs(cards []Poker) (bool, *Poker) {
	planeCount := len(cards) / 5
	if len(cards)%5 != 0 || planeCount < 2 {
		return false, nil
	}

	bomb := GetMaxSamePointCards(cards)

	var remain []Poker
	if len(bomb) == 4 { // 3334445555视为飞机带对子
		remain = RemoveAll(cards, bomb)
	} else {
		remain = cards
	}
	shunZi := make([]Poker, 0, planeCount)
	for i := 0; i < planeCount; i++ {
		triple := GetMaxSamePointCards(remain)
		if len(triple) != 3 {
			return false, nil
		}
		shunZi = append(shunZi, triple[0])
		remain = RemoveAll(remain, triple)
	}

	if len(bomb) != 4 {
		for i := 0; i < planeCount; i++ {
			pair := GetMaxSamePointCards(remain)
			if len(pair) < 2 {
				return false, nil
			}
			remain = RemoveAll(remain, pair[0:2])
		}
	}
	return isMinShunZi(shunZi, 2)
}

// 飞机带单张
func IsTriplesAndSingles(cards []Poker) (bool, *Poker) {
	planeCount := len(cards) / 4
	if len(cards)%4 != 0 || planeCount < 2 {
		return false, nil
	}

	shunZi := make([]Poker, 0, planeCount)
	remain := cards
	for i := 0; i < planeCount; i++ {
		triple := GetMaxSamePointCards(remain)
		if len(triple) < 3 { // 333344445555视为飞机带翅膀
			return false, nil
		}
		shunZi = append(shunZi, triple[0])
		remain = RemoveAll(remain, triple[0:3])
	}

	return isMinShunZi(shunZi, 2)
}

// 连对
func IsPairs(cards []Poker) (bool, *Poker) {
	pairs := len(cards) / 2
	if len(cards)%2 != 0 || pairs < 3 {
		return false, nil
	}

	shunZi := make([]Poker, 0, pairs)
	remain := cards
	for i := 0; i < pairs; i++ {
		pair := GetMaxSamePointCards(remain)
		if len(pair) != 2 {
			return false, nil
		}
		shunZi = append(shunZi, pair[0])
		remain = RemoveAll(remain, pair)
	}

	return isMinShunZi(shunZi, 3)
}

// 顺子
func IsShunZi(cards []Poker) (bool, *Poker) {
	return isMinShunZi(cards, 5)
}

// 带最小长度的顺子判断
func isMinShunZi(cards []Poker, minLen int) (bool, *Poker) {
	if len(cards) < minLen {
		return false, nil
	}

	if ContainsPoint(cards, p2) || ContainsPoint(cards, pRedJoker) || ContainsPoint(cards, pBlackJoker) { //有2或着大小王直接返回
		return false, nil
	}

	DDZPointSort(cards)
	for i := 0; i < len(cards)-1; i++ {
		if cards[i+1].PointWeight-cards[i].PointWeight != 1 {
			return false, nil
		}
	}

	return true, GetMaxCard(cards)
}

// 三张
func IsTriple(cards []Poker) (bool, *Poker) {
	if len(cards) != 3 {
		return false, nil
	}
	if !IsAllSamePoint(cards) {
		return false, nil
	}
	return true, GetMaxCard(cards)
}

// 三张带对子
func IsTripleAndPair(cards []Poker) (bool, *Poker) {
	if len(cards) != 5 {
		return false, nil
	}

	triple := GetMaxSamePointCards(cards)
	if len(triple) != 3 {
		return false, nil
	}

	remain := RemoveAll(cards, triple)
	if yes, _ := IsPair(remain); !yes {
		return false, nil
	}
	return true, GetMaxCard(triple)
}

// 三张带单张
func IsTripleAndSingle(cards []Poker) (bool, *Poker) {
	if len(cards) != 4 {
		return false, nil
	}

	triple := GetMaxSamePointCards(cards)
	if len(triple) != 3 {
		return false, nil
	}

	return true, GetMaxCard(triple)
}

// 对子
func IsPair(cards []Poker) (bool, *Poker) {
	if len(cards) != 2 {
		return false, nil
	}
	if !IsAllSamePoint(cards) {
		return false, nil
	}
	return true, GetMaxCard(cards)
}

// 单张
func IsSingle(cards []Poker) (bool, *Poker) {
	if len(cards) != 1 {
		return false, nil
	}
	return true, &cards[0]
}

func GetMaxCard(cards []Poker) *Poker {
	DDZPokerSort(cards)
	return &cards[len(cards)-1]
}

// GetMinCard 获取一组牌中最小的那张牌
func GetMinCard(cards []Poker) *Poker {
	DDZPokerSort(cards)
	return &cards[0]
}

// 获取最大相同点数的牌, 如 444555533 返回 5555
func GetMaxSamePointCards(cards []Poker) []Poker {
	pointWeight, count := GetMaxSamePoint(cards)
	maxSamePointCards := make([]Poker, 0, count)
	for _, card := range cards {
		if card.PointWeight == pointWeight {
			maxSamePointCards = append(maxSamePointCards, card)
		}
	}
	return maxSamePointCards
}

// 获取最大相同点数, 如 444555533 返回 pointWeight(5), 4。 KKKKAAAA返回pointWeight(A), 4
func GetMaxSamePoint(cards []Poker) (maxCountPointWeight uint32, maxCount uint32) {
	counts := make(map[uint32]uint32) //Map<Point, count>
	for _, card := range cards {
		pointWeight := card.PointWeight
		count, exists := counts[pointWeight]
		if !exists {
			counts[pointWeight] = 1
		} else {
			counts[pointWeight] = count + 1
		}
	}

	for pointWeight, count := range counts {
		if count > maxCount || (count == maxCount && pointWeight > maxCountPointWeight) {
			maxCount = count
			maxCountPointWeight = pointWeight
		}
	}
	return
}

func IsAllSamePoint(cards []Poker) bool {
	for i := 0; i < len(cards)-1; i++ {
		if cards[i].Point != cards[i+1].Point {
			return false
		}
	}
	return true
}

// ContainsPointWeightCount cards中包含指定无花色权重点数的牌的个数
func ContainsPointWeightCount(cards []Poker, pointWeight uint32) uint32 {
	var count uint32 = 0
	for _, card := range cards {
		if card.PointWeight == pointWeight {
			count++
		}
	}
	return count
}
