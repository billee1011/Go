package common

//适用麻将：四川血流
//前置条件：取麻将现场的杠的牌和最后杠玩家
//处理的事件请求：杠完成请求
//处理请求的过程：设置麻将现场的摸牌玩家
//处理请求的结果：返回摸牌状态ID
//状态退出行为：无
//状态进入行为：设置自动触发杠完成事件,处理杠逻辑，和广播通知客户端杠消息通知，该消息包含来自的玩家，去的玩家，杠的牌，杠类型。
//并设置杠玩家Properties["gang"]为[]byte("true")，最后进行明杠结算
//约束条件：无
import (
	msgid "steve/client_pb/msgid"
	"steve/client_pb/room"
	"steve/majong/global"
	"steve/majong/interfaces"
	"steve/majong/interfaces/facade"
	"steve/majong/utils"
	majongpb "steve/server_pb/majong"

	"github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
)

//MingGangState 明杠状态
type MingGangState struct {
}

var _ interfaces.MajongState = new(MingGangState)

// ProcessEvent 处理事件
func (s *MingGangState) ProcessEvent(eventID majongpb.EventID, eventContext []byte, flow interfaces.MajongFlow) (newState majongpb.StateID, err error) {
	if eventID == majongpb.EventID_event_gang_finish {
		s.setMopaiPlayer(flow)
		return majongpb.StateID(majongpb.StateID_state_gang_settle), nil
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
	mjContext := flow.GetMajongContext()
	playerID := mjContext.GetLastGangPlayer()
	player := utils.GetMajongPlayer(playerID, mjContext)
	card := mjContext.GetGangCard()
	srcPlayerID := mjContext.GetLastChupaiPlayer()
	srcPlayer := utils.GetMajongPlayer(srcPlayerID, mjContext)

	logEntry := logrus.WithFields(logrus.Fields{
		"func_name":      "MingGangState.doMinggang",
		"gang_player_id": playerID,
	})
	logEntry = utils.WithMajongContext(logEntry, mjContext)

	// 从被杠玩家的outCards移除被杠牌
	srcOutCards := srcPlayer.GetOutCards()
	srcPlayer.OutCards = removeLastCard(logEntry, srcOutCards, card)
	// 从杠牌玩家的handCards移除杠牌
	newCards, ok := utils.RemoveCards(player.GetHandCards(), card, 3)
	if !ok {
		logEntry.Errorln("移除玩家手牌失败")
		return
	}
	player.HandCards = newCards

	s.addGangCard(card, player, mjContext.GetLastChupaiPlayer())
	s.notifyPlayers(flow, card, player, mjContext.GetLastChupaiPlayer())
	return
}

// notifyPlayers 广播暗杠消息
func (s *MingGangState) notifyPlayers(flow interfaces.MajongFlow, card *majongpb.Card, player *majongpb.Player, srcPlayerID uint64) {
	intCard := uint32(utils.ServerCard2Number(card))
	body := room.RoomGangNtf{
		ToPlayerId:   proto.Uint64(player.GetPalyerId()),
		FromPlayerId: proto.Uint64(srcPlayerID),
		Card:         proto.Uint32(intCard),
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
	mjContext.MopaiType = majongpb.MopaiType_MT_GANG
}

//	doMingGangSettle 明杠结算
func (s *MingGangState) doMingGangSettle(mjContext *majongpb.MajongContext, player *majongpb.Player, srcPlayerID uint64) {
	allPlayers := make([]uint64, 0)
	for _, player := range mjContext.Players {
		allPlayers = append(allPlayers, player.GetPalyerId())
	}
	param := interfaces.GangSettleParams{
		GangPlayer: player.GetPalyerId(),
		SrcPlayer:  srcPlayerID,
		AllPlayers: allPlayers,
		GangType:   majongpb.GangType_gang_minggang,
		SettleID:   mjContext.CurrentSettleId,
	}

	f := global.GetGameSettlerFactory()
	gameID := int(mjContext.GetGameId())
	settleInfo := facade.SettleGang(f, gameID, param)
	if settleInfo != nil {
		mjContext.SettleInfos = append(mjContext.SettleInfos, settleInfo)
		mjContext.CurrentSettleId++
	}
}
