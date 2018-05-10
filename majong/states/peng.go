package states

import (
	"fmt"
	"steve/majong/interfaces"
	"steve/majong/utils"
	majongpb "steve/server_pb/majong"

	"steve/client_pb/room"

	"steve/client_pb/msgId"

	"github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
)

// PengState 碰状态 @Author:wuhongwei
type PengState struct {
}

var _ interfaces.MajongState = new(PengState)

// ProcessEvent 处理事件
// 碰牌成功后，接受到出牌请求，处理出牌请求，处理完成，进入出牌状态
func (s *PengState) ProcessEvent(eventID majongpb.EventID, eventContext []byte, flow interfaces.MajongFlow) (newState majongpb.StateID, err error) {
	if eventID == majongpb.EventID_event_chupai_request {
		if err := s.chupai(eventContext, flow); err != nil {
			return majongpb.StateID(majongpb.StateID_state_peng), err
		}
		return majongpb.StateID(majongpb.StateID_state_chupai), nil
	}
	return majongpb.StateID_state_peng, nil
}

//Chupai 出牌操作 @Author:wuhongwei
func (s *PengState) chupai(eventContext []byte, flow interfaces.MajongFlow) error {
	// 序列化
	chupaiEvent := new(majongpb.ChupaiRequestEvent)
	if err := proto.Unmarshal(eventContext, chupaiEvent); err != nil {
		return fmt.Errorf("出牌 : %v", errUnmarshalEvent)
	}
	//麻将牌局现场
	mjContext := flow.GetMajongContext()
	// 获取出牌玩家和出的牌
	playerID := chupaiEvent.Head.PlayerId
	outCardPlayer := utils.GetPlayerByID(mjContext.Players, playerID)
	if outCardPlayer == nil {
		return fmt.Errorf("出牌事件失败-出牌玩家ID不存在: %v", playerID)
	}
	// 是否是当前玩家
	if mjContext.GetActivePlayer() != playerID {
		return fmt.Errorf("出牌事件失败-player不为当前行动玩家")
	}
	// 获取出的牌
	outCard := chupaiEvent.Cards
	// 判断出的牌是否存在
	isExist := false
	handCard := outCardPlayer.GetHandCards()
	for _,card := range handCard {
		if utils.CardEqual(card,outCard) {
			isExist = true
			 break
		}
	}
	if !isExist {
		return fmt.Errorf("出牌事件失败-请求出的牌不存在：%v" ,outCard)
	}

	// 删除手中要出牌
	handCards, flag := utils.DeleteCardFromLast(handCard, outCard)
	if !flag {
		return fmt.Errorf("出牌事件-删除牌失败: %v", outCard)
	}
	
	// 修改玩家手牌
	outCardPlayer.HandCards = handCards
	// 将出的牌添加到玩家出牌数组中
	outCardPlayer.OutCards = append(outCardPlayer.OutCards, outCard)
	// 修改麻将牌局现场最后出的牌
	mjContext.LastOutCard = outCard
	// 设置最后出牌玩家ID
	mjContext.LastChupaiPlayer = outCardPlayer.PalyerId
	// 清空玩家可能动作
	outCardPlayer.PossibleActions = outCardPlayer.PossibleActions[:0]

	// 出牌广播通知
	playerIDs := make([]uint64, 0, 0)
	for _, player := range mjContext.Players {
		playerIDs = append(playerIDs, player.GetPalyerId())
	}
	roomCard, err := utils.CardToRoomCard(outCard)
	if err != nil {
		return fmt.Errorf("出牌事件-outCard to roomCard - 失败: %v", outCard)
	}
	toClient := interfaces.ToClientMessage{
		MsgID: int(msgid.MsgID_ROOM_CHUPAI_NTF),
		Msg: &room.RoomChupaiNtf{
			Player: proto.Uint64(outCardPlayer.PalyerId),
			Card:   roomCard,
		},
	}
	flow.PushMessages(playerIDs, toClient)

	// 出过非定缺颜色的牌 TODO
	if len(outCardPlayer.Properties[utils.IsOutNoDingQueColorCard]) == 0 && outCard.Color != outCardPlayer.DingqueColor {
		outCardPlayer.Properties[utils.IsOutNoDingQueColorCard] = []byte{1}
	}

	// 日志
	logrus.WithFields(logrus.Fields{
		"outPlayer_id": outCardPlayer.GetPalyerId(),
		"outCard":      outCard,
		"chupaiEvent":  chupaiEvent,
		"toClient":     toClient,
		"HandCards":    outCardPlayer.GetHandCards(),
		"OutCards":     outCardPlayer.GetOutCards(),
	}).Info("出牌成功")
	return nil
}

// OnEntry 进入状态
func (s *PengState) OnEntry(flow interfaces.MajongFlow) {
	s.doPeng(flow)
}

// OnExit 退出状态
func (s *PengState) OnExit(flow interfaces.MajongFlow) {

}

// doPeng 执行碰操作
func (s *PengState) doPeng(flow interfaces.MajongFlow) {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "PengState.doPeng",
	})

	mjContext := flow.GetMajongContext()
	logEntry = utils.WithMajongContext(logEntry, mjContext)

	pengPlayer := mjContext.GetLastPengPlayer()

	player := utils.GetMajongPlayer(pengPlayer, mjContext)

	card := mjContext.GetLastOutCard()
	logEntry = logEntry.WithFields(logrus.Fields{
		"peng_player_id": pengPlayer,
	})

	newCards, ok := utils.RemoveCards(player.GetHandCards(), card, 2)
	if !ok {
		logEntry.Errorln("移除玩家手牌失败")
		return
	}
	player.HandCards = newCards

	s.addPengCard(card, player, mjContext.GetActivePlayer())
	return
}

// addPengCard 添加碰的牌
func (s *PengState) addPengCard(card *majongpb.Card, player *majongpb.Player, srcPlayerID uint64) {
	player.PengCards = append(player.GetPengCards(), &majongpb.PengCard{
		Card:      card,
		SrcPlayer: srcPlayerID,
	})
}
