package fantype

import majongpb "steve/server_pb/majong"

// checkQiangGangHu 检测抢杠胡 其他玩家补杠，当前玩家抢补杠胡
func checkQiangGangHu(tc *typeCalculator) bool {
	huCard := tc.getHuCard()
	if huCard != nil && huCard.GetType() == majongpb.HuType_hu_qiangganghu {
		return true
	}
	return false
}
