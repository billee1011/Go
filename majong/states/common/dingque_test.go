package common

import (
	"fmt"
	"steve/majong/interfaces"
	majongpb "steve/entity/majong"
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/golang/mock/gomock"
	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
)

// 测试定缺状态未完成，还是定缺状态
func TestDingQueState_WeiWangCheng(t *testing.T) {
	mc := gomock.NewController(t)
	defer mc.Finish()

	//初始化牌局信息
	flow := interfaces.NewMockMajongFlow(mc)
	mjContext := majongpb.MajongContext{
		Players:      []*majongpb.Player{},
		ActivePlayer: 1,
	}

	// 初始玩家信息
	mjContext.Players = mjContext.Players[0:0]
	for i := 0; i < 4; i++ {
		mjContext.Players = append(mjContext.Players, &majongpb.Player{
			PalyerId:  uint64(i),
			HandCards: []*majongpb.Card{&Card1W, &Card1T, &Card1W, &Card1W},
		})
	}
	// 序列化消息
	dingqueEvent := &majongpb.DingqueRequestEvent{
		Head:  &majongpb.RequestEventHead{PlayerId: 1},
		Color: majongpb.CardColor_ColorWan,
	}
	eventContext, err := proto.Marshal(dingqueEvent)
	if err != nil {
		fmt.Println(err)
	}
	flow.EXPECT().GetMajongContext().Return(&mjContext).AnyTimes()

	start := "定缺状态"
	logrus.WithFields(logrus.Fields{
		"状态":                 start,
		"DingqueColor（默认为W）": mjContext.Players[1].DingqueColor,
		"HasDingque":         mjContext.Players[1].HasDingque,
	}).Info("前")
	// 定缺状态接受到定缺消息
	d := new(DingqueState)
	newStateID, err := d.ProcessEvent(majongpb.EventID_event_dingque_request, eventContext, flow)
	logrus.WithFields(logrus.Fields{
		"状态":           start,
		"DingqueColor": mjContext.Players[1].DingqueColor,
		"HasDingque":   mjContext.Players[1].HasDingque,
	}).Info("后")
	assert.Equal(t, mjContext.Players[1].DingqueColor, majongpb.CardColor_ColorWan, "定缺")
	assert.Equal(t, mjContext.Players[1].HasDingque, true, "定缺")
	assert.Equal(t, majongpb.StateID_state_dingque, newStateID, "定缺")
}

// 测试定缺状态完成，进入自询状态
func TestDingQueState_WangCheng_ZiXun(t *testing.T) {
	mc := gomock.NewController(t)
	defer mc.Finish()

	//初始化牌局信息
	flow := interfaces.NewMockMajongFlow(mc)
	mjContext := majongpb.MajongContext{
		Players:      []*majongpb.Player{},
		ActivePlayer: 1,
	}

	// 初始玩家信息，未每个玩家都设置一遍定缺
	mjContext.Players = mjContext.Players[0:0]
	for i := 0; i < 4; i++ {
		mjContext.Players = append(mjContext.Players, &majongpb.Player{
			PalyerId:  uint64(i),
			HandCards: []*majongpb.Card{&Card1W, &Card1T, &Card1W, &Card1W},
		})
	}
	// 玩家一个个去定缺，最后一个人定缺完后，返回自询
	for i := 0; i < len(mjContext.Players); i++ {
		// 序列化消息
		dingqueEvent := &majongpb.DingqueRequestEvent{
			Head:  &majongpb.RequestEventHead{PlayerId: mjContext.Players[i].PalyerId},
			Color: majongpb.CardColor_ColorTong,
		}
		eventContext, err := proto.Marshal(dingqueEvent)
		assert.Nil(t, err)

		flow.EXPECT().GetMajongContext().Return(&mjContext).AnyTimes()
		start := "定缺状态"
		logrus.WithFields(logrus.Fields{
			"状态":                 start,
			"DingqueColor（默认为W）": mjContext.Players[i].DingqueColor,
			"HasDingque":         mjContext.Players[i].HasDingque,
		}).Info("前")
		// 定缺状态接受到定缺消息
		d := new(DingqueState)
		newStateID, err := d.ProcessEvent(majongpb.EventID_event_dingque_request, eventContext, flow)
		if newStateID == majongpb.StateID_state_zixun {
			start = "自询状态"
		}
		logrus.WithFields(logrus.Fields{
			"状态":           start,
			"DingqueColor": mjContext.Players[i].DingqueColor,
			"HasDingque":   mjContext.Players[i].HasDingque,
		}).Info("后")

		assert.Nil(t, err)
		if newStateID == majongpb.StateID_state_zixun {
			assert.Equal(t, mjContext.Players[i].DingqueColor, majongpb.CardColor_ColorTong, "定缺转自询")
			assert.Equal(t, mjContext.Players[i].HasDingque, true, "定缺转自询")
			assert.Equal(t, majongpb.StateID_state_zixun, newStateID, "定缺转自询")
		} else {
			assert.Equal(t, mjContext.Players[i].DingqueColor, majongpb.CardColor_ColorTong, "定缺")
			assert.Equal(t, mjContext.Players[i].HasDingque, true, "定缺")
			assert.Equal(t, majongpb.StateID_state_dingque, newStateID, "定缺")
		}
	}

}
