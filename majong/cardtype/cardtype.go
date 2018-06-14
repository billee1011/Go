package cardtype

import (
	"steve/majong/cardtype/fan"
	"steve/majong/interfaces"
	majongpb "steve/server_pb/majong"

	"github.com/Sirupsen/logrus"
)

//scxlCardTypeCalculator 血流卡牌类型计算器
type scxlCardTypeCalculator struct {
}

func init() {
}

//Calculate 获取能胡所有番型，及根，最小平胡
func (ctc *scxlCardTypeCalculator) Calculate(params interfaces.CardCalcParams) (cardTypes []majongpb.CardType, gengCount uint32) {
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
	logrus.WithFields(logrus.Fields{
		"func_name": "scxlCardTypeCalculator.Calculate",
		"cardTypes": cardTypes,
		"params":    params,
		"genCount":  gengCount,
	}).Debugln("计算胡牌倍数")
	return cardTypes, gengCount
}

//CardTypeValue 获取总倍数及根数	（注：总倍数已包括根的倍数了）
func (ctc *scxlCardTypeCalculator) CardTypeValue(gameID int, cardTypes []majongpb.CardType, gengCount uint32) (uint32, uint32) {
	total := uint32(1)
	// 叠乘番型
	for _, cardType := range cardTypes {
		fanCardType := majongpb.CardType(cardType)
		if multiple, isExist := fan.ScxlFanValue[fanCardType]; isExist {
			total = total * multiple
		}
	}
	// 根平方
	genTotoal := uint32(1 << gengCount)
	// 根乘总番型倍数
	total = total * genTotoal
	return total, gengCount
}
