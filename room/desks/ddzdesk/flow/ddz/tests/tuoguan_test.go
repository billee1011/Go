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

// 测试托管的连对
func Test_Tuoguan_Pairs(t *testing.T) {

	// 手中的牌（10998877665224）
	handCards := states.ToDDZCards([]uint32{0x2A, 0x19, 0x29, 0x28, 0x18, 0x27, 0x17, 0x16, 0x26, 0x45, 0x22, 0x32, 0x24})

	// 上次出的牌（443355）
	lastCards := states.ToDDZCards([]uint32{0x34, 0x14, 0x43, 0x13, 0x45, 0x35})

	// 是否成功
	bSuc, outCards := states.GetMinBiggerPairs(handCards, lastCards)

	// 应该是成功的
	assert.Equal(t, true, bSuc)

	logrus.Errorf("Test_Tuoguan_Pairs()::outCards = %v", outCards)

	// 出的牌
	outInts := states.ToInts(outCards)

	// 应该是667788
	assert.Equal(t, int(0x16), int(outInts[0]))
	assert.Equal(t, int(0x26), int(outInts[1]))
	assert.Equal(t, int(0x17), int(outInts[2]))
	assert.Equal(t, int(0x27), int(outInts[3]))
	assert.Equal(t, int(0x18), int(outInts[4]))
	assert.Equal(t, int(0x28), int(outInts[5]))
}

// 测试托管的三张
func Test_Tuoguan_Triple(t *testing.T) {

	// 手中的牌（999886665224）
	handCards := states.ToDDZCards([]uint32{0x19, 0x29, 0x39, 0x28, 0x18, 0x16, 0x26, 0x36, 0x45, 0x22, 0x32, 0x24})

	// 上次出的牌（444）
	lastCards := states.ToDDZCards([]uint32{0x34, 0x14, 0x44})

	// 是否成功
	bSuc, outCards := states.GetMinBiggerTriple(handCards, lastCards)

	// 应该是成功的
	assert.Equal(t, true, bSuc)

	logrus.Errorf("Test_Tuoguan_Triple()::outCards = %v", outCards)

	// 出的牌
	outInts := states.ToInts(outCards)

	// 应该是667788
	assert.Equal(t, int(0x16), int(outInts[0]))
	assert.Equal(t, int(0x26), int(outInts[1]))
	assert.Equal(t, int(0x36), int(outInts[2]))
}

// 测试托管的3带1
func Test_Tuoguan_3And1(t *testing.T) {

	// 手中的牌（999886665224）
	handCards := states.ToDDZCards([]uint32{0x19, 0x29, 0x39, 0x28, 0x18, 0x16, 0x26, 0x36, 0x45, 0x22, 0x32, 0x24})

	// 上次出的牌（4445）
	lastCards := states.ToDDZCards([]uint32{0x34, 0x15, 0x14, 0x44})

	// 是否成功
	bSuc, outCards := states.GetMinBigger3And1(handCards, lastCards)

	// 应该是成功的
	assert.Equal(t, true, bSuc)

	logrus.Errorf("Test_Tuoguan_3And1()::outCards = %v", outCards)

	// 出的牌
	outInts := states.ToInts(outCards)

	// 应该是6664
	assert.Equal(t, int(0x16), int(outInts[0]))
	assert.Equal(t, int(0x26), int(outInts[1]))
	assert.Equal(t, int(0x36), int(outInts[2]))
	assert.Equal(t, int(0x24), int(outInts[3]))
}

// 测试托管的3带2
func Test_Tuoguan_3And2(t *testing.T) {

	// 手中的牌（999886665224）
	handCards := states.ToDDZCards([]uint32{0x19, 0x29, 0x39, 0x28, 0x18, 0x16, 0x26, 0x36, 0x45, 0x22, 0x32, 0x24})

	// 上次出的牌（44455）
	lastCards := states.ToDDZCards([]uint32{0x34, 0x15, 0x14, 0x25, 0x44})

	// 是否成功
	bSuc, outCards := states.GetMinBigger3And2(handCards, lastCards)

	// 应该是成功的
	assert.Equal(t, true, bSuc)

	logrus.Errorf("Test_Tuoguan_3And2()::outCards = %v", outCards)

	// 出的牌
	outInts := states.ToInts(outCards)

	// 应该是66688
	assert.Equal(t, int(0x16), int(outInts[0]))
	assert.Equal(t, int(0x26), int(outInts[1]))
	assert.Equal(t, int(0x36), int(outInts[2]))
	assert.Equal(t, int(0x18), int(outInts[3]))
	assert.Equal(t, int(0x28), int(outInts[4]))
}

