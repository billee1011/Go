package utils

import (
	"math/rand"
	"sort"
	majongpb "steve/server_pb/majong"
	"time"

	"github.com/Sirupsen/logrus"
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
	sortPriority, colorPrioMap := GetSortPrioAndColorPrioMapByColorCardMap(colorCardsMap)
	minPriority := sortPriority[len(sortPriority)-1] // 获取最小优先级 1为最大，18为最小优先级
	colors := GetColorByPriority(colorPrioMap, minPriority)
	if len(colors) != 1 {
		// 最小优先级不止一个的情况
		rd := rand.New(rand.NewSource(time.Now().UnixNano())) // 随机出颜色
		towards := rd.Intn(len(colors))
		return colors[towards]
	}
	return colors[0]
}

//IsCanGetColor 是否能获取颜色
func IsCanGetColor(colorCardsMap map[majongpb.CardColor][]*majongpb.Card) (bool, majongpb.CardColor) {
	min, mid, max := ColorSort(colorCardsMap)
	if min == 0 || mid-min > 1 { //最少花色数为0，或居中花色与最少的数差值大于1
		colors := GetCardColorsByLen(colorCardsMap, min)
		if len(colors) == 1 {
			return true, colors[0]
		}
		// 可能存在mid与min都为0的情况，随机
		rd := rand.New(rand.NewSource(time.Now().UnixNano())) // 随机出颜色
		towards := rd.Intn(len(colors))
		return true, colors[towards]
	}
	// 判断最多的花色是否参与比较
	if max-min > 1 {
		// 最多花色的牌，与最少的花色的牌数差值大于1，不比较
		colors := GetCardColorsByLen(colorCardsMap, max)
		delete(colorCardsMap, colors[0])
	}
	return false, majongpb.CardColor_ColorWan
}

//GetSortPrioAndColorPrioMapByColorCardMap 根据颜色与牌的映射，获取排序后的优先级，和优先级Map
func GetSortPrioAndColorPrioMapByColorCardMap(colorCardsMap map[majongpb.CardColor][]*majongpb.Card) ([]int, map[majongpb.CardColor]int) {
	// 优先级MAP
	colorPrioMap := make(map[majongpb.CardColor]int)
	sortPriority := make([]int, 0)
	for color, cards := range colorCardsMap {
		priority := GetPriorityByColorCard(cards)
		colorPrioMap[color] = priority
		sortPriority = append(sortPriority, priority)
		logrus.WithFields(logrus.Fields{"func_name": "GetSortPrioAndColorPrioMapByColorCardMap",
			"color": color, "priority": priority}).Info("获取推荐定缺颜色的优先级")
	}
	sort.Ints(sortPriority) // 升序，排序优先级
	return sortPriority, colorPrioMap
}

//GetPriorityByColorCard 获取优先级
func GetPriorityByColorCard(colorCards []*majongpb.Card) int {
	cardCountMap := make(map[int32]int)
	for _, card := range colorCards {
		cardCountMap[card.GetPoint()] = cardCountMap[card.GetPoint()] + 1
	}
	// 优先级从高到低查
	if flag, prio := CheckGangGroup(cardCountMap, len(colorCards)); flag {
		return prio
	}

	if flag, prio := CheckKeGroup(cardCountMap); flag {
		return prio
	}

	duiNum := GetAssignTypeNum(cardCountMap, 2)
	if duiNum >= 3 { //三对
		return 10
	}

	if flag, prio := CheckShunGroup(cardCountMap); flag {
		return prio
	}

	if flag, prio := CheckRemainGroup(duiNum, cardCountMap); flag {
		return prio
	}
	return 18 //单牌
}

// CheckRemainGroup 查剩余组合
func CheckRemainGroup(duiNum int, cardCountMap map[int32]int) (bool, int) {
	switch {
	case GetAssignTypeNum(cardCountMap, 3) > 0: //刻
		return true, 14
	case duiNum >= 2: //俩对
		return true, 15
	case len(GetMinShuns(cardCountMap)) > 0: //顺
		return true, 16
	case duiNum > 0: //对子
		return true, 17
	}
	return false, 18
}

//CheckGangGroup 查杠组合
func CheckGangGroup(cardCountMap map[int32]int, cardLen int) (bool, int) {
	for point, count := range cardCountMap {
		if count == 4 {
			newMap := CopyMap(cardCountMap)
			delete(newMap, point) // 删除杠
			switch {
			case cardLen == 4: //杠（只有四张）
				return true, 1
			case GetAssignTypeNum(newMap, 3) > 0: //杠+刻
				return true, 2
			case len(GetMinShuns(newMap)) > 0: //杠+顺
				return true, 3
			case GetAssignTypeNum(newMap, 2) > 0: //杠+对
				return true, 4
			default:
				return true, 5
			}
		}
	}
	return false, 18
}

