package fantype

import majongpb "steve/server_pb/majong"

// checkAnGang 检测暗杠
func checkAnGang(tc *typeCalculator) bool {
	gangCards := tc.getGangCards()
	for _, gangCard := range gangCards {
		if gangCard.Type == majongpb.GangType_gang_angang {
			return true
		}
	}
	return false
}
