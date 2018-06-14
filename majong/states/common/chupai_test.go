package common

import (
	"steve/majong/interfaces"
	majongpb "steve/server_pb/majong"
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

//TestChupaiState_chupaiwenxun 出牌后，其他玩家有特殊操作，进入出牌问询状态
func TestChupaiState_chupaiwenxun(t *testing.T) {
	mc := gomock.NewController(t)
	defer mc.Finish()
	flow := interfaces.NewMockMajongFlow(mc)
	flow.EXPECT().PushMessages(gomock.Any(), gomock.Any()).AnyTimes()
	flow.EXPECT().GetMajongContext().Return(
		&majongpb.MajongContext{
			GameId: 1,
			Players: []*majongpb.Player{
				&majongpb.Player{
					PalyerId:        1,
					HandCards:       []*majongpb.Card{&Card1W, &Card1W, &Card1W, &Card2W, &Card2W, &Card2W, &Card2W, &Card3W, &Card3W, &Card3W, &Card3W, &Card4W, &Card4W},
					DingqueColor:    majongpb.CardColor_ColorTiao,
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
			ActivePlayer:     1,
			WallCards:        []*majongpb.Card{&Card9B, &Card9B},
			LastOutCard:      &Card1T,
			LastChupaiPlayer: 1,
		},
	).AnyTimes()
	s := ChupaiState{}
	context := flow.GetMajongContext()
	players := context.GetPlayers()
	for _, player := range players {
		logrus.WithFields(FmtPlayerInfo(player)).Info("出牌查询前")
		// beforeResults := ""
		// beforeResults += fmt.Sprintln("出牌查询前:")
		// beforeResults += FmtPlayerInfo(player)
		// logrus.Info(beforeResults)

	}
	logrus.WithFields(FmtMajongContxt(context)).Info("出牌前的麻将现场")
	state, err := s.ProcessEvent(majongpb.EventID_event_chupai_finish, nil, flow)
	assert.Nil(t, err)
	assert.Equal(t, majongpb.StateID_state_chupaiwenxun, state, "查询成功后有特殊操作，应该进入出牌问询状态")
	for _, player := range players {
		logrus.WithFields(FmtPlayerInfo(player)).Info("出牌查询后")
		// afterResults := ""
		// afterResults += fmt.Sprintln("出牌查询后:")
		// afterResults += FmtPlayerInfo(player)
		// logrus.Info(afterResults)
	}
	logrus.WithFields(FmtMajongContxt(context)).Info("出牌后的麻将现场")
}

//TestChupaiState_mopai 出牌后，其他玩家没有特殊操作，进入摸牌状态（下家摸牌）
func TestChupaiState_mopai(t *testing.T) {
	mc := gomock.NewController(t)
	defer mc.Finish()
	flow := interfaces.NewMockMajongFlow(mc)
	flow.EXPECT().PushMessages(gomock.Any(), gomock.Any()).AnyTimes()
	flow.EXPECT().GetMajongContext().Return(
		&majongpb.MajongContext{
			GameId: 1,
			Players: []*majongpb.Player{
				&majongpb.Player{
					PalyerId:        1,
					HandCards:       []*majongpb.Card{&Card1T, &Card1W, &Card1W, &Card2W, &Card2W, &Card2W, &Card2W, &Card3W, &Card3W, &Card3W, &Card3W, &Card4W, &Card4W},
					DingqueColor:    majongpb.CardColor_ColorTiao,
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
			ActivePlayer:     1,
			WallCards:        []*majongpb.Card{&Card9B, &Card9B},
			LastOutCard:      &Card1W,
			LastChupaiPlayer: 1,
		},
	).AnyTimes()
	s := ChupaiState{}
	context := flow.GetMajongContext()
	players := context.GetPlayers()
	for _, player := range players {
		logrus.WithFields(FmtPlayerInfo(player)).Info("出牌查询前")
		// beforeResults := ""
		// beforeResults += fmt.Sprintln("出牌查询前:")
		// beforeResults += FmtPlayerInfo(player)
		// logrus.Info(beforeResults)

	}
	logrus.WithFields(FmtMajongContxt(context)).Info("出牌前的麻将现场")
	state, err := s.ProcessEvent(majongpb.EventID_event_chupai_finish, nil, flow)
	assert.Nil(t, err)
	assert.Equal(t, majongpb.StateID_state_mopai, state, "查询成功后无特殊操作，应该进入摸牌状态")
	for _, player := range players {
		logrus.WithFields(FmtPlayerInfo(player)).Info("出牌查询后")
		// afterResults := ""
		// afterResults += fmt.Sprintln("出牌查询后:")
		// afterResults += FmtPlayerInfo(player)
		// logrus.Info(afterResults)
	}
	logrus.WithFields(FmtMajongContxt(context)).Info("出牌后的麻将现场")
}
