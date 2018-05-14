package states

import (
	"math/rand"
	msgid "steve/client_pb/msgId"
	"steve/client_pb/room"
	"steve/majong/interfaces"
	"steve/majong/interfaces/facade"
	"steve/majong/utils"
	"steve/peipai"
	majongpb "steve/server_pb/majong"
	"strconv"
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
	gameName := getGameName(flow)
	PeiPai(cards, flow.GetMajongContext(), gameName)
	length := peipai.GetLensOfWallCards(gameName)
	if length != 0 {
		cards = cards[:length]
	}
	//暂时只能配牌和长度，换三张方向的代码已写，具体怎么操作要和向xuzhang了解下
	return cards
}

func getGameName(flow interfaces.MajongFlow) string {
	gameID := flow.GetMajongContext().GetGameId()
	var gameName string
	if gameID == 1 {
		gameName = "scxl"
	}
	return gameName
}

// PeiPai 配牌工具
func PeiPai(wallCards []*majongpb.Card, context *majongpb.MajongContext, gameName string) (bool, []*majongpb.Card) {
	value, err := peipai.GetPeiPai(gameName)
	if err != nil {
		return false, wallCards
	}
	var cards []*majongpb.Card
	for i := 0; i < len(value); i = i + 3 {
		ca, err := strconv.Atoi(value[i : i+2])
		if err != nil {
			return false, wallCards
		}
		card, err := utils.IntToCard(int32(ca))
		if err != nil {
			return false, wallCards
		}
		cards = append(cards, card)
	}
	for i := 0; i < len(cards); i++ {
		for j := len(wallCards) - 1; j >= 0; j-- {
			if utils.CardEqual(cards[i], wallCards[j]) {
				wallCards[i], wallCards[j] = wallCards[j], wallCards[i]
			}
		}
	}
	return true, wallCards
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
