package fantype

import (
	"steve/gutils"
)

// checkXiaoSanYuan 检查小三元，含有“中发白” 2副箭刻子,1副箭将牌
func checkXiaoSanYuan(tc *typeCalculator) bool {
	num := 0
	for _, peng := range tc.getPengCards() {
		pengCard := gutils.ServerCard2Number(peng.GetCard())
		if pengCard >= gutils.Zhong && pengCard <= gutils.Bai {
			num++
		}
	}
	for _, combine := range tc.combines {
		// 将是箭
		if combine.jiang >= gutils.Zhong && combine.jiang <= gutils.Bai {
			count := 0
			// 刻子是箭
			for _, ke := range combine.kes {
				if ke >= gutils.Zhong && ke <= gutils.Bai {
					count++
				}
			}
			//含有2个箭牌的刻子
			if num+count >= 2 {
				return true
			}
		}
	}
	return false
}
