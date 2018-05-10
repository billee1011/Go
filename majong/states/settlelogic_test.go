package states

import (
	"fmt"
	"steve/majong/utils"
	majongpb "steve/server_pb/majong"
	"testing"

	"github.com/stretchr/testify/assert"
)

var players = make([]*majongpb.Player, 0)

var mingGang1w = &majongpb.GangCard{
	Card:      &Card1W,
	Type:      majongpb.GangType_gang_minggang,
	SrcPlayer: uint64(2),
}

var anGang1t = &majongpb.GangCard{
	Card: &Card1T,
	Type: majongpb.GangType_gang_angang,
}

var buGang1b = &majongpb.GangCard{
	Card: &Card1B,
	Type: majongpb.GangType_gang_bugang,
}

var player0 = &majongpb.Player{PalyerId: uint64(0), HandCards: []*majongpb.Card{}}

var player1 = &majongpb.Player{PalyerId: uint64(1), HandCards: []*majongpb.Card{}}

var player2 = &majongpb.Player{PalyerId: uint64(2), HandCards: []*majongpb.Card{}}

var player3 = &majongpb.Player{PalyerId: uint64(3), HandCards: []*majongpb.Card{}}

var card0util = []utils.Card{11, 11, 12, 12, 13, 13, 14, 14, 15, 15, 16, 17, 18}

var card1util = []utils.Card{21, 21, 22, 22, 23, 23, 24, 24, 25, 25, 19, 26, 29}

var card2util = []utils.Card{11, 11, 12, 12, 13, 13, 14, 14, 15, 15, 16, 16, 17}

var card3util = []utils.Card{31, 31, 32, 32, 33, 33, 34, 34, 35, 35, 36, 37, 38}

// 初始化玩家,默认玩家0是花猪，玩家1是未听玩家，玩家2是听玩家，玩家3是胡玩家
func init() {
	prop := map[string][]byte{utils.IsOutNoDingQueColorCard: []byte{1}}

	player0.Properties = prop
	player0.HandCards, _ = utils.CheckHuUtilCardsToHandCards(card0util)
	player0.DingqueColor = majongpb.CardColor_ColorWan

	player1.Properties = prop
	player1.HandCards, _ = utils.CheckHuUtilCardsToHandCards(card1util)
	player1.DingqueColor = majongpb.CardColor_ColorTong

	player2.Properties = prop
	player2.HandCards, _ = utils.CheckHuUtilCardsToHandCards(card2util)
	player2.DingqueColor = majongpb.CardColor_ColorTiao

	player3.Properties = prop
	player3.HandCards, _ = utils.CheckHuUtilCardsToHandCards(card3util)
	player3.DingqueColor = majongpb.CardColor_ColorWan

	players = append(players, player0, player1, player2, player3)
}

//  0个花猪的情况,不进行花猪
func TestCheckFlowerPigSettle(t *testing.T) {
	for i := 0; i < len(players); i++ {
		player := players[i]
		player.HandCards, _ = utils.CheckHuUtilCardsToHandCards(card1util)
		player.HuCards = []*majongpb.HuCard{}
		player.DingqueColor = majongpb.CardColor_ColorTong
	}
	settleInfos, err := CheckFlowerPigSettle(players)
	assert.Nil(t, err)
	assert.Equal(t, len(settleInfos), 0)
	fmt.Println(settleInfos)
}

//  1个花猪的情况,1个听，1个未听，1个胡
func TestCheckFlowerPigSettleB(t *testing.T) {
	player3.HuCards = append(player3.HuCards, &majongpb.HuCard{Card: &Card1B, SrcPlayer: uint64(0), Type: majongpb.HuType_hu_dianpao})
	settleInfos, err := CheckFlowerPigSettle(players)
	assert.Nil(t, err)
	assert.Equal(t, len(settleInfos), 1)
	assert.Equal(t, settleInfos[0].Id, uint64(0))
	fmt.Println(settleInfos)
}

// 1个未听玩家,2个听玩家
func TestCheckYellSettleA(t *testing.T) {
	settleInfos, err := CheckYellSettle(players)
	assert.Nil(t, err)
	assert.Equal(t, len(settleInfos), 1)
	for k := range settleInfos {
		assert.Equal(t, settleInfos[k].Id, uint64(1))
		for ID := range settleInfos[k].Scores {
			assert.True(t, ID == uint64(2) || ID == uint64(3) || ID == uint64(1))
		}
	}
	fmt.Println(settleInfos)
}

func TestCallDivertSettle(t *testing.T) {
	// CallDivertSettle
}
