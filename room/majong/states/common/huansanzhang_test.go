package common

import (
	majongpb "steve/entity/majong"
	"steve/room/majong/interfaces"
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/golang/mock/gomock"
	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
)

func TestHuanSanZhangState_huansanzhang(t *testing.T) {
	mc := gomock.NewController(t)
	defer mc.Finish()

	wallCards := getOriginCards(0)

	flow := interfaces.NewMockMajongFlow(mc)

	mjContext := majongpb.MajongContext{
		Players:        []*majongpb.Player{},
		WallCards:      wallCards,
		ZhuangjiaIndex: 0,
	}
	initPlayers(&mjContext)

	makeCard(&mjContext)

	flow.EXPECT().GetMajongContext().Return(&mjContext).AnyTimes()

	f := new(HuansanzhangState)

	evenContext := &majongpb.HuansanzhangRequestEvent{
		Head: &majongpb.RequestEventHead{
			PlayerId: 0,
		},
		Cards: []*majongpb.Card{&Card1W, &Card1W, &Card1W},
		Sure:  true,
	}
	newState, _ := f.ProcessEvent(majongpb.EventID_event_huansanzhang_request, evenContext, flow)

	event = &majongpb.HuansanzhangRequestEvent{
		Head: &majongpb.RequestEventHead{
			PlayerId: 1,
		},
		Cards: []*majongpb.Card{&Card3W, &Card3W, &Card3W},
		Sure:  true,
	}
	evenContext, _ = proto.Marshal(event)
	newState, _ = f.ProcessEvent(majongpb.EventID_event_huansanzhang_request, evenContext, flow)

	event = &majongpb.HuansanzhangRequestEvent{
		Head: &majongpb.RequestEventHead{
			PlayerId: 2,
		},
		Cards: []*majongpb.Card{&Card5W, &Card5W, &Card5W},
		Sure:  true,
	}
	evenContext, _ = proto.Marshal(event)
	newState, _ = f.ProcessEvent(majongpb.EventID_event_huansanzhang_request, evenContext, flow)

	event = &majongpb.HuansanzhangRequestEvent{
		Head: &majongpb.RequestEventHead{
			PlayerId: 3,
		},
		Cards: []*majongpb.Card{&Card7W, &Card7W, &Card7W},
		Sure:  true,
	}
	evenContext, _ = proto.Marshal(event)
	newState, _ = f.ProcessEvent(majongpb.EventID_event_huansanzhang_request, evenContext, flow)

	logrus.Info(newState)
	assert.Equal(t, newState, majongpb.StateID_state_dingque)
}

func makeCard(mjContext *majongpb.MajongContext) {
	mjContext.Players[0].HandCards = []*majongpb.Card{&Card1W, &Card1W, &Card1W, &Card2W, &Card2W, &Card2W, &Card2W}
	mjContext.Players[1].HandCards = []*majongpb.Card{&Card3W, &Card3W, &Card3W, &Card4W, &Card4W, &Card4W, &Card4W}
	mjContext.Players[2].HandCards = []*majongpb.Card{&Card5W, &Card5W, &Card5W, &Card6W, &Card6W, &Card6W, &Card6W}
	mjContext.Players[3].HandCards = []*majongpb.Card{&Card7W, &Card7W, &Card7W, &Card8W, &Card8W, &Card8W, &Card8W}
}
