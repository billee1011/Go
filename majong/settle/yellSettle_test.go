package settle

import (
	"fmt"
	"steve/majong/states"
	"steve/majong/utils"
	majongpb "steve/server_pb/majong"
	"testing"

	"github.com/stretchr/testify/assert"
)

var players = make([]*majongpb.Player, 0)

var mingGang1w = &majongpb.GangCard{
	Card:      &states.Card1W,
	Type:      majongpb.GangType_gang_minggang,
	SrcPlayer: uint64(2),
}

var anGang1t = &majongpb.GangCard{
	Card: &states.Card1T,
	Type: majongpb.GangType_gang_angang,
}

var buGang1b = &majongpb.GangCard{
	Card: &states.Card1B,
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

// 初始化玩家,默认玩家0是花猪，玩家1是未听玩家，玩家2，3是听玩家
func init() {

	player0.HandCards, _ = utils.CheckHuUtilCardsToHandCards(card0util)
	player0.OutCards = []*majongpb.Card{
		&majongpb.Card{
			Color: majongpb.CardColor_ColorTong,
			Point: 1,
		},
	}
	player0.DingqueColor = majongpb.CardColor_ColorWan

	player1.HandCards, _ = utils.CheckHuUtilCardsToHandCards(card1util)
	player1.DingqueColor = majongpb.CardColor_ColorTong

	player2.HandCards, _ = utils.CheckHuUtilCardsToHandCards(card2util)
	player2.PengCards = []*majongpb.PengCard{
		&majongpb.PengCard{
			Card:      &states.Card1W,
			SrcPlayer: 0,
		},
	}
	player2.DingqueColor = majongpb.CardColor_ColorTiao

	player3.HandCards, _ = utils.CheckHuUtilCardsToHandCards(card3util)
	player3.DingqueColor = majongpb.CardColor_ColorWan

	players = append(players, player0, player1, player2, player3)
}

// 1个未听玩家,2个听玩家
func TestCheckYellSettleA(t *testing.T) {
	yellSettle := new(YellSettle)
	context := &majongpb.MajongContext{
		GameId:  1,
		Players: players,
	}
	settleInfos := yellSettle.SettleYell(context)
	assert.Equal(t, len(settleInfos), 1)
	for k := range settleInfos {
		assert.Equal(t, settleInfos[k].Id, uint64(1))
		total := int64(0)
		for pid, score := range settleInfos[k].Scores {
			assert.True(t, pid == uint64(2) || pid == uint64(3) || pid == uint64(1) || pid == uint64(0))
			total = total + score
		}
		assert.Equal(t, total, int64(0))
	}
	fmt.Println(settleInfos)
}
