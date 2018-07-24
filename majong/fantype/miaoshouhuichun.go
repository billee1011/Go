package fantype

import (
	"steve/common/mjoption"
	"steve/majong/utils"
	majongpb "steve/entity/majong"
)

//miaoshouhuichun 检测妙手回春，最后一张牌自摸
func checkMiaoShouHuiChun(tc *typeCalculator) bool {
	huCard := tc.getHuCard()
	if huCard != nil && huCard.GetType() == majongpb.HuType_hu_zimo {
		mjContext := tc.mjContext
		if !IsWallCanMoPai(mjContext) {
			return true
		}
	}
	return false
}

// IsWallCanMoPai 判断牌墙是否能摸牌
func IsWallCanMoPai(context *majongpb.MajongContext) bool {
	if len(context.WallCards) == 0 {
		return false
	}
	// 由配牌控制是否gameover,配牌长度为0走正常gameover,配牌长度不为0走配牌长度流局
	length := context.GetOption().GetWallcardsLength()
	maxCount := 0
	if mjoption.GetXingpaiOption(int(context.GetXingpaiOptionId())).EnableKaijuAddflower {
		maxCount = int(length) - (len(context.GetPlayers()) * 13)
	} else {
		maxCount = int(length) - (len(context.GetPlayers())*13 + 1)
	}
	if utils.GetAllMopaiCount(context) == maxCount {
		return false
	}
	return true
}
