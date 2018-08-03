package states

import (
	"sort"
	"strconv"

	"github.com/Sirupsen/logrus"
)

var (
	RedJoker    = ToDDZCard(0x0F) //大王
	BlackJoker  = ToDDZCard(0x0E) //小王
	sDiamond    = uint32(0x10)    //方块
	sClub       = uint32(0x20)    //梅花
	sHeart      = uint32(0x30)    //红桃
	sSpade      = uint32(0x40)    //黑桃
	pA          = uint32(0x01)
	p2          = uint32(0x02)
	p3          = uint32(0x03)
	p4          = uint32(0x04)
	p5          = uint32(0x05)
	p6          = uint32(0x06)
	p7          = uint32(0x07)
	p8          = uint32(0x08)
	p9          = uint32(0x09)
	p10         = uint32(0x0A)
	pJ          = uint32(0x0B)
	pQ          = uint32(0x0C)
	pK          = uint32(0x0D)
	pBlackJoker = uint32(0x0E)
	pRedJoker   = uint32(0x0F)
)

type Poker struct {
	Suit        uint32 //花色 0x00,0x10,0x20,0x30,xx40
	Point       uint32 //点数 0x01-0x0D(A-K), 0x0E(小王), 0x0F(大王)
	Weight      uint32 //带花色权重,用于带花色大小比较，同点数的在一起
	PointWeight uint32 //无花色权重，用于无花色大小比较
}

func (c Poker) String() string {
	if c.Suit == sDiamond {
		return "♦" + c.GetPointString()
	} else if c.Suit == sClub {
		return "♣" + c.GetPointString()
	} else if c.Suit == sHeart {
		return "♥" + c.GetPointString()
	} else if c.Suit == sSpade {
		return "♠" + c.GetPointString()
	} else {
		return c.GetPointString()
	}
}

func (c Poker) GetPointString() string {
	if c.Point == pA {
		return "A"
	} else if c.Point == pJ {
		return "J"
	} else if c.Point == pQ {
		return "Q"
	} else if c.Point == pK {
		return "K"
	} else if c.Point == pBlackJoker {
		return "小王"
	} else if c.Point == pRedJoker {
		return "大王"
	} else {
		return strconv.Itoa(int(c.Point))
	}
}

func (c Poker) ToInt() uint32 {
	return c.Suit + c.Point
}

func (c Poker) Equals(other Poker) bool {
	return c.Suit == other.Suit && c.Point == other.Point
}

// 带花色比较，黑桃A 和 方块A比较返回true
func (c Poker) BiggerThan(other Poker) bool {
	return c.Weight > other.Weight
}

// 无花色比较，黑桃A 和 方块A比较返回false
func (c Poker) PointBiggerThan(other Poker) bool {
	return c.PointWeight > other.PointWeight
}

type DDZCardSlice []Poker

func (cs DDZCardSlice) Len() int           { return len(cs) }
func (cs DDZCardSlice) Swap(i, j int)      { cs[i], cs[j] = cs[j], cs[i] }
func (cs DDZCardSlice) Less(i, j int) bool { return cs[i].Weight < cs[j].Weight }

func ToDDZCard(card uint32) Poker {
	result := Poker{}
	result.Suit = card / 16 * 16
	result.Point = card % 16

	// 计算无花色权重
	if result.Point == pA {
		result.PointWeight = pK + pA //A为K加1
	} else if result.Point == p2 {
		result.PointWeight = pK + p2 + 1 //2为A加1,方便断开顺子,连对等
	} else if result.Point == pBlackJoker || result.Point == pRedJoker {
		result.PointWeight = sSpade + pK + result.Point //大小王，加大权重
	} else {
		result.PointWeight = result.Point
	}
	result.Weight = result.PointWeight*5 + result.Suit/16 //带花色权重
	return result
}

func ToDDZCards(cards []uint32) []Poker {
	result := make([]Poker, 0, len(cards))
	for _, card := range cards {
		result = append(result, ToDDZCard(card))
	}
	return result
}

func ToInts(cards []Poker) []uint32 {
	result := make([]uint32, 0, len(cards))
	for _, card := range cards {
		result = append(result, card.ToInt())
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
	cs := DDZCardSlice(ToDDZCards(cards))
	if reverse {
		sort.Sort(sort.Reverse(cs))
	} else {
		sort.Sort(cs)
	}
	result := make([]uint32, 0, cs.Len())
	for i := range cs {
		result = append(result, cs[i].ToInt())
	}
	logrus.WithFields(logrus.Fields{"in": cards, "out:": result}).Debug("斗地主排序")
	return result
}

func DDZPokerSort(cards []Poker) {
	cs := DDZCardSlice(cards)
	sort.Sort(cs)
}

func DDZPokerSortDesc(cards []Poker) {
	cs := DDZCardSlice(cards)
	sort.Sort(sort.Reverse(cs))
}
