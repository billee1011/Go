package common

import (
	"steve/majong/interfaces"
	majongpb "steve/server_pb/majong"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func initPlayers(mjContext *majongpb.MajongContext) {
	mjContext.Players = mjContext.Players[0:0]
	for i := 0; i < 4; i++ {
		mjContext.Players = append(mjContext.Players, &majongpb.Player{
			PalyerId:  uint64(i),
			HandCards: []*majongpb.Card{},
		})
	}
}

func TestFapaiState_fapai(t *testing.T) {
	mc := gomock.NewController(t)
	defer mc.Finish()

	wallCards := getOriginCards(0)
	originWallCardsCount := len(wallCards)

	flow := interfaces.NewMockMajongFlow(mc)
	mjContext := majongpb.MajongContext{
		Players:        []*majongpb.Player{},
		WallCards:      wallCards,
		ZhuangjiaIndex: 0,
	}
	initPlayers(&mjContext)

	flow.EXPECT().GetMajongContext().Return(&mjContext).AnyTimes()

	f := new(FapaiState)
	f.fapai(flow)

	player0CardsCount := len(mjContext.GetPlayers()[0].HandCards)
	player1CardsCount := len(mjContext.GetPlayers()[1].HandCards)
	player2CardsCount := len(mjContext.GetPlayers()[2].HandCards)
	player3CardsCount := len(mjContext.GetPlayers()[3].HandCards)

	assert.Equal(t, player0CardsCount-1, player1CardsCount)
	assert.Equal(t, player2CardsCount, player1CardsCount)
	assert.Equal(t, player3CardsCount, player1CardsCount)

	newWallCardsCount := len(mjContext.WallCards)
	assert.Equal(t, originWallCardsCount-player0CardsCount-player1CardsCount-player2CardsCount-player3CardsCount, newWallCardsCount)
}
