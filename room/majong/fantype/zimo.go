package fantype

import majongpb "steve/entity/majong"

// checkZiMo 检测自摸胡 当前玩家摸牌后胡
func checkZiMo(tc *typeCalculator) bool {
	huCard := tc.getHuCard()
	if huCard != nil && huCard.GetType() == majongpb.HuType_hu_zimo {
		return true
	}
	return false
}
