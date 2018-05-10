package states

import (
	"fmt"
	"steve/majong/interfaces"
	"testing"

	majongpb "steve/server_pb/majong"

	"github.com/Sirupsen/logrus"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

// Test_HuState_hu 胡状态-->摸牌状态
func Test_HuState_hu(t *testing.T) {
	mc := gomock.NewController(t)
	flow := interfaces.NewMockMajongFlow(mc)
	flow.EXPECT().GetMajongContext().Return(
		&majongpb.MajongContext{
			GameId:    1,
			WallCards: []*majongpb.Card{&Card1W},
			Players: []*majongpb.Player{
				&majongpb.Player{
					PalyerId:        1,
					HandCards:       []*majongpb.Card{&Card1W, &Card1W, &Card1W, &Card2W, &Card2W, &Card2W, &Card2W, &Card3W, &Card3W, &Card3W, &Card3W, &Card4W, &Card4W},
					DingqueColor:    majongpb.CardColor_ColorTiao,
					OutCards:        []*majongpb.Card{&Card1T},
					PossibleActions: []majongpb.Action{},
				},
				&majongpb.Player{
					PalyerId:        2,
					HandCards:       []*majongpb.Card{&Card1T, &Card1T, &Card1B, &Card1B, &Card2B, &Card2B, &Card2B, &Card2B, &Card3B, &Card3B, &Card3B, &Card3B, &Card4B},
					DingqueColor:    majongpb.CardColor_ColorTong,
					PossibleActions: []majongpb.Action{},
				},
				&majongpb.Player{
					PalyerId:        3,
					HandCards:       []*majongpb.Card{&Card2T, &Card3T, &Card5T, &Card5T, &Card5T, &Card6T, &Card6T, &Card6T, &Card7T, &Card7T, &Card7T, &Card8T, &Card8T},
					DingqueColor:    majongpb.CardColor_ColorWan,
					PossibleActions: []majongpb.Action{},
				},
			},
			LastChupaiPlayer: 1,
			LastOutCard:      &Card1T,
			LastHuPlayers:    []uint64{2},
		},
	).AnyTimes()

	s := HuState{}
	context := flow.GetMajongContext()
	beforeResults := ""
	beforeResults += fmt.Sprintln("点炮胡状态：")
	for _, player := range context.Players {
		beforeResults += FmtPlayerInfo(player)
	}
	logrus.Info(beforeResults)
	states, err := s.ProcessEvent(majongpb.EventID_event_hu_finish, nil, flow)
	assert.Nil(t, err)
	assert.Equal(t, majongpb.StateID_state_mopai, states, "点炮胡状态自动跳转到摸牌状态")
	afterResults := ""
	afterResults += fmt.Sprintln("摸牌状态：")
	for _, player := range context.Players {
		afterResults += FmtPlayerInfo(player)
	}
	logrus.Info(afterResults)
}
