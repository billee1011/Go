package gutils

import (
	"math/rand"
	majongpb "steve/server_pb/majong"
	"time"
)

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
		return true, 13
	}
	return false, 18
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

//CopyColorCardMap 复制ColorCardMap
func CopyColorCardMap(currMap map[majongpb.CardColor][]*majongpb.Card) map[majongpb.CardColor][]*majongpb.Card {
	newMap := make(map[majongpb.CardColor][]*majongpb.Card)
	for key, value := range currMap {
		newMap[key] = value
	}
	return newMap
}

// NotDuplicatesCards 去除重复的牌
func NotDuplicatesCards(cards []*majongpb.Card) []*majongpb.Card {
	newCards := make([]*majongpb.Card, 0)
	for _, card := range cards {
		flag := true
		for _, nCard := range newCards {
			if CardEqual(card, nCard) {
				flag = false
				break
			}
		}
		if flag {
			newCards = append(newCards, card)
		}
	}
	return newCards
}

// CardTypeIsSame 牌型一样
func CardTypeIsSame(colors []majongpb.CardColor, colorCardsMap map[majongpb.CardColor][]*majongpb.Card) []*majongpb.Card {
	// 牌数不一样，选择最小牌数
	if flag, minCards := IsCardNumEqualAndMinCards(colors, colorCardsMap); !flag {
		return minCards
	}
	//牌数一样，随机
	rd := rand.New(rand.NewSource(time.Now().UnixNano())) // 随机出颜色
	towards := rd.Intn(len(colors))
	cards := colorCardsMap[colors[towards]]
	return cards
}
