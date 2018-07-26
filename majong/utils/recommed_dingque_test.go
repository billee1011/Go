package utils

import (
	"fmt"
	"steve/gutils"
	majongpb "steve/server_pb/majong"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test_LiangMengWeiLing 俩门花色为0的情况下，随机推荐一门
func Test_LiangMengWeiLing(t *testing.T) {
	hanCard := []Card{21, 21, 21, 22, 22, 22, 23, 23, 23, 24, 24, 24, 25, 25}
	cards, _ := CheckHuUtilCardsToHandCards(hanCard)
	color := gutils.GetRecommedDingQueColor(cards)
	assert.Contains(t, []majongpb.CardColor{majongpb.CardColor_ColorWan, majongpb.CardColor_ColorTong}, color)
	fmt.Println(color)
}

// Test_yiMengWeiLing 1门花色为0的情况下，推荐为0的那张
func Test_yiMengWeiLing(t *testing.T) {
	hanCard := []Card{21, 21, 21, 22, 22, 22, 23, 23, 23, 34, 34, 34, 35, 35}
	cards, _ := CheckHuUtilCardsToHandCards(hanCard)
	color := gutils.GetRecommedDingQueColor(cards)
	assert.Equal(t, majongpb.CardColor_ColorWan, color)
	fmt.Println(color)
}

//Test_LiangMengPaiShuChaDaYuDeng2  两门牌数相差>=2张
func Test_LiangMengPaiShuChaDaYuDeng2(t *testing.T) {
	hanCard := []Card{21, 21, 21, 22, 22, 22, 11, 11, 12, 34, 34, 34, 35, 35}
	cards, _ := CheckHuUtilCardsToHandCards(hanCard)
	color := gutils.GetRecommedDingQueColor(cards)
	assert.Equal(t, majongpb.CardColor_ColorWan, color)
	fmt.Println(color)
}

//Test_paiXing_3vs4_gangVsKe  3张vs4张,gang vs ke
func Test_paiXing_3vs4_gangVsKe(t *testing.T) {
	hanCard := []Card{12, 12, 12, 22, 22, 22, 22, 31, 31, 31, 32, 32, 32, 33, 33}
	cards, _ := CheckHuUtilCardsToHandCards(hanCard)
	color := gutils.GetRecommedDingQueColor(cards)
	assert.Equal(t, majongpb.CardColor_ColorWan, color)
	fmt.Println(color)
}

//Test_paiXing_3vs4_keVsShunJiaDan 刻vs顺+单
func Test_paiXing_3vs4_keVsShunJiaDan(t *testing.T) {
	hanCard := []Card{11, 11, 11, 22, 23, 24, 28, 31, 32, 33, 34, 35, 36, 37}
	cards, _ := CheckHuUtilCardsToHandCards(hanCard)
	color := gutils.GetRecommedDingQueColor(cards)
	assert.Equal(t, majongpb.CardColor_ColorTiao, color)
	fmt.Println(color)
}

//Test_paiXing_3vs4_keVsShunJiaDan2 刻vs顺+单
func Test_paiXing_3vs4_keVsShunJiaDan2(t *testing.T) {
	hanCard := []Card{11, 11, 11, 22, 23, 24, 25, 31, 32, 33, 34, 35, 36, 37}
	cards, _ := CheckHuUtilCardsToHandCards(hanCard)
	color := gutils.GetRecommedDingQueColor(cards)
	assert.Equal(t, majongpb.CardColor_ColorTiao, color)
	fmt.Println(color)
}

//Test_paiXing_3vs4_ShunJiaDan2VsDuiJiaDan 顺+单vs对+单
func Test_paiXing_3vs4_ShunJiaDan2VsDuiJiaDan(t *testing.T) {
	hanCard := []Card{12, 13, 14, 18, 22, 22, 23, 31, 32, 33, 34, 35, 36, 37}
	cards, _ := CheckHuUtilCardsToHandCards(hanCard)
	color := gutils.GetRecommedDingQueColor(cards)
	assert.Equal(t, majongpb.CardColor_ColorTiao, color)
	fmt.Println(color)
}

//Test_paiXing_3vs4_DuiJiaDanVsDanJiaDan 对+单vs单
func Test_paiXing_3vs4_DuiJiaDanVsDanJiaDan(t *testing.T) {
	hanCard := []Card{11, 13, 15, 15, 22, 23, 25, 31, 32, 33, 34, 35, 36, 37}
	cards, _ := CheckHuUtilCardsToHandCards(hanCard)
	color := gutils.GetRecommedDingQueColor(cards)
	assert.Equal(t, majongpb.CardColor_ColorTiao, color)
	fmt.Println(color)
}

//Test_paiXing_3vs4_KeVsLiangDui 刻vs两对 222vs3344
func Test_paiXing_3vs4_KeVsLiangDui(t *testing.T) {
	hanCard := []Card{12, 12, 12, 23, 23, 24, 24, 31, 32, 33, 34, 35, 36, 37}
	cards, _ := CheckHuUtilCardsToHandCards(hanCard)
	color := gutils.GetRecommedDingQueColor(cards)
	assert.Equal(t, majongpb.CardColor_ColorTiao, color)
	fmt.Println(color)
}

//Test_paiXing_3vs4_LiangDuiVsShun 两对vs顺 5588vs123
func Test_paiXing_3vs4_LiangDuiVsShun(t *testing.T) {
	hanCard := []Card{15, 15, 18, 18, 21, 22, 23, 31, 32, 33, 34, 35, 36, 37}
	cards, _ := CheckHuUtilCardsToHandCards(hanCard)
	color := gutils.GetRecommedDingQueColor(cards)
	assert.Equal(t, majongpb.CardColor_ColorTiao, color)
	fmt.Println(color)
}

//Test_paiXing_3vs4_LiangDuiVsShun2 两对vs顺 4466vs123
func Test_paiXing_3vs4_LiangDuiVsShun2(t *testing.T) {
	hanCard := []Card{14, 14, 16, 16, 21, 22, 23, 31, 32, 33, 34, 35, 36, 37}
	cards, _ := CheckHuUtilCardsToHandCards(hanCard)
	color := gutils.GetRecommedDingQueColor(cards)
	assert.Equal(t, majongpb.CardColor_ColorTiao, color)
	fmt.Println(color)
}

//Test_paiXing_3vs4_DuiJiaDanVsShun 对+单vs顺 1124vs123
func Test_paiXing_3vs4_DuiJiaDanVsShun(t *testing.T) {
	hanCard := []Card{11, 11, 12, 14, 21, 22, 23, 31, 32, 33, 34, 35, 36, 37}
	cards, _ := CheckHuUtilCardsToHandCards(hanCard)
	color := gutils.GetRecommedDingQueColor(cards)
	assert.Equal(t, majongpb.CardColor_ColorWan, color)
	fmt.Println(color)
}

//Test_LiangMengXiangTong_3Vs3_keVsShun 3张vs3张 刻vs顺 111vs234
func Test_LiangMengXiangTong_3Vs3_keVsShun(t *testing.T) {
	hanCard := []Card{11, 11, 11, 22, 23, 24, 31, 32, 33, 34, 35, 36, 37, 38}
	cards, _ := CheckHuUtilCardsToHandCards(hanCard)
	color := gutils.GetRecommedDingQueColor(cards)
	assert.Equal(t, majongpb.CardColor_ColorTiao, color)
	fmt.Println(color)
}

//Test_LiangMengXiangTong_3Vs3_ShunVsDui 3张vs3张 顺vs对 123vs112
func Test_LiangMengXiangTong_3Vs3_ShunVsDui(t *testing.T) {
	hanCard := []Card{11, 12, 13, 21, 21, 22, 31, 32, 33, 34, 35, 36, 37, 38}
	cards, _ := CheckHuUtilCardsToHandCards(hanCard)
	color := gutils.GetRecommedDingQueColor(cards)
	assert.Equal(t, majongpb.CardColor_ColorTiao, color)
	fmt.Println(color)
}

//Test_LiangMengXiangTong_3Vs3_DuiJiaDanVsDan 3张vs3张 对+单vs单 223vs147
func Test_LiangMengXiangTong_3Vs3_DuiJiaDanVsDan(t *testing.T) {
	hanCard := []Card{12, 12, 13, 21, 14, 17, 31, 32, 33, 34, 35, 36, 37, 38}
	cards, _ := CheckHuUtilCardsToHandCards(hanCard)
	color := gutils.GetRecommedDingQueColor(cards)
	assert.Equal(t, majongpb.CardColor_ColorTiao, color)
	fmt.Println(color)
}

//Test_LiangMengXiangTong_3Vs3_KeVsDuiJiaDan 3张vs3张 刻vs对+单 333vs233
func Test_LiangMengXiangTong_3Vs3_KeVsDuiJiaDan(t *testing.T) {
	hanCard := []Card{13, 13, 13, 22, 23, 23, 31, 32, 33, 34, 35, 36, 37, 38}
	cards, _ := CheckHuUtilCardsToHandCards(hanCard)
	color := gutils.GetRecommedDingQueColor(cards)
	assert.Equal(t, majongpb.CardColor_ColorTiao, color)
	fmt.Println(color)
}

//Test_LiangMengXiangTong_3Vs3_KeVsDan 3张vs3张 刻vs单 333vs135
func Test_LiangMengXiangTong_3Vs3_KeVsDan(t *testing.T) {
	hanCard := []Card{13, 13, 13, 21, 23, 25, 31, 32, 33, 34, 35, 36, 37, 38}
	cards, _ := CheckHuUtilCardsToHandCards(hanCard)
	color := gutils.GetRecommedDingQueColor(cards)
	assert.Equal(t, majongpb.CardColor_ColorTiao, color)
	fmt.Println(color)
}

//Test_LiangMengXiangTong_3Vs3_ShunVsDan 3张vs3张 顺vs单 234v156
func Test_LiangMengXiangTong_3Vs3_ShunVsDan(t *testing.T) {
	hanCard := []Card{12, 13, 14, 21, 25, 26, 31, 32, 33, 34, 35, 36, 37, 38}
	cards, _ := CheckHuUtilCardsToHandCards(hanCard)
	color := gutils.GetRecommedDingQueColor(cards)
	assert.Equal(t, majongpb.CardColor_ColorTiao, color)
	fmt.Println(color)
}

//Test_LiangMengXiangTong_4Vs4_KeJiaDanVsLiangDui 4张vs4张 刻+单vs两对 1888vs1188
func Test_LiangMengXiangTong_4Vs4_KeJiaDanVsLiangDui(t *testing.T) {
	hanCard := []Card{11, 18, 18, 18, 21, 21, 28, 28, 31, 32, 33, 34, 35, 36}
	cards, _ := CheckHuUtilCardsToHandCards(hanCard)
	color := gutils.GetRecommedDingQueColor(cards)
	assert.Equal(t, majongpb.CardColor_ColorTiao, color)
	fmt.Println(color)
}

//Test_LiangMengXiangTong_4Vs4_KeJiaDanVsShunJiaDan 4张vs4张 刻+单vs顺+单 2224vs1456
func Test_LiangMengXiangTong_4Vs4_KeJiaDanVsShunJiaDan(t *testing.T) {
	hanCard := []Card{12, 12, 12, 14, 21, 24, 25, 26, 31, 32, 33, 34, 35, 36}
	cards, _ := CheckHuUtilCardsToHandCards(hanCard)
	color := gutils.GetRecommedDingQueColor(cards)
	assert.Equal(t, majongpb.CardColor_ColorTiao, color)
	fmt.Println(color)
}

//Test_LiangMengXiangTong_4Vs4_KeJiaDanVsShunJiaDan2 4张vs4张 刻+单vs顺+单 1222vs4568
func Test_LiangMengXiangTong_4Vs4_KeJiaDanVsShunJiaDan2(t *testing.T) {
	hanCard := []Card{12, 12, 12, 11, 24, 25, 26, 28, 31, 32, 33, 34, 35, 36}
	cards, _ := CheckHuUtilCardsToHandCards(hanCard)
	color := gutils.GetRecommedDingQueColor(cards)
	assert.Equal(t, majongpb.CardColor_ColorTiao, color)
	fmt.Println(color)
}

//Test_LiangMengXiangTong_4Vs4_GangVsKeJiaDan 4张vs4张 杠vs刻+单 2222vs7888
func Test_LiangMengXiangTong_4Vs4_GangVsKeJiaDan(t *testing.T) {
	hanCard := []Card{12, 12, 12, 12, 27, 28, 28, 28, 31, 32, 33, 34, 35, 36}
	cards, _ := CheckHuUtilCardsToHandCards(hanCard)
	color := gutils.GetRecommedDingQueColor(cards)
	assert.Equal(t, majongpb.CardColor_ColorTiao, color)
	fmt.Println(color)
}

//Test_LiangMengXiangTong_4Vs4_GangVsKeJiaDan2 4张vs4张 杠vs刻+单 2222vs7888
func Test_LiangMengXiangTong_4Vs4_GangVsKeJiaDan2(t *testing.T) {
	hanCard := []Card{17, 18, 18, 18, 22, 22, 22, 22, 31, 32, 33, 34, 35, 36}
	cards, _ := CheckHuUtilCardsToHandCards(hanCard)
	color := gutils.GetRecommedDingQueColor(cards)
	assert.Equal(t, majongpb.CardColor_ColorWan, color)
	fmt.Println(color)
}

//Test_LiangMengXiangTong_4Vs4_GangVsLiangDui 4张vs4张 杠vs两对 3333vs5577
func Test_LiangMengXiangTong_4Vs4_GangVsLiangDui(t *testing.T) {
	hanCard := []Card{25, 25, 27, 27, 23, 23, 23, 23, 31, 32, 33, 34, 35, 36}
	cards, _ := CheckHuUtilCardsToHandCards(hanCard)
	color := gutils.GetRecommedDingQueColor(cards)
	assert.Equal(t, majongpb.CardColor_ColorWan, color)
	fmt.Println(color)
}

//Test_LiangMengXiangTong_4Vs4_LiangDuiVsShunJiaDan 4张vs4张 两对vs顺+单 4455vs3456
func Test_LiangMengXiangTong_4Vs4_LiangDuiVsShunJiaDan(t *testing.T) {
	hanCard := []Card{14, 14, 15, 15, 23, 24, 25, 26, 31, 32, 33, 34, 35, 36}
	cards, _ := CheckHuUtilCardsToHandCards(hanCard)
	color := gutils.GetRecommedDingQueColor(cards)
	assert.Equal(t, majongpb.CardColor_ColorTiao, color)
	fmt.Println(color)
}

//Test_LiangMengXiangTong_4Vs4_ShunJiaDanVsDan 4张vs4张 顺+单vs单 1234vs1458
func Test_LiangMengXiangTong_4Vs4_ShunJiaDanVsDan(t *testing.T) {
	hanCard := []Card{11, 12, 13, 14, 21, 24, 25, 28, 31, 32, 33, 34, 35, 36}
	cards, _ := CheckHuUtilCardsToHandCards(hanCard)
	color := gutils.GetRecommedDingQueColor(cards)
	assert.Equal(t, majongpb.CardColor_ColorTiao, color)
	fmt.Println(color)
}

//Test_LiangMengXiangTong_4Vs4_GangVsDan 4张vs4张 杠vs单 2222vs5689
func Test_LiangMengXiangTong_4Vs4_GangVsDan(t *testing.T) {
	hanCard := []Card{12, 12, 12, 12, 25, 26, 28, 29, 31, 32, 33, 34, 35, 36}
	cards, _ := CheckHuUtilCardsToHandCards(hanCard)
	color := gutils.GetRecommedDingQueColor(cards)
	assert.Equal(t, majongpb.CardColor_ColorTiao, color)
	fmt.Println(color)
}

//Test_LiangMengXiangTong_4Vs4_shunJiaDanVsDuiJiaDan 4张vs4张 顺+单vs对+单 1123vs1134
func Test_LiangMengXiangTong_4Vs4_shunJiaDanVsDuiJiaDan(t *testing.T) {
	hanCard := []Card{11, 11, 12, 13, 21, 21, 23, 24, 31, 32, 33, 34, 35, 36}
	cards, _ := CheckHuUtilCardsToHandCards(hanCard)
	color := gutils.GetRecommedDingQueColor(cards)
	assert.Equal(t, majongpb.CardColor_ColorTiao, color)
	fmt.Println(color)
}

//Test_LiangMengXiangTong_4Vs4_shunJiaDanVsDuiJiaDan2 4张vs4张 顺+单vs对+单 3455vs6689
func Test_LiangMengXiangTong_4Vs4_shunJiaDanVsDuiJiaDan2(t *testing.T) {
	hanCard := []Card{13, 14, 15, 15, 26, 26, 28, 29, 31, 32, 33, 34, 35, 36}
	cards, _ := CheckHuUtilCardsToHandCards(hanCard)
	color := gutils.GetRecommedDingQueColor(cards)
	assert.Equal(t, majongpb.CardColor_ColorTiao, color)
	fmt.Println(color)
}

//Test_chaYiAndEqual_4Vs4Vs5_GangVsKeJiaDangVsShunJiaDui  同时存在相差一张和相同的 4张vs4张vs5张
//杠vs刻+单vs顺+对 2222vs4448vs12344
func Test_chaYiAndEqual_4Vs4Vs5_GangVsKeJiaDangVsShunJiaDui(t *testing.T) {
	hanCard := []Card{12, 12, 12, 12, 24, 24, 24, 28, 31, 32, 33, 34, 34}
	cards, _ := CheckHuUtilCardsToHandCards(hanCard)
	color := gutils.GetRecommedDingQueColor(cards)
	assert.Equal(t, majongpb.CardColor_ColorTiao, color)
	fmt.Println(color)
}

//Test_chaYiAndEqual_4Vs4Vs5_KeJiaDuiVsLiangDuiVsShunJiaDan  同时存在相差一张和相同的 4张vs4张vs5张
//刻+对vs两对vs顺+单 33355vs3344vs1233
func Test_chaYiAndEqual_4Vs4Vs5_KeJiaDuiVsLiangDuiVsShunJiaDan(t *testing.T) {
	hanCard := []Card{13, 13, 13, 15, 15, 23, 23, 24, 24, 31, 32, 33, 33}
	cards, _ := CheckHuUtilCardsToHandCards(hanCard)
	color := gutils.GetRecommedDingQueColor(cards)
	assert.Equal(t, majongpb.CardColor_ColorTong, color)
	fmt.Println(color)
}

//Test_chaYiAndEqual_4Vs4Vs5_GangVsGangJiaDangVsLiangDui  同时存在相差一张和相同的 4张vs4张vs5张
//杠vs杠+单vs两对 1111vs11112vs3355
func Test_chaYiAndEqual_4Vs4Vs5_GangVsGangJiaDangVsLiangDui(t *testing.T) {
	hanCard := []Card{11, 11, 11, 11, 21, 21, 21, 21, 22, 33, 33, 35, 35}
	cards, _ := CheckHuUtilCardsToHandCards(hanCard)
	color := gutils.GetRecommedDingQueColor(cards)
	assert.Equal(t, majongpb.CardColor_ColorTong, color)
	fmt.Println(color)
}

//Test_chaYiAndEqual_4Vs4Vs5_LiangDuiVsKeJiaDanVsShunJiaDan  同时存在相差一张和相同的 4张vs4张vs5张
//两对vs刻子+单牌vs顺子+单牌 2233vs1444vs3455
func Test_chaYiAndEqual_4Vs4Vs5_LiangDuiVsKeJiaDanVsShunJiaDan(t *testing.T) {
	hanCard := []Card{12, 12, 13, 13, 21, 24, 24, 24, 33, 34, 35, 35}
	cards, _ := CheckHuUtilCardsToHandCards(hanCard)
	color := gutils.GetRecommedDingQueColor(cards)
	assert.Equal(t, majongpb.CardColor_ColorTong, color)
	fmt.Println(color)
}

//Test_chaYiAndEqual_4Vs4Vs5_LiangDuiVsShunJiaDanVsDuiJiaDan  同时存在相差一张和相同的 4张vs4张vs5张
//两对vs顺子+单vs对子+单 6677vs1345vs2245
func Test_chaYiAndEqual_4Vs4Vs5_LiangDuiVsShunJiaDanVsDuiJiaDan(t *testing.T) {
	hanCard := []Card{16, 16, 17, 17, 21, 23, 24, 25, 32, 32, 34, 45}
	cards, _ := CheckHuUtilCardsToHandCards(hanCard)
	color := gutils.GetRecommedDingQueColor(cards)
	assert.Equal(t, majongpb.CardColor_ColorTong, color)
	fmt.Println(color)
}

//Test_chaYiAndEqual_4Vs4Vs5_ShunJiaDuiVsKeJiaDanVsGang  同时存在相差一张和相同的 4张vs4张vs5张
//顺+对vs刻+单vs杠 34555vs6668vs4444
func Test_chaYiAndEqual_4Vs4Vs5_ShunJiaDuiVsKeJiaDanVsGang(t *testing.T) {
	hanCard := []Card{13, 14, 15, 15, 15, 26, 26, 26, 28, 34, 34, 34, 34}
	cards, _ := CheckHuUtilCardsToHandCards(hanCard)
	color := gutils.GetRecommedDingQueColor(cards)
	assert.Equal(t, majongpb.CardColor_ColorTiao, color)
	fmt.Println(color)
}

//Test_chaYiAndEqual_4Vs4Vs5_LiangJiaDangVsLianDuiVsGang  同时存在相差一张和相同的 4张vs4张vs5张
//俩对+单vs俩对vs杠 11223vs4455vs4444
func Test_chaYiAndEqual_4Vs4Vs5_LiangJiaDangVsLianDuiVsGang(t *testing.T) {
	hanCard := []Card{11, 11, 12, 12, 13, 24, 24, 25, 25, 34, 34, 34, 34}
	cards, _ := CheckHuUtilCardsToHandCards(hanCard)
	color := gutils.GetRecommedDingQueColor(cards)
	assert.Contains(t, []majongpb.CardColor{majongpb.CardColor_ColorWan, majongpb.CardColor_ColorTiao}, color)
	fmt.Println(color)
}

//Test_chaYiAndEqual_4Vs4Vs5_GangVsShunJiaDuiVsKeJiaDan  同时存在相差一张和相同的 4张vs5张vs5张
//杠vs顺子+对子vs刻子+单牌 1111vs56788vs4445
func Test_chaYiAndEqual_4Vs4Vs5_GangVsShunJiaDuiVsKeJiaDan(t *testing.T) {
	hanCard := []Card{11, 11, 11, 11, 25, 26, 27, 28, 28, 34, 34, 34, 35}
	cards, _ := CheckHuUtilCardsToHandCards(hanCard)
	color := gutils.GetRecommedDingQueColor(cards)
	assert.Equal(t, majongpb.CardColor_ColorTong, color)
	fmt.Println(color)
}

//Test_chaYiAndEqual_4Vs4Vs5_LiangDuiVsKeJiaDuiVsShunJiaDui  同时存在相差一张和相同的 4张vs5张vs5张
//两对vs刻子+对子vs顺子+对子 3344vs22244vs55678
func Test_chaYiAndEqual_4Vs4Vs5_LiangDuiVsKeJiaDuiVsShunJiaDui(t *testing.T) {
	hanCard := []Card{13, 13, 14, 14, 22, 22, 22, 24, 24, 35, 35, 36, 37, 38}
	cards, _ := CheckHuUtilCardsToHandCards(hanCard)
	color := gutils.GetRecommedDingQueColor(cards)
	assert.Equal(t, majongpb.CardColor_ColorWan, color)
	fmt.Println(color)
}

//Test_chaYiAndEqual_4Vs4Vs5_GangVsGangJiaDangVsKeJiaDui  同时存在相差一张和相同的 4张vs5张vs5张
//杠vs杠+单牌vs刻+对 5555vs14444vs66677
func Test_chaYiAndEqual_4Vs4Vs5_GangVsGangJiaDangVsKeJiaDui(t *testing.T) {
	hanCard := []Card{15, 15, 15, 15, 21, 24, 24, 24, 24, 36, 36, 36, 37, 37}
	cards, _ := CheckHuUtilCardsToHandCards(hanCard)
	color := gutils.GetRecommedDingQueColor(cards)
	assert.Equal(t, majongpb.CardColor_ColorTong, color)
	fmt.Println(color)
}

//Test_chaYiAndEqual_4Vs4Vs5_ShunJiaDanVsShunJiaDuiVsKeJiaDui  同时存在相差一张和相同的 4张vs5张vs5张
//顺子+单牌vs顺子+对子vs刻子+对子 1567vs34555vs66677
func Test_chaYiAndEqual_4Vs4Vs5_ShunJiaDanVsShunJiaDuiVsKeJiaDui(t *testing.T) {
	hanCard := []Card{11, 15, 16, 17, 23, 24, 25, 25, 25, 36, 36, 36, 37, 37}
	cards, _ := CheckHuUtilCardsToHandCards(hanCard)
	color := gutils.GetRecommedDingQueColor(cards)
	assert.Equal(t, majongpb.CardColor_ColorWan, color)
	fmt.Println(color)
}

//Test_chaYiAndEqual_4Vs4Vs5_DuiJiaDangVsShunJiaDanVsLiangDuiJiaDan  同时存在相差一张和相同的 4张vs5张vs5张
//对子+单牌vs顺子+单牌vs两对+单牌 5578vs3567vs33445?/55889
func Test_chaYiAndEqual_4Vs4Vs5_DuiJiaDangVsShunJiaDanVsLiangDuiJiaDan(t *testing.T) {
	hanCard := []Card{15, 15, 17, 18, 23, 25, 26, 27, 33, 33, 34, 34, 35}
	cards, _ := CheckHuUtilCardsToHandCards(hanCard)
	color := gutils.GetRecommedDingQueColor(cards)
	assert.Equal(t, majongpb.CardColor_ColorWan, color)
	fmt.Println(color)
	hanCard2 := []Card{15, 15, 17, 18, 23, 25, 26, 27, 35, 35, 38, 38, 39}
	cards2, _ := CheckHuUtilCardsToHandCards(hanCard2)
	color2 := gutils.GetRecommedDingQueColor(cards2)
	assert.Equal(t, majongpb.CardColor_ColorWan, color2)
	fmt.Println(color)
}

//Test_chaYiAndEqual_4Vs4Vs5_DuiJiaDanVsShunJiaDanVsDan  同时存在相差一张和相同的 4张vs5张vs5张
//对子+单牌vs顺子+单牌vs单牌 3356vs45678vs14769
func Test_chaYiAndEqual_4Vs4Vs5_DuiJiaDanVsShunJiaDanVsDan(t *testing.T) {
	hanCard := []Card{13, 13, 15, 16, 24, 25, 26, 27, 28, 31, 34, 37, 36, 39}
	cards, _ := CheckHuUtilCardsToHandCards(hanCard)
	color := gutils.GetRecommedDingQueColor(cards)
	assert.Equal(t, majongpb.CardColor_ColorTong, color)
	fmt.Println(color)
}

//Test_chaYiAndEqual_4Vs4Vs5_ShunJiaDuiVsShunJiaDuiVsGang  同时存在相差一张和相同的 4张vs5张vs5张
//顺+对vs顺+对vs杠 11234vs22234vs4444
func Test_chaYiAndEqual_4Vs4Vs5_ShunJiaDuiVsShunJiaDuiVsGang(t *testing.T) {
	hanCard := []Card{11, 11, 12, 13, 14, 22, 22, 22, 23, 24, 34, 34, 34, 34}
	cards, _ := CheckHuUtilCardsToHandCards(hanCard)
	color := gutils.GetRecommedDingQueColor(cards)
	assert.Contains(t, []majongpb.CardColor{majongpb.CardColor_ColorWan, majongpb.CardColor_ColorTiao}, color)
	fmt.Println(color)
}

//Test_TeShu_siZhangWanVsJiuTiao  特殊：有一门花色只有四张不定该花色 有一门牌只有四张且构成杠的
//手牌4张万+9张条
func Test_TeShu_siZhangWanVsJiuTiao(t *testing.T) {
	hanCard := []Card{11, 11, 11, 11, 21, 22, 23, 24, 25, 26, 27, 28, 29}
	cards, _ := CheckHuUtilCardsToHandCards(hanCard)
	color := gutils.GetRecommedDingQueColor(cards)
	assert.Equal(t, color, majongpb.CardColor_ColorTong)
	fmt.Println(color)
}

//Test_TeShu_siTiaoVsWuWanVsWuTong  特殊：有一门花色只有四张不定该花色 有一门牌只有四张且构成杠的
//手牌中有4张条+5张万+5张筒
func Test_TeShu_siTiaoVsWuWanVsWuTong(t *testing.T) {
	hanCard := []Card{21, 21, 21, 21, 11, 11, 11, 11, 15, 31, 31, 31, 31, 35}
	cards, _ := CheckHuUtilCardsToHandCards(hanCard)
	color := gutils.GetRecommedDingQueColor(cards)
	assert.Contains(t, []majongpb.CardColor{majongpb.CardColor_ColorWan, majongpb.CardColor_ColorTong}, color)
	fmt.Println(color)
}

//Test_TeShu_siTiaoVsSiWanVsWuTong  特殊：有2门的牌只有四张且构成杠
//手牌中有4张条+4张万+5张筒
func Test_TeShu_siTiaoVsSiWanVsWuTong(t *testing.T) {
	hanCard := []Card{21, 21, 21, 21, 11, 11, 11, 11, 31, 31, 31, 31, 32}
	cards, _ := CheckHuUtilCardsToHandCards(hanCard)
	color := gutils.GetRecommedDingQueColor(cards)
	assert.Equal(t, majongpb.CardColor_ColorTong, color)
	fmt.Println(color)
}

//Test_TeShu_siTiaoVsSiWanVsLiuTong  特殊：有2门的牌只有四张且构成杠
//手牌中有4张条+4张万+6张筒
func Test_TeShu_siTiaoVsSiWanVsLiuTong(t *testing.T) {
	hanCard := []Card{21, 21, 21, 21, 11, 11, 11, 11, 31, 31, 31, 31, 32, 32}
	cards, _ := CheckHuUtilCardsToHandCards(hanCard)
	color := gutils.GetRecommedDingQueColor(cards)
	assert.Contains(t, []majongpb.CardColor{majongpb.CardColor_ColorWan, majongpb.CardColor_ColorTiao}, color)
	fmt.Println(color)
}

func Test_B(t *testing.T) {
	hanCard := []Card{11, 12, 13, 15, 21, 22, 23, 31, 32, 33, 34, 35, 36, 37, 38}
	cards, _ := CheckHuUtilCardsToHandCards(hanCard)
	color := gutils.GetRecommedDingQueColor(cards)
	assert.Equal(t, majongpb.CardColor_ColorTiao, color)
	fmt.Println(color)
}
