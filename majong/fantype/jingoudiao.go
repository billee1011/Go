package fantype

import (
	"steve/gutils"
)

// checkJinGouDiao 检测金钩钓 胡牌时手里只剩一张，并且单钓一这张，其他的牌都被杠或碰了,不计碰碰胡。
func checkJinGouDiao(tc *typeCalculator) bool {
	if huCard := tc.getHuCard(); huCard != nil {
		if len(tc.getHandCards()) == 1 && gutils.CardEqual(tc.getHandCards()[0], huCard.GetCard()) {
			return true
		}
	}
	return false
}
