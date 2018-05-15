package states

import (
	"fmt"
	"steve/client_pb/msgId"
	"steve/client_pb/room"
	"steve/majong/global"
	"steve/majong/interfaces"
	"steve/majong/utils"
	majongpb "steve/server_pb/majong"

	"github.com/Sirupsen/logrus"

	"github.com/golang/protobuf/proto"
)

// ZiXunState 摸牌状态
type ZiXunState struct {
}

var _ interfaces.MajongState = new(ZiXunState)

// ProcessEvent 处理事件
func (s *ZiXunState) ProcessEvent(eventID majongpb.EventID, eventContext []byte, flow interfaces.MajongFlow) (newState majongpb.StateID, err error) {
	switch eventID {
	case majongpb.EventID_event_hu_request:
		{
			message := &majongpb.HuRequestEvent{}
			err := proto.Unmarshal(eventContext, message)
			if err != nil {
				return majongpb.StateID_state_zixun, global.ErrInvalidEvent
			}
			return s.zimo(flow, message)

		}
	case majongpb.EventID_event_chupai_request:
		{
			message := &majongpb.ChupaiRequestEvent{}
			err := proto.Unmarshal(eventContext, message)
			if err != nil {
				return majongpb.StateID_state_zixun, global.ErrInvalidEvent
			}
			return s.chupai(flow, message)
		}
	case majongpb.EventID_event_gang_request:
		{
			message := &majongpb.GangRequestEvent{}
			err := proto.Unmarshal(eventContext, message)
			if err != nil {
				return majongpb.StateID_state_zixun, global.ErrInvalidEvent
			}
			return s.gang(flow, message)

		}
	default:
		{
			return majongpb.StateID_state_zixun, global.ErrInvalidEvent
		}
	}
}

//angang 决策杠
func (s *ZiXunState) gang(flow interfaces.MajongFlow, message *majongpb.GangRequestEvent) (majongpb.StateID, error) {
	canAnGang, _ := s.canAnGang(flow, message)
	if canAnGang {
		return majongpb.StateID_state_angang, nil
	}
	canBuGang, _ := s.canBuGang(flow, message)
	if canBuGang {
		//补杠的时候，可能会有玩家抢杠胡，所以此处也将胡牌玩家列表清空
		flow.GetMajongContext().LastHuPlayers = flow.GetMajongContext().LastHuPlayers[:0]
		hasQGH := s.hasQiangGangHu(flow)
		if hasQGH {
			//TODO: 可以在这里广播补杠的消息（也可以在waitqiangganghu的entry中进行广播）
			return majongpb.StateID_state_waitqiangganghu, nil
		}
		//在这里判断是否可以抢杠胡，可以抢杠胡进入抢杠胡状态，否则进入补杠状态
		return majongpb.StateID_state_bugang, nil
	}
	logrus.Errorln("当前玩家的请求有问题，无法进入杠的状态，执行杠的逻辑")
	return majongpb.StateID_state_zixun, errInvalidEvent
}

//angang 决策暗杠
func (s *ZiXunState) angang(flow interfaces.MajongFlow) {
	context := flow.GetMajongContext()
	activePlayer := utils.GetPlayerByID(context.GetPlayers(), context.LastGangPlayer)
	card := context.GetGangCard()
	activePlayer.HandCards, _ = utils.RemoveCards(activePlayer.HandCards, card, 4)
	angangType := &majongpb.GangCard{
		Card:      card,
		Type:      majongpb.GangType_gang_angang,
		SrcPlayer: context.GetActivePlayer(),
	}
	activePlayer.GangCards = append(activePlayer.GangCards, angangType)
	//TODO:广播消息
	playerIDs := make([]uint64, 0, 0)
	for _, player := range context.Players {
		playerIDs = append(playerIDs, player.GetPalyerId())
	}
	angangCard, _ := utils.CardToRoomCard(card)
	angang := &room.RoomGangNtf{
		ToPlayerId:   proto.Uint64(activePlayer.PalyerId),
		Card:         angangCard,
		GangType:     room.GangType_AnGang.Enum(),
		FromPlayerId: proto.Uint64(activePlayer.PalyerId),
	}
	toClient := interfaces.ToClientMessage{
		MsgID: int(msgid.MsgID_ROOM_GANG_NTF),
		Msg:   angang,
	}
	flow.PushMessages(playerIDs, toClient)
}

