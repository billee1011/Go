package fantype

import (
	"steve/gutils"
	majongpb "steve/server_pb/majong"
)

// checkDaSiXi 检查大四喜 含有“东南西北”4副风刻或杠
func checkDaSiXi(tc *typeCalculator) bool {
	// 是碰碰胡的 不能有吃和顺子，做多只有4个组合，风牌不可能是顺
	if !tc.callCheckFunc(pengpenghuFuncID) {
		return false
	}
	fCardCountMap := getCardsToFengCardMap(tc)
	// 必须都有东南西北
	if len(fCardCountMap) < 4 {
		return false
	}
	// 每种必须>=3
	for _, cardNum := range fCardCountMap {
		if cardNum < 3 {
			return false
		}
	}
	return true
}

// getCardsToFengCardMap 杠，碰，手牌转风牌数量映射
func getCardsToFengCardMap(tc *typeCalculator) map[uint32]int {
	currCard := make([]*majongpb.Card, 0)
	// 杠
	currCard = append(currCard, gangToCards(tc.getGangCards())...)
	// 碰
	currCard = append(currCard, pengToCards(tc.getPengCards())...)
	// 手，胡牌
	currCard = append(currCard, huJoinHandCard(tc.getHandCards(), tc.getHuCard())...)
	return getAssignCardMap(currCard, gutils.Zhong, gutils.Bai)
}
