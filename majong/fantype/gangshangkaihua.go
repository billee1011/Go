package fantype

import (
	majongpb "steve/server_pb/majong"
)

// checkGangShangKaiHua 杠上开花
func checkGangShangKaiHua(tc *typeCalculator) bool {
	return tc.huCard.Type == majongpb.HuType_hu_gangkai
}
