package states

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

// ChupaiState 初始化状态
type ChupaiState struct {
}

var _ interfaces.MajongState = new(ChupaiState)

// ProcessEvent 处理事件
func (s *ChupaiState) ProcessEvent(eventID majongpb.EventID, eventContext []byte, flow interfaces.MajongFlow) (newState majongpb.StateID, err error) {
	if eventID == majongpb.EventID_event_chupai_finish {
		s.chupai(flow)
		context := flow.GetMajongContext()
		players := context.GetPlayers()
		card := context.GetLastOutCard()
		var hasChupaiwenxun bool
		//出完牌后，将上轮添加的胡牌玩家列表重置
		context.LastHuPlayers = context.LastHuPlayers[:0]
		for _, player := range players {
			if context.GetLastChupaiPlayer() == player.GetPalyerId() {
				continue
			}
			ntf, need := checkActions(context, player, card)
			if need {
				flow.PushMessages([]uint64{player.GetPalyerId()}, interfaces.ToClientMessage{
					MsgID: int(msgid.MsgID_ROOM_CHUPAIWENXUN_NTF),
					Msg:   ntf,
				})
				hasChupaiwenxun = true
			}
		}
		if hasChupaiwenxun {
			return majongpb.StateID_state_chupaiwenxun, nil
		}
		player := utils.GetNextPlayerByID(context.GetPlayers(), context.GetLastChupaiPlayer())
		context.MopaiPlayer = player.GetPalyerId()
		context.MopaiType = majongpb.MopaiType_MT_NORMAL
		return majongpb.StateID_state_mopai, nil
	}
	return majongpb.StateID_state_init, global.ErrInvalidEvent
}

//checkActions 检查玩家可以有哪些操作
func checkActions(context *majongpb.MajongContext, player *majongpb.Player, card *majongpb.Card) (*room.RoomChupaiWenxunNtf, bool) {
	player.PossibleActions = player.PossibleActions[:0]

	chupaiWenxunNtf := &room.RoomChupaiWenxunNtf{}
	chupaiWenxunNtf.Card = proto.Uint32(uint32(utils.ServerCard2Number(card)))
	canMingGang := checkMingGang(context, player, card)
	chupaiWenxunNtf.EnableMinggang = proto.Bool(canMingGang)
	if canMingGang {
		player.PossibleActions = append(player.PossibleActions, majongpb.Action_action_gang)
	}
	canDianPao := checkDianPao(context, player, card)
	chupaiWenxunNtf.EnableDianpao = proto.Bool(canDianPao)
	if canDianPao {
		context.LastHuPlayers = append(context.LastHuPlayers, player.GetPalyerId())
		player.PossibleActions = append(player.PossibleActions, majongpb.Action_action_hu)
	}
	canPeng := checkPeng(context, player, card)
	chupaiWenxunNtf.EnablePeng = proto.Bool(canPeng)
	if canPeng {
		player.PossibleActions = append(player.PossibleActions, majongpb.Action_action_peng)
	}
	chupaiWenxunNtf.EnableQi = proto.Bool(true)
	return chupaiWenxunNtf, canDianPao || canMingGang || canPeng
}

//checkMingGang 查明杠
func checkMingGang(context *majongpb.MajongContext, player *majongpb.Player, card *majongpb.Card) bool {
	// 没有墙牌不能明杠
	if len(context.WallCards) == 0 {
		return false
	}
	outCard := context.GetLastOutCard()
	color := player.GetDingqueColor()
	//定缺牌不查
	if outCard.Color == color {
		return false
	}
	cards := player.HandCards
	num := 0
	for _, card := range cards {
		if utils.CardEqual(card, outCard) {
			num++
		}
	}
	if num == 3 {
		if len(player.GetHuCards()) > 0 {
			//创建副本，移除相应的杠牌进行查胡
			newcards := make([]*majongpb.Card, 0, len(cards))
			newcards = append(newcards, cards...)
			newcards, _ = utils.RemoveCards(newcards, outCard, num)
			newcardsI, _ := utils.CardsToInt(newcards)
			cardsI := utils.IntToUtilCard(newcardsI)
			laizi := make(map[utils.Card]bool)
			huCards := utils.FastCheckTingV2(cardsI, laizi)
			if utils.ContainHuCards(huCards, utils.HuCardsToUtilCards(player.HuCards)) {
				return true
			}
		} else {
			return true
		}
	}
	return false
}

//checkPeng 查碰
func checkPeng(context *majongpb.MajongContext, player *majongpb.Player, card *majongpb.Card) bool {
	color := player.GetDingqueColor()
	//胡牌后不能碰了
	if len(player.GetHuCards()) > 0 || card.Color == color {
		return false
	}
	num := 0
	for _, handCard := range player.GetHandCards() {
		if utils.CardEqual(handCard, card) {
			num++
		}
	}
	logrus.WithFields(logrus.Fields{
		"func_name":    "checkPeng",
		"check_player": player.GetPalyerId(),
		"check_card":   card,
		"hand_cards":   player.GetHandCards(),
		"count":        num,
	}).Debugln("检查是否可碰")
	return num >= 2
}

//checkDianPao 查点炮
func checkDianPao(context *majongpb.MajongContext, player *majongpb.Player, card *majongpb.Card) bool {
	cpCard := context.GetLastOutCard()
	color := player.GetDingqueColor()
	hasDingQueCard := utils.CheckHasDingQueCard(player.HandCards, color)
	if hasDingQueCard {
		return false
	}
	handCard := player.GetHandCards() // 当前点炮胡玩家手牌
	cardI, _ := utils.CardToInt(*cpCard)
	flag := utils.CheckHu(handCard, uint32(*cardI))
	if flag {
		return true
	}
	return false
}

//chupai 决策出牌
func (s *ChupaiState) chupai(flow interfaces.MajongFlow) {
	context := flow.GetMajongContext()
	activePlayer := utils.GetPlayerByID(context.GetPlayers(), context.GetLastChupaiPlayer())
	card := context.GetLastOutCard()
	activePlayer.HandCards, _ = utils.RemoveCards(activePlayer.HandCards, card, 1)
	activePlayer.OutCards = append(activePlayer.OutCards, card)
	facade.BroadcaseMessage(flow, msgid.MsgID_ROOM_CHUPAI_NTF, &room.RoomChupaiNtf{
		Player: proto.Uint64(activePlayer.GetPalyerId()),
		Card:   proto.Uint32(utils.ServerCard2Uint32(card)),
	})

}

// OnEntry 进入状态
func (s *ChupaiState) OnEntry(flow interfaces.MajongFlow) {
	flow.SetAutoEvent(majongpb.AutoEvent{
		EventId:      majongpb.EventID_event_chupai_finish,
		EventContext: nil,
	})
}

// OnExit 退出状态
func (s *ChupaiState) OnExit(flow interfaces.MajongFlow) {

}
