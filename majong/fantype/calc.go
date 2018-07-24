package fantype

import (
	"steve/majong/interfaces"
	majongpb "steve/entity/majong"
	"steve/majong/bus"
)

type fanTypeCalculator struct {
}

func (fc fanTypeCalculator) Calculate(fanParams interfaces.FantypeParams) (fanTypes []int, gengCount int, huaCount int) {
	return CalculateFanTypes(fanParams.MjContext, fanParams.PlayerID, fanParams.HandCard, fanParams.HuCard)
}

func (fc fanTypeCalculator) CardTypeValue(mjContext *majongpb.MajongContext, fanTypes []int, gengCount int, huaCount int) uint64 {
	return CalculateScore(mjContext, fanTypes, int(gengCount), int(huaCount))
}

func init() {
	fc := fanTypeCalculator{}
	bus.SetFanTypeCalculator(fc)

}
