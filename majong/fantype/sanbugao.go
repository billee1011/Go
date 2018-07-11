package fantype

import (
	majongpb "steve/server_pb/majong"
)

//checkSanBuGao 检测三步高 含有一种花色3副依次递增一位数或二位数的顺子
func checkSanBuGao(tc *typeCalculator) bool {
	for _, combine := range tc.combines {
		colorPointMap := make(map[majongpb.CardColor][]int32)
		// 吃
		for _, chi := range tc.getChiCards() {
			chiCard := chi.GetOprCard()
			colorPointMap[chiCard.GetColor()] = append(colorPointMap[chiCard.GetColor()], chiCard.GetPoint())
		}
		for _, shun := range combine.shuns {
			shunCard := intToCard(shun)
			colorPointMap[shunCard.GetColor()] = append(colorPointMap[shunCard.GetColor()], shunCard.GetPoint())
		}
		for _, cardPoints := range colorPointMap {
			if len(cardPoints) >= 3 {
				// 差值
				one, two := diff(cardPoints, 1), diff(cardPoints, 2)
				if one >= 3 || two >= 3 {
					return true
				}
			}
		}
	}
	return false
}
