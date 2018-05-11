package fan

import (
	"steve/majong/utils"
	"steve/server_pb/majong"
)

// Fan 番型
type Fan struct {
	name      string
	value     uint32
	Condition condition
}

type condition func(context majong.MajongContext, huType majong.HuType, player *majong.Player) bool

// GetFanName 获取番型名字
func (f *Fan) GetFanName() string {
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
var FanName = map[majong.Fan]string{
	majong.Fan_PingHu:          "平胡",
	majong.Fan_FZiMo:           "自摸",
	majong.Fan_QiangGangHu:     "抢杠胡",
	majong.Fan_DianPaoHu:       "点炮",
	majong.Fan_GangHouPao:      "杠后炮",
	majong.Fan_GangKai:         "杠上开花",
	majong.Fan_HaiDiLao:        "海底捞",
	majong.Fan_QingYiSe:        "清一色",
	majong.Fan_QiDui:           "七对",
	majong.Fan_QingQiDui:       "清七对",
	majong.Fan_LongQiDui:       "龙七对",
	majong.Fan_QingLongQiDui:   "清龙七对",
	majong.Fan_PengPengHu:      "碰碰胡",
	majong.Fan_QingPeng:        "清碰",
	majong.Fan_JingGouDiao:     "金钩钓",
	majong.Fan_QingJingGouDiao: "清金钩钓",
	majong.Fan_TianHu:          "天胡",
	majong.Fan_DiHu:            "地胡",
	majong.Fan_ShiBaLuoHan:     "十八罗汉",
	majong.Fan_QingShiBaLuoHan: "清十八罗汉",
}

// FanValue 番型对应的倍数
var FanValue = map[majong.Fan]uint32{
	majong.Fan_PingHu:          1,
	majong.Fan_FZiMo:           2,
	majong.Fan_QiangGangHu:     2,
	majong.Fan_DianPaoHu:       1,
	majong.Fan_GangHouPao:      2,
	majong.Fan_GangKai:         2,
	majong.Fan_HaiDiLao:        2,
	majong.Fan_QingYiSe:        4,
	majong.Fan_QiDui:           4,
	majong.Fan_QingQiDui:       16,
	majong.Fan_LongQiDui:       8,
	majong.Fan_QingLongQiDui:   32,
	majong.Fan_PengPengHu:      2,
	majong.Fan_QingPeng:        8,
	majong.Fan_JingGouDiao:     4,
	majong.Fan_QingJingGouDiao: 16,
	majong.Fan_TianHu:          32,
	majong.Fan_DiHu:            32,
	majong.Fan_ShiBaLuoHan:     64,
	majong.Fan_QingShiBaLuoHan: 256,
}

func init() {

	fanPinghu := Fan{name: FanName[majong.Fan_PingHu], value: 1, Condition: func(context majong.MajongContext, huType majong.HuType, player *majong.Player) bool {
		return checkPingHu(context, huType, player)
	}}

	fanZimo := Fan{name: FanName[majong.Fan_FZiMo], value: 2, Condition: func(context majong.MajongContext, huType majong.HuType, player *majong.Player) bool {
		return checkZiMo(context, huType, player)
	}}

	fanQiangGangHu := Fan{name: FanName[majong.Fan_QiangGangHu], value: 2, Condition: func(context majong.MajongContext, huType majong.HuType, player *majong.Player) bool {
		return checkQiangGangHu(context, huType, player)
	}}

	fanDianPaoHu := Fan{name: FanName[majong.Fan_DianPaoHu], value: 1, Condition: func(context majong.MajongContext, huType majong.HuType, player *majong.Player) bool {
		return checkDianPao(context, huType, player)
	}}

	fanGangHouPao := Fan{name: FanName[majong.Fan_GangHouPao], value: 2, Condition: func(context majong.MajongContext, huType majong.HuType, player *majong.Player) bool {
		return checkGangHouPao(context, huType, player)
	}}

	fanGangKai := Fan{name: FanName[majong.Fan_GangKai], value: 2, Condition: func(context majong.MajongContext, huType majong.HuType, player *majong.Player) bool {
		return checkGangKai(context, huType, player)
	}}

	fanHaiDiLao := Fan{name: FanName[majong.Fan_HaiDiLao], value: 2, Condition: func(context majong.MajongContext, huType majong.HuType, player *majong.Player) bool {
		return checkHaiDiLao(context, huType, player)
	}}

	fanQingYiSe := Fan{name: FanName[majong.Fan_QingYiSe], value: 4, Condition: func(context majong.MajongContext, huType majong.HuType, player *majong.Player) bool {
		return checkQingYiSe(context, huType, player)
	}}

	fanQiDui := Fan{name: FanName[majong.Fan_QiDui], value: 4, Condition: func(context majong.MajongContext, huType majong.HuType, player *majong.Player) bool {
		return checkQiDui(context, huType, player)
	}}

	fanQingQiDui := Fan{name: FanName[majong.Fan_QingQiDui], value: 16, Condition: func(context majong.MajongContext, huType majong.HuType, player *majong.Player) bool {
		return checkQingQiDui(context, huType, player)
	}}

	fanLongQiDui := Fan{name: FanName[majong.Fan_LongQiDui], value: 8, Condition: func(context majong.MajongContext, huType majong.HuType, player *majong.Player) bool {
		return checkLongQiDui(context, huType, player)
	}}

	fanQingLongQiDui := Fan{name: FanName[majong.Fan_QingLongQiDui], value: 32, Condition: func(context majong.MajongContext, huType majong.HuType, player *majong.Player) bool {
		return checkQingLongQiDui(context, huType, player)
	}}

	fanPengPengHu := Fan{name: FanName[majong.Fan_PengPengHu], value: 2, Condition: func(context majong.MajongContext, huType majong.HuType, player *majong.Player) bool {
		return checkPengPengHu(context, huType, player)
	}}

	fanQingPeng := Fan{name: FanName[majong.Fan_QingPeng], value: 8, Condition: func(context majong.MajongContext, huType majong.HuType, player *majong.Player) bool {
		return checkQingPeng(context, huType, player)
	}}

	fanJingGouDiao := Fan{name: FanName[majong.Fan_JingGouDiao], value: 4, Condition: func(context majong.MajongContext, huType majong.HuType, player *majong.Player) bool {
		return checkJingGouDiao(context, huType, player)
	}}

	fanQingJingGouDiao := Fan{name: FanName[majong.Fan_QingJingGouDiao], value: 16, Condition: func(context majong.MajongContext, huType majong.HuType, player *majong.Player) bool {
		return checkQingJingGouDiao(context, huType, player)
	}}

	fanTianHu := Fan{name: FanName[majong.Fan_TianHu], value: 32, Condition: func(context majong.MajongContext, huType majong.HuType, player *majong.Player) bool {
		return checkTianHu(context, huType, player)
	}}

	fanDiHu := Fan{name: FanName[majong.Fan_DiHu], value: 32, Condition: func(context majong.MajongContext, huType majong.HuType, player *majong.Player) bool {
		return checkDiHu(context, huType, player)
	}}

	fanShiBaLuoHan := Fan{name: FanName[majong.Fan_ShiBaLuoHan], value: 64, Condition: func(context majong.MajongContext, huType majong.HuType, player *majong.Player) bool {
		return checkShiBaLuoHan(context, huType, player)
	}}

	fanQingShiBaLuoHan := Fan{name: FanName[majong.Fan_QingShiBaLuoHan], value: 256, Condition: func(context majong.MajongContext, huType majong.HuType, player *majong.Player) bool {
		return checkQingShiBaLuoHan(context, huType, player)
	}}

	AllFan = append(AllFan, fanPinghu, fanZimo, fanQiangGangHu, fanDianPaoHu, fanGangHouPao, fanGangKai, fanHaiDiLao, fanQingYiSe, fanQiDui, fanQingQiDui, fanLongQiDui, fanQingLongQiDui, fanPengPengHu, fanQingPeng, fanJingGouDiao, fanQingJingGouDiao, fanTianHu, fanDiHu, fanShiBaLuoHan, fanQingShiBaLuoHan)

	ScxlFan = append(ScxlFan, fanPinghu, fanQiangGangHu, fanGangHouPao, fanGangKai, fanHaiDiLao, fanQingYiSe, fanQiDui, fanLongQiDui, fanPengPengHu, fanJingGouDiao, fanTianHu, fanDiHu, fanShiBaLuoHan)
}

// checkPingHu 平胡-不包含其他番型
func checkPingHu(context majong.MajongContext, huType majong.HuType, player *majong.Player) bool {
	if checkQingYiSe(context, huType, player) {
		return false
	} else if checkQiDui(context, huType, player) {
		return false
	} else if checkPengPengHu(context, huType, player) {
		return false
	} else if checkJingGouDiao(context, huType, player) {
		return false
	} else if checkShiBaLuoHan(context, huType, player) {
		return false
	}
	return true
}

// checkZiMo 自摸
func checkZiMo(context majong.MajongContext, huType majong.HuType, player *majong.Player) bool {
	if huType == majong.HuType_hu_zimo {
		return true
	}
	return false
}

// checkQiangGangHu 抢杠胡-输家补杠时，赢抢杠胡
func checkQiangGangHu(context majong.MajongContext, huType majong.HuType, player *majong.Player) bool {
	if huType == majong.HuType_hu_qiangganghu {
		return true
	}
	return false
}

// checkDianPao 点炮
func checkDianPao(context majong.MajongContext, huType majong.HuType, player *majong.Player) bool {
	if huType == majong.HuType_hu_dianpao {
		return true
	}
	return false
}

// checkGangHouPao 杠后炮-其他玩家杠后摸牌出牌后被赢家点炮胡
func checkGangHouPao(context majong.MajongContext, huType majong.HuType, player *majong.Player) bool {
	if huType == majong.HuType_hu_ganghoupao {
		return true
	}
	return false
}

// checkGangKai 杠上开花-杠后摸牌胡
func checkGangKai(context majong.MajongContext, huType majong.HuType, player *majong.Player) bool {
	if huType == majong.HuType_hu_gangkai {
		return true
	}
	return false
}

// checkQiDui 海底捞-自摸胡的是最后一张
func checkHaiDiLao(context majong.MajongContext, huType majong.HuType, player *majong.Player) bool {
	if (huType == majong.HuType_hu_zimo) && len(context.WallCards) == 0 {
		return true
	}
	return false
}

// checkQingYiSe 清一色-所有牌同一花色
func checkQingYiSe(context majong.MajongContext, huType majong.HuType, player *majong.Player) bool {
	checkCards := getCheckCards(player.HandCards, player.HuCards)
	for _, pengCard := range player.PengCards {
		checkCards = append(checkCards, pengCard.Card)
	}
	for _, gangCard := range player.GangCards {
		checkCards = append(checkCards, gangCard.Card)
	}
	color := majong.CardColor(-1)
	for _, card := range checkCards {
		if color == 1 {
			color = card.Color
		} else if color != card.Color {
			return false
		}
	}
	return true
}

// checkQiDui 七对-由七个对子组成--不能有碰杠
func checkQiDui(context majong.MajongContext, huType majong.HuType, player *majong.Player) bool {
	checkCards := getCheckCards(player.HandCards, player.HuCards)

	if len(checkCards) != 14 {
		return false
	}
	cardCount := make(map[int32]int)
	for _, card := range checkCards {
		cardCount[card.Point] = cardCount[card.Point] + 1
	}
	for _, v := range cardCount {
		if v%2 != 0 {
			return false
		}
	}
	return true
}

// checkQingQiDui 清七对-清一色+七对
func checkQingQiDui(context majong.MajongContext, huType majong.HuType, player *majong.Player) bool {
	if checkQiDui(context, huType, player) && checkQingYiSe(context, huType, player) {
		return true
	}
	return false
}

// checkLongQiDui 龙七对-至少有一个根的七对
func checkLongQiDui(context majong.MajongContext, huType majong.HuType, player *majong.Player) bool {
	if checkQiDui(context, huType, player) && GetGenCount(player) > 0 {
		return true
	}
	return false
}

// checkQingLongQiDui 清龙七对-清一色+龙七对
func checkQingLongQiDui(context majong.MajongContext, huType majong.HuType, player *majong.Player) bool {
	if checkLongQiDui(context, huType, player) && checkQingYiSe(context, huType, player) {
		return true
	}
	return false
}

// checkPengPengHu 对对(碰碰)胡-刻子或碰或杠，加將牌组成就是没有顺子
func checkPengPengHu(context majong.MajongContext, huType majong.HuType, player *majong.Player) bool {
	checkCards := getCheckCards(player.HandCards, player.HuCards)
	//开牌，即碰杠这些
	openCardSum := len(player.PengCards) + len(player.GangCards)
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
func checkQingPeng(context majong.MajongContext, huType majong.HuType, player *majong.Player) bool {
	if checkPengPengHu(context, huType, player) && checkQingYiSe(context, huType, player) {
		return true
	}
	return false
}

// checkJingGouDiao 金钩钓-胡牌时手里只剩一张，并且单钓一这张，其他的牌都被杠或碰了,不计碰碰胡。
func checkJingGouDiao(context majong.MajongContext, huType majong.HuType, player *majong.Player) bool {
	handCards := player.HandCards
	if len(handCards) == 1 {
		huCard := player.HuCards[len(player.HuCards)-1]
		if utils.CardEqual(huCard.Card, handCards[0]) {
			if len(player.PengCards) != 0 {
				return true
			}
		}
	}
	return false
}

// checkQingJingGouDiao 清金钩钓-清一色+金钩钓
func checkQingJingGouDiao(context majong.MajongContext, huType majong.HuType, player *majong.Player) bool {
	if checkJingGouDiao(context, huType, player) && checkQingYiSe(context, huType, player) {
		return true
	}
	return false
}

// checkTianHu 天胡-庄家发完牌自摸
func checkTianHu(context majong.MajongContext, huType majong.HuType, player *majong.Player) bool {
	if len(context.WallCards) == (108-14-(len(context.Players)-1)*13) && huType == majong.HuType_hu_zimo {
		return true
	}
	return false
}

// checkDiHu 地胡-闲家摸第一张牌自摸
func checkDiHu(context majong.MajongContext, huType majong.HuType, player *majong.Player) bool {
	if len(player.PengCards) != 0 || len(player.GangCards) != 0 {
		return false
	}
	if len(player.OutCards) == 0 && huType == majong.HuType_hu_zimo {
		return true
	}
	return false
}

// checkShiBaLuoHan 十八罗汉-胡牌时手上只剩一张牌单吊，其他手牌形成四个杠，此时不计四根和碰碰胡。
func checkShiBaLuoHan(context majong.MajongContext, huType majong.HuType, player *majong.Player) bool {
	l := len(player.GangCards)
	if l == 4 {
		return true
	}
	return false
}

//checkQingShiBaLuoHan 清十八罗汉-清一色+十八罗汉
func checkQingShiBaLuoHan(context majong.MajongContext, huType majong.HuType, player *majong.Player) bool {
	if checkShiBaLuoHan(context, huType, player) && checkQingYiSe(context, huType, player) {
		return true
	}
	return false
}

// GetGenCount 获取玩家牌型根的数目
func GetGenCount(player *majong.Player) uint32 {
	var gCount uint32
	gangCards := player.GangCards
	pengCards := player.PengCards

	checkCards := getCheckCards(player.HandCards, player.HuCards)

	cardCount := make(map[int32]int)
	for _, card := range checkCards {
		cardCount[card.Point] = cardCount[card.Point] + 1
	}
	for card, sum := range cardCount {
		if sum >= 4 {
			gCount++
		} else if sum == 1 {
			for _, pengCard := range pengCards {
				if pengCard.Card.Point == card {
					gCount++
				}
			}
		}
	}
	gCount = gCount + uint32(len(gangCards))
	return gCount
}

// getCheckCards 获取校验的牌组
func getCheckCards(handCards []*majong.Card, huCards []*majong.HuCard) []*majong.Card {
	checkCard := handCards
	if len(huCards) > 0 {
		checkCard = append(checkCard, huCards[len(huCards)-1].Card)
	}
	return checkCard
}
