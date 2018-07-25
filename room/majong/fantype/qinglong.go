package fantype

import (
	majongpb "steve/entity/majong"
)

//checkQingLong 检测清龙 含有一种花色1-9相连的序数牌
func checkQingLong(tc *typeCalculator) bool {
	// 所有牌
	cardAll := getPlayerCardAll(tc)
	// 颜色对应牌点数映射
	colorPointMap := make(map[majongpb.CardColor][]int32)
	for _, card := range cardAll {
		if IsXuShuCard(card) {
			colorPointMap[card.GetColor()] = append(colorPointMap[card.GetColor()], card.GetPoint())
		}
	}
	// 相同颜色的序数牌所有相连的牌数量为9张
	for _, points := range colorPointMap {
		if len(sortRemoveDuplicate(points)) == 9 {
			return true
		}
	}
	return false
}
