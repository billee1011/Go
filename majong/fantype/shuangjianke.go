package fantype

import (
	"steve/gutils"
)

// checkShuangJianKe 检测双箭刻 含有2副箭(中发白)刻或杠
func checkShuangJianKe(tc *typeCalculator) bool {
	count := 0
	for _, gang := range tc.getGangCards() {
		gangCard := gutils.ServerCard2Number(gang.GetCard())
		if gangCard >= gutils.Zhong && gangCard <= gutils.Bai {
			count++
		}
	}
	for _, peng := range tc.getPengCards() {
		pengCard := gutils.ServerCard2Number(peng.GetCard())
		if pengCard >= gutils.Zhong && pengCard <= gutils.Bai {
			count++
		}
	}
	for _, combine := range tc.combines {
		keCount := 0
		for _, ke := range combine.kes {
			if ke >= gutils.Zhong && ke <= gutils.Bai {
				keCount++
			}
		}
		if keCount+count >= 2 {
			return true
		}
	}
	return false
}
