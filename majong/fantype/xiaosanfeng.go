package fantype

import (
	majongpb "steve/server_pb/majong"
)

//checkDaSanFeng 检测小三风 含有2个风牌的刻子或杠，以及1对风将牌
func checkXiaoSanFeng(tc *typeCalculator) bool {
	currCard := make([]*majongpb.Card, 0)
	// 杠是风
	currCard = append(currCard, gangToCards(tc.getGangCards())...)
	// 碰是风
	currCard = append(currCard, pengToCards(tc.getPengCards())...)
	fCardMap := getAssignCardMap(currCard, 41, 44)
	num := 0
	for _, count := range fCardMap {
		if count >= 3 {
			num++
		}
	}
	for _, combine := range tc.combines {
		// 将是风
		if combine.jiang > 40 && combine.jiang%10 <= 4 {
			count := 0
			// 刻子是风
			for _, ke := range combine.kes {
				if ke > 40 && ke%10 <= 4 {
					count++
				}
			}
			//含有2个风牌的刻子或杠
			if num+count >= 2 {
				return true
			}
		}
	}
	return false
}
