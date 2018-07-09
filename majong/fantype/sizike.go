package fantype

//checkSiZiKe 检测四字刻 含有4副字牌的刻或杠
func checkSiZiKe(tc *typeCalculator) bool {
	// 是碰碰胡的 不能有吃和顺子,做多只有4个组合，字牌不可能是顺
	if !tc.callCheckFunc(pengpenghuFuncID) {
		return false
	}
	// 所有牌
	cardAll := getPlayerCardAll(tc)
	pointCountMap := make(map[int32]int)
	for _, card := range cardAll {
		// 字牌
		if !IsNotFlowerCard(card) {
			pointCountMap[card.GetPoint()] = pointCountMap[card.GetPoint()] + 1
		}
	}
	// 每种字牌数量>=3的数量
	count := 0
	for _, num := range pointCountMap {
		if num >= 3 {
			count++
		}
	}
	return count == 4
}
