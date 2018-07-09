package fantype

import (
	majongpb "steve/server_pb/majong"
)

// checkQuanDaiYao 检测全带幺 每副牌，将牌都有幺九(1,9,字牌)
func checkQuanDaiYao(tc *typeCalculator) bool {
	// 碰杠吃
	for _, peng := range tc.getPengCards() {
		pengCard := peng.GetCard()
		if !isYaoJiuByCard(pengCard) {
			return false
		}
	}
	for _, gang := range tc.getGangCards() {
		gangCard := gang.GetCard()
		if !isYaoJiuByCard(gangCard) {
			return false
		}
	}
	for _, chi := range tc.getChiCards() {
		chiCard := chi.GetOprCard()
		if chiCard.GetPoint() > 1 && chiCard.GetPoint() < 7 {
			return false
		}
	}
	// 手上的牌判断
Next:
	for _, combine := range tc.combines {
		// 将
		if !isYaoJiuByInt(combine.jiang) {
			continue
		}
		//顺
		for _, shun := range combine.shuns {
			shunValue := shun % 10
			if shunValue > 1 && shunValue < 7 {
				continue Next
			}
		}
		//刻
		for _, ke := range combine.kes {
			if !isYaoJiuByInt(ke) {
				continue Next
			}
		}
		return true
	}
	return false
}

// isYaoJiuByCard 判断是否是幺九(1,9,字)
func isYaoJiuByCard(card *majongpb.Card) bool {
	if !IsNotFlowerCard(card) {
		//又不是1和9的序数牌
		if card.GetPoint() > 1 && card.GetPoint() < 9 {
			return false
		}
	}
	return true
}
