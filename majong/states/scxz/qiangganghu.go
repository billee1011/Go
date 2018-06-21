package scxz

import (
	msgid "steve/client_pb/msgId"
	"steve/client_pb/room"
	"steve/majong/global"
	"steve/majong/interfaces"
	"steve/majong/interfaces/facade"
	"steve/majong/states/common"
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
		return majongpb.StateID_state_qiangganghu_settle, nil
	}
	return majongpb.StateID_state_qiangganghu, global.ErrInvalidEvent
}

// OnEntry 进入状态
func (s *QiangganghuState) OnEntry(flow interfaces.MajongFlow) {
	s.doHu(flow)
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

	mjContext.MopaiPlayer = common.CalcMopaiPlayer(logEntry, huPlayers, srcPlayer, players)
	mjContext.MopaiType = majongpb.MopaiType_MT_NORMAL
}

// addHuCard 添加胡的牌
func (s *QiangganghuState) addHuCard(card *majongpb.Card, player *majongpb.Player, srcPlayerID uint64, isReal bool) {
	common.AddHuCard(card, player, srcPlayerID, majongpb.HuType_hu_dianpao, isReal)
}

func (s *QiangganghuState) removeSrcCard(card *majongpb.Card, srcPlayer *majongpb.Player) {
	var succ bool
	srcPlayer.HandCards, succ = utils.RemoveCards(srcPlayer.GetHandCards(), card, 1)
	if !succ {
		logrus.WithFields(logrus.Fields{
			"func_name":      "QiangganghuState.removeSrcCard",
			"hand_cards":     srcPlayer.GetHandCards(),
			"gang_player_id": srcPlayer.GetPalyerId(),
		}).Errorln("移除杠者的杠牌失败")
	}
}

// doHu 执行胡操作
func (s *QiangganghuState) doHu(flow interfaces.MajongFlow) {
	mjContext := flow.GetMajongContext()

	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "QiangganghuState.doHu",
	})
	logEntry = utils.WithMajongContext(logEntry, mjContext)
	players := mjContext.GetLastHuPlayers()
	srcPlayerID := mjContext.GetLastGangPlayer()
	srcPlayer := utils.GetPlayerByID(mjContext.GetPlayers(), srcPlayerID)
	card := mjContext.GetGangCard() // 杠的牌为抢杠胡的牌

	isReal := true
	for _, playerID := range players {
		player := utils.GetMajongPlayer(playerID, mjContext)
		s.addHuCard(card, player, srcPlayerID, isReal)
		isReal = false
	}
	s.removeSrcCard(card, srcPlayer)
	s.notifyHu(flow)
	return
}

// QiangganghuState 广播胡
func (s *QiangganghuState) notifyHu(flow interfaces.MajongFlow) {
	mjContext := flow.GetMajongContext()
	card := mjContext.GetGangCard()
	body := room.RoomHuNtf{
		Players:      mjContext.GetLastHuPlayers(),
		FromPlayerId: proto.Uint64(mjContext.GetLastMopaiPlayer()),
		Card:         proto.Uint32(uint32(utils.ServerCard2Number(card))),
		HuType:       room.HuType_HT_QIANGGANGHU.Enum(),
	}
	facade.BroadcaseMessage(flow, msgid.MsgID_ROOM_HU_NTF, &body)
}
