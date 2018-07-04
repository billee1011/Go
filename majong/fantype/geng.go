package fantype

import "steve/majong/utils"

// calcGengCount 计算根的数量
func (tc *typeCalculator) calcGengCount(fantypes []int) int {
	total := tc.calcTotalGengCount()
	if total == 0 {
		return 0
	}
	subCount := tc.calcSubGengCount(fantypes)
	if total >= subCount {
		return total - subCount
	}
	return 0
}

// calcSubGengCount 计算要减去的根数量
func (tc *typeCalculator) calcSubGengCount(fantypes []int) int {
	subCount := 0
	option := tc.getOption()

	for _, fantype := range fantypes {
		fan := option.Fantypes[fantype]
		subCount += fan.SubGeng
	}
	return subCount
}

// calcTotalGengCount 计算总根数量
func (tc *typeCalculator) calcTotalGengCount() int {
	countMap := map[int]int{}

	handCards := tc.getHandCards()
	for _, card := range handCards {
		c := utils.ServerCard2Number(card)
		countMap[c]++
	}

	huCard := tc.getHuCard()
	if huCard != nil {
		c := utils.ServerCard2Number(huCard.GetCard())
		countMap[c]++
	}

	chiCards := tc.getChiCards()
	for _, chiCard := range chiCards {
		c := utils.ServerCard2Number(chiCard.GetCard())
		countMap[c]++
		countMap[c+1]++
		countMap[c+2]++
	}

	pengCards := tc.getPengCards()
	for _, pengCard := range pengCards {
		c := utils.ServerCard2Number(pengCard.GetCard())
		countMap[c] += 3
	}

	gangCards := tc.getGangCards()
	for _, gangCard := range gangCards {
		c := utils.ServerCard2Number(gangCard.GetCard())
		countMap[c] += 4
	}

	total := 0
	for _, count := range countMap {
		total += count / 4
	}
	return total
}
