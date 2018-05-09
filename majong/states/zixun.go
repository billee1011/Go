package states

import (
	"fmt"
	"steve/client_pb/msgId"
	"steve/client_pb/room"
	"steve/majong/interfaces"
	"steve/majong/utils"
	majongpb "steve/server_pb/majong"

	"github.com/golang/protobuf/proto"
)

// ZiXunState 摸牌状态
type ZiXunState struct {
}

var _ interfaces.MajongState = new(ZiXunState)

// ProcessEvent 处理事件
func (s *ZiXunState) ProcessEvent(eventID majongpb.EventID, eventContext []byte, flow interfaces.MajongFlow) (newState majongpb.StateID, err error) {
	switch eventID {
	case majongpb.EventID_event_angang_request:
		{
			// message := &clientpb.GameActionReq{}
			message := &majongpb.AngangRequestEvent{}
			err := proto.Unmarshal(eventContext, message)
			if err != nil {
				return majongpb.StateID_state_zixun, errInvalidEvent
			}
			return s.angang(flow, message)
		}
	case majongpb.EventID_event_zimo_request:
		{
			// message := &clientpb.GameActionReq{}
			message := &majongpb.ZimoRequestEvent{}
			err := proto.Unmarshal(eventContext, message)
			if err != nil {
				return majongpb.StateID_state_zixun, errInvalidEvent
			}
			return s.zimo(flow, message)

		}
	case majongpb.EventID_event_chupai_request:
		{
			// message := &clientpb.GameChuPaiReq{}
			message := &majongpb.ChupaiRequestEvent{}
			err := proto.Unmarshal(eventContext, message)
			if err != nil {
				return majongpb.StateID_state_zixun, errInvalidEvent
			}
			return s.chupai(flow, message)
		}
	case majongpb.EventID_event_bugang_request:
		{
			// message := &clientpb.GameActionReq{}
			message := &majongpb.BugangRequestEvent{}
			err := proto.Unmarshal(eventContext, message)
			if err != nil {
				return majongpb.StateID_state_zixun, errInvalidEvent
			}
			return s.bugang(flow, message)
		}
	default:
		{
			return majongpb.StateID_state_zixun, errInvalidEvent
		}
	}
}

//angang 决策暗杠
func (s *ZiXunState) angang(flow interfaces.MajongFlow, message *majongpb.AngangRequestEvent) (majongpb.StateID, error) {
	can, err := s.canAnGang(flow, message)
	if err != nil {
		return majongpb.StateID_state_zixun, err
	}
	if can {
		context := flow.GetMajongContext()
		activePlayer := utils.GetPlayerByID(context.GetPlayers(), context.GetActivePlayer())
		card := message.GetCards()
		activePlayer.HandCards, _ = utils.DeleteCardFromLast(activePlayer.HandCards, card)
		activePlayer.HandCards, _ = utils.DeleteCardFromLast(activePlayer.HandCards, card)
		activePlayer.HandCards, _ = utils.DeleteCardFromLast(activePlayer.HandCards, card)
		activePlayer.HandCards, _ = utils.DeleteCardFromLast(activePlayer.HandCards, card)
		// activePlayer.PossibleActions = activePlayer.PossibleActions[:0]
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
		angang := &room.RoomAngangNtf{
			Player: proto.Uint64(activePlayer.PalyerId),
			Card:   angangCard,
		}
		toClient := interfaces.ToClientMessage{
			MsgID: int(msgid.MsgID_room_angang_ntf),
			Msg:   angang,
		}
		flow.PushMessages(playerIDs, toClient)
		return majongpb.StateID_state_angang, nil
	}
	return majongpb.StateID_state_zixun, errInvalidEvent
}

