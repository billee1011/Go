package fantype

import majongpb "steve/server_pb/majong"

//checkHaiDiLaoYue 检测海底捞月 胡最后打出的牌，必须是最后摸牌的人点炮
func checkHaiDiLaoYue(tc *typeCalculator) bool {
	huCard := tc.getHuCard()
	if huCard != nil && huCard.GetType() == majongpb.HuType_hu_dianpao {
		mjContext := tc.mjContext
		if len(mjContext.GetWallCards()) == 0 && huCard.GetSrcPlayer() == mjContext.GetLastMopaiPlayer() {
			return true
		}
	}
	return false
}
