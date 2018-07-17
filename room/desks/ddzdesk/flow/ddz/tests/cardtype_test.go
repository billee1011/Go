package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"steve/room/desks/ddzdesk/flow/ddz/states"
	"steve/server_pb/ddz"
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
	cards := states.ToDDZCards([]uint32{0x13, 0x23, 0x33, 0x14, 0x24, 0x34, 0x15, 0x25, 0x35, 0x45}) // 3334445555视为飞机带对子
	is, pivot := states.IsTriplesAndPairs(cards)
	assert.Equal(t, true, is)
	assert.Equal(t, uint32(4), pivot.PointWeight)
}

func Test_IsTriplesAndPairs2(t *testing.T) {
	cards := states.ToDDZCards([]uint32{0x13, 0x23, 0x33, 0x14, 0x24, 0x34, 0x15, 0x25, 0x17, 0x37})
	is, pivot := states.IsTriplesAndPairs(cards)
	assert.Equal(t, true, is)
	assert.Equal(t, uint32(4), pivot.PointWeight)
}

func Test_IsTriplesAndSingles1(t *testing.T) {
	cards := states.ToDDZCards([]uint32{0x13, 0x23, 0x33, 0x43, 0x14, 0x24, 0x34, 0x44, 0x15, 0x25, 0x35, 0x45}) // 333344445555视为飞机带翅膀
	is, pivot := states.IsTriplesAndSingles(cards)
	assert.Equal(t, true, is)
	assert.Equal(t, uint32(5), pivot.PointWeight)
}

func Test_IsTriplesAndSingles2(t *testing.T) {
	cards := states.ToDDZCards([]uint32{0x13, 0x23, 0x33, 0x14, 0x24, 0x34, 0x15, 0x25, 0x35, 0x17, 0x29, 0x43}) //333444555793
	is, pivot := states.IsTriplesAndSingles(cards)
	assert.Equal(t, true, is)
	assert.Equal(t, uint32(5), pivot.PointWeight)
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
	assert.Equal(t, ddz.CardType_CT_NONE, cardType)
}

func Test_ShunZi(t *testing.T) {
	cards := states.ToDDZCards([]uint32{0x13, 0x24, 0x35, 0x46, 0x17, 0x28, 0x39, 0x4A, 0x1B, 0x2C, 0x3D, 0x41})
	cardType, _ := states.GetCardType(cards)
	assert.Equal(t, ddz.CardType_CT_SHUNZI, cardType)
}
