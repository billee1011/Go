package tests

import (
	"testing"

	"math/rand"
	"steve/room/flows/ddzflow/ddz/states"

	"github.com/stretchr/testify/assert"
	"steve/entity/poker"
)

func Test_IsKingBomb(t *testing.T) {
	is, pivot := states.IsKingBomb([]states.Poker{states.BlackJoker, states.RedJoker})
	assert.Equal(t, true, is)
	assert.Equal(t, &states.RedJoker, pivot)
}

func Test_IsBomb(t *testing.T) {
	cards := states.ToDDZCards([]uint32{0x13, 0x23, 0x33, 0x43})
	is, pivot := states.IsBomb(cards)
	assert.Equal(t, true, is)
	assert.Equal(t, &cards[3], pivot)
}

func Test_IsBombAndPairs(t *testing.T) {
	cards := states.ToDDZCards([]uint32{0x13, 0x23, 0x33, 0x43, 0x14, 0x24, 0x34, 0x44})
	is, pivot := states.IsBombAndPairs(cards)
	assert.Equal(t, true, is)
	assert.Equal(t, &cards[7], pivot)
}

func Test_IsBombAndSingles(t *testing.T) {
	cards := states.ToDDZCards([]uint32{0x13, 0x23, 0x33, 0x43, 0x14, 0x24})
	is, pivot := states.IsBombAndSingles(cards)
	assert.Equal(t, true, is)
	assert.Equal(t, &cards[3], pivot)
}

func Test_IsTriples(t *testing.T) {
	cards := states.ToDDZCards([]uint32{0x13, 0x23, 0x33, 0x14, 0x24, 0x34})
	is, pivot := states.IsTriples(cards)
	assert.Equal(t, true, is)
	assert.Equal(t, uint32(4), pivot.PointWeight)
}

func Test_IsTriplesAndPairs1(t *testing.T) {
	cards := states.ToDDZCards([]uint32{0x13, 0x23, 0x33, 0x14, 0x24, 0x34, 0x15, 0x25, 0x17, 0x37}) // 3334445577 常规牌型
	is, pivot := states.IsTriplesAndPairs(cards)
	assert.Equal(t, true, is)
	assert.Equal(t, uint32(4), pivot.PointWeight)
}

func Test_IsTriplesAndPairs2(t *testing.T) {
	cards := states.ToDDZCards([]uint32{0x13, 0x23, 0x33, 0x14, 0x24, 0x34, 0x15, 0x25, 0x35, 0x45}) // 3334445555 含一个炸弹牌型
	is, pivot := states.IsTriplesAndPairs(cards)
	assert.Equal(t, true, is)
	assert.Equal(t, uint32(4), pivot.PointWeight)
}

func Test_IsTriplesAndPairs3(t *testing.T) {
	cards := states.ToDDZCards([]uint32{0x13, 0x23, 0x33, 0x14, 0x24, 0x34, 0x15, 0x25, 0x35, 0x16, 0x26, 0x36, 0x18, 0x28, 0x38, 0x48, 0x19, 0x29, 0x39, 0x49}) // 33344455566688889999 含两个炸弹牌型
	is, pivot := states.IsTriplesAndPairs(cards)
	assert.Equal(t, true, is)
	assert.Equal(t, uint32(6), pivot.PointWeight)
}

func Test_IsTriplesAndPairs4(t *testing.T) {
	cards := states.ToDDZCards([]uint32{0x13, 0x23, 0x33, 0x14, 0x24, 0x34, 0x15, 0x25, 0x35, 0x16, 0x26, 0x36, 0x18, 0x28, 0x38, 0x48, 0x17, 0x29, 0x3B, 0x4C}) // 333444555666 8888 79JQ 不符牌型
	is, _ := states.IsTriplesAndPairs(cards)
	assert.Equal(t, false, is)
}

func Test_IsTriplesAndPairs5(t *testing.T) {
	cards := states.ToDDZCards([]uint32{0x13, 0x23, 0x33, 0x14, 0x24, 0x34, 0x15, 0x25, 0x35, 0x16, 0x26, 0x36, 0x18, 0x28, 0x38, 0x48, 0x19, 0x29, 0x3B, 0x4B}) // 333444555666 8888 99JJ 符牌型
	is, _ := states.IsTriplesAndPairs(cards)
	assert.Equal(t, true, is)
}

func Test_IsTriplesAndSingles1(t *testing.T) {
	cards := states.ToDDZCards([]uint32{0x13, 0x23, 0x33, 0x14, 0x24, 0x34, 0x17, 0x29}) // 33344479 常规牌型
	is, pivot := states.IsTriplesAndSingles(cards)
	assert.Equal(t, true, is)
	assert.Equal(t, uint32(4), pivot.PointWeight)
}

func Test_IsTriplesAndSingles2(t *testing.T) {
	cards := states.ToDDZCards([]uint32{0x13, 0x23, 0x33, 0x43, 0x14, 0x24, 0x34, 0x44, 0x15, 0x25, 0x35, 0x45}) // 333344445555 全是炸弹牌型
	is, pivot := states.IsTriplesAndSingles(cards)
	assert.Equal(t, true, is)
	assert.Equal(t, uint32(5), pivot.PointWeight)
}