//bugang 决策补杠
func (s *ZiXunState) bugang(flow interfaces.MajongFlow, message *majongpb.BugangRequestEvent) (majongpb.StateID, error) {
	can, err := s.canBuGang(flow, message)
	if err != nil {
		return majongpb.StateID_state_zixun, err
	}
	//TODO: 检查是否能可以进行抢杠胡
	//可以抢杠胡的话，进入等待抢杠胡的状态,否则是补杠状态
	ctx := flow.GetMajongContext()
	card := message.GetCards()
	cardI, _ := utils.CardToInt(*card)
	var hasQGanghu bool
	for _, player := range ctx.GetPlayers() {
		player.PossibleActions = player.PossibleActions[:0]
		if player.GetPalyerId() != ctx.ActivePlayer {
			flag := utils.CheckHu(player.HandCards, uint32(*cardI))
			if flag {
				hasQGanghu = true
				playersID := make([]uint64, 0, 0)
				playersID = append(playersID, player.PalyerId)
				qianggangCard, _ := utils.CardToRoomCard(card)
				angang := &room.RoomWaitQianggangHuNtf{
					Card: qianggangCard,
				}
				toClientMessage := interfaces.ToClientMessage{
					MsgID: int(msgid.MsgID_room_wait_qiangganghu_ntf),
					Msg:   angang,
				}
				player.PossibleActions = append(player.PossibleActions, majongpb.Action_action_qiangganghu)
				flow.PushMessages(playersID, toClientMessage)
			}
		}
	}
	if can {
		activePlayer := utils.GetPlayerByID(ctx.Players, ctx.ActivePlayer)
		//先进行补杠的处理，将碰的action改为补杠的action，如果之后有人执行了抢杠胡，再将对应的补杠action改为碰
		activePlayer.HandCards, _ = utils.DeleteCardFromLast(activePlayer.HandCards, card)
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
		// activePlayer.PossibleActions = activePlayer.PossibleActions[:0]
		if hasQGanghu {
			return majongpb.StateID_state_waitqiangganghu, nil
		}
		//广播补杠消息
		playerIDs := make([]uint64, 0, 0)
		for _, player := range ctx.Players {
			playerIDs = append(playerIDs, player.GetPalyerId())
		}
		bugangCard, _ := utils.CardToRoomCard(card)
		bugang := &room.RoomBugangNtf{
			Player: proto.Uint64(activePlayer.PalyerId),
			Card:   bugangCard,
		}
		toClient := interfaces.ToClientMessage{
			MsgID: int(msgid.MsgID_room_bugang_ntf),
			Msg:   bugang,
		}
		flow.PushMessages(playerIDs, toClient)
		return majongpb.StateID_state_bugang, nil
	}
	return majongpb.StateID_state_zixun, errInvalidEvent
}

//zimo 决策自摸
func (s *ZiXunState) zimo(flow interfaces.MajongFlow, message *majongpb.ZimoRequestEvent) (majongpb.StateID, error) {
	can, err := s.canZiMo(flow, message)
	if err != nil {
		return majongpb.StateID_state_zixun, err
	}
	if can {
		ctx := flow.GetMajongContext()
		//自摸决策通过后，移除胡的那张牌，并且添加到胡的牌中
		head := message.GetHead()
		activePlayer := utils.GetPlayerByID(ctx.Players, head.GetPlayerId())
		card := message.GetCards()
		activePlayer.HandCards, _ = utils.DeleteCardFromLast(activePlayer.HandCards, card)
		huCard := &majongpb.HuCard{
			Card:      card,
			SrcPlayer: message.GetHead().GetPlayerId(),
			Type:      majongpb.HuType_hu_zimo,
		}
		activePlayer.HuCards = append(activePlayer.HuCards, huCard)
		// activePlayer.PossibleActions = activePlayer.PossibleActions[:0]
		toclientCard, _ := utils.CardToRoomCard(card)
		playersID := make([]uint64, 0, 0)
		for _, player := range flow.GetMajongContext().GetPlayers() {
			playersID = append(playersID, player.PalyerId)
		}
		ntf := &room.RoomZimoNtf{
			Player: proto.Uint64(activePlayer.PalyerId),
			Card:   toclientCard,
		}
		toClientMessage := interfaces.ToClientMessage{
			MsgID: int(msgid.MsgID_room_zimo_ntf),
			Msg:   ntf,
		}
		flow.PushMessages(playersID, toClientMessage)
		return majongpb.StateID_state_zimo, nil
	}
	return majongpb.StateID_state_zixun, errInvalidEvent
}

//chupai 决策出牌
func (s *ZiXunState) chupai(flow interfaces.MajongFlow, message *majongpb.ChupaiRequestEvent) (majongpb.StateID, error) {
	//检查玩家收牌中是否包含出的牌
	context := flow.GetMajongContext()
	pid := message.GetHead().GetPlayerId()
	if context.GetActivePlayer() != pid {
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
		activePlayer.HandCards, _ = utils.DeleteCardFromLast(activePlayer.HandCards, card)
		context.LastOutCard = card
		activePlayer.OutCards = append(activePlayer.OutCards, card)
		playersID := make([]uint64, 0, 0)
		for _, player := range context.GetPlayers() {
			playersID = append(playersID, player.GetPalyerId())
		}
		cardToClient, _ := utils.CardToRoomCard(card)
		toClientMessage := interfaces.ToClientMessage{
			MsgID: int(msgid.MsgID_room_chupai_ntf),
			Msg: &room.RoomChupaiNtf{
				Player: proto.Uint64(activePlayer.GetPalyerId()),
				Card:   cardToClient,
			},
		}
		flow.PushMessages(playersID, toClientMessage)
		// activePlayer.PossibleActions = activePlayer.PossibleActions[:0]
		return majongpb.StateID_state_chupai, nil
	}
	return majongpb.StateID_state_zixun, errInvalidEvent
}

