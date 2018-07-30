package fantype

//checkDaSanFeng 检测大三风 含有3个风（东南西北）刻
func checkDaSanFeng(tc *typeCalculator) bool {
	//获取风牌映射，碰，杠，胡，手牌
	fengCardMap := getCardsToFengCardMap(tc)
	// 风牌种类必须>=3
	if len(fengCardMap) < 3 {
		return false
	}
	// 每种风刻数量为3的数量
	count := 0
	for _, num := range fengCardMap {
		if num == 3 {
			count++
		}
	}
	return count == 3
}
