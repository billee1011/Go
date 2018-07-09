package fantype

import (
	"steve/majong/utils"
	majongpb "steve/server_pb/majong"
)

// checkSanLianKe 三连刻:胡牌时,含有一种花色 3 副依次递增一位数字的刻子
func checkSanLianKe(tc *typeCalculator) bool {
	pengCount := len(tc.getPengCards())
	keCombines := make([]Combine, 0)
	for _, combine := range tc.combines {
		if len(combine.kes)+pengCount >= 3 {
			keCombines = append(keCombines, combine)
		}
	}
	for _, keCombine := range keCombines {
		colorCount, cardCount := getChiCardsDetails(tc.getChiCards())
		for _, ke := range keCombine.kes {
			keCard := intToCard(ke)
			cardCount[ke] = cardCount[ke] + 1
			colorCount[keCard.Color] = colorCount[keCard.Color] + 1
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

func getPengCardDetails(pengCards []*majongpb.PengCard) (colorCount map[majongpb.CardColor]int, cardCount map[int]int) {
	colorCount = make(map[majongpb.CardColor]int, 0)
	cardCount = make(map[int]int, 0)
	for _, pengCard := range pengCards {
		chiValue := utils.ServerCard2Number(pengCard.Card)
		cardCount[chiValue] = cardCount[chiValue] + 1
		colorCount[pengCard.Card.Color] = colorCount[pengCard.Card.Color] + 1
	}
	return
}
