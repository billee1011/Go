package fantype

import (
	majongpb "steve/server_pb/majong"
)

// checkDianPao 点炮
func checkDianPao(tc *typeCalculator) bool {
	if tc.huCard.Type == majongpb.HuType_hu_dianpao {
		return true
	}
	return false
}
