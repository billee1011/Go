package fantype

import (
	majongpb "steve/server_pb/majong"
)

// checkTianHu 天胡
func checkTianHu(tc *typeCalculator) bool {
	if tc.huCard.Type == majongpb.HuType_hu_tianhu {
		return true
	}
	return false
}

// calcHuType 计算胡牌类型
// func calcHuType(mjContext *majongpb.MajongContext, huPlayerID uint64) majongpb.HuType {
// 	afterGang := isAfterGang(mjContext)
// 	isLast := noCardsToTake(mjContext)
// 	if afterGang && isLast {
// 		return majongpb.HuType_hu_gangshanghaidilao
// 	} else if afterGang {
// 		return majongpb.HuType_hu_gangkai
// 	} else if isLast {
// 		return majongpb.HuType_hu_haidilao
// 	}
// 	huPlayer := utils.GetMajongPlayer(huPlayerID, mjContext)
// 	if len(huPlayer.PengCards) == 0 && len(huPlayer.GangCards) == 0 && len(huPlayer.HuCards) == 0 {
// 		if huPlayer.MopaiCount == 0 && huPlayerID == mjContext.Players[mjContext.ZhuangjiaIndex].GetPalyerId() {
// 			return majongpb.HuType_hu_tianhu
// 		}
// 		if huPlayer.MopaiCount == 1 && huPlayerID != mjContext.Players[mjContext.ZhuangjiaIndex].GetPalyerId() {
// 			return majongpb.HuType_hu_dihu
// 		}
// 	}
// 	return majongpb.HuType_hu_zimo
// }

// // isAfterGang 是否为杠后炮
// // 杠后摸牌、自询出牌则为杠后炮
// func isAfterGang(mjContext *majongpb.MajongContext) bool {
// 	zxType := mjContext.GetZixunType()
// 	mpType := mjContext.GetMopaiType()
// 	return mpType == majongpb.MopaiType_MT_GANG && zxType == majongpb.ZixunType_ZXT_NORMAL
// }

// func noCardsToTake(context *majongpb.MajongContext) bool {
// 	length := context.GetOption().GetWallcardsLength()
// 	if utils.GetAllMopaiCount(context) == int(length)-53 {
// 		return true
// 	}
// 	if len(context.WallCards) == 0 {
// 		return true
// 	}
// 	return false
// }
