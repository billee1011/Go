package fantype

import (
	"fmt"
	majongpb "steve/server_pb/majong"
)

//checkHaiDiLaoYue 检测海底捞月 胡最后打出的牌，必须是最后摸牌的人点炮
func checkHaiDiLaoYue(tc *typeCalculator) bool {
	huCard := tc.getHuCard()
	fmt.Println("----------------------------")
	fmt.Println(huCard)
	if huCard != nil && huCard.GetType() == majongpb.HuType_hu_dianpao {
		mjContext := tc.mjContext
		fmt.Println(!IsWallCanMoPai(mjContext))
		fmt.Println(huCard.GetSrcPlayer())
		fmt.Println(mjContext.GetLastMopaiPlayer())
		if !IsWallCanMoPai(mjContext) && huCard.GetSrcPlayer() == mjContext.GetLastMopaiPlayer() {
			return true
		}
	}
	fmt.Println("-========================")
	return false
}
