package states

import (
	"math/rand"
	msgid "steve/client_pb/msgId"
	"steve/client_pb/room"
	"steve/majong/interfaces"
	"steve/majong/interfaces/facade"
	majongpb "steve/server_pb/majong"
	"time"

	"github.com/golang/protobuf/proto"
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
	rand.Seed(int64(time.Now().Nanosecond()))
	rand.Shuffle(len(cards), func(i, j int) {
		cards[i], cards[j] = cards[j], cards[i]
	})
	return cards
}

// randDices 随机筛子
func (s *XipaiState) randDices() [2]int {
	rand.Seed(int64(time.Now().Nanosecond()))
	return [2]int{
		rand.Intn(6) + 1,
		rand.Intn(6) + 1,
	}
}

// selectZhuangjia 选择庄家
func (s *XipaiState) selectZhuangjia(mjContext *majongpb.MajongContext, dices [2]int) int {
	totalDice := dices[0] + dices[1]

	mjContext.ZhuangjiaIndex = uint32(totalDice % len(mjContext.Players))
	return int(mjContext.ZhuangjiaIndex)
}

// pushMessages 发送消息给玩家
func (s *XipaiState) pushMessages(cardCount int, dices [2]int, zjIndex int, flow interfaces.MajongFlow) {
	protoDices := []uint32{uint32(dices[0]), uint32(dices[1])}
	facade.BroadcaseMessage(flow, msgid.MsgID_ROOM_XIPAI_NTF, &room.RoomXipaiNtf{
		Dices:          protoDices,
		TotalCard:      proto.Uint32(uint32(cardCount)),
		ZhuangjiaIndex: proto.Uint32(uint32(zjIndex)),
	})
}

// OnEntry 进入状态
func (s *XipaiState) OnEntry(flow interfaces.MajongFlow) {
	mjContext := flow.GetMajongContext()

	mjContext.WallCards = s.xipai(flow)
	dices := s.randDices()
	zjIndex := s.selectZhuangjia(mjContext, dices)

	s.pushMessages(len(mjContext.WallCards), dices, zjIndex, flow)

	flow.SetAutoEvent(majongpb.AutoEvent{
		EventId:      majongpb.EventID_event_xipai_finish,
		EventContext: nil,
	})
}

// OnExit 退出状态
func (s *XipaiState) OnExit(flow interfaces.MajongFlow) {
}
