package fantype

import (
	"steve/majong/global"
	"steve/majong/utils"
	majongpb "steve/entity/majong"
)

// checkChunXiaQiuDong 检测春夏秋冬
func checkChunXiaQiuDong(tc *typeCalculator) bool {
	chuXiaQiuDong := []majongpb.Card{
		global.Card1H, global.Card2H, global.Card3H, global.Card4H,
	}
	cardCount := make(map[int]int)

	for _, card := range tc.getHuaCards() {
		cardValue := utils.ServerCard2Number(card)
		cardCount[cardValue] = cardCount[cardValue] + 1
	}

	for _, huaPaiCard := range chuXiaQiuDong {
		cardValue := utils.ServerCard2Number(&huaPaiCard)
		if cardCount[cardValue] == 0 {
			return false
		}
	}
	return true
}
