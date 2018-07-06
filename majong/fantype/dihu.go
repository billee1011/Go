package fantype

import (
	majongpb "steve/server_pb/majong"
)

// checkDiHu 地胡
func checkDiHu(tc *typeCalculator) bool {
	if tc.huCard.Type == majongpb.HuType_hu_dihu {
		return true
	}
	return false
}
