package fantype

import (
	majongpb "steve/server_pb/majong"
)

//checkBuQiuRen 检测不求人 4副牌及将牌中，没有吃，碰，明杠的自摸胡牌
func checkBuQiuRen(tc *typeCalculator) bool {
	huCard := tc.getHuCard()
	// 必须是自摸
	if huCard != nil && huCard.GetType() != majongpb.HuType_hu_zimo {
		return false
	}
	// 不能有吃碰
	chiGangNum := len(tc.getChiCards()) + len(tc.getPengCards())
	if chiGangNum != 0 {
		return false
	}
	// 不能有明杠
	for _, gangCard := range tc.getGangCards() {
		if gangCard.GetType() != majongpb.GangType_gang_angang {
			return false
		}
	}
	return true
}
