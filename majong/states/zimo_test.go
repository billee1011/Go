package states

import (
	"steve/majong/interfaces"
	"steve/majong/utils"
	"testing"

	majongpb "steve/server_pb/majong"

	"github.com/Sirupsen/logrus"
	"github.com/golang/mock/gomock"
	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
)

// Test_InitState_ProcessEventStartGame 测试初始状态接收游戏开始事件
func Test_ZimoState_zimo(t *testing.T) {
	// mc := gomock.NewController(t)
	// flow := interfaces.NewMockMajongFlow(mc)
	// flow.EXPECT().GetMajongContext().Return(
	// 	&majongpb.MajongContext{
	// 		GameId:       1,
	// 		WallCards:    []*majongpb.Card{&Card1W},
	// 		Players:      []*majongpb.Player{&majongpb.Player{PalyerId: 1}, &majongpb.Player{PalyerId: 2}},
	// 		ActivePlayer: 1,
	// 	},
	// ).AnyTimes()

	// s := ZimoState{}
	// context := flow.GetMajongContext()
	// beforeResults := ""
	// beforeResults += fmt.Sprintln("自摸状态：")
	// for _, player := range context.Players {
	// 	beforeResults += FmtPlayerInfo(player)
	// }
	// logrus.Info(beforeResults)
	// states, err := s.ProcessEvent(majongpb.EventID_event_zimo_finish, nil, flow)
	// assert.Nil(t, err)
	// assert.Equal(t, majongpb.StateID_state_mopai, states, "自摸状态自动跳转到摸牌状态")
	// afterResults := ""
	// afterResults += fmt.Sprintln("摸牌状态：")
	// for _, player := range context.Players {
	// 	afterResults += FmtPlayerInfo(player)
	// }
	// logrus.Info(afterResults)
	mc := gomock.NewController(t)
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
			MopaiPlayer:   1,
			LastMopaiCard: &Card4W,
			WallCards:     []*majongpb.Card{&Card1T},
		},
	).AnyTimes()

	s := ZimoState{}
	autoEvent := majongpb.AutoEvent{
		EventId:      majongpb.EventID_event_zimo_finish,
		EventContext: nil,
	}
	requestEvent, err := proto.Marshal(&autoEvent)
	context := flow.GetMajongContext()
	player := utils.GetPlayerByID(context.GetPlayers(), context.GetMopaiPlayer())
	logrus.WithFields(FmtPlayerInfo(player)).Info("自摸前")
	// beforeResults := ""
	// beforeResults += fmt.Sprintln("before自摸：")
	// beforeResults += FmtPlayerInfo(player)
	// logrus.Info(beforeResults)

	stateID, err := s.ProcessEvent(majongpb.EventID_event_zimo_finish, requestEvent, flow)
	assert.Nil(t, err)
	assert.Equal(t, majongpb.StateID_state_mopai, stateID, "执行自摸操作成功后，状态应该为自摸状态")
	logrus.WithFields(FmtPlayerInfo(player)).Info("自摸后")
	// results := ""
	// results += fmt.Sprintln("after自摸：")
	// results += FmtPlayerInfo(player)
	// logrus.Info(results)
}
