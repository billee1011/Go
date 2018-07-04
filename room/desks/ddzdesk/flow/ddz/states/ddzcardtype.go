package states

import (
	"steve/server_pb/ddz"
	"math"
)

func getCardType(cards []DDZCard) ddz.CardType {
	if isKingBomb(cards) {
		return ddz.CardType_CT_KINGBOMB
	} else if isBomb(cards) {
		return ddz.CardType_CT_BOMB
	} else if isBombAndPairs(cards) {
		return ddz.CardType_CT_4SAND2S
	} else if isBombAndSingles(cards) {
		return ddz.CardType_CT_4SAND1S
	} else if isTriples(cards) {
		return ddz.CardType_CT_TRIPLES
	} else if isTriplesAndPairs(cards) {
		return ddz.CardType_CT_3SAND2S
	} else if isTriplesAndSingles(cards){
		return ddz.CardType_CT_3SAND1S
	} else if isPairs(cards) {
		return ddz.CardType_CT_PAIRS
	} else if isShunZi(cards) {
		return ddz.CardType_CT_SHUNZI
	} else if isTriple(cards) {
		return ddz.CardType_CT_TRIPLE
	} else if isTripleAndPair(cards) {
		return ddz.CardType_CT_3AND2
	} else if isTripleAndSingle(cards) {
		return ddz.CardType_CT_3AND1
	} else if isPair(cards) {
		return ddz.CardType_CT_PAIR
	} else if isSingle(cards) {
		return ddz.CardType_CT_SINGLE
	}

	return ddz.CardType_CT_NONE
}

// 火箭
func isKingBomb(cards []DDZCard) bool {
	if len(cards) != 2 {
		return false
	}
	return Contains(cards, blackJoker) && Contains(cards, redJoker)
}

// 炸弹
func isBomb(cards []DDZCard) bool {
	if len(cards) != 4 {
		return false
	}
	return isAllSamePoint(cards)
}

// 四带两对
func isBombAndPairs(cards []DDZCard) bool {
	if len(cards) != 8 {
		return false
	}

	bomb := getMaxSamePointCards(cards)
	if len(bomb) != 4 {
		return false //没有炸弹
	}

	remain := RemoveAll(cards, bomb)
	firstPair := getMaxSamePointCards(remain)
	if len(firstPair) != 2 {
		return false
	}

	remain = RemoveAll(remain, firstPair)
	if remain[0].point != remain[1].point {
		return false
	}

	return true
}

// 四带二
func isBombAndSingles(cards []DDZCard) bool {
	if len(cards) != 6 {
		return false
	}

	_, count := getMaxSamePoint(cards)
	if count != 4 {
		return false
	}

	return true
}

// 飞机
func isTriples(cards []DDZCard) bool {
	planeCount := len(cards)/3
	if planeCount < 2 {
		return false
	}

	remain := cards
	for i:=0; i<planeCount; i++ {
		triple := getMaxSamePointCards(cards)
		if len(triple) != 3 {
			return false
		}
		remain = RemoveAll(remain, triple)
	}
	return true
}

// 飞机带对子
func isTriplesAndPairs(cards []DDZCard) bool {
	planeCount := len(cards)/5
	if planeCount < 2 {
		return false
	}

	remain := cards
	for i:=0; i<planeCount; i++ {
		triple := getMaxSamePointCards(cards)
		if len(triple) != 3 {
			return false
		}
		remain = RemoveAll(remain, triple)
	}

	for i:=0; i<planeCount; i++ {
		pair := getMaxSamePointCards(cards)
		if len(pair) != 3 {
			return false
		}
		remain = RemoveAll(remain, pair)
	}
	return true
}

// 飞机带单张
func isTriplesAndSingles(cards []DDZCard) bool {
	planeCount := len(cards)/4
	if planeCount < 2 {
		return false
	}

	remain := cards
	for i:=0; i<planeCount; i++ {
		triple := getMaxSamePointCards(cards)
		if len(triple) != 3 {
			return false
		}
		remain = RemoveAll(remain, triple)
	}

	return true
}

// 连对
func isPairs(cards []DDZCard) bool {
	pairs := len(cards)/2
	if pairs < 3 {
		return false
	}

	remain := cards
	for i:=0; i< pairs; i++ {
		pair := getMaxSamePointCards(cards)
		if len(pair) != 2 {
			return false
		}
		remain = RemoveAll(remain, pair)
	}
	return true
}

// 顺子
func isShunZi(cards []DDZCard) bool {
	if len(cards) < 5 {
		return false
	}

	if ContainsPoint(cards, 0x02) || ContainsPoint(cards, redJoker.point) || ContainsPoint(cards, blackJoker.point) {//有2或着大小王直接返回
		return false
	}

	if ContainsPoint(cards, 0x01) && !ContainsPoint(cards, 0x0D) {//有A没K直接返回
		return false
	}

	if ContainsPoint(cards, 0x01) {
		remain, deleted := RemovePoint(cards, 0x01)
		if len(deleted) != 1 {//有多张A，直接返回
			return false
		}
		cards = remain
	}

	cards = ddzPointSort(cards)
	for i:=0; i<len(cards)-1; i++ {
		if math.Abs(float64(cards[i+1].point - cards[i].point)) != 1 { //TODO:如果确定是升序排列，不需要取绝对值
			return false
		}
	}

	return true
}

// 三张
func isTriple(cards []DDZCard) bool {
	if len(cards) != 3 {
		return false
	}
	return isAllSamePoint(cards)
}

// 三张带对子
func isTripleAndPair(cards []DDZCard) bool {
	if len(cards) != 5 {
		return false
	}

	triple := getMaxSamePointCards(cards)
	if len(triple) != 3 {
		return false
	}

	remain := RemoveAll(cards, triple)
	return isAllSamePoint(remain)
}

// 三张带单张
func isTripleAndSingle(cards []DDZCard) bool {
	if len(cards) != 4 {
		return false
	}

	_, count := getMaxSamePoint(cards)
	if count != 3 {
		return false
	}

	return true
}

// 对子
func isPair(cards []DDZCard) bool {
	if len(cards) != 2 {
		return false
	}
	return isAllSamePoint(cards)
}

// 单张
func isSingle(cards []DDZCard) bool {
	if len(cards) != 1 {
		return false
	}
	return true
}

// 获取最大相同点数的牌, 如 444555533 返回 5555
func getMaxSamePointCards(cards []DDZCard) []DDZCard {
	point, count := getMaxSamePoint(cards)
	maxSamePointCards := make([]DDZCard, count)
	for _, card := range cards {
		if card.point == point {
			maxSamePointCards = append(maxSamePointCards, card)
		}
	}
	return maxSamePointCards
}

// 获取最大相同点数, 如 444555533 返回 5,4
func getMaxSamePoint(cards []DDZCard) (maxCountPoint uint32, maxCount uint32) {
	counts := make(map[uint32]uint32) //Map<point, count>
	for _, card := range cards {
		point := card.point
		count, exists := counts[point]
		if !exists {
			counts[point] = 1
		} else {
			counts[point] = count + 1
		}
	}

	for point, count := range counts {
		if count > maxCount {
			maxCount = count
			maxCountPoint = point
		}
	}
	return
}

func isAllSamePoint(cards []DDZCard) bool {
	for i:=0;i<len(cards)-1;i++ {
		if(cards[i].point != cards[i+1].point){
			return false
		}
	}
	return true
}