package fantype

import (
	"steve/gutils"
	majongpb "steve/server_pb/majong"
)

// checkHuJueZhang 检测胡绝张 胡牌池，桌面已亮明的3张牌所剩的第4张牌,抢扛胡不算
func checkHuJueZhang(tc *typeCalculator) bool {
	//抢杠胡不算
	huCard := tc.getHuCard()
	if huCard != nil && huCard.GetType() != majongpb.HuType_hu_qiangganghu {
		// 牌墙
		wall := tc.mjContext.GetWallCards()
		// 所有玩家手牌
		players := tc.mjContext.GetPlayers()
		cards := make([]*majongpb.Card, 0, len(wall))
		cards = append(cards, wall...)
		for _, palyer := range players {
			cards = append(cards, palyer.GetHandCards()...)
		}
		for _, card := range cards {
			if gutils.CardEqual(card, huCard.GetCard()) {
				return false
			}
		}
	}
	return true
}
