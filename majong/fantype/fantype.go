package fantype

import (
	"steve/common/mjoption"
	"steve/server_pb/majong"
)

// combine 牌型组合
type combine struct {
	jiang *majong.Card   // 将牌
	shuns []*majong.Card // 顺牌
	kes   []*majong.Card // 刻牌
}

// typeCalculator 牌型计算器
type typeCalculator struct {
	mjContext *majong.MajongContext
	playerID  uint64

	combines  []combine
	player    *majong.Player
	handCards []*majong.Card // 手牌
	huCard    *majong.HuCard
	cache     map[int]bool // 函数执行结果缓存, 函数ID->bool
}

// getOption 获取牌型选项
func (tc *typeCalculator) getOption(mjContext *majong.MajongContext) *mjoption.CardTypeOption {
	optionID := mjContext.GetCardtypeOptionId()
	return mjoption.GetCardTypeOption(int(optionID))
}

// makeCombines 计算所有组合
func (tc *typeCalculator) makeCombines() {
	tc.combines = []combine{}
	// TODO
}

// typeCalculator 计算出玩家胡牌所有的番型
func (tc *typeCalculator) calclate() []int {
	result := []int{}
	tc.cache = make(map[int]bool)
	tc.makeCombines()

	option := tc.getOption(tc.mjContext)
	for ID, fantype := range option.Fantypes {
		match := tc.callCheckFunc(fantype.FuncID)
		if match {
			result = append(result, ID)
		}
	}
	return result
}

// callCheckFunc 调用检测函数，如果有缓存，从缓存中取出结果
// 如果没有缓存，重新调用函数计算，并记录缓存
func (tc *typeCalculator) callCheckFunc(funcID int) bool {
	if result, ok := tc.cache[funcID]; ok {
		return result
	}
	f, ok := checkFuncs[funcID]
	if !ok {
		tc.cache[funcID] = false
		return false
	}
	tc.cache[funcID] = f(tc)
	return tc.cache[funcID]
}

// getPlayer 获取玩家
func (tc *typeCalculator) getPlayer() *majong.Player {
	if tc.player != nil {
		return tc.player
	}
	for _, player := range tc.mjContext.Players {
		if player.GetPalyerId() == tc.playerID {
			tc.player = player
			return player
		}
	}
	return nil
}

// getChiCards 获取吃的牌
func (tc *typeCalculator) getChiCards() []*majong.ChiCard {
	return tc.getPlayer().GetChiCards()
}

// getGangCards 获取杠的牌
func (tc *typeCalculator) getGangCards() []*majong.GangCard {
	return tc.getPlayer().GetGangCards()
}

// getPengCards 获取碰的牌
func (tc *typeCalculator) getPengCards() []*majong.PengCard {
	return tc.getPlayer().GetPengCards()
}

// getHandCards 获取手牌
func (tc *typeCalculator) getHandCards() []*majong.Card {
	if tc.handCards == nil {
		return tc.getPlayer().GetHandCards()
	}
	return tc.handCards
}

// getHuCard 获取胡的牌
func (tc *typeCalculator) getHuCard() *majong.HuCard {
	if tc.huCard == nil {
		huCards := tc.getPlayer().GetHuCards()
		if len(huCards) == 0 {
			return nil
		}
		return huCards[len(huCards)-1]
	}
	return tc.huCard
}

// CalculateFanTypes 计算番型
// handCards 手牌，如果为nil，则使用玩家自己的手牌
// huCard : 胡的牌，如果为 nil ，则使用玩家最后一次胡的牌
func CalculateFanTypes(mjContext *majong.MajongContext, playerID uint64, handCards []*majong.Card, huCard *majong.HuCard) []int {
	tc := typeCalculator{
		mjContext: mjContext,
		playerID:  playerID,
		handCards: handCards,
		huCard:    huCard,
	}
	return tc.calclate()
}
