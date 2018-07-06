package fantype

import (
	"math"
	"steve/common/mjoption"
	"steve/majong/utils"
	"steve/server_pb/majong"
)

// Combine 牌型组合
type Combine struct {
	jiang int   // 将牌
	shuns []int // 顺牌
	kes   []int // 刻牌
}

func min(values ...int) int {
	t := values[0]
	for _, v := range values {
		if t > v {
			t = v
		}
	}
	return t
}

func minCard(values ...utils.Card) utils.Card {
	t := values[0]
	for _, v := range values {
		if t > v {
			t = v
		}
	}
	return t
}

// newCombines 根据 utils 的 Combine 创建 Combine 数组
func newCombines(combines utils.Combines) []Combine {
	result := make([]Combine, 0, len(combines))
	for _, combine := range combines {
		localCombine := Combine{
			shuns: []int{},
			kes:   []int{},
		}
		// TODO 先作一个能暂时满足需求的， 目前牌中没有癞子
		for _, group := range combine {
			if group.GroupType == utils.TypeJiang {
				localCombine.jiang = int(group.Cards[1])
			}
			if group.GroupType == utils.TypeShun {
				localCombine.shuns = append(localCombine.shuns, int(minCard(group.Cards...)))
			}
			if group.GroupType == utils.TypeKe {
				localCombine.kes = append(localCombine.kes, int(group.Cards[0]))
			}
		}
		result = append(result, localCombine)
	}
	return result
}

// typeCalculator 牌型计算器
type typeCalculator struct {
	mjContext *majong.MajongContext
	playerID  uint64
	handCards []*majong.Card // 手牌
	huCard    *majong.HuCard // 胡的牌

	combines []Combine
	player   *majong.Player // 玩家
	cache    map[int]bool   // 函数执行结果缓存, 函数ID->bool
}

// getOption 获取选项
func getOption(mjContext *majong.MajongContext) *mjoption.CardTypeOption {
	optionID := mjContext.GetCardtypeOptionId()
	return mjoption.GetCardTypeOption(int(optionID))
}

// getOption 获取牌型选项
func (tc *typeCalculator) getOption() *mjoption.CardTypeOption {
	return getOption(tc.mjContext)
}

func (tc *typeCalculator) makeCombines() {
	handCards := tc.getHandCards()
	cards := utils.ServerCards2UtilsCards(handCards)
	huCard := tc.getHuCard()
	if huCard != nil {
		card := utils.ServerCard2UtilCard(huCard.GetCard())
		cards = append(cards, card)
	}
	_, combines := utils.FastCheckHuV2(cards, nil, true)
	tc.combines = newCombines(combines)
}

// typeCalculator 计算出玩家胡牌所有的番型
func (tc *typeCalculator) calclate() (fantypes []int, gengcount int, huacount int) {
	tc.makeCombines()
	tc.cache = make(map[int]bool)

	fantypes = []int{}
	mutexs := map[int]struct{}{} // 当前哪些牌型被互斥了的

	option := tc.getOption()
	for ID, fantype := range option.Fantypes {
		// 已经互斥的不检测
		if _, ok := mutexs[ID]; ok {
			continue
		}
		match := tc.callCheckFunc(fantype.FuncID)
		if match {
			fantypes = append(fantypes, ID)
			for _, m := range fantype.Mutex {
				mutexs[m] = struct{}{}
			}
		}
	}
	// 移除掉互斥的番型
	tmpFantypes := []int{}
	for _, ID := range fantypes {
		if _, ok := mutexs[ID]; !ok {
			tmpFantypes = append(tmpFantypes, ID)
		}
	}
	fantypes = tmpFantypes

	if option.EnableGeng {
		gengcount = tc.calcGengCount(fantypes)
	}
	if option.EnableHua {
		huacount = tc.calcHuaCount()
	}
	return
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
// 返回值：
//	fantypes 番型列表
//	gengcount 根数量
//  huacount 花数量
func CalculateFanTypes(mjContext *majong.MajongContext, playerID uint64, handCards []*majong.Card, huCard *majong.HuCard) (fantypes []int, gengcount int, huacount int) {
	tc := typeCalculator{
		mjContext: mjContext,
		playerID:  playerID,
		handCards: handCards,
		huCard:    huCard,
	}
	return tc.calclate()
}

// CalculateScore 计算总番数
func CalculateScore(mjContext *majong.MajongContext, fantypes []int, gengcount int, huacount int) uint64 {
	var sumScore uint64     // 相加值
	var mulScore uint64 = 1 // 相乘值
	option := getOption(mjContext)

	if option.EnableGeng {
		calcScoreByMethod(option.GengMethod, &sumScore, &mulScore, gengcount, option.GengScore)
	}
	if option.EnableHua {
		calcScoreByMethod(option.HuaMethod, &sumScore, &mulScore, huacount, option.HuaScore)
	}
	for _, fan := range fantypes {
		fantype, ok := option.Fantypes[fan]
		if !ok {
			continue
		}
		calcScoreByMethod(fantype.Method, &sumScore, &mulScore, 1, fantype.Score)
	}
	if sumScore == 0 {
		sumScore = 1
	}
	return sumScore * mulScore
}

func calcScoreByMethod(method int, sumScore *uint64, mulScore *uint64, count int, score int) {
	if count == 0 {
		return
	}
	if method == 0 {
		*sumScore += uint64(score * count)
	} else if method == 1 {
		*mulScore *= uint64(score * count)
	} else if method == 2 {
		*mulScore *= uint64(math.Pow(float64(score), float64(count)))
	}
}
