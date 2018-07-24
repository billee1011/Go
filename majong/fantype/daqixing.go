package fantype

import (
	"steve/gutils"
	majongpb "steve/entity/majong"
)

// checkDaQiXing 大七星:胡牌为七对,并且由“东南西北中发白”其中的字牌构成
func checkDaQiXing(tc *typeCalculator) bool {
	if !tc.callCheckFunc(qiduiFuncID) {
		return false
	}
	cards := make([]*majongpb.Card, 0)
	cards = append(cards, tc.getHandCards()...)
	cards = append(cards, tc.getHuCard().GetCard())
	count := 0
	ziCards := []int{gutils.Dong, gutils.Nan, gutils.Xi, gutils.Bei, gutils.Zhong, gutils.Fa, gutils.Bai}
	for _, ziCard := range ziCards {
		for _, card := range cards {
			if IsXuShuCard(card) {
				return false
			}
			cardValue := gutils.ServerCard2Number(card)
			if cardValue == uint32(ziCard) {
				count++
				break
			}
		}
	}
	return count == 7
}
