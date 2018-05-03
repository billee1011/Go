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

var (
	// Card1W 1 万
	Card1W = majongpb.Card{Color: majongpb.CardColor_ColorWan, Point: 1}
	// Card2W 2 万
	Card2W = majongpb.Card{Color: majongpb.CardColor_ColorWan, Point: 2}
	// Card3W 3 万
	Card3W = majongpb.Card{Color: majongpb.CardColor_ColorWan, Point: 3}
	// Card4W 4 万
	Card4W = majongpb.Card{Color: majongpb.CardColor_ColorWan, Point: 4}
	// Card5W 5 万
	Card5W = majongpb.Card{Color: majongpb.CardColor_ColorWan, Point: 5}
	// Card6W 6 万
	Card6W = majongpb.Card{Color: majongpb.CardColor_ColorWan, Point: 6}
	// Card7W 7 万
	Card7W = majongpb.Card{Color: majongpb.CardColor_ColorWan, Point: 7}
	// Card8W 8 万
	Card8W = majongpb.Card{Color: majongpb.CardColor_ColorWan, Point: 8}
	// Card9W 9 万
	Card9W = majongpb.Card{Color: majongpb.CardColor_ColorWan, Point: 9}

	// Card1T 1 条
	Card1T = majongpb.Card{Color: majongpb.CardColor_ColorTiao, Point: 1}
	// Card2T 2 条
	Card2T = majongpb.Card{Color: majongpb.CardColor_ColorTiao, Point: 2}
	// Card3T 3 条
	Card3T = majongpb.Card{Color: majongpb.CardColor_ColorTiao, Point: 3}
	// Card4T 4 条
	Card4T = majongpb.Card{Color: majongpb.CardColor_ColorTiao, Point: 4}
	// Card5T 5 条
	Card5T = majongpb.Card{Color: majongpb.CardColor_ColorTiao, Point: 5}
	// Card6T 6 条
	Card6T = majongpb.Card{Color: majongpb.CardColor_ColorTiao, Point: 6}
	// Card7T 7 条
	Card7T = majongpb.Card{Color: majongpb.CardColor_ColorTiao, Point: 7}
	// Card8T 8 条
	Card8T = majongpb.Card{Color: majongpb.CardColor_ColorTiao, Point: 8}
	// Card9T 9 条
	Card9T = majongpb.Card{Color: majongpb.CardColor_ColorTiao, Point: 9}

	// Card1B 1 筒
	Card1B = majongpb.Card{Color: majongpb.CardColor_ColorTong, Point: 1}
	// Card2B 2 筒
	Card2B = majongpb.Card{Color: majongpb.CardColor_ColorTong, Point: 2}
	// Card3B 3 筒
	Card3B = majongpb.Card{Color: majongpb.CardColor_ColorTong, Point: 3}
	// Card4B 4 筒
	Card4B = majongpb.Card{Color: majongpb.CardColor_ColorTong, Point: 4}
	// Card5B 5 筒
	Card5B = majongpb.Card{Color: majongpb.CardColor_ColorTong, Point: 5}
	// Card6B 6 筒
	Card6B = majongpb.Card{Color: majongpb.CardColor_ColorTong, Point: 6}
	// Card7B 7 筒
	Card7B = majongpb.Card{Color: majongpb.CardColor_ColorTong, Point: 7}
	// Card8B 8 筒
	Card8B = majongpb.Card{Color: majongpb.CardColor_ColorTong, Point: 8}
	// Card9B 9 筒
	Card9B = majongpb.Card{Color: majongpb.CardColor_ColorTong, Point: 9}
)

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
