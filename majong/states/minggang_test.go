package states

import (
	"steve/majong/interfaces"
	majongpb "steve/server_pb/majong"
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

// 测试明杠状态接收到非摸牌消息
func TestMingGangState_MoPai(t *testing.T) {
	mc := gomock.NewController(t)
	defer mc.Finish()

	//初始化牌局信息
	flow := interfaces.NewMockMajongFlow(mc)
	mjContext := majongpb.MajongContext{
		Players:      []*majongpb.Player{},
		ActivePlayer: 1,
	}

	flow.EXPECT().GetMajongContext().Return(&mjContext).AnyTimes()
	start := "明杠状态"
	logrus.WithFields(logrus.Fields{
		"状态": start,
	}).Info("前")
	// 明杠状态接受到明杠完成消息
	gangState := new(MingGangState)
	newStateID, err := gangState.ProcessEvent(majongpb.EventID_event_gang_finish, nil, flow)
	if newStateID == majongpb.StateID_state_mopai {
		start = "摸牌状态"
	}
	logrus.WithFields(logrus.Fields{
		"状态": start,
	}).Info("后")
	assert.Nil(t, err)
	assert.Equal(t, majongpb.StateID_state_mopai, newStateID)
}
