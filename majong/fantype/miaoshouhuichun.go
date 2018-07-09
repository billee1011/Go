package fantype

import majongpb "steve/server_pb/majong"

//miaoshouhuichun 检测妙手回春，最后一张牌自摸
func checkMiaoShouHuiChun(tc *typeCalculator) bool {
	huCard := tc.getHuCard()
	if huCard != nil && huCard.GetType() == majongpb.HuType_hu_zimo {
		mjContext := tc.mjContext
		if len(mjContext.GetWallCards()) == 0 {
			return true
		}
	}
	return false
}
