package hustates

import (
	msgid "steve/client_pb/msgId"
	"steve/client_pb/room"
	"steve/majong/global"
	"steve/majong/interfaces"
	"steve/majong/interfaces/facade"
	"steve/majong/utils"
	majongpb "steve/server_pb/majong"

	"github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
)

// QiangganghuState 抢杠胡状态
// 执行抢杠胡操作，并广播
// 从上个摸牌的玩家算起，最后胡的玩家的下家摸牌
type QiangganghuState struct {
}

var _ interfaces.MajongState = new(QiangganghuState)

// ProcessEvent 处理事件
func (s *QiangganghuState) ProcessEvent(eventID majongpb.EventID, eventContext []byte, flow interfaces.MajongFlow) (newState majongpb.StateID, err error) {
	if eventID == majongpb.EventID_event_qiangganghu_finish {
		s.setMopaiPlayer(flow)
		return majongpb.StateID_state_mopai, nil
	}
	return majongpb.StateID_state_qiangganghu, global.ErrInvalidEvent
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

// setMopaiPlayer 设置摸牌玩家
func (s *QiangganghuState) setMopaiPlayer(flow interfaces.MajongFlow) {
	mjContext := flow.GetMajongContext()
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "QiangganghuState.setMopaiPlayer",
	})
	logEntry = utils.WithMajongContext(logEntry, mjContext)
	huPlayers := mjContext.GetLastHuPlayers()
	srcPlayer := mjContext.GetLastMopaiPlayer()
	players := mjContext.GetPlayers()

	mjContext.MopaiPlayer = calcMopaiPlayer(logEntry, huPlayers, srcPlayer, players)
}

// addHuCard 添加胡的牌
func (s *QiangganghuState) addHuCard(card *majongpb.Card, player *majongpb.Player, srcPlayerID uint64) {
	addHuCard(card, player, srcPlayerID, majongpb.HuType_hu_dianpao)
}

// doHu 执行胡操作
func (s *QiangganghuState) doHu(flow interfaces.MajongFlow) {
	mjContext := flow.GetMajongContext()

	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "QiangganghuState.doHu",
	})
	logEntry = utils.WithMajongContext(logEntry, mjContext)
	players := mjContext.GetLastHuPlayers()

	for _, playerID := range players {
		player := utils.GetMajongPlayer(playerID, mjContext)
		card := mjContext.GetGangCard() // 杠的牌为抢杠胡的牌
		s.addHuCard(card, player, playerID)
	}
	s.notifyHu(flow)
	return
}

// QiangganghuState 广播胡
func (s *QiangganghuState) notifyHu(flow interfaces.MajongFlow) {
	mjContext := flow.GetMajongContext()
	huCard, _ := utils.CardToRoomCard(mjContext.GetLastOutCard())
	body := room.RoomHuNtf{
		Players:      mjContext.GetLastHuPlayers(),
		FromPlayerId: proto.Uint64(mjContext.GetLastMopaiPlayer()),
		Card:         huCard,
		HuType:       room.HuType_QiangGangHu.Enum(),
	}
	facade.BroadcaseMessage(flow, msgid.MsgID_ROOM_HU_NTF, &body)
}
