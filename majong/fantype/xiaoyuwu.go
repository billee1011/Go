package fantype

//checkXiaoYuWu 检查小于五 只能有序数牌并且序数牌<5
func checkXiaoYuWu(tc *typeCalculator) bool {
	// 所有牌
	cardAll := getPlayerCardAll(tc)
	for _, card := range cardAll {
		// 只能有序数牌并且序数牌<5
		if !IsFlowerCard(card) || card.GetPoint() >= 5 {
			return false
		}
	}
	return true
}
