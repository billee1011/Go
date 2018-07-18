package tests

import (
	"testing"

	"steve/room/desks/ddzdesk/flow/ddz/states"

	"github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

// 牌的花色
//PS_NONE 		= 0x00;//无花色(大小王)
//PS_DIAMOND 	= 0x10;//方块
//PS_CLUB 		= 0x20;//梅花
//PS_HEART 		= 0x30;//红桃
//PS_SPADE 		= 0x40;//黑桃

// 牌的点数
//PV_A 			= 0x01;//A
//PV_2 			= 0x02;//2
//PV_3 			= 0x03;//3
//PV_4 			= 0x04;//4
//PV_5 			= 0x05;//5
//PV_6 			= 0x06;//6
//PV_7 			= 0x07;//7
//PV_8 			= 0x08;//8
//PV_9 			= 0x09;//9
//PV_10 		= 0x0A;//10
//PV_J 			= 0x0B;//J
//PV_Q 			= 0x0C;//Q
//PV_K 			= 0x0D;//K
//PV_BLACK_JOKER= 0x0E;//小王
//PV_RED_JOKER 	= 0x0F;//大王

// 测试托管的单张
func Test_Tuoguan_Single(t *testing.T) {

	// 手中的牌（96652244）
	handCards := states.ToDDZCards([]uint32{0x29, 0x16, 0x26, 0x45, 0x22, 0x32, 0x24, 0x34})

	// 上次出的牌（3）
	lastCards := states.ToDDZCards([]uint32{0x33})

	// 是否成功
	bSuc, outCards := states.GetMinBiggerSingle(handCards, lastCards)

	// 应该是成功的
	assert.Equal(t, true, bSuc)

	logrus.Errorf("Test_Tuoguan_Single()::outCards = %v", outCards)

	// 出的牌
	outInts := states.ToInts(outCards)

	// 应该是4
	assert.Equal(t, int(0x24), int(outInts[0]))
}

// 测试托管的对子
func Test_Tuoguan_Pair(t *testing.T) {

	// 手中的牌（9665224）
	handCards := states.ToDDZCards([]uint32{0x29, 0x16, 0x26, 0x45, 0x22, 0x32, 0x24})

	// 上次出的牌（33）
	lastCards := states.ToDDZCards([]uint32{0x33, 0x43})

	// 是否成功
	bSuc, outCards := states.GetMinBiggerPair(handCards, lastCards)

	// 应该是成功的
	assert.Equal(t, true, bSuc)

	logrus.Errorf("Test_Tuoguan_Pair()::outCards = %v", outCards)

	// 出的牌
	outInts := states.ToInts(outCards)

	// 应该是66
	assert.Equal(t, int(0x16), int(outInts[0]))
	assert.Equal(t, int(0x26), int(outInts[1]))
}

// 测试托管的顺子
func Test_Tuoguan_Shunzi(t *testing.T) {

	// 手中的牌（987665224）
	handCards := states.ToDDZCards([]uint32{0x29, 0x28, 0x27, 0x16, 0x26, 0x45, 0x22, 0x32, 0x24})

	// 上次出的牌（43567）
	lastCards := states.ToDDZCards([]uint32{0x34, 0x43, 0x45, 0x46, 0x47})

	// 是否成功
	bSuc, outCards := states.GetMinBiggerShunzi(handCards, lastCards)

	// 应该是成功的
	assert.Equal(t, true, bSuc)

	logrus.Errorf("Test_Tuoguan_Shunzi()::outCards = %v", outCards)

	// 出的牌
	outInts := states.ToInts(outCards)

	// 应该是45678
	assert.Equal(t, int(0x24), int(outInts[0]))
	assert.Equal(t, int(0x45), int(outInts[1]))
	assert.Equal(t, int(0x16), int(outInts[2]))
	assert.Equal(t, int(0x27), int(outInts[3]))
	assert.Equal(t, int(0x28), int(outInts[4]))
}
