package fantype

import (
	majongpb "steve/entity/majong"
)

//checkLaoShaoFu 检查老少副 含有一种花色的123,789的俩副顺子
func checkLaoShaoFu(tc *typeCalculator) bool {
	cards := make([]*majongpb.Card, 0, len(tc.getChiCards()))
	for _, chi := range tc.getChiCards() {
		cards = append(cards, chi.GetCard())
	}
	// 顺+吃
	for _, combine := range tc.combines {
		newCards := append(cards, intsToCards(combine.shuns)...)
		if isLaoShaoFu(newCards) {
			return true
		}
	}
	return false
}

func isLaoShaoFu(newCards []*majongpb.Card) bool {
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
		_, isExist1 := pointMap[one]
		_, isExist2 := pointMap[seven]
		if isExist1 && isExist2 {
			return true
		}
	}
	return false
}
