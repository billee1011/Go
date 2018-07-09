package fantype

import (
	majongpb "steve/server_pb/majong"
)

//checkLianLiu 检测连六 含有一种花色6张相连的序数牌
func checkLianLiu(tc *typeCalculator) bool {
	// 所有牌：手，胡,碰，杠，吃牌
	cardAll := getPlayerCardAll(tc)
	colorPointMap := make(map[majongpb.CardColor][]int32)
	for _, card := range cardAll {
		// 不字牌
		if !IsNotFlowerCard(card) {
			color := card.GetColor()
			colorPointMap[color] = append(colorPointMap[color], card.GetPoint())
		}
	}
	// 每种花色的序数牌中是否连六
	for _, points := range colorPointMap {
		newPoints := sortRemoveDuplicate(points)
		if l := len(newPoints); l >= 6 {
			// 从最大值递减下去，是否连续6张
			count := 1
			for i := l - 1; i > 0; i-- {
				if newPoints[i]-newPoints[i-1] == 1 {
					count++
				} else {
					count = 1
				}
			}
			if count >= 6 {
				return true
			}
		}
	}
	return false
}
