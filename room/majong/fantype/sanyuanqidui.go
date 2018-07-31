package fantype

import "steve/gutils"

//checkSanYuanQiDui 检测三元七对  胡牌为七对，并且包含“中发白”
func checkSanYuanQiDui(tc *typeCalculator) bool {
	// 是否七对
	if tc.callCheckFunc(qiduiFuncID) {
		handCards := huJoinHandCard(tc.getHandCards(), tc.getHuCard())
		//含“中发白”
		if len(getAssignCardMap(handCards, gutils.Zhong, gutils.Bai)) == 3 {
			return true
		}
	}
	return false
}
