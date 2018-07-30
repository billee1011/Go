package fantype

import (
	majongpb "steve/entity/majong"
)

//checkSiAnKe 检查四暗刻 含有4个暗刻子或暗杠
func checkSiAnKe(tc *typeCalculator) bool {
	num := 0
	// 暗杠
	for _, gangCard := range tc.getGangCards() {
		if gangCard.GetType() == majongpb.GangType_gang_angang {
			num++
		}
	}
	// 暗刻
	for _, combine := range tc.combines {
		keLen := len(combine.kes)
		if num+keLen >= 4 {
			return true
		}
	}
	return false
}
