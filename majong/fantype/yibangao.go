package fantype

import (
	"steve/gutils"
)

//checkYiBangGao 检测一般高 由一种花色2副相同的顺子组成的胡牌
func checkYiBangGao(tc *typeCalculator) bool {
	for _, combine := range tc.combines {
		cards := make([]int, 0)
		// 吃
		for _, chi := range tc.getChiCards() {
			cardValue := gutils.ServerCard2Number(chi.GetCard())
			cards = append(cards, int(cardValue))
		}
		// 顺子+吃
		newCards := append(cards, combine.shuns...)
		if len(newCards) < 2 {
			continue
		}
		cardMap := make(map[int]int)
		// 判断顺子+吃中是否有相同的
		for i := 0; i < len(newCards); i++ {
			cardMap[newCards[i]] = cardMap[newCards[i]] + 1
			if cardMap[newCards[i]] >= 2 {
				return true
			}
		}
	}
	return false
}
