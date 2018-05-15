package gangstates

import (
	msgid "steve/client_pb/msgId"
	"steve/client_pb/room"
	"steve/majong/interfaces"
	"steve/majong/interfaces/facade"
	"steve/majong/settle"
	"steve/majong/utils"
	majongpb "steve/server_pb/majong"

	"github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
)

// TODO 结算

//AnGangState 暗杠状态 @Author:wuhongwei
type AnGangState struct {
}

var _ interfaces.MajongState = new(AnGangState)

// ProcessEvent 处理事件
// 暗杠逻辑执行完后，进入暗杠状态，确认接收到暗杠完成请求，返回摸牌状态
func (s *AnGangState) ProcessEvent(eventID majongpb.EventID, eventContext []byte, flow interfaces.MajongFlow) (newState majongpb.StateID, err error) {
	if eventID == majongpb.EventID_event_angang_finish {
		s.setMopaiPlayer(flow)
		return majongpb.StateID(majongpb.StateID_state_mopai), nil
	}
	return majongpb.StateID(majongpb.StateID_state_angang), nil
}

// OnEntry 进入状态
func (s *AnGangState) OnEntry(flow interfaces.MajongFlow) {
	s.doAngang(flow)
	flow.SetAutoEvent(majongpb.AutoEvent{
		EventId:      majongpb.EventID_event_angang_finish,
		EventContext: nil,
	})
}

// OnExit 退出状态
func (s *AnGangState) OnExit(flow interfaces.MajongFlow) {
}

// doAngang 执行暗杠操作
func (s *AnGangState) doAngang(flow interfaces.MajongFlow) {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "AnGangState.doAngang",
	})

	mjContext := flow.GetMajongContext()
	logEntry = utils.WithMajongContext(logEntry, mjContext)

	playerID := mjContext.GetLastGangPlayer()
	player := utils.GetMajongPlayer(playerID, mjContext)
	card := mjContext.GetGangCard()

	logEntry = logEntry.WithFields(logrus.Fields{
		"gang_player_id": playerID,
	})

	newCards, ok := utils.RemoveCards(player.GetHandCards(), card, 4)
	if !ok {
		logEntry.Errorln("移除玩家手牌失败")
		return
	}
	player.HandCards = newCards
	player.Properties["gang"] = []byte("true")
	s.addGangCard(card, player, player.GetPalyerId())
	s.notifyPlayers(flow, card, player)
	s.doAnGangSettle(mjContext, player)
	return
}

// notifyPlayers 广播暗杠消息
func (s *AnGangState) notifyPlayers(flow interfaces.MajongFlow, card *majongpb.Card, player *majongpb.Player) {
	roomCard, _ := utils.CardToRoomCard(card)
	body := room.RoomGangNtf{
		ToPlayerId:   proto.Uint64(player.GetPalyerId()),
		FromPlayerId: proto.Uint64(player.GetPalyerId()),
		Card:         roomCard,
		GangType:     room.GangType_AnGang.Enum(),
	}
	facade.BroadcaseMessage(flow, msgid.MsgID_ROOM_GANG_NTF, &body)
}

// addGangCard 添加暗杠的牌
func (s *AnGangState) addGangCard(card *majongpb.Card, player *majongpb.Player, srcPlayerID uint64) {
	player.GangCards = append(player.GetGangCards(), &majongpb.GangCard{
		Card:      card,
		Type:      majongpb.GangType_gang_angang,
		SrcPlayer: srcPlayerID,
	})
}

// setMopaiPlayer 设置摸牌玩家
func (s *AnGangState) setMopaiPlayer(flow interfaces.MajongFlow) {
	mjContext := flow.GetMajongContext()
	mjContext.MopaiPlayer = mjContext.GetLastGangPlayer()
}

//	doAnGangSettle 暗杠结算
func (s *AnGangState) doAnGangSettle(mjContext *majongpb.MajongContext, player *majongpb.Player) {
	allPlayers := make([]uint64, 0)
	for _, player := range mjContext.Players {
		allPlayers = append(allPlayers, player.GetPalyerId())
	}
	param := interfaces.GangSettleParams{
		GangPlayer: player.GetPalyerId(),
		SrcPlayer:  player.GetPalyerId(),
		AllPlayers: allPlayers,
		GangType:   majongpb.GangType_gang_angang,
		SettleID:   mjContext.CurrentSettleId,
	}

	anGangSettle := new(settle.GangSettle)
	settleInfo := anGangSettle.Settle(param)
	mjContext.SettleInfos = append(mjContext.SettleInfos, settleInfo)
	mjContext.CurrentSettleId++
}