func Test_IsTriplesAndSingles3(t *testing.T) {
	cards := states.ToDDZCards([]uint32{0x13, 0x23, 0x33, 0x43, 0x14, 0x24, 0x34, 0x44, 0x15, 0x25, 0x35, 0x45, 0x16, 0x26, 0x36, 0x46}) // 3333444455556666
	is, pivot := states.IsTriplesAndSingles(cards)
	assert.Equal(t, true, is)
	assert.Equal(t, uint32(6), pivot.PointWeight)
}

func Test_IsTriplesAndSingles4(t *testing.T) {
	cards := states.ToDDZCards([]uint32{0x13, 0x23, 0x33, 0x14, 0x24, 0x34, 0x15, 0x25, 0x35, 0x17, 0x29, 0x43}) //333444555793 一个炸弹牌型
	is, pivot := states.IsTriplesAndSingles(cards)
	assert.Equal(t, true, is)
	assert.Equal(t, uint32(5), pivot.PointWeight)
}

func Test_IsTriplesAndSingles5(t *testing.T) {
	cards := states.ToDDZCards([]uint32{0x15, 0x25, 0x35, 0x16, 0x26, 0x36, 0x17, 0x27, 0x37, 0x18, 0x28, 0x38, 0x1D, 0x2D, 0x4D, 0x4D}) //555666777888KKKK 特殊炸弹牌型
	is, pivot := states.IsTriplesAndSingles(cards)
	assert.Equal(t, true, is)
	assert.Equal(t, uint32(8), pivot.PointWeight)
}

func Test_IsTriplesAndSingles6(t *testing.T) {
	cards := states.ToDDZCards([]uint32{0x15, 0x25, 0x35, 0x17, 0x27, 0x37, 0x16, 0x26, 0x36, 0x2D, 0x4D, 0x4D}) //555666777KKK 全三牌型
	is, pivot := states.IsTriplesAndSingles(cards)
	assert.Equal(t, true, is)
	assert.Equal(t, uint32(7), pivot.PointWeight)
}

func Test_IsPairs(t *testing.T) {
	cards := states.ToDDZCards([]uint32{0x13, 0x23, 0x14, 0x24, 0x35, 0x45})
	is, pivot := states.IsPairs(cards)
	assert.Equal(t, true, is)
	assert.Equal(t, uint32(5), pivot.PointWeight)
}

func Test_IsShunZi1(t *testing.T) {
	cards := states.ToDDZCards([]uint32{0x39, 0x2A, 0x3B, 0x2C, 0x1D, 0x11}) // 9 10 J Q K A
	is, pivot := states.IsShunZi(cards)
	assert.Equal(t, true, is)
	assert.Equal(t, &cards[5], pivot)
}

func Test_IsShunZi2(t *testing.T) {
	cards := states.ToDDZCards([]uint32{0x39, 0x2A, 0x3B, 0x2C, 0x1D, 0x11, 0x12}) // 9 10 J Q K A 2
	is, pivot := states.IsShunZi(cards)
	assert.Equal(t, false, is)
	assert.Nil(t, pivot)
}

func Test_IsTriple(t *testing.T) {
	cards := states.ToDDZCards([]uint32{0x13, 0x23, 0x33})
	is, pivot := states.IsTriple(cards)
	assert.Equal(t, true, is)
	assert.Equal(t, &cards[2], pivot)
}

func Test_IsTripleAndPair(t *testing.T) {
	cards := states.ToDDZCards([]uint32{0x13, 0x23, 0x33, 0x41, 0x31})
	is, pivot := states.IsTripleAndPair(cards)
	assert.Equal(t, true, is)
	assert.Equal(t, &cards[2], pivot)
}

func Test_IsTripleAndSingle(t *testing.T) {
	cards := states.ToDDZCards([]uint32{0x13, 0x23, 0x33, 0x49})
	is, pivot := states.IsTripleAndSingle(cards)
	assert.Equal(t, true, is)
	assert.Equal(t, &cards[2], pivot)
}

func Test_IsPair(t *testing.T) {
	cards := states.ToDDZCards([]uint32{0x13, 0x23})
	is, pivot := states.IsPair(cards)
	assert.Equal(t, true, is)
	assert.Equal(t, &cards[1], pivot)
}

func Test_IsSingle(t *testing.T) {
	cards := states.ToDDZCards([]uint32{0x13})
	is, pivot := states.IsSingle(cards)
	assert.Equal(t, true, is)
	assert.Equal(t, &cards[0], pivot)
}

func Test_GetCardType(t *testing.T) {
	cards := states.ToDDZCards([]uint32{0x13, 0x33, 0x28, 0x38, 0x48, 0x4B})
	cardType, _ := states.GetCardType(cards)
	assert.Equal(t, poker.CardType_CT_NONE, cardType)
}

func Test_ShunZi(t *testing.T) {
	cards := states.ToDDZCards([]uint32{0x13, 0x24, 0x35, 0x46, 0x17, 0x28, 0x39, 0x4A, 0x1B, 0x2C, 0x3D, 0x41})
	cardType, _ := states.GetCardType(cards)
	assert.Equal(t, poker.CardType_CT_SHUNZI, cardType)
}

func randCard() uint32 {
	suit := rand.Intn(4) + 1
	point := rand.Intn(13) + 1
	return uint32(suit*16 + point)
}

func Benchmark_GetCardType(b *testing.B) {
	for i := 0; i < b.N; i++ {
		cards := []uint32{}
		for j := 0; j < 13; j++ {
			cards = append(cards, randCard())
		}
		states.GetCardType(states.ToDDZCards(cards))
	}

}
