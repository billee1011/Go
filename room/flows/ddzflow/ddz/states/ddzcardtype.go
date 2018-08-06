package states

import "steve/entity/poker"

func CanBiggerThan(mine poker.CardType, other poker.CardType) bool {
	if mine == poker.CardType_CT_NONE {
		return false
	}
	if mine == poker.CardType_CT_KINGBOMB {
		return true
	} else if mine == poker.CardType_CT_BOMB && other != poker.CardType_CT_KINGBOMB {
		return true
	} else {
		return mine == other
	}
}

func GetCardType(cards []Poker) (poker.CardType, *Poker) {
	if yes, pivot := IsKingBomb(cards); yes {
		return poker.CardType_CT_KINGBOMB, pivot
	} else if yes, pivot := IsBomb(cards); yes {
		return poker.CardType_CT_BOMB, pivot
	} else if yes, pivot := IsBombAndPairs(cards); yes {
		return poker.CardType_CT_4SAND2S, pivot
	} else if yes, pivot := IsBombAndSingles(cards); yes {
		return poker.CardType_CT_4SAND1S, pivot
	} else if yes, pivot := IsTriples(cards); yes {
		return poker.CardType_CT_TRIPLES, pivot
	} else if yes, pivot := IsTriplesAndPairs(cards); yes {
		return poker.CardType_CT_3SAND2S, pivot
	} else if yes, pivot := IsTriplesAndSingles(cards); yes {
		return poker.CardType_CT_3SAND1S, pivot
	} else if yes, pivot := IsPairs(cards); yes {
		return poker.CardType_CT_PAIRS, pivot
	} else if yes, pivot := IsShunZi(cards); yes {
		return poker.CardType_CT_SHUNZI, pivot
	} else if yes, pivot := IsTriple(cards); yes {
		return poker.CardType_CT_TRIPLE, pivot
	} else if yes, pivot := IsTripleAndPair(cards); yes {
		return poker.CardType_CT_3AND2, pivot
	} else if yes, pivot := IsTripleAndSingle(cards); yes {
		return poker.CardType_CT_3AND1, pivot
	} else if yes, pivot := IsPair(cards); yes {
		return poker.CardType_CT_PAIR, pivot
	} else if yes, pivot := IsSingle(cards); yes {
		return poker.CardType_CT_SINGLE, pivot
	}

	return poker.CardType_CT_NONE, nil
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

	bombs := GetSpecificCountCards(cards, 4)
	if len(bombs) > 0 {
		// 333444 5555牌型 和 333444555666 8888 99JJ牌型
		planes := GetSpecificCountCards(cards, 3)
		if len(planes) != planeCount {
			return false, nil
		}

		remain := RemoveByPoint(cards, planes)
		for i := 0; i < planeCount; i++ {
			pair := GetMaxSamePointCards(remain)
			if len(pair) < 2 {
				return false, nil
			}
			remain = RemoveAll(remain, pair[0:2])
		}

		return isMinShunZi(planes, 2)
	} else {
		remain := cards
		shunZi := make([]Poker, 0, planeCount)
		for i := 0; i < planeCount; i++ {
			triple := GetMaxSamePointCards(remain)
			if len(triple) != 3 {
				return false, nil
			}
			shunZi = append(shunZi, triple[0])
			remain = RemoveAll(remain, triple)
		}

		for i := 0; i < planeCount; i++ {
			pair := GetMaxSamePointCards(remain)
			if len(pair) < 2 {
				return false, nil
			}
			remain = RemoveAll(remain, pair[0:2])
		}
		return isMinShunZi(shunZi, 2)
	}
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

	yes, pivot := isMinShunZi(shunZi, 2)
	if yes {
		return yes, pivot
	}

	// 555666777888 KKKK 牌型
	planes := GetSpecificCountCards(cards, 3)
	if planeCount%3 == 0 && len(planes) == planeCount+planeCount/3 { //555666777KKK 全三牌型
		DDZPokerSort(planes)
		return isMinShunZi(planes[0:len(planes)-planeCount/3], 2)
	}
	if len(planes) != planeCount {
		return false, nil
	}

	return isMinShunZi(planes, 2)
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

	DDZPokerSort(cards)
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
	DDZPokerSortDesc(cards)
	return &cards[0]
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
	DDZPokerSortDesc(maxSamePointCards)
	return maxSamePointCards
}

// 获取指定数量点数相同的牌，如 555666777888 KKKK，num = 3，返回5678
func GetSpecificCountCards(cards []Poker, num int) []Poker {
	return getCountCards(cards, num, false)
}

func getCountCards(cards []Poker, num int, canMore bool) []Poker {
	if num < 1 {
		return []Poker{}
	}
	sameCards := make(map[uint32][]Poker) //Map<PointWeight, []Poker>
	for _, card := range cards {
		pointWeight := card.PointWeight
		sameCard := sameCards[pointWeight]
		sameCards[pointWeight] = append(sameCard, card)
	}

	var result []Poker
	for _, sameCard := range sameCards {
		if (!canMore && len(sameCard) == num) || (canMore && len(sameCard) >= num) {
			result = append(result, sameCard[0])
		}
	}
	return result
}

// 获取最大相同点数, 如 444555533 返回 pointWeight(5), 4。 KKKKAAAA返回pointWeight(A), 4
func GetMaxSamePoint(cards []Poker) (maxCountPointWeight uint32, maxCount uint32) {
	counts := make(map[uint32]uint32) //Map<PointWeight, count>
	for _, card := range cards {
		pointWeight := card.PointWeight
		counts[pointWeight]++
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
