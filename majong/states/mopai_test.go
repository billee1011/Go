package states

import (
	"fmt"
	"steve/majong/interfaces"
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	majongpb "steve/server_pb/majong"

	"github.com/golang/mock/gomock"
)

// Test_MopaiState_mopai 摸牌状态
func Test_MopaiState_mopai(t *testing.T) {
	mc := gomock.NewController(t)
	defer mc.Finish()
	flow := interfaces.NewMockMajongFlow(mc)
	flow.EXPECT().PushMessages(gomock.Any(), gomock.Any()).AnyTimes()
	flow.EXPECT().GetMajongContext().Return(
		&majongpb.MajongContext{
			GameId: 1,
			Players: []*majongpb.Player{
				&majongpb.Player{
					PalyerId:     1,
					HandCards:    []*majongpb.Card{&Card1W, &Card1W, &Card1W, &Card1W, &Card2W, &Card2W, &Card2W, &Card2W, &Card3W, &Card3W, &Card3W, &Card3W, &Card4W, &Card4W},
					DingqueColor: majongpb.CardColor_ColorTiao,
				},
			},
			ActivePlayer: 1,
			WallCards:    []*majongpb.Card{&Card1T},
		},
	).AnyTimes()
	s := MoPaiState{}
	context := flow.GetMajongContext()
	players := context.GetPlayers()
	player := players[0]
	beforeResults := ""
	beforeResults += fmt.Sprintln("查询前:")
	beforeResults += FmtPlayerInfo(player)
	logrus.Info(beforeResults)
	state, err := s.ProcessEvent(majongpb.EventID_event_mopai_finish, nil, flow)
	afterResults := ""
	afterResults += fmt.Sprintln("查询后:")
	afterResults += FmtPlayerInfo(player)
	logrus.Info(afterResults)
	assert.Nil(t, err)
	assert.Equal(t, majongpb.StateID_state_zixun, state, "查询成功后有特殊操作，应该进入自询状态")
}

func Test_MopaiState_mopai1(t *testing.T) {
	mc := gomock.NewController(t)
	defer mc.Finish()
	flow := interfaces.NewMockMajongFlow(mc)
	flow.EXPECT().PushMessages(gomock.Any(), gomock.Any()).AnyTimes()
	flow.EXPECT().GetMajongContext().Return(
		&majongpb.MajongContext{
			GameId: 1,
			Players: []*majongpb.Player{
				&majongpb.Player{
					PalyerId:     1,
					HandCards:    []*majongpb.Card{&Card1W, &Card2W, &Card2W, &Card2W, &Card2W, &Card3W, &Card3W, &Card3W, &Card3W, &Card4W, &Card4W},
					DingqueColor: majongpb.CardColor_ColorTiao,
					PengCards: []*majongpb.PengCard{
						&majongpb.PengCard{
							Card:      &Card1W,
							SrcPlayer: 2,
						},
					},
				},
			},
			ActivePlayer: 1,
			WallCards:    []*majongpb.Card{&Card1T},
		},
	).AnyTimes()
	s := MoPaiState{}
	context := flow.GetMajongContext()
	players := context.GetPlayers()
	player := players[0]
	beforeResults := ""
	beforeResults += fmt.Sprintln("查询前:")
	beforeResults += FmtPlayerInfo(player)
	logrus.Info(beforeResults)
	state, err := s.ProcessEvent(majongpb.EventID_event_mopai_finish, nil, flow)
	afterResults := ""
	afterResults += fmt.Sprintln("查询后:")
	afterResults += FmtPlayerInfo(player)
	logrus.Info(afterResults)
	assert.Nil(t, err)
	assert.Equal(t, majongpb.StateID_state_zixun, state, "查询成功后有特殊操作，应该进入自询状态")
}
