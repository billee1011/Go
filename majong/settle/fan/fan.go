package fan

import (
	"steve/majong/interfaces"
	"steve/majong/utils"
	majongpb "steve/server_pb/majong"
)

// Fan 番型
type Fan struct {
	name      majongpb.CardType
	value     uint32
	Condition condition
}

type condition func(cardCalcParams interfaces.CardCalcParams) bool

// GetFanName 获取番型名字
func (f *Fan) GetFanName() majongpb.CardType {
	return f.name
}

// GetFanValue 获取番型倍数
func (f *Fan) GetFanValue() uint32 {
	return f.value
}

// ScxlFan 血流麻将所有番型
var ScxlFan []Fan

// AllFan 所有番型
var AllFan []Fan

// FanName 番型对应的名字
var FanName = map[majongpb.CardType]string{
	majongpb.CardType_PingHu:          "平胡",
	majongpb.CardType_QingYiSe:        "清一色",
	majongpb.CardType_QiDui:           "七对",
	majongpb.CardType_QingQiDui:       "清七对",
	majongpb.CardType_LongQiDui:       "龙七对",
	majongpb.CardType_QingLongQiDui:   "清龙七对",
	majongpb.CardType_PengPengHu:      "碰碰胡",
	majongpb.CardType_QingPeng:        "清碰",
	majongpb.CardType_JingGouDiao:     "金钩钓",
	majongpb.CardType_QingJingGouDiao: "清金钩钓",
	majongpb.CardType_ShiBaLuoHan:     "十八罗汉",
	majongpb.CardType_QingShiBaLuoHan: "清十八罗汉",
}

// FanValue 番型对应的倍数
var FanValue = map[majongpb.CardType]uint32{
	majongpb.CardType_PingHu:          1,
	majongpb.CardType_QingYiSe:        4,
	majongpb.CardType_QiDui:           4,
	majongpb.CardType_QingQiDui:       16,
	majongpb.CardType_LongQiDui:       8,
	majongpb.CardType_QingLongQiDui:   32,
	majongpb.CardType_PengPengHu:      2,
	majongpb.CardType_QingPeng:        8,
	majongpb.CardType_JingGouDiao:     4,
	majongpb.CardType_QingJingGouDiao: 16,
	majongpb.CardType_ShiBaLuoHan:     64,
	majongpb.CardType_QingShiBaLuoHan: 256,
}

func init() {

	fanPinghu := Fan{name: majongpb.CardType_PingHu, value: 1, Condition: func(cardCalcParams interfaces.CardCalcParams) bool {
		return checkPingHu(cardCalcParams)
	}}

	fanQingYiSe := Fan{name: majongpb.CardType_QingYiSe, value: 4, Condition: func(cardCalcParams interfaces.CardCalcParams) bool {
		return checkQingYiSe(cardCalcParams)
	}}

	fanQiDui := Fan{name: majongpb.CardType_QiDui, value: 4, Condition: func(cardCalcParams interfaces.CardCalcParams) bool {
		return checkQiDui(cardCalcParams)
	}}

	fanQingQiDui := Fan{name: majongpb.CardType_QingQiDui, value: 16, Condition: func(cardCalcParams interfaces.CardCalcParams) bool {
		return checkQingQiDui(cardCalcParams)
	}}

	fanQingLongQiDui := Fan{name: majongpb.CardType_QingLongQiDui, value: 32, Condition: func(cardCalcParams interfaces.CardCalcParams) bool {
		return checkQingLongQiDui(cardCalcParams)
	}}

	fanPengPengHu := Fan{name: majongpb.CardType_PengPengHu, value: 2, Condition: func(cardCalcParams interfaces.CardCalcParams) bool {
		return checkPengPengHu(cardCalcParams)
	}}

	fanQingPeng := Fan{name: majongpb.CardType_QingPeng, value: 8, Condition: func(cardCalcParams interfaces.CardCalcParams) bool {
		return checkQingPeng(cardCalcParams)
	}}

	fanJingGouDiao := Fan{name: majongpb.CardType_JingGouDiao, value: 4, Condition: func(cardCalcParams interfaces.CardCalcParams) bool {
		return checkJingGouDiao(cardCalcParams)
	}}

	fanQingJingGouDiao := Fan{name: majongpb.CardType_QingJingGouDiao, value: 16, Condition: func(cardCalcParams interfaces.CardCalcParams) bool {
		return checkQingJingGouDiao(cardCalcParams)
	}}

	fanShiBaLuoHan := Fan{name: majongpb.CardType_ShiBaLuoHan, value: 64, Condition: func(cardCalcParams interfaces.CardCalcParams) bool {
		return checkShiBaLuoHan(cardCalcParams)
	}}

	fanQingShiBaLuoHan := Fan{name: majongpb.CardType_QingShiBaLuoHan, value: 256, Condition: func(cardCalcParams interfaces.CardCalcParams) bool {
		return checkQingShiBaLuoHan(cardCalcParams)
	}}

	AllFan = append(AllFan, fanPinghu, fanQingYiSe, fanQiDui, fanQingQiDui, fanQingLongQiDui, fanPengPengHu, fanQingPeng, fanJingGouDiao, fanQingJingGouDiao, fanShiBaLuoHan, fanQingShiBaLuoHan)

	ScxlFan = append(ScxlFan, fanPinghu, fanQingYiSe, fanQiDui, fanPengPengHu, fanJingGouDiao, fanShiBaLuoHan)
}

