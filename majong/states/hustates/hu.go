package hustates

import (
	"steve/client_pb/msgId"
	"steve/client_pb/room"
	"steve/majong/global"
	"steve/majong/interfaces"
	"steve/majong/interfaces/facade"
	"steve/majong/utils"
	majongpb "steve/server_pb/majong"

	"github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
)

// HuState 胡状态
// 进入胡状态时， 执行胡操作。设置胡完成事件
// 收到胡完成事件时，设置摸牌玩家，返回摸牌状态
type HuState struct {
}

var _ interfaces.MajongState = new(HuState)

// ProcessEvent 处理事件
func (s *HuState) ProcessEvent(eventID majongpb.EventID, eventContext []byte, flow interfaces.MajongFlow) (newState majongpb.StateID, err error) {
	if eventID == majongpb.EventID_event_hu_finish {
		s.setMopaiPlayer(flow)
		return majongpb.StateID_state_mopai, nil
	}
	return majongpb.StateID_state_hu, global.ErrInvalidEvent
}

// OnEntry 进入状态
func (s *HuState) OnEntry(flow interfaces.MajongFlow) {
	s.doHu(flow)
	flow.SetAutoEvent(majongpb.AutoEvent{
		EventId:      majongpb.EventID_event_hu_finish,
		EventContext: nil,
	})
}

// OnExit 退出状态
func (s *HuState) OnExit(flow interfaces.MajongFlow) {

}

// addHuCard 添加胡的牌
func (s *HuState) addHuCard(card *majongpb.Card, player *majongpb.Player, srcPlayerID uint64) {
	addHuCard(card, player, srcPlayerID, majongpb.HuType_hu_dianpao)
}

// doHu 执行胡操作
func (s *HuState) doHu(flow interfaces.MajongFlow) {
	mjContext := flow.GetMajongContext()

	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "HuState.doHu",
	})
	logEntry = utils.WithMajongContext(logEntry, mjContext)
	players := mjContext.GetLastHuPlayers()

	for _, playerID := range players {
		player := utils.GetMajongPlayer(playerID, mjContext)
		card := mjContext.GetLastOutCard()
		s.addHuCard(card, player, playerID)
	}
	s.notifyHu(flow)
	return
}

// HuState 广播胡
func (s *HuState) notifyHu(flow interfaces.MajongFlow) {
	mjContext := flow.GetMajongContext()
	card := mjContext.GetLastOutCard()
	body := room.RoomHuNtf{
		Players:      mjContext.GetLastHuPlayers(),
		FromPlayerId: proto.Uint64(mjContext.GetLastChupaiPlayer()),
		Card:         proto.Uint32(uint32(utils.ServerCard2Number(card))),
		HuType:       room.HuType_DianPao.Enum(),
	}
	facade.BroadcaseMessage(flow, msgid.MsgID_ROOM_HU_NTF, &body)
}

// setMopaiPlayer 设置摸牌玩家
func (s *HuState) setMopaiPlayer(flow interfaces.MajongFlow) {
	mjContext := flow.GetMajongContext()
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "QiangganghuState.setMopaiPlayer",
	})
	logEntry = utils.WithMajongContext(logEntry, mjContext)
	huPlayers := mjContext.GetLastHuPlayers()
	srcPlayer := mjContext.GetLastChupaiPlayer()
	players := mjContext.GetPlayers()

	mjContext.MopaiPlayer = calcMopaiPlayer(logEntry, huPlayers, srcPlayer, players)
}
