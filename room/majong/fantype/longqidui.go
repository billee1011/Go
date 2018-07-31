package fantype

import (
	"steve/gutils"
)

//checkLongQiDui 检查龙七对，七对+1根（4张一样的牌）
func checkLongQiDui(tc *typeCalculator) bool {
	if tc.callCheckFunc(qiduiFuncID) {
		cardMap := make(map[uint32]int)
		for _, card := range tc.getHandCards() {
			cardValue := gutils.ServerCard2Number(card)
			cardMap[cardValue] = cardMap[cardValue] + 1
		}
		for _, count := range cardMap {
			if count == 4 {
				return true
			}
		}
	}
	return false
}
