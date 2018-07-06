package fantype

import (
	"steve/gutils"
	"steve/majong/global"
	majongpb "steve/server_pb/majong"
)

// checkMeiLanZhuJiu 检测梅兰竹菊
func checkMeiLanZhuJiu(tc *typeCalculator) bool {
	meiLangZhuJu := []majongpb.Card{
		global.Card5H, global.Card6H, global.Card7H, global.Card8H,
	}
	cardCount := make(map[uint32]int)

	for _, card := range tc.getHuaCards() {
		cardValue := gutils.ServerCard2Number(card)
		cardCount[cardValue] = cardCount[cardValue] + 1
	}

	for _, huaPaiCard := range meiLangZhuJu {
		cardValue := gutils.ServerCard2Number(&huaPaiCard)
		if cardCount[cardValue] == 0 {
			return false
		}
	}
	return true
}
