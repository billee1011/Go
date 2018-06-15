package fan

import (
	"steve/majong/interfaces"
	majongpb "steve/server_pb/majong"
)

// ScxlFan 血流麻将所有番型
var ScxlFan []Fan

//ScxlFanValue 番型对应的倍数
var ScxlFanValue = map[majongpb.CardType]uint32{
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

	fanPengPengHu := Fan{name: majongpb.CardType_PengPengHu, value: 2, Condition: func(cardCalcParams interfaces.CardCalcParams) bool {
		return checkPengPengHu(cardCalcParams)
	}}

	fanJingGouDiao := Fan{name: majongpb.CardType_JingGouDiao, value: 4, Condition: func(cardCalcParams interfaces.CardCalcParams) bool {
		return checkJingGouDiao(cardCalcParams)
	}}

	fanShiBaLuoHan := Fan{name: majongpb.CardType_ShiBaLuoHan, value: 64, Condition: func(cardCalcParams interfaces.CardCalcParams) bool {
		return checkShiBaLuoHan(cardCalcParams)
	}}

	ScxlFan = append(ScxlFan, fanPinghu, fanQingYiSe, fanQiDui, fanPengPengHu, fanJingGouDiao, fanShiBaLuoHan)
}

//ScxlFanMutex 番型和根处理
func ScxlFanMutex(fans []majongpb.CardType, gen uint32) ([]majongpb.CardType, uint32) {
	// 翻型只有1个，并且根为0,直接返回
	if len(fans) == 1 && gen == 0 {
		return []majongpb.CardType{majongpb.CardType(fans[0])}, 0
	}
	fansMap := make(map[majongpb.CardType]majongpb.CardType)
	for _, fanCardType := range fans {
		fansMap[fanCardType] = majongpb.CardType(fanCardType)
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
			fansMap[majongpb.CardType_QingShiBaLuoHan] = majongpb.CardType_QingShiBaLuoHan
			flag = true
		}
		if value, ok := fansMap[majongpb.CardType_JingGouDiao]; ok && value > 0 { // 添加清金钩钓,移除金钩钓
			delete(fansMap, majongpb.CardType_JingGouDiao)
			fansMap[majongpb.CardType_QingJingGouDiao] = majongpb.CardType_QingJingGouDiao
			flag = true
		}
		if value, ok := fansMap[majongpb.CardType_PengPengHu]; ok && value > 0 { // 添加清碰,移除碰碰胡
			delete(fansMap, majongpb.CardType_PengPengHu)
			fansMap[majongpb.CardType_QingPeng] = majongpb.CardType_QingPeng
			flag = true
		}
		if value, ok := fansMap[majongpb.CardType_QiDui]; ok && value > 0 {
			delete(fansMap, majongpb.CardType_QiDui)
			if gen > 0 { // 添加清龙七对,移除七对/清七对/龙七对
				fansMap[majongpb.CardType_QingLongQiDui] = majongpb.CardType_QingLongQiDui
				gen--
			} else { // 添加清七对,移除七对
				fansMap[majongpb.CardType_QingQiDui] = majongpb.CardType_QingQiDui
			}
			flag = true
		}
		if flag { // 存在可以跟清一色可以合组的牌型，移除清一色
			delete(fansMap, majongpb.CardType_QingYiSe)
		}
	}
	if value, ok := fansMap[majongpb.CardType_QiDui]; ok && gen > 0 && value > 0 { // 龙七对
		delete(fansMap, majongpb.CardType_QiDui)                           // 移除七对
		fansMap[majongpb.CardType_LongQiDui] = majongpb.CardType_LongQiDui // 添加龙七对
		gen--                                                              // 根减1
	}
	// fan的cardType，转为卡牌的cardType
	cardTypes := make([]majongpb.CardType, 0)
	for _, cardType := range fansMap {
		cardTypes = append(cardTypes, cardType)
	}
	return cardTypes, gen
}
