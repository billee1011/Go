package states

import (
	"math/rand"
	"steve/majong/interfaces"
	majongpb "steve/server_pb/majong"
	"time"
)

// XipaiState 洗牌状态
type XipaiState struct {
}

var _ interfaces.MajongState = new(XipaiState)

// ProcessEvent 处理事件
func (s *XipaiState) ProcessEvent(eventID majongpb.EventID, eventContext []byte, flow interfaces.MajongFlow) (newState majongpb.StateID, err error) {
	if eventID == majongpb.EventID_event_xipai_finish {
		return majongpb.StateID(majongpb.StateID_state_fapai), nil
	}
	return majongpb.StateID(majongpb.StateID_state_xipai), nil
}

func (s *XipaiState) genOriginCards(flow interfaces.MajongFlow) []*majongpb.Card {
	gameID := flow.GetMajongContext().GetGameId()
	return getOriginCards(int(gameID))
}

func (s *XipaiState) xipai(flow interfaces.MajongFlow) []*majongpb.Card {
	cards := s.genOriginCards(flow)
	rand.Seed(time.Now().Unix())
	rand.Shuffle(len(cards), func(i, j int) {
		tmp := cards[i]
		cards[i] = cards[j]
		cards[j] = tmp
	})
	return cards
}

// OnEntry 进入状态
func (s *XipaiState) OnEntry(flow interfaces.MajongFlow) {
	flow.GetMajongContext().WallCards = s.xipai(flow)

	flow.SetAutoEvent(majongpb.AutoEvent{
		EventId:      majongpb.EventID_event_xipai_finish,
		EventContext: nil,
	})
}

// OnExit 退出状态
func (s *XipaiState) OnExit(flow interfaces.MajongFlow) {
}