// 测试托管的飞机
func Test_Tuoguan_Triples(t *testing.T) {

	// 手中的牌（9998886665224）
	handCards := states.ToDDZCards([]uint32{0x19, 0x29, 0x39, 0x28, 0x18, 0x38, 0x16, 0x26, 0x36, 0x45, 0x22, 0x32, 0x24})

	// 上次出的牌（444333）
	lastCards := states.ToDDZCards([]uint32{0x34, 0x14, 0x43, 0x44, 0x23, 0x13})

	// 是否成功
	bSuc, outCards := states.GetMinBiggerTriples(handCards, lastCards)

	// 应该是成功的
	assert.Equal(t, true, bSuc)

	logrus.Errorf("Test_Tuoguan_Triples()::outCards = %v", outCards)

	// 出的牌
	outInts := states.ToInts(outCards)

	// 应该是888999
	assert.Equal(t, int(0x18), int(outInts[0]))
	assert.Equal(t, int(0x28), int(outInts[1]))
	assert.Equal(t, int(0x38), int(outInts[2]))
	assert.Equal(t, int(0x19), int(outInts[3]))
	assert.Equal(t, int(0x29), int(outInts[4]))
	assert.Equal(t, int(0x39), int(outInts[5]))
}

// 测试托管的飞机带单张
func Test_Tuoguan_3sAnd1s(t *testing.T) {

	// 手中的牌（9998886665224）
	handCards := states.ToDDZCards([]uint32{0x19, 0x29, 0x39, 0x28, 0x18, 0x38, 0x16, 0x26, 0x36, 0x45, 0x22, 0x32, 0x24})

	// 上次出的牌（4443339J）
	lastCards := states.ToDDZCards([]uint32{0x34, 0x14, 0x43, 0x44, 0x23, 0x13, 0x49, 0x1B})

	// 是否成功
	bSuc, outCards := states.GetMinBigger3sAnd1s(handCards, lastCards)

	// 应该是成功的
	assert.Equal(t, true, bSuc)

	logrus.Errorf("Test_Tuoguan_3sAnd1s()::outCards = %v", outCards)

	// 出的牌
	outInts := states.ToInts(outCards)

	// 应该是88899945
	assert.Equal(t, int(0x18), int(outInts[0]))
	assert.Equal(t, int(0x28), int(outInts[1]))
	assert.Equal(t, int(0x38), int(outInts[2]))
	assert.Equal(t, int(0x19), int(outInts[3]))
	assert.Equal(t, int(0x29), int(outInts[4]))
	assert.Equal(t, int(0x39), int(outInts[5]))
	assert.Equal(t, int(0x24), int(outInts[6]))
	assert.Equal(t, int(0x45), int(outInts[7]))
}

// 测试托管的飞机带对子
func Test_Tuoguan_3sAnd2s(t *testing.T) {

	// 手中的牌（9998886665224）
	handCards := states.ToDDZCards([]uint32{0x19, 0x29, 0x39, 0x28, 0x18, 0x38, 0x16, 0x26, 0x36, 0x45, 0x22, 0x32, 0x24})

	// 上次出的牌（44433377JJ）
	lastCards := states.ToDDZCards([]uint32{0x34, 0x14, 0x43, 0x44, 0x23, 0x13, 0x17, 0x1B, 0x27, 0x2B})

	// 是否成功
	bSuc, outCards := states.GetMinBigger3sAnd2s(handCards, lastCards)

	// 应该是成功的
	assert.Equal(t, true, bSuc)

	logrus.Errorf("Test_Tuoguan_3sAnd2s()::outCards = %v", outCards)

	// 出的牌
	outInts := states.ToInts(outCards)

	// 应该是8889996622
	assert.Equal(t, int(0x18), int(outInts[0]))
	assert.Equal(t, int(0x28), int(outInts[1]))
	assert.Equal(t, int(0x38), int(outInts[2]))
	assert.Equal(t, int(0x19), int(outInts[3]))
	assert.Equal(t, int(0x29), int(outInts[4]))
	assert.Equal(t, int(0x39), int(outInts[5]))
	assert.Equal(t, int(0x16), int(outInts[6]))
	assert.Equal(t, int(0x26), int(outInts[7]))
	assert.Equal(t, int(0x22), int(outInts[8]))
	assert.Equal(t, int(0x32), int(outInts[9]))
}