// checkPingHu 平胡-不包含其他番型
func checkPingHu(cardCalcParams interfaces.CardCalcParams) bool {
	if checkQingYiSe(cardCalcParams) {
		return false
	} else if checkQiDui(cardCalcParams) {
		return false
	} else if checkPengPengHu(cardCalcParams) {
		return false
	} else if checkJingGouDiao(cardCalcParams) {
		return false
	} else if checkShiBaLuoHan(cardCalcParams) {
		return false
	}
	return true
}

// checkQingYiSe 清一色-所有牌同一花色
func checkQingYiSe(cardCalcParams interfaces.CardCalcParams) bool {
	checkCards := getCheckCards(cardCalcParams.HandCard, cardCalcParams.HuCard)
	for _, pengCard := range cardCalcParams.PengCard {
		checkCards = append(checkCards, pengCard)
	}
	for _, gangCard := range cardCalcParams.GangCard {
		checkCards = append(checkCards, gangCard)
	}
	color := majongpb.CardColor(-1)
	for _, card := range checkCards {
		if color == -1 {
			color = card.Color
		} else if color != card.Color {
			return false
		}
	}
	return true
}

// checkQiDui 七对-由七个对子组成--不能有碰杠
func checkQiDui(cardCalcParams interfaces.CardCalcParams) bool {
	checkCards := getCheckCards(cardCalcParams.HandCard, cardCalcParams.HuCard)

	if len(checkCards) != 14 {
		return false
	}
	cardCount := make(map[int32]int)
	for _, card := range checkCards {
		cardValue, _ := utils.CardToInt(*card)
		cardCount[*cardValue] = cardCount[*cardValue] + 1
	}
	for _, v := range cardCount {
		if v%2 != 0 {
			return false
		}
	}
	return true
}

// checkQingQiDui 清七对-清一色+七对
func checkQingQiDui(cardCalcParams interfaces.CardCalcParams) bool {
	if checkQiDui(cardCalcParams) && checkQingYiSe(cardCalcParams) {
		return true
	}
	return false
}

// checkLongQiDui 龙七对-至少有一个根的七对
func checkLongQiDui(cardCalcParams interfaces.CardCalcParams) bool {
	if checkQiDui(cardCalcParams) && GetGenCount(cardCalcParams) > 0 {
		return true
	}
	return false
}

// checkQingLongQiDui 清龙七对-清一色+龙七对
func checkQingLongQiDui(cardCalcParams interfaces.CardCalcParams) bool {
	if checkLongQiDui(cardCalcParams) && checkQingYiSe(cardCalcParams) {
		return true
	}
	return false
}

