package fantype

import (
	"steve/majong/utils"
	majongpb "steve/server_pb/majong"
)

// checkSanTongShun 三同顺:胡牌时,含有一种花色 3 副序数相同的顺子
func checkSanTongShun(tc *typeCalculator) bool {
	shunCombines := make([]Combine, 0)
	chiCount := len(tc.getChiCards())
	for _, combine := range tc.combines {
		if len(combine.shuns)+chiCount >= 3 {
			shunCombines = append(shunCombines, combine)
		}
	}
	for _, shunCombine := range shunCombines {
		colorCount, cardCount := getChiCardsDetails(tc.getChiCards())
		for _, shun := range shunCombine.shuns {
			cardCount[shun] = cardCount[shun] + 1
			sColor := shun / 10
			colorCount[sColor] = colorCount[sColor] + 1
		}
		hasColor := false
		for _, count := range colorCount {
			if count >= 3 {
				hasColor = true
				break
			}
		}
		hasValue := false
		for _, count := range cardCount {
			if count >= 3 {
				hasValue = true
			}
		}
		if hasColor && hasValue {
			return true
		}
	}
	return false
}

func getChiCardsDetails(chiCards []*majongpb.ChiCard) (colorCount map[int]int, cardCount map[int]int) {
	colorCount = make(map[int]int, 0)
	cardCount = make(map[int]int, 0)
	for _, chiCard := range chiCards {
		chiValue := utils.ServerCard2Number(chiCard.Card)
		chiColor := chiValue / 10
		cardCount[chiValue] = cardCount[chiValue] + 1
		colorCount[chiColor] = colorCount[chiColor] + 1
	}
	return
}
