package common

import (
	"steve/majong/interfaces"
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	majongpb "steve/entity/majong"

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
					HandCards:    []*majongpb.Card{&Card1W, &Card1W, &Card1W, &Card1W, &Card2W, &Card2W, &Card2W, &Card2W, &Card3W, &Card3W, &Card3W, &Card3W, &Card4W},
					DingqueColor: majongpb.CardColor_ColorTiao,
				},
			},
			MopaiPlayer: 1,
			WallCards:   []*majongpb.Card{&Card4W, &Card1T},
		},
	).AnyTimes()
	s := MoPaiState{}
	context := flow.GetMajongContext()
	players := context.GetPlayers()
	player := players[0]
	logrus.WithFields(FmtPlayerInfo(player)).Info("查询前")
	logrus.WithFields(FmtMajongContxt(context)).Info("摸牌前的麻将现场")
	state, err := s.ProcessEvent(majongpb.EventID_event_mopai_finish, nil, flow)
	assert.Nil(t, err)
	assert.Equal(t, majongpb.StateID_state_zixun, state, "查询成功后有特殊操作，应该进入自询状态")
	logrus.WithFields(FmtPlayerInfo(player)).Info("查询后")
	logrus.WithFields(FmtMajongContxt(context)).Info("摸牌后的麻将现场")
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
					HandCards:    []*majongpb.Card{&Card1W, &Card2W, &Card2W, &Card2W, &Card2W, &Card3W, &Card3W, &Card3W, &Card3W, &Card4W},
					DingqueColor: majongpb.CardColor_ColorTiao,
					PengCards: []*majongpb.PengCard{
						&majongpb.PengCard{
							Card:      &Card1W,
							SrcPlayer: 2,
						},
					},
				},
			},
			MopaiPlayer: 1,
			WallCards:   []*majongpb.Card{&Card4W, &Card1T},
		},
	).AnyTimes()
	s := MoPaiState{}
	context := flow.GetMajongContext()
	players := context.GetPlayers()
	player := players[0]
	logrus.WithFields(FmtPlayerInfo(player)).Info("查询前")
	logrus.WithFields(FmtMajongContxt(context)).Info("摸牌前的麻将现场")
	state, err := s.ProcessEvent(majongpb.EventID_event_mopai_finish, nil, flow)
	assert.Nil(t, err)
	assert.Equal(t, majongpb.StateID_state_zixun, state, "查询成功后有特殊操作，应该进入自询状态")
	logrus.WithFields(FmtPlayerInfo(player)).Info("查询后")
	logrus.WithFields(FmtMajongContxt(context)).Info("摸牌后的麻将现场")
}

func Test_MopaiState_mopai2(t *testing.T) {
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
					HandCards:    []*majongpb.Card{&Card1W, &Card2W, &Card2W, &Card2W, &Card2W, &Card3W, &Card3W, &Card3W, &Card3W, &Card4W},
					DingqueColor: majongpb.CardColor_ColorTiao,
					PengCards: []*majongpb.PengCard{
						&majongpb.PengCard{
							Card:      &Card1W,
							SrcPlayer: 2,
						},
					},
				},
			},
			MopaiPlayer: 1,
			WallCards:   []*majongpb.Card{},
		},
	).AnyTimes()
	s := MoPaiState{}
	context := flow.GetMajongContext()
	players := context.GetPlayers()
	player := players[0]
	logrus.WithFields(FmtPlayerInfo(player)).Info("查询前")
	logrus.WithFields(FmtMajongContxt(context)).Info("摸牌前的麻将现场")

	state, err := s.ProcessEvent(majongpb.EventID_event_mopai_finish, nil, flow)
	assert.Nil(t, err)
	assert.Equal(t, majongpb.StateID_state_gameover, state, "无牌可摸，应该进入gameover状态")
	logrus.WithFields(FmtPlayerInfo(player)).Info("查询后")
	logrus.WithFields(FmtMajongContxt(context)).Info("摸牌后的麻将现场")

}
