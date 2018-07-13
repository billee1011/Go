package fantype

import majongpb "steve/server_pb/majong"

// checkMengQianQing 检测门前清 ：没有吃、碰、杠(暗杠可以)，不能是自摸
func checkMengQianQing(tc *typeCalculator) bool {
	huCard := tc.getHuCard()
	if huCard.GetType() == majongpb.HuType_hu_zimo {
		return false
	}
	pengCards := tc.getPengCards()
	chiCards := tc.getChiCards()
	if len(pengCards) != 0 || len(chiCards) != 0 {
		return false
	}
	for _, gang := range tc.getGangCards() {
		if gang.GetType() != majongpb.GangType_gang_angang {
			return false
		}
	}
	return true
}
