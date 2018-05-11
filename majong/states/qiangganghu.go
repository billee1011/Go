package states

import (
	"steve/client_pb/msgId"
	"steve/client_pb/room"
	"steve/majong/interfaces"
	"steve/majong/utils"
	majongpb "steve/server_pb/majong"

	"github.com/golang/protobuf/proto"
)

// QiangganghuState 抢杠胡状态
type QiangganghuState struct {
}

var _ interfaces.MajongState = new(QiangganghuState)

// ProcessEvent 处理事件
func (s *QiangganghuState) ProcessEvent(eventID majongpb.EventID, eventContext []byte, flow interfaces.MajongFlow) (newState majongpb.StateID, err error) {
	if eventID == majongpb.EventID_event_qiangganghu_finish {
		s.qiangganghu(flow)
		return majongpb.StateID_state_mopai, nil
	}
	return majongpb.StateID_state_qiangganghu, errInvalidEvent
}

func (s *QiangganghuState) qiangganghu(flow interfaces.MajongFlow) {
	context := flow.GetMajongContext()
	card := context.GetLastMopaiCard()
	huPlayers := context.GetLastHuPlayers()
	//抢杠胡的玩家，添加胡的牌
	for _, huPlayerID := range huPlayers {
		huPlayer := utils.GetPlayerByID(context.GetPlayers(), huPlayerID)
		huPlayer.HuCards = append(huPlayer.HuCards, &majongpb.HuCard{
			Card:      card,
			Type:      majongpb.HuType_hu_qiangganghu,
			SrcPlayer: context.GetLastGangPlayer(),
		})
	}
	//被抢杠的玩家移除手牌
	srcPlayer := utils.GetPlayerByID(context.GetPlayers(), context.GetLastMopaiPlayer())
	srcPlayer.HandCards, _ = utils.RemoveCards(srcPlayer.HandCards, card, 1)
	//广播抢杠胡成功的消息
	roomCard, _ := utils.CardToRoomCard(card)
	huNtf := &room.RoomHuNtf{
		Players:      huPlayers,
		FromPlayerId: proto.Uint64(context.GetLastMopaiPlayer()),
		Card:         roomCard,
		HuType:       room.HuType_QiangGangHu.Enum(),
	}
	toClientMessage := interfaces.ToClientMessage{
		MsgID: int(msgid.MsgID_ROOM_HU_NTF),
		Msg:   huNtf,
	}
	playerIDs := []uint64{}
	for _, p := range context.GetPlayers() {
		playerIDs = append(playerIDs, p.GetPalyerId())
	}
	flow.PushMessages(playerIDs, toClientMessage)
	//通知完成后，将进入摸牌状态，这里将重置mopaiPlayer
	//TODO:先暂时取胡牌玩家列表中的最后一个玩家
	lastHuPlayerID := huPlayers[len(huPlayers)-1]
	nextMopaiPlayer := utils.GetNextPlayerByID(context.GetPlayers(), lastHuPlayerID)
	context.MopaiPlayer = nextMopaiPlayer.GetPalyerId()
}

// OnEntry 进入状态
func (s *QiangganghuState) OnEntry(flow interfaces.MajongFlow) {
	flow.SetAutoEvent(majongpb.AutoEvent{
		EventId:      majongpb.EventID_event_qiangganghu_finish,
		EventContext: nil,
	})
}

// OnExit 退出状态
func (s *QiangganghuState) OnExit(flow interfaces.MajongFlow) {

}
