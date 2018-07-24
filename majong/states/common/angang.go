package common

//适用麻将：四川血流
//前置条件：取麻将现场的杠的牌和最后杠玩家
//处理的事件请求：暗杠完成请求
//处理请求的过程：设置麻将现场的摸牌玩家
//处理请求的结果：返回摸牌状态ID
//状态退出行为：无
//状态进入行为：设置自动触发杠完成事件，处理暗杠逻辑，和广播通知客户端杠消息通知，该消息包含来自的玩家，去的玩家，杠的牌，还有杠类型。
//并设置杠玩家Properties["gang"]为[]byte("true")，最后进行暗杠结算
//约束条件：无
import (
	"steve/client_pb/msgid"
	"steve/client_pb/room"
	"steve/common/mjoption"
	majongpb "steve/entity/majong"
	"steve/majong/interfaces"
	"steve/majong/interfaces/facade"
	"steve/majong/utils"

	"github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
)

// TODO 结算

//AnGangState 暗杠状态
type AnGangState struct {
}

var _ interfaces.MajongState = new(AnGangState)

// ProcessEvent 处理事件
// 暗杠逻辑执行完后，进入暗杠结算状态，确认接收到暗杠结算完成请求，返回摸牌状态
func (s *AnGangState) ProcessEvent(eventID majongpb.EventID, eventContext interface{}, flow interfaces.MajongFlow) (newState majongpb.StateID, err error) {
	if eventID == majongpb.EventID_event_angang_finish {
		s.setMopaiPlayer(flow)
		xpOption := mjoption.GetXingpaiOption(int(flow.GetMajongContext().GetXingpaiOptionId()))
		if xpOption.EnableGangSettle {
			return majongpb.StateID(majongpb.StateID_state_gang_settle), nil
		}
		return majongpb.StateID_state_mopai, nil
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
	s.addGangCard(card, player, player.GetPalyerId())
	s.notifyPlayers(flow, card, player)
	return
}

// notifyPlayers 广播暗杠消息
func (s *AnGangState) notifyPlayers(flow interfaces.MajongFlow, card *majongpb.Card, player *majongpb.Player) {
	intCard := uint32(utils.ServerCard2Number(card))
	body := room.RoomGangNtf{
		ToPlayerId:   proto.Uint64(player.GetPalyerId()),
		FromPlayerId: proto.Uint64(player.GetPalyerId()),
		Card:         proto.Uint32(intCard),
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
	mjContext.MopaiType = majongpb.MopaiType_MT_GANG
}
