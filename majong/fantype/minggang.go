package fantype

import majongpb "steve/server_pb/majong"

// checkMingGang 检测明杠
func checkMingGang(tc *typeCalculator) bool {
	gangCards := tc.getGangCards()
	for _, gangCard := range gangCards {
		if gangCard.Type == majongpb.GangType_gang_bugang || gangCard.Type == majongpb.GangType_gang_minggang {
			return true
		}
	}
	return false
}
