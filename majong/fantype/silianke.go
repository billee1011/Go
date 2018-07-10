package fantype

import (
	"steve/majong/utils"
	majongpb "steve/server_pb/majong"
)

// checkSiLianKe 四连刻:胡牌时,含有一种花色 4 副依次递增一位数的刻子;
func checkSiLianKe(tc *typeCalculator) bool {
	pengCount := len(tc.getPengCards())
	keCombines := make([]Combine, 0)
	for _, combine := range tc.combines {
		if len(combine.kes)+pengCount >= 4 {
			keCombines = append(keCombines, combine)
		}
	}
	for _, keCombine := range keCombines {
		colorCount, cardCount, minValue := getPengCardsDetails(tc.getPengCards())
		for _, ke := range keCombine.kes {
			cardCount[ke] = cardCount[ke] + 1
			kcolor := ke / 10
			colorCount[kcolor] = colorCount[kcolor] + 1
			if minValue == 0 || ke < minValue {
				minValue = ke
			}
		}
		if len(colorCount) > 1 {
			return false
		}
		for i := 0; i < 4; i++ {
			if cardCount[minValue+i] == 0 {
				return false
			}
		}
	}
	return true
}

func getPengCardsDetails(pengCards []*majongpb.PengCard) (colorCount map[int]int, cardCount map[int]int, minValue int) {
	colorCount = make(map[int]int, 0)
	cardCount = make(map[int]int, 0)
	if len(pengCards) != 0 {
		minValue = utils.ServerCard2Number(pengCards[0].Card)
	}
	for _, pengCard := range pengCards {
		pengValue := utils.ServerCard2Number(pengCard.Card)
		pengColor := pengValue / 10
		cardCount[pengValue] = cardCount[pengValue] + 1
		colorCount[pengColor] = colorCount[pengColor] + 1
		if pengValue < minValue {
			minValue = pengValue
		}
	}
	return
}
