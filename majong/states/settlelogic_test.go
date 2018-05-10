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

	player0.HandCards, _ = utils.CheckHuUtilCardsToHandCards(card0util)
	player0.OutCards = []*majongpb.Card{&Card1T}
	player0.DingqueColor = majongpb.CardColor_ColorWan

	player1.HandCards, _ = utils.CheckHuUtilCardsToHandCards(card1util)
	player1.OutCards = []*majongpb.Card{&Card2T}
	player1.DingqueColor = majongpb.CardColor_ColorTong

	player2.HandCards, _ = utils.CheckHuUtilCardsToHandCards(card2util)
	player2.OutCards = []*majongpb.Card{&Card1W}
	player2.DingqueColor = majongpb.CardColor_ColorTiao

	player3.HandCards, _ = utils.CheckHuUtilCardsToHandCards(card3util)
	player3.OutCards = []*majongpb.Card{&Card1B}
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

// 明杠后炮呼叫转移，1人胡
func TestCallDivertSettle(t *testing.T) {
	player0.GangCards = append(player0.GangCards, mingGang1w)
	// player3为胡家，player0输家-明杠
	settleInfos, err := CallDivertSettle(majongpb.HuType_hu_ganghoupao, players, []*majongpb.Player{player3}, player0)
	assert.Nil(t, err)
	assert.Equal(t, len(settleInfos), 1)
	// 明杠转移2分
	scoreMap := map[uint64]int64{
		player0.PalyerId: -2,
		player3.PalyerId: 2,
	}
	for k := range settleInfos {
		assert.Equal(t, settleInfos[k].Id, uint64(0))
		for ID, score := range settleInfos[k].Scores {
			assert.Equal(t, score, scoreMap[ID])
		}
	}
	fmt.Println(settleInfos)
}

// 明杠后炮呼叫转移，2人胡,胡家不包含明杠
func TestCallDivertSettleB(t *testing.T) {
	mingGang1w.SrcPlayer = uint64(1)
	player0.GangCards = append(player0.GangCards, mingGang1w)
	// player3为胡家，player0输家-明杠
	settleInfos, err := CallDivertSettle(majongpb.HuType_hu_ganghoupao, players, []*majongpb.Player{player3, player2}, player0)
	assert.Nil(t, err)
	assert.Equal(t, len(settleInfos), 1)
	// 明杠转移2分
	scoreMap := map[uint64]int64{
		player0.PalyerId: -2,
		player3.PalyerId: 1,
		player2.PalyerId: 1,
	}
	for k := range settleInfos {
		assert.Equal(t, settleInfos[k].Id, uint64(0))
		for ID, score := range settleInfos[k].Scores {
			assert.Equal(t, score, scoreMap[ID])
		}
	}
	fmt.Println(settleInfos)
}

// 明杠后炮呼叫转移，2人胡,胡家包含明杠
func TestCallDivertSettleC(t *testing.T) {
	mingGang1w.SrcPlayer = uint64(2)
	player0.GangCards = append(player0.GangCards, mingGang1w)
	// player3为胡家，player0输家-明杠
	settleInfos, err := CallDivertSettle(majongpb.HuType_hu_ganghoupao, players, []*majongpb.Player{player3, player2}, player0)
	assert.Nil(t, err)
	assert.Equal(t, len(settleInfos), 1)
	// 明杠转移2分
	scoreMap := map[uint64]int64{
		player0.PalyerId: -2,
		player2.PalyerId: 2,
	}
	for k := range settleInfos {
		assert.Equal(t, settleInfos[k].Id, uint64(0))
		for ID, score := range settleInfos[k].Scores {
			assert.Equal(t, score, scoreMap[ID])
		}
	}
	fmt.Println(settleInfos)
}

// 明杠后炮呼叫转移，3人胡,胡家包含明杠
func TestCallDivertSettleD(t *testing.T) {
	mingGang1w.SrcPlayer = uint64(2)
	player0.GangCards = append(player0.GangCards, mingGang1w)
	// player3为胡家，player0输家-明杠
	settleInfos, err := CallDivertSettle(majongpb.HuType_hu_ganghoupao, players, []*majongpb.Player{player3, player2, player1}, player0)
	assert.Nil(t, err)
	assert.Equal(t, len(settleInfos), 1)
	// 明杠转移2分
	scoreMap := map[uint64]int64{
		player0.PalyerId: -2,
		player2.PalyerId: 2,
	}
	for k := range settleInfos {
		assert.Equal(t, settleInfos[k].Id, uint64(0))
		for ID, score := range settleInfos[k].Scores {
			assert.Equal(t, score, scoreMap[ID])
		}
	}
	fmt.Println(settleInfos)
}

// 补杠后炮呼叫转移，1人胡
func TestCallDivertSettleE(t *testing.T) {
	player0.GangCards = append(player0.GangCards, buGang1b)
	// player3为胡家，player0输家-明杠
	settleInfos, err := CallDivertSettle(majongpb.HuType_hu_ganghoupao, players, []*majongpb.Player{player3}, player0)
	assert.Nil(t, err)
	assert.Equal(t, len(settleInfos), 1)
	// 补杠转移3分
	scoreMap := map[uint64]int64{
		player0.PalyerId: -3,
		player3.PalyerId: 3,
	}
	for k := range settleInfos {
		assert.Equal(t, settleInfos[k].Id, uint64(0))
		for ID, score := range settleInfos[k].Scores {
			assert.Equal(t, score, scoreMap[ID])
		}
	}
	fmt.Println(settleInfos)
}

