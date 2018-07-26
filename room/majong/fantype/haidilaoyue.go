package fantype

import (
	majongpb "steve/entity/majong"
)

//checkHaiDiLaoYue 检测海底捞月 胡最后打出的牌，必须是最后摸牌的人点炮
func checkHaiDiLaoYue(tc *typeCalculator) bool {
	huCard := tc.getHuCard()
	if huCard != nil && huCard.GetType() == majongpb.HuType_hu_dianpao {
		mjContext := tc.mjContext
		if !IsWallCanMoPai(mjContext) && huCard.GetSrcPlayer() == mjContext.GetLastMopaiPlayer() {
			return true
		}
	}
	return false
}
