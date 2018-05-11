package gangstates

import (
	msgid "steve/client_pb/msgId"
	"steve/client_pb/room"
	"steve/majong/interfaces"
	"steve/majong/interfaces/facade"
	"steve/majong/utils"
	majongpb "steve/server_pb/majong"

	"github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
)

// TODO 结算

//MingGangState 明杠状态 @Author:wuhongwei
type MingGangState struct {
}

var _ interfaces.MajongState = new(MingGangState)

// ProcessEvent 处理事件
func (s *MingGangState) ProcessEvent(eventID majongpb.EventID, eventContext []byte, flow interfaces.MajongFlow) (newState majongpb.StateID, err error) {
	if eventID == majongpb.EventID_event_gang_finish {
		s.setMopaiPlayer(flow)
		return majongpb.StateID(majongpb.StateID_state_mopai), nil
	}
	return majongpb.StateID(majongpb.StateID_state_gang), nil
}

// OnEntry 进入状态
func (s *MingGangState) OnEntry(flow interfaces.MajongFlow) {
	s.doMinggang(flow)
	flow.SetAutoEvent(majongpb.AutoEvent{
		EventId:      majongpb.EventID_event_gang_finish,
		EventContext: nil,
	})
}

// OnExit 退出状态
func (s *MingGangState) OnExit(flow interfaces.MajongFlow) {
}

// doMinggang 执行明杠操作
func (s *MingGangState) doMinggang(flow interfaces.MajongFlow) {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "MingGangState.doMinggang",
	})

	mjContext := flow.GetMajongContext()
	logEntry = utils.WithMajongContext(logEntry, mjContext)

	playerID := mjContext.GetLastGangPlayer()
	player := utils.GetMajongPlayer(playerID, mjContext)
	card := mjContext.GetGangCard()

	logEntry = logEntry.WithFields(logrus.Fields{
		"gang_player_id": playerID,
	})

	newCards, ok := utils.RemoveCards(player.GetHandCards(), card, 3)
	if !ok {
		logEntry.Errorln("移除玩家手牌失败")
		return
	}
	player.HandCards = newCards

	s.addGangCard(card, player, player.GetPalyerId())
	s.notifyPlayers(flow, card, player, mjContext.GetLastChupaiPlayer())
	return
}

// notifyPlayers 广播暗杠消息
func (s *MingGangState) notifyPlayers(flow interfaces.MajongFlow, card *majongpb.Card, player *majongpb.Player, srcPlayerID uint64) {
	roomCard, _ := utils.CardToRoomCard(card)
	body := room.RoomGangNtf{
		ToPlayerId:   proto.Uint64(player.GetPalyerId()),
		FromPlayerId: proto.Uint64(srcPlayerID),
		Card:         roomCard,
		GangType:     room.GangType_MingGang.Enum(),
	}
	facade.BroadcaseMessage(flow, msgid.MsgID_ROOM_GANG_NTF, &body)
}

// addGangCard 添加明杠的牌
func (s *MingGangState) addGangCard(card *majongpb.Card, player *majongpb.Player, srcPlayerID uint64) {
	player.GangCards = append(player.GetGangCards(), &majongpb.GangCard{
		Card:      card,
		Type:      majongpb.GangType_gang_minggang,
		SrcPlayer: srcPlayerID,
	})
}

// setMopaiPlayer 设置摸牌玩家
func (s *MingGangState) setMopaiPlayer(flow interfaces.MajongFlow) {
	mjContext := flow.GetMajongContext()
	mjContext.MopaiPlayer = mjContext.GetLastGangPlayer()
}
