package fantype

import (
	"steve/majong/global"
	"steve/majong/utils"
	majongpb "steve/server_pb/majong"
)

// checkQuanHua 检测全花
func checkQuanHua(tc *typeCalculator) bool {
	countMap := map[int]int{}
	huaCards := tc.getHuaCards()
	for _, card := range huaCards {
		c := utils.ServerCard2Number(card)
		countMap[c] = countMap[c] + 1
	}

	for _, count := range countMap {
		if count%2 != 0 {
			return false
		}
	}
	huaPaiCards := []majongpb.Card{
		global.Card1H, global.Card2H, global.Card3H, global.Card4H,
		global.Card5H, global.Card6H, global.Card7H, global.Card8H,
	}
	for _, huaPaiCard := range huaPaiCards {
		cardValue := utils.ServerCard2Number(&huaPaiCard)
		if countMap[cardValue] == 0 {
			return false
		}
	}
	return true
}
