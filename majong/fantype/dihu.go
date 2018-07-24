package fantype

import (
	majongpb "steve/entity/majong"
)

//checkDiHu 检测地胡 闲家摸到第一张牌就胡牌，此为地胡，若闲家抓的第一张牌是花牌，那么补花之后胡牌也算地胡；若闲家抓牌前有人吃碰杠（包括暗杠），那么不算地胡
func checkDiHu(tc *typeCalculator) bool {
	huCard := tc.getHuCard()
	if huCard != nil && huCard.GetType() == majongpb.HuType_hu_zimo {
		currPlayer := tc.getPlayer()
		// 当前玩家不是庄家
		mjContext := tc.mjContext
		if mjContext.Players[mjContext.ZhuangjiaIndex].GetPalyerId() != currPlayer.GetPalyerId() {
			if currPlayer.GetZixunCount() == 1 && len(currPlayer.GetGangCards())+len(currPlayer.GetHuCards()) == 0 {
				return true
			}
		}
	}
	return false
}
