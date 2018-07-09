package fantype

import (
	"steve/majong/utils"
	majongpb "steve/server_pb/majong"
)

// checkSanTongShun 三同顺:胡牌时,含有一种花色 3 副序数相同的顺子
func checkSanTongShun(tc *typeCalculator) bool {
	shunCombines := make([]combine, 0)
	chiCount := len(tc.getChiCards())
	for _, combine := range tc.combines {
		if len(combine.shuns)+chiCount >= 3 {
			shunCombines = append(shunCombines, combine)
		}
	}
	for _, shunCombine := range shunCombines {
		colorCount, cardCount := getChiCardsDetails(tc.getChiCards())
		for _, shun := range shunCombine.shuns {
			shunValue := utils.ServerCard2Number(shun)
			cardCount[shunValue] = cardCount[shunValue] + 1
			colorCount[shun.Color] = colorCount[shun.Color] + 1
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
			if cardCount[count+1] != 0 && cardCount[count+2] != 0 {
				hasValue = true
			}
		}
		if !hasColor || !hasValue {
			return false
		}
	}
	return false
}

func getChiCardsDetails(chiCards []*majongpb.ChiCard) (colorCount map[majongpb.CardColor]int, cardCount map[int]int) {
	colorCount = make(map[majongpb.CardColor]int, 0)
	cardCount = make(map[int]int, 0)
	for _, chiCard := range chiCards {
		chiValue := utils.ServerCard2Number(chiCard.Card)
		cardCount[chiValue] = cardCount[chiValue] + 1
		colorCount[chiCard.Card.Color] = colorCount[chiCard.Card.Color] + 1
	}
	return
}
