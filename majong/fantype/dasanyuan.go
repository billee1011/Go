package fantype

import (
	"steve/gutils"
)

// checkDaSanYuan 检查大三元，含有“中发白” 3副刻子
func checkDaSanYuan(tc *typeCalculator) bool {
	num := 0
	for _, peng := range tc.getPengCards() {
		pengCard := gutils.ServerCard2Number(peng.GetCard())
		if pengCard >= gutils.Zhong && pengCard <= gutils.Bai {
			num++
		}
	}
	for _, combine := range tc.combines {
		count := 0
		for _, ke := range combine.kes {
			if ke >= gutils.Zhong && ke <= gutils.Bai {
				count++
			}
		}
		if count+num == 3 {
			return true
		}
	}
	return true
}
