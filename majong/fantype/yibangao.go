package fantype

import (
	"steve/gutils"
	majongpb "steve/server_pb/majong"
)

//checkYiBangGao 检测一般高 由一种花色2副相同的顺子组成的胡牌
func checkYiBangGao(tc *typeCalculator) bool {
	cards := make([]*majongpb.Card, 0)
	// 吃
	for _, chi := range tc.getChiCards() {
		cards = append(cards, chi.GetCard())
	}
	for _, combine := range tc.combines {
		shunCards := intsToCards(combine.shuns)
		// 顺子+吃
		newCards := append(cards, shunCards...)
		if len(newCards) < 2 {
			continue
		}
		// 判断顺子+吃中是否有相同的
		for i := 0; i < len(newCards); i++ {
			comCard := append(newCards[:i], newCards[i+1:]...)
			for _, card := range comCard {
				if gutils.CardEqual(card, newCards[i]) {
					return true
				}
			}
		}
	}
	return false
}
