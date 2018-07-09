package fantype

import (
	majongpb "steve/server_pb/majong"
)

// checkMengQianQing 检测门前清
func checkMengQianQing(tc *typeCalculator) bool {
	gangCards := tc.getGangCards()
	pengCards := tc.getPengCards()
	chiCards := tc.getChiCards()
	huCard := tc.getHuCard()

	if len(gangCards) != 0 || len(pengCards) != 0 || len(chiCards) != 0 {
		return false
	}

	dianPaoHu := map[majongpb.HuType]bool{
		majongpb.HuType_hu_dianpao:     true,
		majongpb.HuType_hu_ganghoupao:  true,
		majongpb.HuType_hu_qiangganghu: true,
	}
	return dianPaoHu[huCard.GetType()]
}
