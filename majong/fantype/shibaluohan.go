package fantype

import (
	"steve/gutils"
)

//checkShiBaLuoHan 检查十八罗汉 胡牌时手上只剩一张牌单吊，其他手牌形成四个杠，此时不计四根和碰碰胡。
func checkShiBaLuoHan(tc *typeCalculator) bool {
	if len(tc.getGangCards()) == 4 {
		if huCard := tc.getHuCard(); huCard != nil && len(tc.getHandCards()) == 1 {
			if gutils.CardEqual(tc.getHandCards()[0], huCard.GetCard()) {
				return true
			}
		}
	}
	return false
}
