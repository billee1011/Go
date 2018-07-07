package fantype

import majongpb "steve/server_pb/majong"

// checkShuanLongHui 双龙会:由一种花色的 2 个老少副,5 为将牌组成的胡牌
func checkShuanLongHui(tc *typeCalculator) bool {
	cards := make([]*majongpb.Card, 0, len(tc.getChiCards()))
	for _, chi := range tc.getChiCards() {
		cards = append(cards, chi.GetCard())
	}
	// 顺+吃
	for _, combine := range tc.combines {
		if isShuanLongHui(cards, intsToCards(combine.shuns)) {
			return true
		}
	}
	return false
}

func isShuanLongHui(cardA, cardB []*majongpb.Card) bool {
	newCards := append(cardA, cardB...)
	colorPointMap := make(map[majongpb.CardColor]map[int32]int)
	for _, card := range newCards {
		if cardMap, isExist := colorPointMap[card.GetColor()]; isExist {
			cardMap[card.GetPoint()] = cardMap[card.GetPoint()] + 1
			colorPointMap[card.GetColor()] = cardMap
		} else {
			colorPointMap[card.GetColor()] = map[int32]int{card.GetPoint(): 1}
		}
	}

	for _, pointMap := range colorPointMap {
		one, seven := int32(1), int32(7)
		oneNum, _ := pointMap[one]
		severnNum, _ := pointMap[seven]
		if oneNum >= 2 && severnNum >= 2 {
			return true
		}
	}
	return false
}
