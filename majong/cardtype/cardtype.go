package cardtype

import (
	"steve/majong/interfaces"
	"steve/majong/settle/fan"
	majongpb "steve/server_pb/majong"
)

//XueLiuCardTypeCalculator 血流卡牌类型计算器
type XueLiuCardTypeCalculator struct {
}

func init() {
}

//Calculate 获取能胡所有番型，及根，最小平胡
func (ctc *XueLiuCardTypeCalculator) Calculate(params interfaces.CardCalcParams) (cardTypes []interfaces.CardType, gengCount uint32) {
	fanCardTypes := make([]majongpb.CardType, 0)
	// 遍历可行番型
	for i := 0; i < len(fan.ScxlFan); i++ {
		if fan.ScxlFan[i].Condition(params) {
			fanName := fan.ScxlFan[i].GetFanName()
			fanCardTypes = append(fanCardTypes, fanName)
		}
	}
	// 番型名和根处理
	cardTypes, gengCount = fan.ScxlFanMutex(fanCardTypes, fan.GetGenCount(params))
	return cardTypes, gengCount
}

//CardTypeValue 获取总倍数
func (ctc *XueLiuCardTypeCalculator) CardTypeValue(cardTypes []interfaces.CardType, gengCount uint32) uint32 {
	total := uint32(1)
	// 叠乘番型
	for _, cardType := range cardTypes {
		fanCardType := majongpb.CardType(cardType)
		if multiple, isExist := fan.FanValue[fanCardType]; isExist {
			total = total * multiple
		}
	}
	// 根平方
	genTotoal := uint32(1 << gengCount)
	// 根乘总番型倍数
	total = total * genTotoal
	return total
}

//CardGenSum 获取根数量，没有做处理的，如是七对，并且是1根，就直接返回1
func (ctc *XueLiuCardTypeCalculator) CardGenSum(params interfaces.CardCalcParams) uint32 {
	return fan.GetGenCount(params)
}
