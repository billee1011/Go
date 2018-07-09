package fantype

import majongpb "steve/server_pb/majong"

// checkRenHu 检测人胡,庄家打出的第一张牌闲家就胡牌，此为人胡，若庄家出牌前有暗杠，那么不算人胡；
func checkRenHu(tc *typeCalculator) bool {
	huCard := tc.getHuCard()
	if huCard != nil && huCard.GetType() == majongpb.HuType_hu_dianpao {
		mjContext := tc.mjContext
		zjPlayer := mjContext.Players[mjContext.ZhuangjiaIndex]
		currPlayer := tc.getPlayer()
		// 不是庄家
		if zjPlayer.GetPalyerId() != currPlayer.GetPalyerId() {
			if zjPlayer.GetZixunCount() == 1 && tc.getPlayer().GetMopaiCount() == 0 {
				return true
			}
		}
	}
	return false
}