//bugang 决策补杠
func (s *ZiXunState) bugang(flow interfaces.MajongFlow) {
	//TODO: 检查是否能可以进行抢杠胡
	//可以抢杠胡的话，进入等待抢杠胡的状态,否则是补杠状态
	ctx := flow.GetMajongContext()
	card := ctx.GetGangCard()
	activePlayer := utils.GetPlayerByID(ctx.Players, ctx.ActivePlayer)
	activePlayer.HandCards, _ = utils.RemoveCards(activePlayer.HandCards, card, 1)
	for k, peng := range activePlayer.PengCards {
		if utils.CardEqual(peng.Card, card) {
			activePlayer.PengCards = append(activePlayer.PengCards[:k], activePlayer.PengCards[k+1:]...)
		}
	}
	activePlayer.GangCards = append(activePlayer.GangCards, &majongpb.GangCard{
		Card:      card,
		Type:      majongpb.GangType_gang_bugang,
		SrcPlayer: activePlayer.PalyerId,
	})
	//广播补杠消息
	playerIDs := make([]uint64, 0, 0)
	for _, player := range ctx.Players {
		playerIDs = append(playerIDs, player.GetPalyerId())
	}
	bugangCard, _ := utils.CardToRoomCard(card)
	bugang := &room.RoomGangNtf{
		ToPlayerId:   proto.Uint64(activePlayer.PalyerId),
		Card:         bugangCard,
		GangType:     room.GangType_BuGang.Enum(),
		FromPlayerId: proto.Uint64(activePlayer.PalyerId),
	}
	toClient := interfaces.ToClientMessage{
		MsgID: int(msgid.MsgID_ROOM_GANG_NTF),
		Msg:   bugang,
	}
	flow.PushMessages(playerIDs, toClient)

}

//zimo 决策自摸
func (s *ZiXunState) zimo(flow interfaces.MajongFlow, message *majongpb.HuRequestEvent) (majongpb.StateID, error) {
	can, err := s.canZiMo(flow, message)
	if err != nil {
		return majongpb.StateID_state_zixun, err
	}
	if can {
		return majongpb.StateID_state_zimo, nil
	}
	return majongpb.StateID_state_zixun, global.ErrInvalidEvent
}

//chupai 决策出牌
func (s *ZiXunState) chupai(flow interfaces.MajongFlow, message *majongpb.ChupaiRequestEvent) (majongpb.StateID, error) {
	//检查玩家收牌中是否包含出的牌
	context := flow.GetMajongContext()
	pid := message.GetHead().GetPlayerId()
	if context.GetMopaiPlayer() != pid {
		return majongpb.StateID_state_zixun, fmt.Errorf("未到玩家：%v 出牌，当前应该出牌的玩家是：%v", pid, context.GetActivePlayer())
	}
	card := message.GetCards()
	activePlayer := utils.GetPlayerByID(context.GetPlayers(), pid)
	var canOutCard bool
	for _, c := range activePlayer.GetHandCards() {
		if utils.CardEqual(c, card) {
			canOutCard = true
			break
		}
	}
	if canOutCard {
		//决策成功，移除手牌，并且广播
		// activePlayer.HandCards, _ = utils.DeleteCardFromLast(activePlayer.HandCards, card)
		// context.LastOutCard = card
		// activePlayer.OutCards = append(activePlayer.OutCards, card)
		// playersID := make([]uint64, 0, 0)
		// for _, player := range context.GetPlayers() {
		// 	playersID = append(playersID, player.GetPalyerId())
		// }
		// cardToClient, _ := utils.CardToRoomCard(card)
		// toClientMessage := interfaces.ToClientMessage{
		// 	MsgID: int(msgid.MsgID_ROOM_CHUPAI_NTF),
		// 	Msg: &room.RoomChupaiNtf{
		// 		Player: proto.Uint64(activePlayer.GetPalyerId()),
		// 		Card:   cardToClient,
		// 	},
		// }
		// flow.PushMessages(playersID, toClientMessage)
		// activePlayer.PossibleActions = activePlayer.PossibleActions[:0]
		context.LastOutCard = card
		context.LastChupaiPlayer = pid
		return majongpb.StateID_state_chupai, nil
	}
	return majongpb.StateID_state_zixun, global.ErrInvalidEvent
}

