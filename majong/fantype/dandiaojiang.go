package fantype

import (
	"steve/gutils"
	majongpb "steve/server_pb/majong"
)

// checkDanDiaoJiang 单钓将:钓单张牌作将成胡,1112 胡 2 算单钓将,1234 胡 1、4 不算单钓将
func checkDanDiaoJiang(tc *typeCalculator) bool {
	huCard := tc.getHuCard()

	player := tc.getPlayer()

	huCardN := gutils.ServerCard2Number(huCard.Card)
	isJiang := false
	for _, combine := range tc.combines {
		jiang := gutils.ServerCard2Number(combine.jiang)
		if jiang == huCardN {
			isJiang = true
			break
		}
	}
	if !isJiang {
		return false
	}
	if isZiMo(huCard) {
		canTingCardInfos := player.GetZixunRecord().CanTingCardInfo
		for _, canTingCardInfo := range canTingCardInfos {
			if canTingCardInfo.OutCard == huCardN && len(canTingCardInfo.TingCardInfo) == 1 {
				return true
			}
		}

	} else if len(player.TingCardInfo) == 1 && player.TingCardInfo[0].TingCard == huCardN {
		return true
	}
	return false
}

func isZiMo(huCard *majongpb.HuCard) bool {
	return map[majongpb.HuType]bool{
		majongpb.HuType_hu_dihu:              true,
		majongpb.HuType_hu_gangkai:           true,
		majongpb.HuType_hu_gangshanghaidilao: true,
		majongpb.HuType_hu_haidilao:          true,
		majongpb.HuType_hu_tianhu:            true,
		majongpb.HuType_hu_zimo:              true,
	}[huCard.GetType()]
}