// CheckKeGroup 查刻组合
func CheckKeGroup(cardCountMap map[int32]int) (bool, int) {
	for point, count := range cardCountMap {
		if count == 3 {
			newMap := CopyMap(cardCountMap)
			delete(newMap, point) // 删除当前刻
			switch {
			case GetAssignTypeNum(newMap, 3) > 0: //两刻
				return true, 6
			case GetAssignTypeNum(newMap, 2) >= 2: //刻+两对
				return true, 7
			case len(GetMinShuns(newMap)) > 0: // 刻+顺
				return true, 8
			case GetAssignTypeNum(newMap, 2) > 0: //刻+对
				return true, 9
			}
		}
	}
	return false, 18
}

//CheckShunGroup 查顺组合
func CheckShunGroup(cardCountMap map[int32]int) (bool, int) {
	duiCount := GetDuiNumByShunJiaDui(cardCountMap)
	if duiCount == 2 { //顺+2对
		return true, 11
	}
	shuns := GetMinShuns(cardCountMap)
	if len(shuns) >= 2 { //两顺
		return true, 12
	}
	if duiCount == 1 { //顺+对
		return true, 11
	}
	return false, 18
}

//GetColorStatistics 统计各花色
func GetColorStatistics(handCards []*majongpb.Card) map[majongpb.CardColor][]*majongpb.Card {
	colorCardsMap := make(map[majongpb.CardColor][]*majongpb.Card)
	for _, card := range handCards {
		if IsXuShuCard(card.GetColor()) {
			colorCardsMap[card.GetColor()] = append(colorCardsMap[card.GetColor()], card)
		}
	}
	return colorCardsMap
}

//ColorSort 对统计后的各个花色进行排序
func ColorSort(colorCardsMap map[majongpb.CardColor][]*majongpb.Card) (min, mid, max int) {
	wanLen := len(colorCardsMap[majongpb.CardColor_ColorWan])
	tiaoLen := len(colorCardsMap[majongpb.CardColor_ColorTiao])
	tongLen := len(colorCardsMap[majongpb.CardColor_ColorTong])
	// 获取 各花色的数量差异
	cardLen := []int{wanLen, tiaoLen, tongLen}
	sort.Ints(cardLen) // 升序
	logrus.WithFields(logrus.Fields{"func_name": "ColorSort", "cardLen": cardLen}).Info("获取推荐定缺颜色的长度")
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

// GetCardColorsByLen 获取指定长度的颜色
func GetCardColorsByLen(colorCardsMap map[majongpb.CardColor][]*majongpb.Card, cardLen int) []majongpb.CardColor {
	res := make([]majongpb.CardColor, 0)
	if len(colorCardsMap[majongpb.CardColor_ColorWan]) == cardLen {
		res = append(res, majongpb.CardColor_ColorWan)
	}
	if len(colorCardsMap[majongpb.CardColor_ColorTiao]) == cardLen {
		res = append(res, majongpb.CardColor_ColorTiao)
	}
	if len(colorCardsMap[majongpb.CardColor_ColorTong]) == cardLen {
		res = append(res, majongpb.CardColor_ColorTong)
	}
	return res
}

// IsXuShuCard 判断是否是序数牌
func IsXuShuCard(color majongpb.CardColor) bool {
	return color == majongpb.CardColor_ColorWan || color == majongpb.CardColor_ColorTiao || color == majongpb.CardColor_ColorTong
}

//GetAssignTypeNum cardType 4=gang 3==ke 2==dui,count=每种类型的数量
func GetAssignTypeNum(cardCountMap map[int32]int, num int) int {
	count := 0
	for _, sum := range cardCountMap {
		if sum == num {
			count++
		}
	}
	return count
}

//GetMinShuns 获取以最小值作为顺子
func GetMinShuns(cardCountMap map[int32]int) []int32 {
	newMap := CopyMap(cardCountMap)
	shun := make([]int32, 0)
	var minPoint int32 = 9
	for point := range newMap {
		if minPoint > point {
			minPoint = point
		}
	}
	if minPoint == 0 {
		return shun
	}
	point := minPoint
	for i := 0; i < 7; i++ {
		if newMap[point] > 0 && newMap[point+1] > 0 && newMap[point+2] > 0 {
			shun = append(shun, point) //存放顺子最大的值
			newMap[point] = newMap[point] - 1
			newMap[point+1] = newMap[point+1] - 1
			newMap[point+2] = newMap[point+2] - 1
		} else { // 不存在才查下一个
			point = point + 1
			if point > 8 {
				break
			}
		}
	}
	return shun
}

//CopyMap 复制Map
func CopyMap(currMap map[int32]int) map[int32]int {
	newMap := make(map[int32]int)
	for key, value := range currMap {
		newMap[key] = value
	}
	return newMap
}

//GetDuiNumByShunJiaDui 获取顺加对,对数量
func GetDuiNumByShunJiaDui(cardCountMap map[int32]int) int {
	duiCount := 0
	newMap := CopyMap(cardCountMap)
	for point, sum := range cardCountMap {
		if sum >= 2 {
			newMap[point] = newMap[point] - 2
			//除去对后，还能组成顺
			if newShuns := GetMinShuns(newMap); len(newShuns) > 0 {
				duiCount++
			} else {
				newMap[point] = newMap[point] + 2
			}
		}
	}
	return duiCount
}