//checkAnGang 检查暗杠 (判断当前事件是否可行)
func (s *ZiXunState) canAnGang(flow interfaces.MajongFlow, message *majongpb.GangRequestEvent) (bool, error) {
	angangCard := message.GetCard()
	mjContext := flow.GetMajongContext()
	wallCards := mjContext.GetWallCards()
	activePlayer := utils.GetPlayerByID(mjContext.Players, mjContext.MopaiPlayer)
	if len(wallCards) == 0 {
		return false, fmt.Errorf("墙牌为0，不允许暗杠")
	}
	if message.GetHead().GetPlayerId() != mjContext.GetLastMopaiPlayer() {
		return false, fmt.Errorf("当前玩家不是可执行玩家，不予操作")
	}
	//检查手牌中是否有足够的暗杠牌
	gangCardsNum := 0
	for _, card := range activePlayer.HandCards {
		if utils.CardEqual(card, angangCard) {
			gangCardsNum++
		}
	}
	if gangCardsNum != 4 {
		return false, fmt.Errorf("暗杠的牌不足4张")
	}
	//判断当前玩家是否胡过牌，胡过牌了，当前玩家需要移除杠牌进行查胡，判断移除后是否会影响胡牌
	if len(activePlayer.HuCards) > 0 {
		//创建副本，移除相应的杠牌进行查胡
		newcards := make([]*majongpb.Card, 0, len(activePlayer.HandCards))
		newcards = append(newcards, activePlayer.HandCards...)
		newcards, _ = utils.RemoveCards(newcards, angangCard, 4)
		newcardsI, _ := utils.CardsToInt(newcards)
		cardsI := utils.IntToUtilCard(newcardsI)
		laizi := make(map[utils.Card]bool)
		huCards := utils.FastCheckTingV2(cardsI, laizi)
		if !utils.ContainHuCards(huCards, utils.HuCardsToUtilCards(activePlayer.HuCards)) {
			return false, fmt.Errorf("当前的明杠操作会影响胡牌后的胡牌牌型，不允许暗杠")
		}
	}
	mjContext.GangCard = angangCard
	mjContext.LastGangPlayer = activePlayer.GetPalyerId()
	return true, nil
}

