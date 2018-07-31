package fantype

import (
	"steve/gutils"
	majongpb "steve/entity/majong"
)

// checkHuJueZhang 检测胡绝张 胡牌池，桌面已亮明(吃碰出)的3张牌所剩的第4张牌,抢扛胡不算
func checkHuJueZhang(tc *typeCalculator) bool {
	//抢杠胡不算
	huCard := tc.getHuCard()
	if huCard != nil && huCard.GetType() != majongpb.HuType_hu_qiangganghu {
		players := tc.mjContext.GetPlayers()
		hCard := huCard.GetCard()
		cards := make([]*majongpb.Card, 0)
		for _, player := range players {
			// 碰
			pengs := player.GetPengCards()
			for _, pengCard := range pengs {
				if gutils.CardEqual(pengCard.GetCard(), hCard) {
					return true
				}
			}
			// 出
			cards = append(cards, player.GetOutCards()...)
			// 吃
			chis := player.GetChiCards()
			for _, chi := range chis {
				minCard := chi.GetCard() //最小牌
				midCard := &majongpb.Card{
					Color: minCard.GetColor(),
					Point: minCard.GetPoint() + 1,
				}
				maxCard := &majongpb.Card{
					Color: minCard.GetColor(),
					Point: minCard.GetPoint() + 2,
				}
				cards = append(cards, minCard, midCard, maxCard)
			}
		}
		return CheckAssignCardNum(cards, hCard, 3)
	}
	return false
}

//CheckAssignCardNum 检测指定牌数量
func CheckAssignCardNum(cards []*majongpb.Card, assCard *majongpb.Card, num int) bool {
	count := 0
	for _, card := range cards {
		if gutils.CardEqual(card, assCard) {
			count++
		}
	}
	if count == num {
		return true
	}
	return false
}
