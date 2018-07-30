package fantype

import "steve/gutils"

// checkXiaoSiXi 检查小四喜 含有“东南西北”风将牌，3个风刻子
func checkXiaoSiXi(tc *typeCalculator) bool {
	num := 0
	for _, peng := range tc.getPengCards() {
		pengCard := gutils.ServerCard2Number(peng.GetCard())
		if pengCard >= gutils.Dong && pengCard <= gutils.Bei {
			num++
		}
	}
	for _, combine := range tc.combines {
		// 将是风
		if combine.jiang >= gutils.Dong && combine.jiang <= gutils.Bei {
			count := 0
			// 刻子是风
			for _, ke := range combine.kes {
				if ke >= gutils.Dong && ke <= gutils.Bei {
					count++
				}
			}
			//含有3个风牌的刻子
			if num+count >= 3 {
				return true
			}
		}
	}
	return false
}
