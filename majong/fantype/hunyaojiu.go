package fantype

import "steve/gutils"

// checkHunYaoJiu 检测混幺九 由序数牌1,9和字牌的刻子，将牌组成
func checkHunYaoJiu(tc *typeCalculator) bool {
	// 吃，杠数量
	chiGangNum := len(tc.getChiCards()) + len(tc.getGangCards())
	if chiGangNum != 0 {
		return false
	}
	existOne, existNine, existZi := false, false, false
	cardAll := getPlayerCardAll(tc)
	for _, card := range cardAll {
		// 不能存在序数2-8
		if !isYaoJiuByCard(card) {
			return false
		}
		cardValue := gutils.ServerCard2Number(card)
		if cardValue >= uint32(gutils.Dong) {
			existZi = true
		} else {
			if card.GetPoint() == 9 {
				existNine = true
			}
			if card.GetPoint() == 1 {
				existOne = true
			}
		}
	}
	// 必须存在1,9和字牌
	if existOne && existNine && existZi {
		return true
	}
	return false
}

//isYaoJiuByInt 判断是否是幺九(1,9,字)
func isYaoJiuByInt(card int) bool {
	if card < gutils.Dong {
		cardValue := card % 10
		if cardValue > 1 && cardValue < 9 {
			return false
		}
	}
	return true
}
