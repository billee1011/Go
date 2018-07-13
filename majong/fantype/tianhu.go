package fantype

import majongpb "steve/server_pb/majong"

//checkTianHu 检测天胡 庄家在发完手牌后就胡牌，此为天胡，若庄家有补花，在补完花后就胡牌也算；若庄家在发完牌后有暗杠杠出，那么不算天胡；
//天胡根据自询次数，自询次数是1的是天胡
func checkTianHu(tc *typeCalculator) bool {
	huCard := tc.getHuCard()
	if huCard != nil && huCard.GetType() == majongpb.HuType_hu_zimo {
		currPlayer := tc.getPlayer()
		// 当前玩家是庄家
		mjContext := tc.mjContext
		if mjContext.Players[mjContext.ZhuangjiaIndex].GetPalyerId() == currPlayer.GetPalyerId() {
			if currPlayer.GetZixunCount() == 1 {
				return true
			}
		}
	}
	return false
}
