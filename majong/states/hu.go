package states

import (
	msgid "steve/client_pb/msgId"
	"steve/client_pb/room"
	"steve/majong/interfaces"
	"steve/majong/utils"
	majongpb "steve/server_pb/majong"

	"github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
)

// HuState 胡状态
type HuState struct {
}

var _ interfaces.MajongState = new(HuState)

// ProcessEvent 处理事件
func (s *HuState) ProcessEvent(eventID majongpb.EventID, eventContext []byte, flow interfaces.MajongFlow) (newState majongpb.StateID, err error) {
	if eventID == majongpb.EventID_event_hu_finish {
		s.hu(flow)
		return majongpb.StateID_state_mopai, nil
	}
	return majongpb.StateID_state_hu, errInvalidEvent
}

// OnEntry 进入状态
func (s *HuState) OnEntry(flow interfaces.MajongFlow) {
	flow.SetAutoEvent(majongpb.AutoEvent{
		EventId:      majongpb.EventID_event_hu_finish,
		EventContext: nil,
	})
}

func (s *HuState) hu(flow interfaces.MajongFlow) {
	mjcontext := flow.GetMajongContext()
	lastPlayer := utils.GetPlayerByID(mjcontext.GetPlayers(), mjcontext.GetLastChupaiPlayer())
	//为胡的玩家们添加huCard
	for _, huPlayerID := range mjcontext.GetLastHuPlayers() {
		huPlayer := utils.GetPlayerByID(mjcontext.GetPlayers(), huPlayerID)
		huPlayer.HuCards = append(huPlayer.HuCards, &majongpb.HuCard{
			Card:      mjcontext.GetLastOutCard(),
			SrcPlayer: mjcontext.GetLastChupaiPlayer(),
			Type:      majongpb.HuType_hu_dianpao,
		})
	}
	//广播点炮胡成功的消息
	roomCard, _ := utils.CardToRoomCard(mjcontext.GetLastOutCard())
	huNtf := &room.RoomHuNtf{
		Players:      mjcontext.GetLastHuPlayers(),
		FromPlayerId: proto.Uint64(mjcontext.GetLastChupaiPlayer()),
		Card:         roomCard,
		HuType:       room.HuType_DianPao.Enum(),
	}
	toClientMessage := interfaces.ToClientMessage{
		MsgID: int(msgid.MsgID_ROOM_HU_NTF),
		Msg:   huNtf,
	}
	playerIDs := []uint64{}
	for _, p := range mjcontext.GetPlayers() {
		playerIDs = append(playerIDs, p.GetPalyerId())
	}
	flow.PushMessages(playerIDs, toClientMessage)
	//通知完成后，将进入摸牌状态，这里将重置mopaiPlayer
	//TODO:先暂时取胡牌玩家列表中的最后一个玩家
	var success bool
	lastPlayer.OutCards, success = utils.RemoveCards(lastPlayer.OutCards, mjcontext.GetLastOutCard(), 1)
	if !success {
		logrus.Errorln("移除outCard失败")
	}
	nextMopaiPlayer := utils.GetNextPlayerByID(mjcontext.GetPlayers(), mjcontext.LastHuPlayers[len(mjcontext.LastHuPlayers)-1])
	mjcontext.MopaiPlayer = nextMopaiPlayer.GetPalyerId()
}

// OnExit 退出状态
func (s *HuState) OnExit(flow interfaces.MajongFlow) {

}
