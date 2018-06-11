package utils

import (
	"sort"
	majongpb "steve/server_pb/majong"
)

type cards []*majongpb.Card

func (c cards) Len() int {
	return len(c)
}

func (c cards) Less(i, j int) bool {
	return cti(c[i]) < cti(c[j])
}

func (c cards) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

// Sort is a convenience method.
func (c cards) Sort() { sort.Sort(c) }

//cti 将Card转换成牌值(int)
func cti(card *majongpb.Card) int {
	var color int
	switch card.GetColor() {
	case majongpb.CardColor_ColorWan:
		color = 1
	case majongpb.CardColor_ColorTiao:
		color = 2
	case majongpb.CardColor_ColorTong:
		color = 3
	}
	tValue := int(card.Point)
	value := color*10 + tValue
	return value
}

// SortCards 排序
func SortCards(a []*majongpb.Card) { sort.Sort(cards(a)) }
