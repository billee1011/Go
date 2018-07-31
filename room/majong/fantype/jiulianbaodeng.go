package fantype

//checkJiuLianBaoDeng 检查九莲宝灯,同种颜色的特定牌型 1112345678999 胡同花色的任意一张牌
func checkJiuLianBaoDeng(tc *typeCalculator) bool {
	handHuCards := huJoinHandCard(tc.getHandCards(), tc.getHuCard())
	// 牌必须都在手上
	if len(handHuCards) != 14 {
		return false
	}
	cardMap := make(map[int32]int)
	// 不能有字牌,只能有一种颜色
	intColor := handHuCards[0].GetColor() //初始颜色
	for _, card := range handHuCards {
		currColor := card.GetColor()
		// 牌不属于万也不属于条筒
		if !IsXuShuCard(card) {
			return false
		}
		if intColor != currColor {
			return false
		}
		cardMap[card.GetPoint()] = cardMap[card.GetPoint()] + 1
	}
	//必须有该颜色的所有序数牌1-9
	if len(cardMap) == 9 {
		for cardPoint, count := range cardMap {
			// 序数牌1和9不能小于3张
			if (cardPoint == 1 || cardPoint == 9) && count < 3 {
				return false
			}
		}
		return true
	}
	return false
}
