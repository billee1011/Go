package utils

import (
	majongpb "steve/entity/majong"
)

// GetCardNum 获取指定牌在指定数组中的数量
func GetCardNum(srcCard *majongpb.Card, cards []*majongpb.Card) (num int) {
	for _, card := range cards {
		if CardEqual(card, srcCard) {
			num++
		}
	}
	return
}