// checkPengPengHu 对对(碰碰)胡-刻子或碰或杠，加將牌组成就是没有顺子
func checkPengPengHu(cardCalcParams interfaces.CardCalcParams) bool {
	checkCards := getCheckCards(cardCalcParams.HandCard, cardCalcParams.HuCard)
	//开牌，即碰杠这些
	openCardSum := len(cardCalcParams.PengCard) + len(cardCalcParams.GangCard)
	if openCardSum >= 4 {
		return true
	}
	// 手牌中重复3个的个数
	handCardSum := 0
	cardCount := make(map[int32]int)
	for _, card := range checkCards {
		cardCount[card.Point] = cardCount[card.Point] + 1
	}
	cards := []int32{}
	for cardPoint, v := range cardCount {
		if v == 4 {
			return false
		} else if v == 3 {
			handCardSum++
		} else if v == 1 {
			cards = append(cards, cardPoint)
		}
	}
	if openCardSum+handCardSum >= 4 && len(cards) == 0 {
		return true
	}
	return false
}

// checkQingPeng 清碰碰胡-清一色+碰碰胡
func checkQingPeng(cardCalcParams interfaces.CardCalcParams) bool {
	if checkPengPengHu(cardCalcParams) && checkQingYiSe(cardCalcParams) {
		return true
	}
	return false
}

// checkJingGouDiao 金钩钓-胡牌时手里只剩一张，并且单钓一这张，其他的牌都被杠或碰了,不计碰碰胡。
func checkJingGouDiao(cardCalcParams interfaces.CardCalcParams) bool {
	handCards := cardCalcParams.HandCard
	if len(handCards) == 1 {
		huCard := cardCalcParams.HuCard
		if utils.CardEqual(huCard, handCards[0]) {
			if len(cardCalcParams.PengCard) != 0 {
				return true
			}
		}
	}
	return false
}

// checkQingJingGouDiao 清金钩钓-清一色+金钩钓
func checkQingJingGouDiao(cardCalcParams interfaces.CardCalcParams) bool {
	if checkJingGouDiao(cardCalcParams) && checkQingYiSe(cardCalcParams) {
		return true
	}
	return false
}

// checkShiBaLuoHan 十八罗汉-胡牌时手上只剩一张牌单吊，其他手牌形成四个杠，此时不计四根和碰碰胡。
func checkShiBaLuoHan(cardCalcParams interfaces.CardCalcParams) bool {
	l := len(cardCalcParams.GangCard)
	if l == 4 {
		return true
	}
	return false
}

//checkQingShiBaLuoHan 清十八罗汉-清一色+十八罗汉
func checkQingShiBaLuoHan(cardCalcParams interfaces.CardCalcParams) bool {
	if checkShiBaLuoHan(cardCalcParams) && checkQingYiSe(cardCalcParams) {
		return true
	}
	return false
}

// GetGenCount 获取玩家牌型根的数目
func GetGenCount(cardCalcParams interfaces.CardCalcParams) uint32 {
	var gCount uint32
	gangCards := cardCalcParams.GangCard
	pengCards := cardCalcParams.PengCard

	checkCards := getCheckCards(cardCalcParams.HandCard, cardCalcParams.HuCard)

	cardCount := make(map[int32]int)
	for _, card := range checkCards {
		cardValue, _ := utils.CardToInt(*card)
		cardCount[*cardValue] = cardCount[*cardValue] + 1
	}
	for card, sum := range cardCount {
		if sum >= 4 {
			gCount++
		} else if sum == 1 {
			for _, pengCard := range pengCards {
				cardValue, _ := utils.CardToInt(*pengCard)
				if *cardValue == card {
					gCount++
				}
			}
		}
	}
	gCount = gCount + uint32(len(gangCards))
	return gCount
}

// getCheckCards 获取校验的牌组
func getCheckCards(handCards []*majongpb.Card, huCard *majongpb.Card) []*majongpb.Card {
	checkCard := handCards
	if huCard != nil {
		checkCard = append(checkCard, huCard)
	}
	return checkCard
}

