package fantype

import (
	majongpb "steve/entity/majong"
)

// checkGangShangKaiHua 检测杠开（杠上开花），杠后摸牌自摸
func checkGangShangKaiHua(tc *typeCalculator) bool {
	huCard := tc.getHuCard()
	if huCard != nil && huCard.GetType() == majongpb.HuType_hu_zimo {
		mjContext := tc.mjContext
		if mjContext.GetMopaiType() == majongpb.MopaiType_MT_GANG {
			return true
		}
	}
	return false
}
