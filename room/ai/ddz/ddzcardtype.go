package ddz

import (
	"steve/server_pb/ddz"
)

func canBiggerThan(mine ddz.CardType, other ddz.CardType) bool {
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

func getCardType(cards []Poker) (ddz.CardType, *Poker) {
	if yes, pivot := isKingBomb(cards); yes {
		return ddz.CardType_CT_KINGBOMB, pivot
	} else if yes, pivot := isBomb(cards); yes {
		return ddz.CardType_CT_BOMB, pivot
	} else if yes, pivot := isBombAndPairs(cards); yes {
		return ddz.CardType_CT_4SAND2S, pivot
	} else if yes, pivot := isBombAndSingles(cards); yes {
		return ddz.CardType_CT_4SAND1S, pivot
	} else if yes, pivot := isTriples(cards); yes {
		return ddz.CardType_CT_TRIPLES, pivot
	} else if yes, pivot := isTriplesAndPairs(cards); yes {
		return ddz.CardType_CT_3SAND2S, pivot
	} else if yes, pivot := isTriplesAndSingles(cards); yes {
		return ddz.CardType_CT_3SAND1S, pivot
	} else if yes, pivot := isPairs(cards); yes {
		return ddz.CardType_CT_PAIRS, pivot
	} else if yes, pivot := isShunZi(cards); yes {
		return ddz.CardType_CT_SHUNZI, pivot
	} else if yes, pivot := isTriple(cards); yes {
		return ddz.CardType_CT_TRIPLE, pivot
	} else if yes, pivot := isTripleAndPair(cards); yes {
		return ddz.CardType_CT_3AND2, pivot
	} else if yes, pivot := isTripleAndSingle(cards); yes {
		return ddz.CardType_CT_3AND1, pivot
	} else if yes, pivot := isPair(cards); yes {
		return ddz.CardType_CT_PAIR, pivot
	} else if yes, pivot := isSingle(cards); yes {
		return ddz.CardType_CT_SINGLE, pivot
	}

	return ddz.CardType_CT_NONE, nil
}

// 火箭，true时返回大王的牌
func isKingBomb(cards []Poker) (bool, *Poker) {
	if len(cards) != 2 {
		return false, nil
	}

	// 包含大小王
	return Contains(cards, blackJoker) && Contains(cards, redJoker), &redJoker
}

// 炸弹，true时返回最大的牌
func isBomb(cards []Poker) (bool, *Poker) {
	if len(cards) != 4 {
		return false, nil
	}
	return IsAllSamePoint(cards), GetMaxCard(cards)
}

// 四带两对
func isBombAndPairs(cards []Poker) (bool, *Poker) {
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
func isBombAndSingles(cards []Poker) (bool, *Poker) {
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
func isTriples(cards []Poker) (bool, *Poker) {

	// 至少2个三张
	planeCount := len(cards) / 3
	if planeCount < 2 {
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
func isTriplesAndPairs(cards []Poker) (bool, *Poker) {
	planeCount := len(cards) / 5
	if planeCount < 2 {
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

	for i := 0; i < planeCount; i++ {
		pair := GetMaxSamePointCards(remain)
		if len(pair) < 2 { // 3334445555视为飞机带对子
			return false, nil
		}
		remain = RemoveAll(remain, pair[0:2])
	}
	return isMinShunZi(shunZi, 2)
}

// 飞机带单张
func isTriplesAndSingles(cards []Poker) (bool, *Poker) {
	planeCount := len(cards) / 4
	if planeCount < 2 {
		return false, nil
	}

	shunZi := make([]Poker, 0, planeCount)
	remain := cards
	for i := 0; i < planeCount; i++ {
		triple := GetMaxSamePointCards(cards)
		if len(triple) < 3 { // 333344445555视为飞机带翅膀
			return false, nil
		}
		shunZi = append(shunZi, triple[0])
		remain = RemoveAll(remain, triple[0:3])
	}

	return isMinShunZi(shunZi, 2)
}

// 连对
func isPairs(cards []Poker) (bool, *Poker) {
	pairs := len(cards) / 2
	if pairs < 3 {
		return false, nil
	}

	shunZi := make([]Poker, 0, pairs)
	remain := cards
	for i := 0; i < pairs; i++ {
		pair := GetMaxSamePointCards(cards)
		if len(pair) != 2 {
			return false, nil
		}
		shunZi = append(shunZi, pair[0])
		remain = RemoveAll(remain, pair)
	}

	return isMinShunZi(shunZi, 3)
}

// 顺子
func isShunZi(cards []Poker) (bool, *Poker) {
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

	ddzPointSort(cards)
	for i := 0; i < len(cards)-1; i++ {
		if cards[i+1].PointWeight-cards[i].PointWeight != 1 {
			return false, nil
		}
	}

	return true, GetMaxCard(cards)
}

// 三张
func isTriple(cards []Poker) (bool, *Poker) {
	if len(cards) != 3 {
		return false, nil
	}
	if !IsAllSamePoint(cards) {
		return false, nil
	}
	return true, GetMaxCard(cards)
}

// 三张带对子
func isTripleAndPair(cards []Poker) (bool, *Poker) {
	if len(cards) != 5 {
		return false, nil
	}

	// 获取最大相同点数的牌，应该是3张
	triple := GetMaxSamePointCards(cards)
	if len(triple) != 3 {
		return false, nil
	}

	// 移除掉这三张
	remain := RemoveAll(cards, triple)

	// 剩下的应该是对子
	if yes, _ := isPair(remain); !yes {
		return false, nil
	}
	return true, GetMaxCard(triple)
}

// 三张带单张
func isTripleAndSingle(cards []Poker) (bool, *Poker) {
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
func isPair(cards []Poker) (bool, *Poker) {
	if len(cards) != 2 {
		return false, nil
	}
	if !IsAllSamePoint(cards) {
		return false, nil
	}
	return true, GetMaxCard(cards)
}

// 单张
func isSingle(cards []Poker) (bool, *Poker) {
	if len(cards) != 1 {
		return false, nil
	}
	return true, &cards[0]
}

// GetMaxCard 获取棋牌数组中最大的那张牌
func GetMaxCard(cards []Poker) *Poker {

	// 按照棋牌的排序权重进行排序
	DdzPokerSort(cards)

	// 最后的那张牌
	return &cards[len(cards)-1]
}

// 获取最大相同点数的牌, 如 444555533 返回 5555
func GetMaxSamePointCards(cards []Poker) []Poker {

	// 先找到牌中最大相同点数的牌的权重，个数
	pointWeight, count := GetMaxSamePoint(cards)
	maxSamePointCards := make([]Poker, count)

	// 再从牌中找到这个权重的所有牌
	for _, card := range cards {
		if card.PointWeight == pointWeight {
			maxSamePointCards = append(maxSamePointCards, card)
		}
	}
	return maxSamePointCards
}

// 获取最大相同点数的牌的无花色权重和牌的个数, 如 444555533 返回 pointWeight(5), 4。 KKKKAAAA返回pointWeight(A), 4
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

// IsAllSamePoint 是否所有牌都是相同点数
func IsAllSamePoint(cards []Poker) bool {
	for i := 0; i < len(cards)-1; i++ {
		if cards[i].Point != cards[i+1].Point {
			return false
		}
	}
	return true
}

// Contains cards是否包含card
func Contains(cards []Poker, card Poker) bool {
	for _, value := range cards {
		if value.equals(card) {
			return true
		}
	}
	return false
}

// ContainsPoint cards是否包含点数
func ContainsPoint(cards []Poker, point uint32) bool {
	for _, card := range cards {
		if card.Point == point {
			return true
		}
	}
	return false
}

// RemoveAll 从cards中删除removeCards
func RemoveAll(cards []Poker, removeCards []Poker) []Poker {
	var result []Poker
	for _, card := range cards {
		if !Contains(removeCards, card) {
			result = append(result, card)
		}
	}
	return result
}

// AppendAll 从cards中添加addCards
func AppendAll(cards []Poker, addCards []Poker) []Poker {
	for _, card := range addCards {
		cards = append(cards, card)
	}
	return cards
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
