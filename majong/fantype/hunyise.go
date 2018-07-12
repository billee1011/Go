package fantype

import majongpb "steve/server_pb/majong"

// checkHunYiSe 检测混一色
func checkHunYiSe(tc *typeCalculator) bool {
	checkCards := getPlayerCardAll(tc)

	//存在字牌
	existZi := false
	existXuShu := false
	cardColor := majongpb.CardColor(-1)
	for _, card := range checkCards {
		if !IsFlowerCard(card) {
			existZi = true
			continue
		}
		if cardColor == -1 {
			cardColor = card.Color
		} else if cardColor != card.Color {
			return false
		}
		existXuShu = true
	}
	return existZi && existXuShu
}