//checkBuGang 检查补杠 (判断当前事件是否可行)
func (s *ZiXunState) canBuGang(flow interfaces.MajongFlow, message *majongpb.GangRequestEvent) (bool, error) {
	context := flow.GetMajongContext()
	activePlayer := utils.GetPlayerByID(context.Players, context.MopaiPlayer)
	if len(context.WallCards) == 0 {
		return false, fmt.Errorf("墙牌为0时，不予补杠")
	}
	//判断是否轮到当前玩家操作
	if activePlayer.PalyerId != message.GetHead().GetPlayerId() {
		return false, fmt.Errorf("当前玩家不允许操作")
	}
	bugangCard := message.GetCard()
	handCards := activePlayer.GetHandCards()
	//检查手牌中是否有补杠牌
	exist := false
	for _, card := range handCards {
		if utils.CardEqual(card, bugangCard) {
			exist = true
		}
	}
	if !exist {
		return false, fmt.Errorf("手中没有请求中可以进行补杠的牌")
	}
	//补杠需要检查当前玩家已经操作的actionCard中是否有相应的碰进行补杠
	exist0 := false
	pengCards := activePlayer.GetPengCards()
	for _, pengCard := range pengCards {
		if utils.CardEqual(pengCard.Card, bugangCard) {
			exist0 = true
		}
	}
	if !exist0 {
		return false, fmt.Errorf("碰的牌中没有请求中可以进行补杠的牌")
	}
	//判断当前玩家是否胡过牌，胡过牌了，当前玩家需要移除杠牌进行查胡，判断移除后是否会影响胡牌
	if len(activePlayer.HuCards) > 0 {
		//创建副本，移除相应的杠牌进行查胡
		newcards := make([]*majongpb.Card, 0, len(handCards))
		newcards = append(newcards, handCards...)
		newcards, _ = utils.RemoveCards(newcards, bugangCard, 1)
		newcardsI, _ := utils.CardsToInt(newcards)
		cardsI := utils.IntToUtilCard(newcardsI)
		laizi := make(map[utils.Card]bool)
		huCards := utils.FastCheckTingV2(cardsI, laizi)
		if !utils.ContainHuCards(huCards, utils.HuCardsToUtilCards(activePlayer.HuCards)) {
			return false, fmt.Errorf("当前的补杠杠操作会影响胡牌后的胡牌牌型，不允许补杠")
		}
	}
	context.GangCard = bugangCard
	context.LastGangPlayer = activePlayer.GetPalyerId()
	return true, nil
}

//checkZiMo 检查自摸 (判断当前事件是否可行)
func (s *ZiXunState) canZiMo(flow interfaces.MajongFlow, message *majongpb.HuRequestEvent) (bool, error) {
	context := flow.GetMajongContext()
	playerID := message.GetHead().GetPlayerId()
	if context.GetMopaiPlayer() != message.GetHead().GetPlayerId() {
		return false, fmt.Errorf("当前玩家不允许操作")
	}
	player := utils.GetPlayerByID(context.Players, playerID)
	handCard := player.GetHandCards()
	if utils.CheckHasDingQueCard(handCard, player.GetDingqueColor()) {
		return false, fmt.Errorf("手中有定缺牌，不能胡牌")
	}
	l := len(handCard)
	if l%3 != 2 {
		return false, fmt.Errorf("手牌数量不正常，不能胡牌")
	}
	flag := utils.CheckHu(handCard, 0)
	if !flag {
		return false, fmt.Errorf("查胡为false，不能胡牌")
	}
	return true, nil
}

func (s *ZiXunState) hasQiangGangHu(flow interfaces.MajongFlow) bool {
	ctx := flow.GetMajongContext()
	card := ctx.GetGangCard()
	cardI, _ := utils.CardToInt(*card)
	var hasQGanghu bool
	for _, player := range ctx.GetPlayers() {
		player.PossibleActions = []majongpb.Action{}
		if player.GetPalyerId() != ctx.GetLastGangPlayer() {
			flag := utils.CheckHu(player.HandCards, uint32(*cardI))
			if flag {
				hasQGanghu = true
				// playersID := make([]uint64, 0, 0)
				// playersID = append(playersID, player.PalyerId)
				// qianggangCard, _ := utils.CardToRoomCard(card)
				// angang := &room.RoomWaitQianggangHuNtf{
				// 	Card: qianggangCard,
				// }
				// toClientMessage := interfaces.ToClientMessage{
				// 	MsgID: int(msgid.MsgID_ROOM_WAIT_QIANGGANGHU_NTF),
				// 	Msg:   angang,
				// }
				// ctx.LastHuPlayers = append(ctx.LastHuPlayers, player.GetPalyerId())
				player.PossibleActions = append(player.PossibleActions, majongpb.Action_action_hu)
				// flow.PushMessages(playersID, toClientMessage)
			}
		}
	}
	return hasQGanghu
}

// OnEntry 进入状态
func (s *ZiXunState) OnEntry(flow interfaces.MajongFlow) {

}

// OnExit 退出状态
func (s *ZiXunState) OnExit(flow interfaces.MajongFlow) {

}
