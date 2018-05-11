package hustates

import (
	"errors"
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

// ZimoState 自摸状态
// 进入状态时，执行自摸动作，并广播给玩家
// 自摸完成事件，进入下家摸牌状态
type ZimoState struct {
}

var _ interfaces.MajongState = new(ZimoState)

// ProcessEvent 处理事件
func (s *ZimoState) ProcessEvent(eventID majongpb.EventID, eventContext []byte, flow interfaces.MajongFlow) (newState majongpb.StateID, err error) {
	if eventID == majongpb.EventID_event_zimo_finish {
		s.setMopaiPlayer(flow)
		return majongpb.StateID_state_mopai, nil
	}
	return majongpb.StateID_state_zimo, global.ErrInvalidEvent
}

// OnEntry 进入状态
func (s *ZimoState) OnEntry(flow interfaces.MajongFlow) {
	s.doZimo(flow)
	flow.SetAutoEvent(majongpb.AutoEvent{
		EventId:      majongpb.EventID_event_zimo_finish,
		EventContext: nil,
	})
}

// OnExit 退出状态
func (s *ZimoState) OnExit(flow interfaces.MajongFlow) {

}

// doZimo 执行自摸操作
func (s *ZimoState) doZimo(flow interfaces.MajongFlow) {
	mjContext := flow.GetMajongContext()

	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "ZimoState.doZimo",
	})
	logEntry = utils.WithMajongContext(logEntry, mjContext)
	player, card, err := s.getZimoInfo(mjContext)
	if err != nil {
		logEntry.Errorln(err)
		return
	}
	mjContext.LastHuPlayers = []uint64{player.GetPalyerId()}

	player.HandCards, _ = utils.RemoveCards(player.GetHandCards(), card, 1)
	addHuCard(card, player, player.GetPalyerId(), majongpb.HuType_hu_zimo)
	s.notifyHu(card, player.GetPalyerId(), flow)
}

// notifyHu 广播胡
func (s *ZimoState) notifyHu(card *majongpb.Card, playerID uint64, flow interfaces.MajongFlow) {
	huCard, _ := utils.CardToRoomCard(card)
	body := room.RoomHuNtf{
		Players:      []uint64{playerID},
		FromPlayerId: proto.Uint64(playerID),
		Card:         huCard,
		HuType:       room.HuType_ZiMo.Enum(),
	}
	facade.BroadcaseMessage(flow, msgid.MsgID_ROOM_HU_NTF, &body)
}

// setMopaiPlayer 设置摸牌玩家
func (s *ZimoState) setMopaiPlayer(flow interfaces.MajongFlow) {
	mjContext := flow.GetMajongContext()
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "ZimoState.doZimo",
	})
	logEntry = utils.WithMajongContext(logEntry, mjContext)

	huPlayers := mjContext.GetLastHuPlayers()
	if len(huPlayers) == 0 {
		logEntry.Errorln("胡牌玩家列表为空")
		return
	}
	players := mjContext.GetPlayers()
	mjContext.MopaiPlayer = calcMopaiPlayer(logEntry, huPlayers, huPlayers[0], players)
}

// getZimoInfo 获取自摸信息
func (s *ZimoState) getZimoInfo(mjContext *majongpb.MajongContext) (player *majongpb.Player, card *majongpb.Card, err error) {
	playerID := mjContext.GetLastMopaiPlayer()
	// 没有上个摸牌的玩家，是为天胡， 取庄家作为胡牌玩家
	if playerID == 0 {
		players := mjContext.GetPlayers()
		zjIndex := mjContext.GetZhuangjiaIndex()
		if zjIndex >= uint32(len(players)) {
			err = errors.New("庄家索引越界")
			return
		}
		player = players[zjIndex]
		card = s.calcTianhuCard(player.GetHandCards())
	} else {
		player = utils.GetPlayerByID(mjContext.GetPlayers(), playerID)
		card = mjContext.GetLastMopaiCard()
	}
	mjContext.LastHuPlayers = []uint64{playerID}
	return
}

// calcTianhuCard 计算天胡胡的牌
func (s *ZimoState) calcTianhuCard(cards []*majongpb.Card) *majongpb.Card {
	// TODO
	return cards[len(cards)-1]
}
