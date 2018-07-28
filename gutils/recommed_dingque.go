package gutils

import (
	"sort"
	majongpb "steve/server_pb/majong"
)

//GetRecommedDingQueColor 获取推荐定却颜色 牌数最少，优先级最低
func GetRecommedDingQueColor(handCards []*majongpb.Card) majongpb.CardColor {
	colorCardsMap := GetColorStatistics(handCards)
	//是否能获取颜色
	flag, color := IsCanGetColor(colorCardsMap)
	if flag {
		return color
	}
	// 获取每种比较的颜色的优先级
	sortPriority, colorPrioMap := GetColorPriorityInfo(colorCardsMap)
	minPriority := sortPriority[len(sortPriority)-1] // 获取最小优先级 1为最大，18为最小优先级
	colors := GetColorByPriority(colorPrioMap, minPriority)
	if len(colors) == 1 {
		return colors[0]
	}
	currCards := CardTypeIsSame(colors, colorCardsMap)
	if len(currCards) <= 0 {
		return majongpb.CardColor_ColorWan
	}
	return currCards[0].GetColor()
}

//GetColorPriorityInfo 获取颜色优先级信息
func GetColorPriorityInfo(colorCardsMap map[majongpb.CardColor][]*majongpb.Card) ([]int, map[majongpb.CardColor]int) {
	// 优先级MAP
	colorPrioMap := make(map[majongpb.CardColor]int)
	sortPriority := make([]int, 0)
	for color, cards := range colorCardsMap {
		priority := GetPriorityByColorCard(cards)
		colorPrioMap[color] = priority
		sortPriority = append(sortPriority, priority)
	}
	sort.Ints(sortPriority) // 升序，排序优先级
	return sortPriority, colorPrioMap
}

//ColorSort 对统计后的各个花色进行排序
func ColorSort(colorCardsMap map[majongpb.CardColor][]*majongpb.Card) (min, mid, max int) {
	wanLen := len(colorCardsMap[majongpb.CardColor_ColorWan])
	tiaoLen := len(colorCardsMap[majongpb.CardColor_ColorTiao])
	tongLen := len(colorCardsMap[majongpb.CardColor_ColorTong])
	// 获取 各花色的数量差异
	cardLen := []int{wanLen, tiaoLen, tongLen}
	sort.Ints(cardLen) // 升序
	return cardLen[0], cardLen[1], cardLen[2]
}

// GetColorByPriority 根据优先级获取颜色
func GetColorByPriority(colorPrioMap map[majongpb.CardColor]int, currPrio int) []majongpb.CardColor {
	colors := make([]majongpb.CardColor, 0)
	for color, priority := range colorPrioMap {
		if currPrio == priority {
			colors = append(colors, color)
		}
	}
	return colors
}