//checkAnGang 检查暗杠 (判断当前事件是否可行)
func (s *ZiXunState) canAnGang(flow interfaces.MajongFlow, message *majongpb.AngangRequestEvent) (bool, error) {
	angangCard := message.Cards
	mjContext := flow.GetMajongContext()
	wallCards := mjContext.GetWallCards()
	activePlayer := utils.GetPlayerByID(mjContext.Players, mjContext.ActivePlayer)
	if len(wallCards) == 0 {
		return false, fmt.Errorf("墙牌为0，不允许暗杠")
	}
	if message.GetHead().GetPlayerId() != mjContext.ActivePlayer {
		return false, fmt.Errorf("当前玩家不是可执行玩家，不予操作")
	}
	//检查手牌中是否有足够的暗杠牌
	gangCardsNum := 0
	// angangCard, err := utils.IntToCard(int32(actionInfo.ActionCards[0]))
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
		newcards, _ = utils.DeleteCardFromLast(newcards, angangCard)
		newcards, _ = utils.DeleteCardFromLast(newcards, angangCard)
		newcards, _ = utils.DeleteCardFromLast(newcards, angangCard)
		newcards, _ = utils.DeleteCardFromLast(newcards, angangCard)
		newcardsI, _ := utils.CardsToInt(newcards)
		cardsI := utils.IntToUtilCard(newcardsI)
		laizi := make(map[utils.Card]bool)
		huCards := utils.FastCheckTingV2(cardsI, laizi)
		if !utils.ContainHuCards(huCards, utils.HuCardsToUtilCards(activePlayer.HuCards)) {
			return false, fmt.Errorf("当前的明杠操作会影响胡牌后的胡牌牌型，不允许暗杠")
		}
	}
	return true, nil
}

//checkBuGang 检查补杠 (判断当前事件是否可行)
func (s *ZiXunState) canBuGang(flow interfaces.MajongFlow, message *majongpb.BugangRequestEvent) (bool, error) {
	context := flow.GetMajongContext()
	activePlayer := utils.GetPlayerByID(context.Players, context.ActivePlayer)
	if len(context.WallCards) == 0 {
		return false, fmt.Errorf("墙牌为0时，不予补杠")
	}
	//判断是否轮到当前玩家操作
	if activePlayer.PalyerId != message.GetHead().GetPlayerId() {
		return false, fmt.Errorf("当前玩家不允许操作")
	}
	bugangCard := message.GetCards()
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
		newcards, _ = utils.DeleteCardFromLast(newcards, bugangCard)
		newcardsI, _ := utils.CardsToInt(newcards)
		cardsI := utils.IntToUtilCard(newcardsI)
		laizi := make(map[utils.Card]bool)
		huCards := utils.FastCheckTingV2(cardsI, laizi)
		if !utils.ContainHuCards(huCards, utils.HuCardsToUtilCards(activePlayer.HuCards)) {
			return false, fmt.Errorf("当前的补杠杠操作会影响胡牌后的胡牌牌型，不允许补杠")
		}
	}
	return true, nil
}

//checkZiMo 检查自摸 (判断当前事件是否可行)
func (s *ZiXunState) canZiMo(flow interfaces.MajongFlow, message *majongpb.ZimoRequestEvent) (bool, error) {
	context := flow.GetMajongContext()
	activePlayer := utils.GetPlayerByID(context.Players, context.ActivePlayer)
	//判断是否轮到当前玩家操作
	if activePlayer.PalyerId != message.GetHead().GetPlayerId() {
		return false, fmt.Errorf("当前玩家不允许操作")
	}
	handCard := activePlayer.GetHandCards()
	if utils.CheckHasDingQueCard(handCard, activePlayer.GetDingqueColor()) {
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

// OnEntry 进入状态
func (s *ZiXunState) OnEntry(flow interfaces.MajongFlow) {

}

// OnExit 退出状态
func (s *ZiXunState) OnExit(flow interfaces.MajongFlow) {

}
