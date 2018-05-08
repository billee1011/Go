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
			GameId:       1,
			WallCards:    []*majongpb.Card{&Card1W},
			Players:      []*majongpb.Player{&majongpb.Player{PalyerId: 1}, &majongpb.Player{PalyerId: 2}},
			ActivePlayer: 1,
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
