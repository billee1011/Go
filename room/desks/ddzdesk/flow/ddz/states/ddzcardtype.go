package states

import (
	"steve/server_pb/ddz"
)

func getCardType(cards []uint32) ddz.CardType {
	if isKingBomb(cards) {
		return ddz.CardType_CT_KINGBOMB
	}
	if isBomb(cards) {
		return ddz.CardType_CT_BOMB
	}

	return ddz.CardType_CT_NONE
}

// 火箭
func isKingBomb(cards []uint32) bool {
	if len(cards) != 2 {
		return false
	}
	return Contains(cards, 0x0E) && Contains(cards, 0x0F)
}

// 炸弹
func isBomb(cards []uint32) bool {
	if len(cards) != 4 {
		return false
	}
	return isAllSamePoint(cards)
}

// 四带二
func isBombAndPairs(cards []uint32) bool {
	if len(cards) != 8 {
		return false
	}

	return true
}

func getMaxSamePointCards(cards []uint32) []uint32 {
	point, count := getMaxSamePoint(cards)
	maxSamePointCards := make([]uint32, count)
	for _, card := range cards {
		if card%16 == point {
			maxSamePointCards = append(maxSamePointCards, card)
		}
	}
	return maxSamePointCards
}

func getMaxSamePoint(cards []uint32) (maxCountPoint uint32, maxCount uint32) {
	counts := make(map[uint32]uint32) //Map<point, count>
	for _, card := range cards {
		point := card % 16
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

func isSamePoint(card1 uint32, card2 uint32) bool {
	return card1%16 == card2%16
}

func isAllSamePoint(cards []uint32) bool {
	for i:=0;i<len(cards)-1;i++ {
		if(!isSamePoint(cards[i], cards[i+1])){
			return false
		}
	}
	return true
}