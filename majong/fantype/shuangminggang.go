package fantype

import (
	majongpb "steve/server_pb/majong"
)

// checkShuangMingGang 检测双明杠 含有两个明杠(直杠，补杠)
func checkShuangMingGang(tc *typeCalculator) bool {
	gangNum := 0
	for _, gangCard := range tc.getGangCards() {
		if gangCard.GetType() != majongpb.GangType_gang_angang {
			gangNum++
		}
	}
	return gangNum >= 2
}
