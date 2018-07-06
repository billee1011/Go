package fantype

import (
	majongpb "steve/server_pb/majong"
)

// checkQiangGangHu 枪杠胡
func checkQiangGangHu(tc *typeCalculator) bool {
	return tc.huCard.Type == majongpb.HuType_hu_qiangganghu
}