// 测试托管的4带2单张
func Test_Tuoguan_4sAnd1s(t *testing.T) {

	// 手中的牌（9988886666224）
	handCards := states.ToDDZCards([]uint32{0x19, 0x29, 0x48, 0x28, 0x18, 0x38, 0x16, 0x26, 0x36, 0x46, 0x22, 0x32, 0x24})

	// 上次出的牌（7777Q9）
	lastCards := states.ToDDZCards([]uint32{0x17, 0x1C, 0x27, 0x37, 0x47, 0x19})

	// 是否成功
	bSuc, outCards := states.GetMinBigger4sAnd1s(handCards, lastCards)

	// 应该是成功的
	assert.Equal(t, true, bSuc)

	logrus.Errorf("Test_Tuoguan_4sAnd1s()::outCards = %v", outCards)

	// 出的牌
	outInts := states.ToInts(outCards)

	// 应该是888846
	assert.Equal(t, int(0x18), int(outInts[0]))
	assert.Equal(t, int(0x28), int(outInts[1]))
	assert.Equal(t, int(0x38), int(outInts[2]))
	assert.Equal(t, int(0x48), int(outInts[3]))
	assert.Equal(t, int(0x24), int(outInts[4]))
	assert.Equal(t, int(0x16), int(outInts[5]))
}

// 测试托管的4带2对子
func Test_Tuoguan_4sAnd2s(t *testing.T) {

	// 手中的牌（9988886666224）
	handCards := states.ToDDZCards([]uint32{0x19, 0x29, 0x48, 0x28, 0x18, 0x38, 0x16, 0x26, 0x36, 0x46, 0x22, 0x32, 0x24})

	// 上次出的牌（77779933）
	lastCards := states.ToDDZCards([]uint32{0x17, 0x27, 0x37, 0x47, 0x39, 0x39, 0x33, 0x23})

	// 是否成功
	bSuc, outCards := states.GetMinBigger4sAnd2s(handCards, lastCards)

	// 应该是成功的
	assert.Equal(t, true, bSuc)

	logrus.Errorf("Test_Tuoguan_4sAnd2s()::outCards = %v", outCards)

	// 出的牌
	outInts := states.ToInts(outCards)

	// 应该是88886699
	assert.Equal(t, int(0x18), int(outInts[0]))
	assert.Equal(t, int(0x28), int(outInts[1]))
	assert.Equal(t, int(0x38), int(outInts[2]))
	assert.Equal(t, int(0x48), int(outInts[3]))
	assert.Equal(t, int(0x16), int(outInts[4]))
	assert.Equal(t, int(0x26), int(outInts[5]))
	assert.Equal(t, int(0x19), int(outInts[6]))
	assert.Equal(t, int(0x29), int(outInts[7]))
}

// 测试托管的炸弹
func Test_Tuoguan_Boom(t *testing.T) {

	// 手中的牌（9988886666224）
	handCards := states.ToDDZCards([]uint32{0x19, 0x29, 0x48, 0x28, 0x18, 0x38, 0x16, 0x26, 0x36, 0x46, 0x22, 0x32, 0x24})

	// 上次出的牌（7777）
	lastCards := states.ToDDZCards([]uint32{0x17, 0x27, 0x37, 0x47})

	// 是否成功
	bSuc, outCards := states.GetMinBiggerBoom(handCards, lastCards)

	// 应该是成功的
	assert.Equal(t, true, bSuc)

	logrus.Errorf("Test_Tuoguan_Boom()::outCards = %v", outCards)

	// 出的牌
	outInts := states.ToInts(outCards)

	// 应该是8888
	assert.Equal(t, int(0x18), int(outInts[0]))
	assert.Equal(t, int(0x28), int(outInts[1]))
	assert.Equal(t, int(0x38), int(outInts[2]))
	assert.Equal(t, int(0x48), int(outInts[3]))
}
