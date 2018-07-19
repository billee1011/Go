package fantype

import (
	majongpb "steve/server_pb/majong"
)

//checkSiBuGao 检测四步高 含有一种花色4副依次递增一位数或二位数的顺子,包括吃
func checkSiBuGao(tc *typeCalculator) bool {
	// 不能有碰杠
	if len(tc.getGangCards())+len(tc.getPengCards()) > 0 {
		return false
	}
	for _, combine := range tc.combines {
		// 刻子为0
		if len(combine.kes) != 0 {
			continue
		}
		colorPointMap := make(map[majongpb.CardColor][]int32)
		// 吃
		for _, chi := range tc.getChiCards() {
			chiCard := chi.GetCard()
			colorPointMap[chiCard.GetColor()] = append(colorPointMap[chiCard.GetColor()], chiCard.GetPoint())
		}
		for _, shun := range combine.shuns {
			shunCard := intToCard(shun)
			colorPointMap[shunCard.GetColor()] = append(colorPointMap[shunCard.GetColor()], shunCard.GetPoint())
		}
		if len(colorPointMap) > 1 {
			return false
		}
		for _, cardPoints := range colorPointMap {
			// 差值
			one, two := diff(cardPoints, 1), diff(cardPoints, 2)
			if one == 4 || two == 4 {
				return true
			}
		}
	}
	return false
}

func diff(cardPoints []int32, diff int32) int {
	cardPoints = sortRemoveDuplicate(cardPoints)
	count := 1
	// 每次的差值1的次数
	for i := len(cardPoints) - 1; i > 0; i-- {
		if cardPoints[i]-cardPoints[i-1] == diff {
			count++
		} else { //断了，重来
			count = 1
		}
	}
	return count
}