// 补杠后炮呼叫转移，2人胡,不够分，第一个胡家，多赢一分
func TestCallDivertSettleF(t *testing.T) {
	player0.GangCards = append(player0.GangCards, buGang1b)
	// player3为胡家，player0输家-明杠
	settleInfos, err := CallDivertSettle(majongpb.HuType_hu_ganghoupao, players, []*majongpb.Player{player3, player2}, player0)
	assert.Nil(t, err)
	assert.Equal(t, len(settleInfos), 1)
	// 补杠转移3分
	scoreMap := map[uint64]int64{
		player0.PalyerId: -3,
		player2.PalyerId: 2,
		player3.PalyerId: 1,
	}
	for k := range settleInfos {
		assert.Equal(t, settleInfos[k].Id, uint64(0))
		for ID, score := range settleInfos[k].Scores {
			assert.Equal(t, score, scoreMap[ID])
		}
	}
	fmt.Println(settleInfos)
}

// 补杠后炮呼叫转移，2人胡,不够分，第一个胡家，多赢一分
func TestCallDivertSettleK(t *testing.T) {
	player3.GangCards = append(player3.GangCards, buGang1b)
	// player3为胡家，player0输家-明杠
	settleInfos, err := CallDivertSettle(majongpb.HuType_hu_ganghoupao, players, []*majongpb.Player{player1, player2}, player3)
	assert.Nil(t, err)
	assert.Equal(t, len(settleInfos), 1)
	// 补杠转移3分
	scoreMap := map[uint64]int64{
		player3.PalyerId: -3,
		player1.PalyerId: 2,
		player2.PalyerId: 1,
	}
	for k := range settleInfos {
		assert.Equal(t, settleInfos[k].Id, uint64(3))
		for ID, score := range settleInfos[k].Scores {
			assert.Equal(t, score, scoreMap[ID])
		}
	}
	fmt.Println(settleInfos)
}

// 补杠后炮呼叫转移，3人胡
func TestCallDivertSettleG(t *testing.T) {
	player0.GangCards = append(player0.GangCards, buGang1b)
	// player3为胡家，player0输家-明杠
	settleInfos, err := CallDivertSettle(majongpb.HuType_hu_ganghoupao, players, []*majongpb.Player{player3, player2, player1}, player0)
	assert.Nil(t, err)
	assert.Equal(t, len(settleInfos), 1)
	// 补杠转移3分
	scoreMap := map[uint64]int64{
		player0.PalyerId: -3,
		player2.PalyerId: 1,
		player3.PalyerId: 1,
		player1.PalyerId: 1,
	}
	for k := range settleInfos {
		assert.Equal(t, settleInfos[k].Id, uint64(0))
		for ID, score := range settleInfos[k].Scores {
			assert.Equal(t, score, scoreMap[ID])
		}
	}
	fmt.Println(settleInfos)
}

// 暗杠后炮呼叫转移，1人胡
func TestCallDivertSettleH(t *testing.T) {
	player0.GangCards = append(player0.GangCards, anGang1t)
	// player3为胡家，player0输家-明杠
	settleInfos, err := CallDivertSettle(majongpb.HuType_hu_ganghoupao, players, []*majongpb.Player{player3}, player0)
	assert.Nil(t, err)
	assert.Equal(t, len(settleInfos), 1)
	// 补杠转移3分
	scoreMap := map[uint64]int64{
		player0.PalyerId: -6,
		player3.PalyerId: 6,
	}
	for k := range settleInfos {
		assert.Equal(t, settleInfos[k].Id, uint64(0))
		for ID, score := range settleInfos[k].Scores {
			assert.Equal(t, score, scoreMap[ID])
		}
	}
	fmt.Println(settleInfos)
}

// 暗杠后炮呼叫转移，2人胡
func TestCallDivertSettleI(t *testing.T) {
	player0.GangCards = append(player0.GangCards, anGang1t)
	// player3为胡家，player0输家-明杠
	settleInfos, err := CallDivertSettle(majongpb.HuType_hu_ganghoupao, players, []*majongpb.Player{player3, player2}, player0)
	assert.Nil(t, err)
	assert.Equal(t, len(settleInfos), 1)
	// 补杠转移3分
	scoreMap := map[uint64]int64{
		player0.PalyerId: -6,
		player3.PalyerId: 3,
		player2.PalyerId: 3,
	}
	for k := range settleInfos {
		assert.Equal(t, settleInfos[k].Id, uint64(0))
		for ID, score := range settleInfos[k].Scores {
			assert.Equal(t, score, scoreMap[ID])
		}
	}
	fmt.Println(settleInfos)
}

// 暗杠后炮呼叫转移，3人胡
func TestCallDivertSettleJ(t *testing.T) {
	player0.GangCards = append(player0.GangCards, anGang1t)
	// player3为胡家，player0输家-明杠
	settleInfos, err := CallDivertSettle(majongpb.HuType_hu_ganghoupao, players, []*majongpb.Player{player3, player2, player1}, player0)
	assert.Nil(t, err)
	assert.Equal(t, len(settleInfos), 1)
	// 补杠转移3分
	scoreMap := map[uint64]int64{
		player0.PalyerId: -6,
		player3.PalyerId: 2,
		player2.PalyerId: 2,
		player1.PalyerId: 2,
	}
	for k := range settleInfos {
		assert.Equal(t, settleInfos[k].Id, uint64(0))
		for ID, score := range settleInfos[k].Scores {
			assert.Equal(t, score, scoreMap[ID])
		}
	}
	fmt.Println(settleInfos)
}
