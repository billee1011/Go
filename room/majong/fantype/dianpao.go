package fantype

import (
	majongpb "steve/entity/majong"
)

// checkDianPao 点炮
func checkDianPao(tc *typeCalculator) bool {
	if tc.huCard.Type == majongpb.HuType_hu_dianpao {
		return true
	}
	return false
}
