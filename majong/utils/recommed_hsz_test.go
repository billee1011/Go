package utils

import (
	"fmt"
	majongpb "steve/entity/majong"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Benchmark_CardNum_sanTiao 牌数比较 有一门牌数=三张，一门>=5张，一门<3张 提起三张条
func Benchmark_CardNum_sanTiao(t *testing.B) {
	//手牌有3张条+11张筒,
	hanCard := []Card{21, 27, 28, 31, 31, 32, 32, 33, 34, 35, 36, 37, 38, 39}
	cards, _ := CheckHuUtilCardsToHandCards(hanCard)
	hszCards := GetRecommedHuanSanZhang(cards)
	assert.Equal(t, 3, len(hszCards))
	comCards, _ := CheckHuUtilCardsToHandCards([]Card{21, 27, 28})
	assert.Equal(t, hszCards, comCards)
	fmt.Println(hszCards)

	//手牌有3张条+11张万
	hanCard2 := []Card{22, 23, 22, 11, 12, 13, 14, 15, 16, 17, 18, 19, 19, 11}
	cards2, _ := CheckHuUtilCardsToHandCards(hanCard2)
	hszCards2 := GetRecommedHuanSanZhang(cards2)
	assert.Equal(t, 3, len(hszCards2))
	comCards2, _ := CheckHuUtilCardsToHandCards([]Card{22, 23, 22})
	assert.Equal(t, hszCards2, comCards2)
	fmt.Println(hszCards2)

	//手牌有2张万+3张条+9张筒
	hanCard3 := []Card{12, 11, 25, 25, 25, 31, 32, 33, 34, 35, 36, 37, 38, 39}
	cards3, _ := CheckHuUtilCardsToHandCards(hanCard3)
	hszCards3 := GetRecommedHuanSanZhang(cards3)
	assert.Equal(t, 3, len(hszCards3))
	comCards3, _ := CheckHuUtilCardsToHandCards([]Card{25, 25, 25})
	assert.Equal(t, hszCards3, comCards3)
	fmt.Println(hszCards3)

	//手牌有1张万+3张条+10张筒
	hanCard4 := []Card{12, 29, 25, 25, 36, 31, 32, 33, 34, 35, 36, 37, 38, 39}
	cards4, _ := CheckHuUtilCardsToHandCards(hanCard4)
	hszCards4 := GetRecommedHuanSanZhang(cards4)
	assert.Equal(t, 3, len(hszCards4))
	comCards4, _ := CheckHuUtilCardsToHandCards([]Card{29, 25, 25})
	assert.Equal(t, hszCards4, comCards4)
	fmt.Println(hszCards4)

	//手牌有1张筒+3张条+10张万
	hanCard5 := []Card{33, 26, 29, 21, 11, 12, 13, 14, 15, 16, 17, 18, 19, 13}
	cards5, _ := CheckHuUtilCardsToHandCards(hanCard5)
	hszCards5 := GetRecommedHuanSanZhang(cards5)
	assert.Equal(t, 3, len(hszCards5))
	comCards5, _ := CheckHuUtilCardsToHandCards([]Card{26, 29, 21})
	assert.Equal(t, hszCards5, comCards5)
	fmt.Println(hszCards5)

	//手牌有2张筒+3张条+9张万
	hanCard6 := []Card{31, 39, 21, 25, 29, 11, 12, 13, 14, 15, 16, 17, 18, 19}
	cards6, _ := CheckHuUtilCardsToHandCards(hanCard6)
	hszCards6 := GetRecommedHuanSanZhang(cards6)
	assert.Equal(t, 3, len(hszCards6))
	comCards6, _ := CheckHuUtilCardsToHandCards([]Card{21, 25, 29})
	assert.Equal(t, hszCards6, comCards6)
	fmt.Println(hszCards6)

	//手牌有2张筒+3张条+8张万
	hanCard7 := []Card{33, 33, 24, 27, 25, 11, 12, 13, 14, 15, 16, 17, 18}
	cards7, _ := CheckHuUtilCardsToHandCards(hanCard7)
	hszCards7 := GetRecommedHuanSanZhang(cards7)
	assert.Equal(t, 3, len(hszCards7))
	comCards7, _ := CheckHuUtilCardsToHandCards([]Card{24, 27, 25})
	assert.Equal(t, hszCards7, comCards7)
	fmt.Println(hszCards7)

	//手牌中有3张条+5张万+5张筒
	hanCard8 := []Card{21, 23, 25, 11, 12, 17, 18, 19, 31, 31, 31, 32, 32}
	cards8, _ := CheckHuUtilCardsToHandCards(hanCard8)
	hszCards8 := GetRecommedHuanSanZhang(cards8)
	assert.Equal(t, 3, len(hszCards8))
	comCards8, _ := CheckHuUtilCardsToHandCards([]Card{21, 23, 25})
	assert.Equal(t, hszCards8, comCards8)
	fmt.Println(hszCards8)
}

// Test_CardNum_sanWan 牌数比较 只有一门牌数大于三张
func Test_CardNum_sanWan(t *testing.T) {
	//手牌14张万
	hanCard := []Card{11, 11, 11, 11, 12, 12, 12, 12, 13, 13, 13, 13, 14, 14}
	cards, _ := CheckHuUtilCardsToHandCards(hanCard)
	hszCards := GetRecommedHuanSanZhang(cards)
	assert.Equal(t, 3, len(hszCards))
	comCards, _ := CheckHuUtilCardsToHandCards([]Card{11, 14, 12})
	if hszCards[2].GetPoint() == 3 {
		comCards, _ = CheckHuUtilCardsToHandCards([]Card{11, 14, 13})
	}
	assert.Equal(t, hszCards, comCards)
	fmt.Println(hszCards)

	//手牌1张条+13张万
	hanCard2 := []Card{21, 11, 11, 11, 12, 12, 12, 13, 13, 13, 14, 14, 14, 15}
	cards2, _ := CheckHuUtilCardsToHandCards(hanCard2)
	hszCards2 := GetRecommedHuanSanZhang(cards2)
	assert.Equal(t, 3, len(hszCards2))
	comCards2, _ := CheckHuUtilCardsToHandCards([]Card{11, 15, 13})
	assert.Equal(t, hszCards2, comCards2)
	fmt.Println(hszCards2)

	//手牌1张条+1张万+11张筒
	hanCard3 := []Card{29, 11, 31, 32, 33, 34, 35, 36, 37, 38, 39, 34, 34, 35}
	cards3, _ := CheckHuUtilCardsToHandCards(hanCard3)
	hszCards3 := GetRecommedHuanSanZhang(cards3)
	assert.Equal(t, 3, len(hszCards3))
	comCards3, _ := CheckHuUtilCardsToHandCards([]Card{31, 39, 35})
	assert.Equal(t, hszCards3, comCards3)
	fmt.Println(hszCards3)

	//手牌1张条+2张万+10张筒
	hanCard4 := []Card{23, 14, 19, 31, 32, 33, 34, 35, 36, 37, 38, 38, 35}
	cards4, _ := CheckHuUtilCardsToHandCards(hanCard4)
	hszCards4 := GetRecommedHuanSanZhang(cards4)
	assert.Equal(t, 3, len(hszCards4))
	comCards4, _ := CheckHuUtilCardsToHandCards([]Card{31, 38, 34})
	if hszCards4[2].GetPoint() == 5 {
		comCards4, _ = CheckHuUtilCardsToHandCards([]Card{31, 38, 35})
	}
	assert.Equal(t, hszCards4, comCards4)
	fmt.Println(hszCards4)

	// 手牌2张条+12张万
	hanCard5 := []Card{21, 23, 11, 12, 13, 14, 14, 16, 17, 18, 19, 11, 12}
	cards5, _ := CheckHuUtilCardsToHandCards(hanCard5)
	hszCards5 := GetRecommedHuanSanZhang(cards5)
	assert.Equal(t, 3, len(hszCards5))
	comCards5, _ := CheckHuUtilCardsToHandCards([]Card{11, 19, 14})
	if hszCards5[2].GetPoint() == 6 {
		comCards5, _ = CheckHuUtilCardsToHandCards([]Card{11, 19, 16})
	}
	assert.Equal(t, hszCards5, comCards5)
	fmt.Println(hszCards5)

	// 手牌有2张条+2张万+9张筒
	hanCard6 := []Card{21, 23, 11, 19, 31, 32, 33, 33, 33, 37, 37, 37, 37}
	cards6, _ := CheckHuUtilCardsToHandCards(hanCard6)
	hszCards6 := GetRecommedHuanSanZhang(cards6)
	assert.Equal(t, 3, len(hszCards6))
	comCards6, _ := CheckHuUtilCardsToHandCards([]Card{31, 37, 33})
	assert.Equal(t, hszCards6, comCards6)
	fmt.Println(hszCards6)
}

// Test_CardNum_YouLianMenPaiDaSanQieChaDaYuDeng2 牌数比较 有两门牌数大于三张且相差>=2张
func Test_CardNum_YouLianMenPaiDaSanQieChaDaYuDeng2(t *testing.T) {
	//手牌中有1张筒+4张万+8张条
	hanCard := []Card{31, 11, 11, 12, 13, 21, 22, 23, 24, 25, 26, 27, 28}
	cards, _ := CheckHuUtilCardsToHandCards(hanCard)
	hszCards := GetRecommedHuanSanZhang(cards)
	assert.Equal(t, 3, len(hszCards))
	comCards, _ := CheckHuUtilCardsToHandCards([]Card{11, 13, 12})
	assert.Equal(t, hszCards, comCards)
	fmt.Println(hszCards)

	// 手牌中有1张筒+5张万+7张条
	hanCard2 := []Card{35, 11, 15, 19, 19, 19, 21, 22, 23, 24, 25, 26, 27}
	cards2, _ := CheckHuUtilCardsToHandCards(hanCard2)
	hszCards2 := GetRecommedHuanSanZhang(cards2)
	assert.Equal(t, 3, len(hszCards2))
	comCards2, _ := CheckHuUtilCardsToHandCards([]Card{11, 19, 15})
	assert.Equal(t, hszCards2, comCards2)
	fmt.Println(hszCards2)

	// 手牌6张万+8张条
	hanCard3 := []Card{11, 11, 12, 13, 14, 15, 21, 21, 22, 22, 23, 23, 24, 24}
	cards3, _ := CheckHuUtilCardsToHandCards(hanCard3)
	hszCards3 := GetRecommedHuanSanZhang(cards3)
	assert.Equal(t, 3, len(hszCards3))
	comCards3, _ := CheckHuUtilCardsToHandCards([]Card{11, 15, 13})
	assert.Equal(t, hszCards3, comCards3)
	fmt.Println(hszCards3)
}

// Test_PaiXing_Xiao2 牌型比较 牌型一样 有两门牌数大于三张且相差<2张
func Test_PaiXing_PaiXingYiYang_Xiao2(t *testing.T) {
	//随机选一种花色定牌
	hanCard := []Card{11, 11, 12, 12, 13, 13, 21, 21, 22, 22, 23, 23, 14, 24}
	cards, _ := CheckHuUtilCardsToHandCards(hanCard)
	hszCards := GetRecommedHuanSanZhang(cards)
	assert.Equal(t, 3, len(hszCards))
	comCards, _ := CheckHuUtilCardsToHandCards([]Card{11, 14, 12})
	if hszCards[0].GetColor() == majongpb.CardColor_ColorTiao {
		comCards, _ = CheckHuUtilCardsToHandCards([]Card{21, 24, 22})
		if hszCards[2].GetPoint() == 3 {
			comCards, _ = CheckHuUtilCardsToHandCards([]Card{21, 24, 23})
		}
	} else if hszCards[2].GetPoint() == 3 {
		comCards, _ = CheckHuUtilCardsToHandCards([]Card{11, 14, 13})
	}
	assert.Equal(t, hszCards, comCards)
	fmt.Println(hszCards)

	//牌数不一样,选择牌数少的一种花色定牌
	hanCard2 := []Card{11, 11, 12, 12, 13, 13, 21, 21, 22, 22, 23, 23, 29, 31}
	cards2, _ := CheckHuUtilCardsToHandCards(hanCard2)
	hszCards2 := GetRecommedHuanSanZhang(cards2)
	assert.Equal(t, 3, len(hszCards2))
	comCards2, _ := CheckHuUtilCardsToHandCards([]Card{11, 13, 12})
	assert.Equal(t, hszCards2, comCards2)
	fmt.Println(hszCards2)
}

// Test_PaiXing_PaiXingBuYiYang_Xiao2_3Vs3 牌型比较 牌型不一样 有两门牌数大于三张且相差<2张 3张vs3张
func Test_PaiXing_PaiXingBuYiYang_Xiao2_3Vs3(t *testing.T) {
	// 刻vs顺 111vs234
	hanCard := []Card{11, 11, 11, 22, 23, 24}
	cards, _ := CheckHuUtilCardsToHandCards(hanCard)
	hszCards := GetRecommedHuanSanZhang(cards)
	assert.Equal(t, 3, len(hszCards))
	comCards, _ := CheckHuUtilCardsToHandCards([]Card{22, 23, 24})
	assert.Equal(t, hszCards, comCards)
	fmt.Println(hszCards)
	//顺vs对 123vs112
	hanCard2 := []Card{22, 23, 24, 31, 31, 32}
	cards2, _ := CheckHuUtilCardsToHandCards(hanCard2)
	hszCards2 := GetRecommedHuanSanZhang(cards2)
	assert.Equal(t, 3, len(hszCards2))
	comCards2, _ := CheckHuUtilCardsToHandCards([]Card{31, 31, 32})
	assert.Equal(t, hszCards2, comCards2)
	fmt.Println(hszCards2)
	// 刻vs对+单 333vs233
	hanCard3 := []Card{23, 23, 23, 12, 13, 13}
	cards3, _ := CheckHuUtilCardsToHandCards(hanCard3)
	hszCards3 := GetRecommedHuanSanZhang(cards3)
	assert.Equal(t, 3, len(hszCards3))
	comCards3, _ := CheckHuUtilCardsToHandCards([]Card{12, 13, 13})
	assert.Equal(t, hszCards3, comCards3)
	fmt.Println(hszCards3)
}

// Test_PaiXing_PaiXingBuYiYang_Xiao2_3Vs4 牌型比较 牌型不一样 有两门牌数大于三张且相差<2张 3张vs4张
func Test_PaiXing_PaiXingBuYiYang_Xiao2_3Vs4(t *testing.T) {
	//杠vs刻 1111vs222
	hanCard := []Card{11, 11, 11, 11, 22, 22, 22}
	cards, _ := CheckHuUtilCardsToHandCards(hanCard)
	hszCards := GetRecommedHuanSanZhang(cards)
	assert.Equal(t, 3, len(hszCards))
	comCards, _ := CheckHuUtilCardsToHandCards([]Card{22, 22, 22})
	assert.Equal(t, hszCards, comCards)
	fmt.Println(hszCards)
	//刻vs两对 222vs3344
	hanCard2 := []Card{22, 22, 22, 33, 33, 34, 34}
	cards2, _ := CheckHuUtilCardsToHandCards(hanCard2)
	hszCards2 := GetRecommedHuanSanZhang(cards2)
	assert.Equal(t, 3, len(hszCards2))
	comCards2, _ := CheckHuUtilCardsToHandCards([]Card{33, 34, 33})
	if hszCards2[2].GetPoint() == 4 {
		comCards2, _ = CheckHuUtilCardsToHandCards([]Card{33, 34, 34})
	}
	assert.Equal(t, hszCards2, comCards2)
	fmt.Println(hszCards2)
	// 对+单vs顺 1123vs123
	hanCard3 := []Card{21, 21, 25, 29, 12, 11, 13}
	cards3, _ := CheckHuUtilCardsToHandCards(hanCard3)
	hszCards3 := GetRecommedHuanSanZhang(cards3)
	assert.Equal(t, 3, len(hszCards3))
	comCards3, _ := CheckHuUtilCardsToHandCards([]Card{21, 29, 25})
	assert.Equal(t, hszCards3, comCards3)
	fmt.Println(hszCards3)
}

// Test_PaiXing_PaiXingBuYiYang_Xiao2_4Vs4 牌型比较 牌型不一样 有两门牌数大于三张且相差<2张 4张vs4张
func Test_PaiXing_PaiXingBuYiYang_Xiao2_4Vs4(t *testing.T) {
	//刻+单vs两对 1888vs1188
	hanCard := []Card{11, 18, 18, 18, 21, 21, 28, 28}
	cards, _ := CheckHuUtilCardsToHandCards(hanCard)
	hszCards := GetRecommedHuanSanZhang(cards)
	assert.Equal(t, 3, len(hszCards))
	comCards, _ := CheckHuUtilCardsToHandCards([]Card{21, 28, 21})
	if hszCards[2].GetPoint() == 8 {
		comCards, _ = CheckHuUtilCardsToHandCards([]Card{21, 28, 28})
	}
	assert.Equal(t, hszCards, comCards)
	fmt.Println(hszCards)
	//刻+单vs顺+单 1222vs4568
	hanCard2 := []Card{11, 12, 12, 12, 24, 25, 26, 28}
	cards2, _ := CheckHuUtilCardsToHandCards(hanCard2)
	hszCards2 := GetRecommedHuanSanZhang(cards2)
	assert.Equal(t, 3, len(hszCards2))
	comCards2, _ := CheckHuUtilCardsToHandCards([]Card{24, 28, 26})
	assert.Equal(t, hszCards2, comCards2)
	fmt.Println(hszCards2)
	// 杠vs刻+单 2222vs7888
	hanCard3 := []Card{22, 22, 22, 22, 17, 18, 18, 18}
	cards3, _ := CheckHuUtilCardsToHandCards(hanCard3)
	hszCards3 := GetRecommedHuanSanZhang(cards3)
	assert.Equal(t, 3, len(hszCards3))
	comCards3, _ := CheckHuUtilCardsToHandCards([]Card{17, 18, 18})
	assert.Equal(t, hszCards3, comCards3)
	fmt.Println(hszCards3)
	//两对vs顺+单 4455vs3456
	hanCard5 := []Card{14, 14, 15, 15, 23, 24, 25, 26}
	cards5, _ := CheckHuUtilCardsToHandCards(hanCard5)
	hszCards5 := GetRecommedHuanSanZhang(cards5)
	assert.Equal(t, 3, len(hszCards5))
	comCards5, _ := CheckHuUtilCardsToHandCards([]Card{23, 26, 24})
	if hszCards5[2].GetPoint() == 5 {
		comCards5, _ = CheckHuUtilCardsToHandCards([]Card{23, 26, 25})
	}
	assert.Equal(t, hszCards5, comCards5)
	fmt.Println(hszCards5)
}

// Test_PaiXing_PaiXingBuYiYang_Xiao2_4Vs4Vs5 牌型比较 牌型不一样 有两门牌数大于三张且相差<2张 4张vs4张vs5张
func Test_PaiXing_PaiXingBuYiYang_Xiao2_4Vs4Vs5(t *testing.T) {
	//杠vs刻+单vs顺+对 2222vs4448vs12344
	hanCard := []Card{12, 12, 12, 12, 24, 24, 24, 28, 31, 32, 33, 34, 34}
	cards, _ := CheckHuUtilCardsToHandCards(hanCard)
	hszCards := GetRecommedHuanSanZhang(cards)
	assert.Equal(t, 3, len(hszCards))
	comCards, _ := CheckHuUtilCardsToHandCards([]Card{24, 28, 24})
	assert.Equal(t, hszCards, comCards)
	fmt.Println(hszCards)
	//两对vs刻子+单牌vs顺子+单牌 2233vs1444vs3455
	hanCard2 := []Card{12, 12, 13, 13, 21, 24, 24, 24, 33, 34, 35, 35}
	cards2, _ := CheckHuUtilCardsToHandCards(hanCard2)
	hszCards2 := GetRecommedHuanSanZhang(cards2)
	assert.Equal(t, 3, len(hszCards2))
	comCards2, _ := CheckHuUtilCardsToHandCards([]Card{33, 35, 34})
	assert.Equal(t, hszCards2, comCards2)
	fmt.Println(hszCards2)
	// 杠vs杠+单vs两对 1111vs11112vs3355
	hanCard3 := []Card{11, 11, 11, 11, 21, 21, 21, 21, 22, 33, 33, 35, 35}
	cards3, _ := CheckHuUtilCardsToHandCards(hanCard3)
	hszCards3 := GetRecommedHuanSanZhang(cards3)
	assert.Equal(t, 3, len(hszCards3))
	comCards3, _ := CheckHuUtilCardsToHandCards([]Card{33, 35, 33})
	if hszCards3[2].GetPoint() == 5 {
		comCards3, _ = CheckHuUtilCardsToHandCards([]Card{33, 35, 35})
	}
	assert.Equal(t, hszCards3, comCards3)
	fmt.Println(hszCards3)
	//两对vs顺子+单vs对子+单 6677vs1345vs2245
	hanCard5 := []Card{16, 16, 17, 17, 21, 23, 24, 25, 32, 32, 34, 35}
	cards5, _ := CheckHuUtilCardsToHandCards(hanCard5)
	hszCards5 := GetRecommedHuanSanZhang(cards5)
	assert.Equal(t, 3, len(hszCards5))
	comCards5, _ := CheckHuUtilCardsToHandCards([]Card{32, 35, 34})
	assert.Equal(t, hszCards5, comCards5)
	fmt.Println(hszCards5)
}

// Test_PaiXing_PaiXingBuYiYang_Xiao2_4Vs5Vs5 牌型比较 牌型不一样 有两门牌数大于三张且相差<2张 4张vs5张vs5张
func Test_PaiXing_PaiXingBuYiYang_Xiao2_4Vs5Vs5(t *testing.T) {
	//杠vs顺子+对子vs刻子+单牌 1111vs56788vs4445
	hanCard := []Card{11, 11, 11, 11, 25, 26, 27, 28, 28, 34, 34, 34, 35}
	cards, _ := CheckHuUtilCardsToHandCards(hanCard)
	hszCards := GetRecommedHuanSanZhang(cards)
	assert.Equal(t, 3, len(hszCards))
	comCards, _ := CheckHuUtilCardsToHandCards([]Card{34, 35, 34})
	assert.Equal(t, hszCards, comCards)
	fmt.Println(hszCards)
	//顺子+单牌vs顺子+对子vs刻子+对子 1567vs34555vs66677
	hanCard2 := []Card{11, 15, 16, 17, 23, 24, 25, 25, 25, 36, 36, 36, 37, 37}
	cards2, _ := CheckHuUtilCardsToHandCards(hanCard2)
	hszCards2 := GetRecommedHuanSanZhang(cards2)
	assert.Equal(t, 3, len(hszCards2))
	comCards2, _ := CheckHuUtilCardsToHandCards([]Card{11, 17, 15})
	assert.Equal(t, hszCards2, comCards2)
	fmt.Println(hszCards2)
	// 对子+单牌vs顺子+单牌vs单牌 3356vs45678vs13579
	hanCard3 := []Card{13, 13, 15, 16, 24, 25, 26, 27, 28, 31, 33, 35, 37, 39}
	cards3, _ := CheckHuUtilCardsToHandCards(hanCard3)
	hszCards3 := GetRecommedHuanSanZhang(cards3)
	assert.Equal(t, 3, len(hszCards3))
	comCards3, _ := CheckHuUtilCardsToHandCards([]Card{31, 39, 35})
	assert.Equal(t, hszCards3, comCards3)
	fmt.Println(hszCards3)
	//两对vs刻子+对子vs顺子+对子 3344vs22244vs55678
	hanCard5 := []Card{13, 13, 14, 14, 22, 22, 22, 24, 24, 35, 35, 36, 37, 38}
	cards5, _ := CheckHuUtilCardsToHandCards(hanCard5)
	hszCards5 := GetRecommedHuanSanZhang(cards5)
	assert.Equal(t, 3, len(hszCards5))
	comCards5, _ := CheckHuUtilCardsToHandCards([]Card{13, 14, 13})
	if hszCards5[2].GetPoint() == 4 {
		comCards5, _ = CheckHuUtilCardsToHandCards([]Card{13, 14, 14})
	}
	assert.Equal(t, hszCards5, comCards5)
	fmt.Println(hszCards5)
}

// Test_PaiXing_PaiXingBuYiYang_Xiao2_5Vs6 牌型比较 牌型不一样 有两门牌数大于三张且相差<2张 5张vs6张
func Test_PaiXing_PaiXingBuYiYang_Xiao2_5Vs6(t *testing.T) {
	//杠+对vs杠+单 222244vs55556
	hanCard := []Card{12, 12, 12, 12, 14, 14, 25, 25, 25, 25, 26}
	cards, _ := CheckHuUtilCardsToHandCards(hanCard)
	hszCards := GetRecommedHuanSanZhang(cards)
	assert.Equal(t, 3, len(hszCards))
	comCards, _ := CheckHuUtilCardsToHandCards([]Card{25, 26, 25})
	assert.Equal(t, hszCards, comCards)
	fmt.Println(hszCards)
	//杠+单vs三对 55556vs556677?
	hanCard2 := []Card{15, 15, 15, 15, 16, 25, 25, 26, 26, 27, 27}
	cards2, _ := CheckHuUtilCardsToHandCards(hanCard2)
	hszCards2 := GetRecommedHuanSanZhang(cards2)
	assert.Equal(t, 3, len(hszCards2))
	comCards2, _ := CheckHuUtilCardsToHandCards([]Card{25, 27, 26})
	assert.Equal(t, hszCards2, comCards2)
	fmt.Println(hszCards2)
	// 三对VS俩顺 335577VS123456
	hanCard3 := []Card{11, 11, 12, 12, 13, 13, 21, 22, 23, 24, 25, 26}
	cards3, _ := CheckHuUtilCardsToHandCards(hanCard3)
	hszCards3 := GetRecommedHuanSanZhang(cards3)
	assert.Equal(t, 3, len(hszCards3))
	comCards3, _ := CheckHuUtilCardsToHandCards([]Card{21, 26, 24})
	if hszCards3[2].GetPoint() == 3 {
		comCards3, _ = CheckHuUtilCardsToHandCards([]Card{21, 26, 23})
	}
	assert.Equal(t, hszCards3, comCards3)
	fmt.Println(hszCards3)
	//三对vs顺+对 113355vs22345
	hanCard5 := []Card{11, 11, 13, 13, 15, 15, 22, 22, 23, 24, 25}
	cards5, _ := CheckHuUtilCardsToHandCards(hanCard5)
	hszCards5 := GetRecommedHuanSanZhang(cards5)
	assert.Equal(t, 3, len(hszCards5))
	comCards5, _ := CheckHuUtilCardsToHandCards([]Card{22, 25, 23})
	if hszCards5[2].GetPoint() == 4 {
		comCards5, _ = CheckHuUtilCardsToHandCards([]Card{22, 25, 24})
	}
	assert.Equal(t, hszCards5, comCards5)
	fmt.Println(hszCards5)

	//三对vs两对+单 557799vs33557
	hanCard6 := []Card{15, 15, 17, 17, 19, 19, 23, 23, 25, 25, 27}
	cards6, _ := CheckHuUtilCardsToHandCards(hanCard6)
	hszCards6 := GetRecommedHuanSanZhang(cards6)
	assert.Equal(t, 3, len(hszCards6))
	comCards6, _ := CheckHuUtilCardsToHandCards([]Card{23, 27, 25})
	assert.Equal(t, hszCards6, comCards6)
	fmt.Println(hszCards6)
}

// Test_PaiXing_PaiXingBuYiYang_Xiao2_6Vs6 牌型比较 牌型不一样 有两门牌数大于三张且相差<2张 6张vs6张
func Test_PaiXing_PaiXingBuYiYang_Xiao2_6Vs6(t *testing.T) {
	//杠+对vs两刻 444455vs333444
	hanCard := []Card{14, 14, 14, 14, 15, 15, 23, 23, 23, 24, 24, 24}
	cards, _ := CheckHuUtilCardsToHandCards(hanCard)
	hszCards := GetRecommedHuanSanZhang(cards)
	assert.Equal(t, 3, len(hszCards))
	comCards, _ := CheckHuUtilCardsToHandCards([]Card{23, 24, 23})
	if hszCards[2].GetPoint() == 4 {
		comCards, _ = CheckHuUtilCardsToHandCards([]Card{23, 24, 24})
	}
	assert.Equal(t, hszCards, comCards)
	fmt.Println(hszCards)
	//三对vs六单 668899vs235689
	hanCard2 := []Card{16, 16, 18, 18, 19, 19, 22, 23, 25, 26, 28, 29}
	cards2, _ := CheckHuUtilCardsToHandCards(hanCard2)
	hszCards2 := GetRecommedHuanSanZhang(cards2)
	assert.Equal(t, 3, len(hszCards2))
	comCards2, _ := CheckHuUtilCardsToHandCards([]Card{22, 29, 26})
	if hszCards2[2].GetPoint() == 5 {
		comCards2, _ = CheckHuUtilCardsToHandCards([]Card{22, 29, 25})
	}
	assert.Equal(t, hszCards2, comCards2)
	fmt.Println(hszCards2)
	// 刻+对+单vs三对 222334vs557799
	hanCard3 := []Card{12, 12, 12, 13, 13, 14, 25, 25, 27, 27, 29, 29}
	cards3, _ := CheckHuUtilCardsToHandCards(hanCard3)
	hszCards3 := GetRecommedHuanSanZhang(cards3)
	assert.Equal(t, 3, len(hszCards3))
	comCards3, _ := CheckHuUtilCardsToHandCards([]Card{25, 29, 27})
	assert.Equal(t, hszCards3, comCards3)
	fmt.Println(hszCards3)
	//刻+顺vs刻+对+单 111234vs555667
	hanCard5 := []Card{11, 11, 11, 12, 13, 14, 25, 25, 25, 26, 26, 27}
	cards5, _ := CheckHuUtilCardsToHandCards(hanCard5)
	hszCards5 := GetRecommedHuanSanZhang(cards5)
	assert.Equal(t, 3, len(hszCards5))
	comCards5, _ := CheckHuUtilCardsToHandCards([]Card{25, 27, 26})
	assert.Equal(t, hszCards5, comCards5)
	fmt.Println(hszCards5)
}

// Test_PaiXing_PaiXingBuYiYang_Xiao2_6Vs7 牌型比较 牌型不一样 有两门牌数大于三张且相差<2张 6张vs7张
func Test_PaiXing_PaiXingBuYiYang_Xiao2_6Vs7(t *testing.T) {
	//杠+刻vs杠+对 2222333vs333344
	hanCard := []Card{12, 12, 12, 12, 13, 13, 13, 23, 23, 23, 23, 24, 24}
	cards, _ := CheckHuUtilCardsToHandCards(hanCard)
	hszCards := GetRecommedHuanSanZhang(cards)
	assert.Equal(t, 3, len(hszCards))
	comCards, _ := CheckHuUtilCardsToHandCards([]Card{23, 24, 23})
	if hszCards[2].GetPoint() == 4 {
		comCards, _ = CheckHuUtilCardsToHandCards([]Card{23, 24, 24})
	}
	assert.Equal(t, hszCards, comCards)
	fmt.Println(hszCards)
	//三对vs两顺+单 224455vs2345678
	hanCard2 := []Card{12, 12, 14, 14, 15, 15, 22, 23, 24, 25, 26, 27, 28}
	cards2, _ := CheckHuUtilCardsToHandCards(hanCard2)
	hszCards2 := GetRecommedHuanSanZhang(cards2)
	assert.Equal(t, 3, len(hszCards2))
	comCards2, _ := CheckHuUtilCardsToHandCards([]Card{22, 28, 25})
	assert.Equal(t, hszCards2, comCards2)
	fmt.Println(hszCards2)
	// 杠+对+单vs刻+对+单 1122223vs333445
	hanCard3 := []Card{11, 11, 12, 12, 12, 12, 13, 23, 23, 23, 24, 24, 25}
	cards3, _ := CheckHuUtilCardsToHandCards(hanCard3)
	hszCards3 := GetRecommedHuanSanZhang(cards3)
	assert.Equal(t, 3, len(hszCards3))
	comCards3, _ := CheckHuUtilCardsToHandCards([]Card{23, 25, 24})
	assert.Equal(t, hszCards3, comCards3)
	fmt.Println(hszCards3)
	//顺+两对vs两顺 7893344vs123456
	hanCard5 := []Card{17, 18, 19, 13, 13, 14, 14, 21, 22, 23, 24, 25, 26}
	cards5, _ := CheckHuUtilCardsToHandCards(hanCard5)
	hszCards5 := GetRecommedHuanSanZhang(cards5)
	assert.Equal(t, 3, len(hszCards5))
	comCards5, _ := CheckHuUtilCardsToHandCards([]Card{21, 26, 23})
	if hszCards5[2].GetPoint() == 4 {
		comCards5, _ = CheckHuUtilCardsToHandCards([]Card{21, 26, 24})
	}
	assert.Equal(t, hszCards5, comCards5)
	fmt.Println(hszCards5)
}

// Test_PaiXing_PaiXingBuYiYang_Xiao2_7Vs7 牌型比较 牌型不一样 有两门牌数大于三张且相差<2张 7张vs7张
func Test_PaiXing_PaiXingBuYiYang_Xiao2_7Vs7(t *testing.T) {
	//杠+刻vs杠+顺 1113333vs1111234
	hanCard := []Card{11, 11, 11, 13, 13, 13, 13, 21, 21, 21, 21, 23, 24, 22}
	cards, _ := CheckHuUtilCardsToHandCards(hanCard)
	hszCards := GetRecommedHuanSanZhang(cards)
	assert.Equal(t, 3, len(hszCards))
	comCards, _ := CheckHuUtilCardsToHandCards([]Card{21, 24, 23})
	if hszCards[2].GetPoint() == 2 {
		comCards, _ = CheckHuUtilCardsToHandCards([]Card{21, 24, 22})
	}
	assert.Equal(t, hszCards, comCards)
	fmt.Println(hszCards)
	//两顺+单vs对子+顺子+单牌 234567vs114569
	hanCard2 := []Card{12, 13, 14, 15, 16, 17, 21, 21, 24, 25, 26, 29}
	cards2, _ := CheckHuUtilCardsToHandCards(hanCard2)
	hszCards2 := GetRecommedHuanSanZhang(cards2)
	assert.Equal(t, 3, len(hszCards2))
	comCards2, _ := CheckHuUtilCardsToHandCards([]Card{21, 29, 25})
	assert.Equal(t, hszCards2, comCards2)
	fmt.Println(hszCards2)
	// 三对+单vs顺+两对 2244667vs5566789
	hanCard3 := []Card{12, 12, 14, 14, 16, 16, 17, 25, 25, 26, 26, 27, 28, 29}
	cards3, _ := CheckHuUtilCardsToHandCards(hanCard3)
	hszCards3 := GetRecommedHuanSanZhang(cards3)
	assert.Equal(t, 3, len(hszCards3))
	comCards3, _ := CheckHuUtilCardsToHandCards([]Card{25, 29, 27})
	assert.Equal(t, hszCards3, comCards3)
	fmt.Println(hszCards3)
}

// Test_DingCard_You1And9
func Test_DingCard_You1And9(t *testing.T) {
	// 含1,9 提起1,9
	hanCard3 := []Card{11, 19, 12, 13, 14, 15}
	cards3, _ := CheckHuUtilCardsToHandCards(hanCard3)
	hszCards3 := GetRecommedHuanSanZhang(cards3)
	assert.Equal(t, 3, len(hszCards3))
	comCards3, _ := CheckHuUtilCardsToHandCards([]Card{11, 19, 15})
	assert.Equal(t, hszCards3, comCards3)
	fmt.Println(hszCards3)
	// 含1，不含9，提起1+最靠近9的一张+
	hanCard2 := []Card{11, 12, 13, 14, 15}
	cards2, _ := CheckHuUtilCardsToHandCards(hanCard2)
	hszCards2 := GetRecommedHuanSanZhang(cards2)
	assert.Equal(t, 3, len(hszCards2))
	comCards2, _ := CheckHuUtilCardsToHandCards([]Card{11, 15, 13})
	assert.Equal(t, hszCards2, comCards2)
	fmt.Println(hszCards2)
	//含9不含1，提起9+最靠近1的一张+
	hanCard1 := []Card{12, 15, 13, 14, 19}
	cards1, _ := CheckHuUtilCardsToHandCards(hanCard1)
	hszCards1 := GetRecommedHuanSanZhang(cards1)
	assert.Equal(t, 3, len(hszCards1))
	comCards1, _ := CheckHuUtilCardsToHandCards([]Card{12, 19, 15})
	assert.Equal(t, hszCards1, comCards1)
	fmt.Println(hszCards1)
	// 不含1,9最靠近1的一张+最靠近9的一张+
	hanCard := []Card{13, 18, 14, 16, 17}
	cards, _ := CheckHuUtilCardsToHandCards(hanCard)
	hszCards := GetRecommedHuanSanZhang(cards)
	assert.Equal(t, 3, len(hszCards))
	comCards, _ := CheckHuUtilCardsToHandCards([]Card{13, 18, 16})
	assert.Equal(t, hszCards, comCards)
	fmt.Println(hszCards)
}

// Test_TeShu 特殊三色三杠
func Test_TeShu(t *testing.T) {
	hanCard := []Card{11, 11, 11, 11, 22, 22, 22, 22, 33, 33, 33, 33}
	cards, _ := CheckHuUtilCardsToHandCards(hanCard)
	hszCards := GetRecommedHuanSanZhang(cards)
	assert.Equal(t, 3, len(hszCards))
	comCards, _ := CheckHuUtilCardsToHandCards([]Card{11, 11, 11})
	if hszCards[0].GetColor() == majongpb.CardColor_ColorTiao {
		comCards, _ = CheckHuUtilCardsToHandCards([]Card{22, 22, 22})
	} else if hszCards[0].GetColor() == majongpb.CardColor_ColorTong {
		comCards, _ = CheckHuUtilCardsToHandCards([]Card{33, 33, 33})
	}
	assert.Equal(t, hszCards, comCards)
	fmt.Println(hszCards)
}

func Test_A(t *testing.T) {
	fmt.Println(19 % 2)
}
