package common

//适用麻将：四川血流
//前置条件：取麻将现场的最后打出的牌，和最后出牌的玩家，和最后碰的玩家
//处理的事件请求：出牌请求
//处理请求的过程：验证出牌是否合法，设置麻将牌局现场最后出的牌和最后出牌玩家，还有清空出牌玩家的可能动作
//处理请求的结果：验证通过返回出牌状态ID，否则还是碰状态
//状态退出行为：无
//状态进入行为：处理碰逻辑，并广播通知客户端碰牌消息通知，该消息包含出的牌和来自的玩家，去的玩家
//约束条件：无
import (
	"steve/client_pb/msgid"
	"steve/client_pb/room"
	majongpb "steve/entity/majong"
	"steve/majong/interfaces"
	"steve/majong/interfaces/facade"
	"steve/majong/utils"

	"github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
)

// PengState 碰状态
type PengState struct {
}

var _ interfaces.MajongState = new(PengState)

// ProcessEvent 处理事件
// 碰牌成功后，接受到出牌请求，处理出牌请求，处理完成，进入出牌状态
func (s *PengState) ProcessEvent(eventID majongpb.EventID, eventContext interface{}, flow interfaces.MajongFlow) (newState majongpb.StateID, err error) {
	if eventID == majongpb.EventID_event_peng_finish {
		mjContext := flow.GetMajongContext()
		mjContext.ZixunType = majongpb.ZixunType_ZXT_PENG
		return majongpb.StateID(majongpb.StateID_state_zixun), nil
	}
	return majongpb.StateID_state_peng, nil
}

// OnEntry 进入状态	"steve/majong/interfaces/facade"
func (s *PengState) OnEntry(flow interfaces.MajongFlow) {
	s.doPeng(flow)
	flow.SetAutoEvent(majongpb.AutoEvent{
		EventId:      majongpb.EventID_event_peng_finish,
		EventContext: nil,
	})
}

// OnExit 退出状态
func (s *PengState) OnExit(flow interfaces.MajongFlow) {

}

func (s *PengState) notifyPeng(flow interfaces.MajongFlow, card *majongpb.Card, from uint64, to uint64) {
	facade.BroadcaseMessage(flow, msgid.MsgID_ROOM_PENG_NTF, &room.RoomPengNtf{
		Card:         proto.Uint32(utils.ServerCard2Uint32(card)),
		FromPlayerId: proto.Uint64(from),
		ToPlayerId:   proto.Uint64(to),
	})
}

// doPeng 执行碰操作
func (s *PengState) doPeng(flow interfaces.MajongFlow) {
	mjContext := flow.GetMajongContext()
	card := mjContext.GetLastOutCard()
	pengPlayerID := mjContext.GetLastPengPlayer()
	pengPlayer := utils.GetMajongPlayer(pengPlayerID, mjContext)
	srcPlayerID := mjContext.GetLastChupaiPlayer()
	srcPlayer := utils.GetMajongPlayer(srcPlayerID, mjContext)

	logEntry := logrus.WithFields(logrus.Fields{
		"func_name":      "PengState.doPeng",
		"peng_player_id": pengPlayer,
	})
	logEntry = utils.WithMajongContext(logEntry, mjContext)

	// 从被碰玩家的outCards移除被碰牌
	srcOutCards := srcPlayer.GetOutCards()
	srcPlayer.OutCards = removeLastCard(logEntry, srcOutCards, card)
	// 从碰牌玩家的handCards移除碰牌
	logEntry = logEntry.WithFields(logrus.Fields{})
	newCards, ok := utils.RemoveCards(pengPlayer.GetHandCards(), card, 2)
	if !ok {
		logEntry.Errorln("移除玩家手牌失败")
		return
	}
	pengPlayer.HandCards = newCards
	s.notifyPeng(flow, card, srcPlayerID, pengPlayerID)
	s.addPengCard(card, pengPlayer, srcPlayerID)
	return
}

// addPengCard 添加碰的牌
func (s *PengState) addPengCard(card *majongpb.Card, player *majongpb.Player, srcPlayerID uint64) {
	player.PengCards = append(player.GetPengCards(), &majongpb.PengCard{
		Card:      card,
		SrcPlayer: srcPlayerID,
	})
}

// TODO:  通知碰
