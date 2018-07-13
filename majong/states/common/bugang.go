package common

//适用麻将：四川血流
//前置条件：取麻将现场的杠的牌和最后杠玩家
//处理的事件请求：补杠完成请求
//处理请求的过程：设置麻将现场的摸牌玩家
//处理请求的结果：返回摸牌状态ID
//状态退出行为：无
//状态进入行为：设置自动触发杠完成事件,处理补杠逻辑，和广播通知客户端杠消息通知，该消息包含来自的玩家，去的玩家，杠的牌，杠类型。
//并设置杠玩家Properties["gang"]为[]byte("true")，最后进行补杠结算
//约束条件：无
import (
	msgid "steve/client_pb/msgid"
	"steve/client_pb/room"
	"steve/common/mjoption"
	"steve/majong/interfaces"
	"steve/majong/interfaces/facade"
	"steve/majong/utils"
	majongpb "steve/server_pb/majong"

	"github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
)

// TODO 结算

//BuGangState 补杠状态
type BuGangState struct {
}

var _ interfaces.MajongState = new(BuGangState)

// ProcessEvent 处理事件
func (s *BuGangState) ProcessEvent(eventID majongpb.EventID, eventContext []byte, flow interfaces.MajongFlow) (newState majongpb.StateID, err error) {
	if eventID == majongpb.EventID_event_bugang_finish {
		s.setMopaiPlayer(flow)
		xpOption := mjoption.GetXingpaiOption(int(flow.GetMajongContext().GetXingpaiOptionId()))
		if xpOption.EnableGangSettle {
			return majongpb.StateID_state_gang_settle, nil
		}
		return majongpb.StateID_state_mopai, nil
	}
	return majongpb.StateID_state_bugang, nil
}

// OnEntry 进入状态
func (s *BuGangState) OnEntry(flow interfaces.MajongFlow) {
	s.doBugang(flow)
	flow.SetAutoEvent(majongpb.AutoEvent{
		EventId:      majongpb.EventID_event_bugang_finish,
		EventContext: nil,
	})
}

// OnExit 退出状态
func (s *BuGangState) OnExit(flow interfaces.MajongFlow) {

}

// doBugang 执行补杠操作
func (s *BuGangState) doBugang(flow interfaces.MajongFlow) {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "BuGangState.doBugang",
	})

	mjContext := flow.GetMajongContext()
	logEntry = utils.WithMajongContext(logEntry, mjContext)

	playerID := mjContext.GetLastGangPlayer()
	player := utils.GetMajongPlayer(playerID, mjContext)
	card := mjContext.GetGangCard()

	logEntry = logEntry.WithFields(logrus.Fields{
		"gang_player_id": playerID,
	})

	newCards, ok := utils.RemoveCards(player.GetHandCards(), card, 1)
	if !ok {
		logEntry.Errorln("移除玩家手牌失败")
		return
	}
	player.HandCards = newCards

	s.removePengCard(card, player)
	s.addGangCard(card, player, player.GetPalyerId())
	s.notifyPlayers(flow, card, player)
	return
}

// notifyPlayers 广播暗杠消息
func (s *BuGangState) notifyPlayers(flow interfaces.MajongFlow, card *majongpb.Card, player *majongpb.Player) {
	intCard := uint32(utils.ServerCard2Number(card))
	body := room.RoomGangNtf{
		ToPlayerId:   proto.Uint64(player.GetPalyerId()),
		FromPlayerId: proto.Uint64(player.GetPalyerId()),
		Card:         proto.Uint32(intCard),
		GangType:     room.GangType_BuGang.Enum(),
	}
	facade.BroadcaseMessage(flow, msgid.MsgID_ROOM_GANG_NTF, &body)
}

// addGangCard 添加补杠的牌
func (s *BuGangState) addGangCard(card *majongpb.Card, player *majongpb.Player, srcPlayerID uint64) {
	player.GangCards = append(player.GetGangCards(), &majongpb.GangCard{
		Card:      card,
		Type:      majongpb.GangType_gang_bugang,
		SrcPlayer: srcPlayerID,
	})
}

// removePengCard 移除碰的牌
func (s *BuGangState) removePengCard(card *majongpb.Card, player *majongpb.Player) {
	newPengCards := []*majongpb.PengCard{}
	pengCards := player.GetPengCards()
	for index, pengCard := range pengCards {
		if utils.CardEqual(card, pengCard.GetCard()) {
			newPengCards = append(newPengCards, pengCards[index+1:]...)
			break
		}
		newPengCards = append(newPengCards, pengCard)
	}
	player.PengCards = newPengCards
}

// setMopaiPlayer 设置摸牌玩家
func (s *BuGangState) setMopaiPlayer(flow interfaces.MajongFlow) {
	mjContext := flow.GetMajongContext()
	mjContext.MopaiPlayer = mjContext.GetLastGangPlayer()
	mjContext.MopaiType = majongpb.MopaiType_MT_GANG
}
