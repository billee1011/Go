package states

import (
	msgid "steve/client_pb/msgId"
	"steve/client_pb/room"
	"steve/majong/interfaces"
	"steve/majong/utils"
	majongpb "steve/server_pb/majong"

	"github.com/golang/protobuf/proto"
)

// ZimoState 自摸状态
type ZimoState struct {
}

var _ interfaces.MajongState = new(ZimoState)

// ProcessEvent 处理事件
func (s *ZimoState) ProcessEvent(eventID majongpb.EventID, eventContext []byte, flow interfaces.MajongFlow) (newState majongpb.StateID, err error) {
	if eventID == majongpb.EventID_event_zimo_finish {
		s.zimo(flow)
		return majongpb.StateID_state_mopai, nil
	}
	return majongpb.StateID_state_zimo, errInvalidEvent
}

func (s *ZimoState) zimo(flow interfaces.MajongFlow) {

	ctx := flow.GetMajongContext()

	activePlayer := utils.GetPlayerByID(ctx.Players, ctx.GetMopaiPlayer())
	card := ctx.GetLastMopaiCard()
	activePlayer.HandCards, _ = utils.DeleteCardFromLast(activePlayer.HandCards, card)
	huCard := &majongpb.HuCard{
		Card:      card,
		SrcPlayer: activePlayer.GetPalyerId(),
		Type:      majongpb.HuType_hu_zimo,
	}
	activePlayer.HuCards = append(activePlayer.HuCards, huCard)
	// activePlayer.PossibleActions = activePlayer.PossibleActions[:0]
	toclientCard, _ := utils.CardToRoomCard(card)
	playersID := make([]uint64, 0, 0)
	for _, player := range flow.GetMajongContext().GetPlayers() {
		playersID = append(playersID, player.PalyerId)
	}
	ntf := &room.RoomHuNtf{
		Players:      []uint64{activePlayer.PalyerId},
		Card:         toclientCard,
		HuType:       room.HuType_ZiMo.Enum(),
		FromPlayerId: proto.Uint64(activePlayer.GetPalyerId()),
	}
	toClientMessage := interfaces.ToClientMessage{
		MsgID: int(msgid.MsgID_ROOM_HU_NTF),
		Msg:   ntf,
	}
	flow.PushMessages(playersID, toClientMessage)
}

// OnEntry 进入状态
func (s *ZimoState) OnEntry(flow interfaces.MajongFlow) {
	flow.SetAutoEvent(majongpb.AutoEvent{
		EventId:      majongpb.EventID_event_zimo_finish,
		EventContext: nil,
	})
}

// OnExit 退出状态
func (s *ZimoState) OnExit(flow interfaces.MajongFlow) {

}