//ScxlFanMutex 番型和根处理
func ScxlFanMutex(fans []majongpb.CardType, gen uint32) ([]interfaces.CardType, uint32) {
	// 翻型只有1个，并且根为0,直接返回
	if len(fans) == 1 && gen == 0 {
		return []interfaces.CardType{interfaces.CardType(fans[0])}, 0
	}
	fansMap := make(map[majongpb.CardType]interfaces.CardType)
	for _, fanCardType := range fans {
		fansMap[fanCardType] = interfaces.CardType(fanCardType)
	}
	if value, ok := fansMap[majongpb.CardType_JingGouDiao]; ok && value > 0 { //金钩钓跟碰碰胡互斥
		delete(fansMap, majongpb.CardType_PengPengHu)
	}
	if value, ok := fansMap[majongpb.CardType_ShiBaLuoHan]; ok && value > 0 {
		if value, ok := fansMap[majongpb.CardType_JingGouDiao]; ok && value > 0 { //十八罗汉跟金钩钓互斥
			delete(fansMap, majongpb.CardType_JingGouDiao)
		}
		if value, ok := fansMap[majongpb.CardType_PengPengHu]; ok && value > 0 { //十八罗汉跟碰碰胡互斥
			delete(fansMap, majongpb.CardType_PengPengHu)
		}
		gen = 0
	}

	if value, ok := fansMap[majongpb.CardType_QingYiSe]; ok && value > 0 {
		flag := false
		if value, ok := fansMap[majongpb.CardType_ShiBaLuoHan]; ok && value > 0 { // 添加清十八罗汉,移除十八罗汉
			delete(fansMap, majongpb.CardType_ShiBaLuoHan)
			fansMap[majongpb.CardType_QingShiBaLuoHan] = interfaces.CardType(majongpb.CardType_QingShiBaLuoHan)
			flag = true
		}
		if value, ok := fansMap[majongpb.CardType_JingGouDiao]; ok && value > 0 { // 添加清金钩钓,移除金钩钓
			delete(fansMap, majongpb.CardType_JingGouDiao)
			fansMap[majongpb.CardType_QingJingGouDiao] = interfaces.CardType(majongpb.CardType_QingJingGouDiao)
			flag = true
		}
		if value, ok := fansMap[majongpb.CardType_PengPengHu]; ok && value > 0 { // 添加清碰,移除碰碰胡
			delete(fansMap, majongpb.CardType_PengPengHu)
			fansMap[majongpb.CardType_QingPeng] = interfaces.CardType(majongpb.CardType_QingPeng)
			flag = true
		}
		if value, ok := fansMap[majongpb.CardType_QiDui]; ok && value > 0 {
			delete(fansMap, majongpb.CardType_QiDui)
			if gen > 0 { // 添加清龙七对,移除七对/清七对/龙七对
				fansMap[majongpb.CardType_QingLongQiDui] = interfaces.CardType(majongpb.CardType_QingLongQiDui)
				gen--
			} else { // 添加清七对,移除七对
				fansMap[majongpb.CardType_QingQiDui] = interfaces.CardType(majongpb.CardType_QingQiDui)
			}
			flag = true
		}
		if flag { // 存在可以跟清一色可以合组的牌型，移除清一色
			delete(fansMap, majongpb.CardType_QingYiSe)
		}
	}
	if value, ok := fansMap[majongpb.CardType_QiDui]; ok && gen > 0 && value > 0 { // 龙七对
		delete(fansMap, majongpb.CardType_QiDui)                                                // 移除七对
		fansMap[majongpb.CardType_LongQiDui] = interfaces.CardType(majongpb.CardType_LongQiDui) // 添加龙七对
		gen--                                                                                   // 根减1
	}
	// fan的cardType，转为卡牌的cardType
	cardTypes := make([]interfaces.CardType, 0)
	for _, cardType := range fansMap {
		cardTypes = append(cardTypes, cardType)
	}
	return cardTypes, gen
}
