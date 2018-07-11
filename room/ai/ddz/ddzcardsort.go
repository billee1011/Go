package ddz

import (
	"sort"
	"strconv"

	"github.com/Sirupsen/logrus"
)

var (
	redJoker    = ToDDZCard(0x0F) //大王
	blackJoker  = ToDDZCard(0x0E) //小王
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

// Poker 单张扑克牌的信息
type Poker struct {
	Suit        uint32 //花色 0x00,0x10,0x20,0x30,xx40
	Point       uint32 //点数 0x01-0x0D(A-K), 0x0E(小王), 0x0F(大王)
	Weight      uint32 //带花色权重,用于带花色大小比较，
	PointWeight uint32 //无花色权重，用于无花色大小比较，只比较点数值，结果：大王>小王>2>A>K>Q>J>10>9>8>7>6>5>4>3
	SortWeight  uint32 //排序权重，用于排序，同点数需要在一起
}

//
func (c Poker) String() string {
	if c.Suit == sDiamond {
		return "♦" + c.getPointString()
	} else if c.Suit == sClub {
		return "♣" + c.getPointString()
	} else if c.Suit == sHeart {
		return "♥" + c.getPointString()
	} else if c.Suit == sSpade {
		return "♠" + c.getPointString()
	} else {
		return c.getPointString()
	}
}

// 获取牌点数的字符串
func (c Poker) getPointString() string {
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

// 转为数字
func (c Poker) toInt() uint32 {
	return c.Suit + c.Point
}

// 两张牌是否相同
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

//DDZCardSlice 棋牌数组（比较时按照排序权重）
type DDZCardSlice []Poker

func (cs DDZCardSlice) Len() int           { return len(cs) }
func (cs DDZCardSlice) Swap(i, j int)      { cs[i], cs[j] = cs[j], cs[i] }
func (cs DDZCardSlice) Less(i, j int) bool { return cs[i].SortWeight < cs[j].SortWeight }

//DDZPointSlice 棋牌数组（比较时按照无花色权重）
type DDZPointSlice []Poker

func (cs DDZPointSlice) Len() int           { return len(cs) }
func (cs DDZPointSlice) Swap(i, j int)      { cs[i], cs[j] = cs[j], cs[i] }
func (cs DDZPointSlice) Less(i, j int) bool { return cs[i].PointWeight < cs[j].PointWeight }

// ToDDZCard uint32类型的卡牌转为Poker
func ToDDZCard(card uint32) Poker {
	result := Poker{}

	// 花色值
	result.Suit = card / 16 * 16

	// 点数值
	result.Point = card % 16

	// 计算无花色权重
	// 大王 小王  2   A    K   Q   J   10   9   8   7   6   5   4   3
	//  92  91  16  14   13  12  11   10   9   8   7   6   5   4   3
	if result.Point == pA {
		result.PointWeight = pK + pA //A为K加1
	} else if result.Point == p2 {
		result.PointWeight = pK + p2 + 1 //2为A加1,方便断开顺子,连对等
	} else if result.Point == pBlackJoker || result.Point == pRedJoker {
		result.PointWeight = sSpade + pK + result.Point //大小王，加大权重，所以加上黑桃（64），这样无论带不带花色，权重都是最大
	} else {
		result.PointWeight = result.Point
	}

	// 带花色权重 = 花色值 + 无花色权重
	// 注意：大小王仍然是最大
	result.Weight = result.Suit + result.PointWeight

	// 排序权重 = 无花色权重 * 5 + 花色枚举值（1,2,3,4）
	// 大王 小王  黑桃2	  红桃2	  梅花2	  方块2    黑桃A    红桃A	  梅花A	  方块A   黑桃K    Q   J   10   9   8   7   6   5   4   3
	//  460  455  84     83      82     81  	74       73       72     71      69
	result.SortWeight = result.PointWeight*5 + result.Suit/16 //点数相同的放在一起
	return result
}

// 把uint32数组类型的棋牌转为Poker数组
func ToDDZCards(cards []uint32) []Poker {
	result := make([]Poker, 0, len(cards))
	for _, card := range cards {
		result = append(result, ToDDZCard(card))
	}
	return result
}

// 把Poker数组转为uint32数组
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

// ddzSort 排序
// reverse 为true表从大到小排序
//         为false表从小到大排序
func ddzSort(cards []uint32, reverse bool) []uint32 {
	cs := DDZCardSlice(ToDDZCards(cards))
	if reverse {
		sort.Sort(sort.Reverse(cs))
	} else {
		sort.Sort(cs)
	}
	result := make([]uint32, 0, cs.Len())
	for i := range cs {
		result = append(result, cs[i].toInt())
	}
	logrus.WithFields(logrus.Fields{"in": cards, "out:": result}).Debug("斗地主排序")
	return result
}

// DdzPokerSort 按排序权重排序后返回，从小到大
func DdzPokerSort(cards []Poker) {
	cs := DDZCardSlice(cards)
	sort.Sort(cs)
}

// 按斗地主点数的大小排序后返回，从小到大
func ddzPointSort(cards []Poker) {
	ps := DDZPointSlice(cards)
	sort.Sort(ps)
}
