package fantype

import majongpb "steve/entity/majong"

//checkGangHouPao 检测杠后炮,摸牌类型是杠后摸牌，点炮
func checkGangHouPao(tc *typeCalculator) bool {
	huCard := tc.getHuCard()
	if huCard != nil && huCard.GetType() == majongpb.HuType_hu_dianpao {
		mjContext := tc.mjContext
		if mjContext.GetMopaiType() == majongpb.MopaiType_MT_GANG {
			return true
		}
	}
	return false
}
