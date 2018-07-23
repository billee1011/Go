package fantype

import (
	majongpb "steve/server_pb/majong"
)

// checkQuanQiuRen 全求人:全靠吃牌、碰牌、单钓别人打出的牌胡牌
func checkQuanQiuRen(tc *typeCalculator) bool {
	if tc.getHuCard().GetType() != majongpb.HuType_hu_dianpao {
		return false
	}
	for _, gangCard := range tc.getGangCards() {
		if gangCard.GetType() == majongpb.GangType_gang_angang {
			return false
		}
	}
	if len(tc.handCards) == 1 && tc.huCard != nil {
		return true
	}
	return false
}
