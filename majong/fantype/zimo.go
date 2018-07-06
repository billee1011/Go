package fantype

import (
	majongpb "steve/server_pb/majong"
)

// checkZiMo 自摸
func checkZiMo(tc *typeCalculator) bool {
	if tc.huCard.Type == majongpb.HuType_hu_zimo {
		return true
	}
	return false
}
