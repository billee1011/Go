package states

import (
	"sort"
	"github.com/Sirupsen/logrus"
)

var (
	redJoker = toDDZCard(0x0F)   //大王
	blackJoker = toDDZCard(0x0E) //小王
	sDiamond = uint32(0x10)                   //方块
	sClub = uint32(0x20)                      //梅花
	sHeart = uint32(0x30)                     //红桃
	sSpade = uint32(0x40)                     //黑桃
	pA = uint32(0x01)
	p2 = uint32(0x02)
	p3 = uint32(0x03)
	p4 = uint32(0x04)
	p5 = uint32(0x05)
	p6 = uint32(0x06)
	p7 = uint32(0x07)
	p8 = uint32(0x08)
	p9 = uint32(0x09)
	p10 = uint32(0x0A)
	pJ = uint32(0x0B)
	pQ = uint32(0x0C)
	pK = uint32(0x0D)
	pBlackJoker = uint32(0x0E)
	pRedJoker = uint32(0x0F)
)

type Poker struct {
	Suit   uint32 //花色 0x00,0x10,0x20,0x30,xx40
	Point  uint32 //点数 0x01-0x0D(A-K), 0x0E(小王), 0x0F(大王)
	Weight uint32 //带花色权重
	PointWeight uint32 //无花色权重
}

func (c Poker) toInt() uint32 {
	return c.Suit + c.Point
}

func (c Poker) equals(other Poker) bool {
	return c.Suit == other.Suit && c.Point == other.Point
}

// 带花色比较，黑桃A 和 方块A比较返回true
func (c Poker) biggerThan(other Poker) bool {
	return c.Weight > other.Weight
}

// 无花色比较，黑桃A 和 方块A比较返回false
func (c Poker) pointBiggerThan(other Poker) bool {
	return c.PointWeight > other.PointWeight
}

type DDZCardSlice []Poker
func (cs DDZCardSlice) Len() int           { return len(cs) }
func (cs DDZCardSlice) Swap(i, j int)      { cs[i], cs[j] = cs[j], cs[i] }
func (cs DDZCardSlice) Less(i, j int) bool { return cs[i].Weight < cs[j].Weight }

type DDZPointSlice []Poker
func (cs DDZPointSlice) Len() int           { return len(cs) }
func (cs DDZPointSlice) Swap(i, j int)      { cs[i], cs[j] = cs[j], cs[i] }
func (cs DDZPointSlice) Less(i, j int) bool { return cs[i].PointWeight < cs[j].PointWeight }

func toDDZCard(card uint32) Poker {
	result := Poker{}
	result.Suit = card / 16 * 16
	result.Point = card % 16

	// 计算无花色权重
	if result.Point == pA || result.Point == p2 {
		result.PointWeight = pK + result.Point //A和2，加大权重
	} else if result.Point == pBlackJoker || result.Point == pRedJoker {
		result.PointWeight = sSpade + result.Point //大小王，加大权重
	} else {
		result.PointWeight = result.Point
	}
	result.Weight = result.Suit + result.PointWeight//带花色权重
	return result
}

func toDDZCards(cards []uint32) []Poker {
	result := make([]Poker, 0, len(cards))
	for _, card := range cards {
		result = append(result, toDDZCard(card))
	}
	return result
}

func toInts(cards []Poker) []uint32 {
	result := make([]uint32, 0, len(cards))
	for _, card := range cards {
		result = append(result, card.toInt())
	}
	return result
}

// DDZSort 从小到大排序后返回
func DDZSort(cards []uint32) []uint32 {
	return ddzSort(cards, false)
}

// DDZSortDescend 从大到小排序后返回
func DDZSortDescend(cards []uint32) []uint32 {
	return ddzSort(cards, true)
}

func ddzSort(cards []uint32, reverse bool) []uint32 {
	cs := DDZCardSlice(toDDZCards(cards))
	if reverse {
		sort.Sort(sort.Reverse(cs))
	} else {
		sort.Sort(cs)
	}
	result := make([]uint32, 0, cs.Len())
	for i := range cs {
		result = append(result, cs[i].toInt())
	}
	logrus.WithFields(logrus.Fields{"in":cards, "out:":result}).Debug("斗地主排序")
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