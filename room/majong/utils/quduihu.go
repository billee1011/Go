package utils

//FastCheckQiDuiHu 判断七小对胡牌，判断手牌三否是14张，且将的总数是7,不适用有赖子的情况
func FastCheckQiDuiHu(cards []Card) bool {
	if len(cards) != gMaxSize {
		return false
	}
	var cm countMap = make(map[Card]int)
	cm.addAll(cards...)
	for _, sum := range cm {
		if sum%2 != 0 {
			return false
		}
	}
	return true
}

//FastCheckQiDuiTing 七小对查听，当手牌为13张牌时，遍历牌墙所有牌是否可胡
func FastCheckQiDuiTing(cards []Card, avalibleCards []Card) []Card {
	if len(cards) != (gMaxSize - 1) {
		return []Card{}
	}
	tingCards := make([]Card, 0)
	for _, avalibleCard := range avalibleCards {
		cards = append(cards, avalibleCard)
		hu := FastCheckQiDuiHu(cards)
		if hu {
			tingCards = append(tingCards, avalibleCard)
		}
		cards = cards[:len(cards)-1]
	}
	return tingCards
}

//FastCheckQiDuiTingInfo 七小对查听-对应的听牌提示，打什么牌可以听什么牌
func FastCheckQiDuiTingInfo(cards []Card, avalibleCards []Card) map[Card][]Card {
	if len(cards) != gMaxSize {
		return nil
	}
	tingInfo := make(map[Card][]Card)
	var cm countMap = make(map[Card]int)
	cm.addAll(cards...)
	for index, card := range cards {
		if cm[card] > 0 {
			cm[card] = 0
			cards = append(cards[:index], cards[index+1:]...)
			tingCards := FastCheckQiDuiTing(cards, avalibleCards)
			if len(tingCards) > 0 {
				tingInfo[card] = tingCards
			}
			cards = append(cards[:index], append([]Card{card}, cards[index:]...)...)
		}
	}
	return tingInfo
}
