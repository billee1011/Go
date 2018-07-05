package states

import "sort"

var (
	blackJoker = toDDZCard(0x0E)
	redJoker = toDDZCard(0x0F)
)

type Poker struct {
	suit uint32
	point uint32
	weight uint32
}

func (c Poker) toInt() uint32 {
	return c.suit + c.point
}

func (c Poker) equals(other Poker) bool {
	return c.suit == other.suit && c.point == other.point
}

// 带花色比较，黑桃A 和 方块A比较返回true
func (c Poker) biggerThan(other Poker) bool {
	return c.weight > other.weight
}

// 无花色比较，黑桃A 和 方块A比较返回true
func (c Poker) pointBiggerThan(other Poker) bool {
	return c.getPointWeight() > other.getPointWeight()
}

func (c Poker) getPointWeight() uint32{//方法的方式节省空间，用时间换空间(因为不常用，常用的weight还是得保存到内存中)
	if c.point == 0x01 || c.point == 0x02 {
		return 0x0D + c.point
	} else if c.point == 0x0E || c.point == 0x0F {
		return 0x10 + c.point
	} else {
		return c.point
	}
}

type DDZCardSlice []Poker
func (cs DDZCardSlice) Len() int           { return len(cs) }
func (cs DDZCardSlice) Swap(i, j int)      { cs[i], cs[j] = cs[j], cs[i] }
func (cs DDZCardSlice) Less(i, j int) bool { return cs[i].weight < cs[j].weight }

type DDZPointSlice []Poker
func (cs DDZPointSlice) Len() int           { return len(cs) }
func (cs DDZPointSlice) Swap(i, j int)      { cs[i], cs[j] = cs[j], cs[i] }
func (cs DDZPointSlice) Less(i, j int) bool { return cs[i].point < cs[j].point }

func toDDZCard(card uint32) Poker {
	result := Poker{}
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

func toDDZCards(cards []uint32) []Poker {
	result := make([]Poker, len(cards))
	for _, card := range cards {
		result = append(result, toDDZCard(card))
	}
	return result
}

func toInts(cards []Poker) []uint32 {
	result := make([]uint32, len(cards))
	for _, card := range cards {
		result = append(result, card.toInt())
	}
	return result
}

// 按斗地主牌的大小排序后返回
func ddzSort(cards []uint32) []uint32 {
	cs := DDZCardSlice(toDDZCards(cards))
	sort.Sort(cs)
	result := make([]uint32, len(cards))
	for _, c := range cs {
		result = append(result, c.toInt())
	}
	return result
}

func ddzPokerSort(cards []Poker) {
	cs := DDZCardSlice(cards)
	sort.Sort(cs)
}

// 按斗地主点数的大小排序后返回
func ddzPointSort(cards []Poker) {
	ps := DDZPointSlice(cards)
	sort.Sort(ps)
}

