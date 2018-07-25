package fantype

import majongpb "steve/entity/majong"

// checkShuanAnGang 检测双暗杠
func checkShuanAnGang(tc *typeCalculator) bool {
	count := 0
	for _, gangCard := range tc.getGangCards() {
		if gangCard.Type == majongpb.GangType_gang_angang {
			count++
		}
	}
	if count >= 2 {
		return true
	}
	return false
}
