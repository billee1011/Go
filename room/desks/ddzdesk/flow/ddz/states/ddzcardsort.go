package states

import "sort"

type DDZCard struct {
	suit uint32
	point uint32
	weight uint32
}

func (c DDZCard) toInt() uint32 {
	return c.suit + c.point
}

type DDZCardSlice []DDZCard
func (cs DDZCardSlice) Len() int           { return len(cs) }
func (cs DDZCardSlice) Swap(i, j int)      { cs[i], cs[j] = cs[j], cs[i] }
func (cs DDZCardSlice) Less(i, j int) bool { return cs[i].weight < cs[j].weight }

func toCard(card uint32) DDZCard {
	result := DDZCard{}
	result.suit = card / 16
	result.point = card % 16

	if result.point == 0x01 || result.point == 0x02 {
		result.weight = result.suit + 0x0D + result.point//A和2，加大权重
	} else if result.point == 0x0E || result.point == 0x0F {
		result.weight = 0x50 + result.point//大小王，加大权重
	} else {
		result.weight = result.suit + result.point
	}
	return result
}

func toCards(cards []uint32) DDZCardSlice {
	result := make([]DDZCard, len(cards))
	for _, card := range cards {
		result = append(result, toCard(card))
	}
	return result
}

// 按斗地主牌的大小排序后返回
func ddzSort(cards []uint32) []uint32 {
	cs := toCards(cards)
	sort.Sort(cs)
	result := make([]uint32, len(cards))
	for _, c := range cs {
		result = append(result, c.toInt())
	}
	return result
}

