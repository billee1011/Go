package states

import (
	"steve/majong/interfaces"
	majongpb "steve/server_pb/majong"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/Sirupsen/logrus"
)

// 测试暗杠状态接收到摸牌消息
func TestAnGangState_MoPai(t *testing.T) {
	mc := gomock.NewController(t)
	defer mc.Finish()

	//初始化牌局信息
	flow := interfaces.NewMockMajongFlow(mc)
	mjContext := majongpb.MajongContext{
		Players:      []*majongpb.Player{},
		ActivePlayer: 1,
	}

	flow.EXPECT().GetMajongContext().Return(&mjContext).AnyTimes()
	start := "暗杠状态"
	logrus.WithFields(logrus.Fields{
		"状态": start,
	}).Info("前")
	// 暗杠状态接受到摸牌消息
	gangState := new(AnGangState)
	newStateID, err := gangState.ProcessEvent(majongpb.EventID_event_mopai_finish, nil, flow)
	if newStateID == majongpb.StateID_state_mopai {
		start = "摸牌状态"
	}
	logrus.WithFields(logrus.Fields{
		"状态": start,
	}).Info("后")
	assert.Nil(t, err)
	assert.Equal(t, majongpb.StateID_state_mopai, newStateID)
}

