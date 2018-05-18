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

// getZixunPlayer 获取自询玩家 ID
// 如果没有上一个摸牌的人，则返回庄家。否则返回上一个摸牌的玩家
func (s *ZiXunState) getZixunPlayer(flow interfaces.MajongFlow) uint64 {
	mjContext := flow.GetMajongContext()
	zxType := mjContext.GetZixunType()
	if zxType == majongpb.ZixunType_ZXT_PENG {
		return mjContext.GetLastPengPlayer()
	}
	return mjContext.GetLastMopaiPlayer()
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
	if s.getZixunPlayer(flow) != pid {
		return majongpb.StateID_state_zixun, fmt.Errorf("未到玩家：%v 出牌，当前应该出牌的玩家是：%v", pid, s.getZixunPlayer(flow))
	}
	card := message.GetCards()
	activePlayer := utils.GetPlayerByID(context.GetPlayers(), pid)
	for _, c := range activePlayer.GetHandCards() {
		if utils.CardEqual(c, card) {
			context.LastOutCard = card
			context.LastChupaiPlayer = pid
			return majongpb.StateID_state_chupai, nil
		}
	}
	return majongpb.StateID_state_zixun, nil
}

//checkAnGang 检查暗杠 (判断当前事件是否可行)
func (s *ZiXunState) canAnGang(flow interfaces.MajongFlow, message *majongpb.GangRequestEvent) (bool, error) {
	angangCard := message.GetCard()
	mjContext := flow.GetMajongContext()
	wallCards := mjContext.GetWallCards()
	playerID := s.getZixunPlayer(flow)
	activePlayer := utils.GetPlayerByID(mjContext.Players, playerID)
	if len(wallCards) == 0 {
		return false, fmt.Errorf("墙牌为0，不允许暗杠")
	}
	if message.GetHead().GetPlayerId() != playerID {
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
	activePlayer := utils.GetPlayerByID(context.Players, context.LastMopaiPlayer)
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

// canZiMo 检查自摸 (判断当前事件是否可行)
func (s *ZiXunState) canZiMo(flow interfaces.MajongFlow, message *majongpb.HuRequestEvent) (bool, error) {
	context := flow.GetMajongContext()
	playerID := message.GetHead().GetPlayerId()
	if s.getZixunPlayer(flow) != message.GetHead().GetPlayerId() {
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
				player.PossibleActions = append(player.PossibleActions, majongpb.Action_action_hu)
			}
		}
	}
	return hasQGanghu
}

// canQi 是否可弃
func (s *ZiXunState) canQi(canAngang bool, canBugang bool, canZimo bool, hasHu bool) bool {
	// 没有可选操作不可弃
	if !canAngang && !canBugang && !canZimo {
		return false
	}
	// 胡过且能自摸则不能弃， 否则可弃
	return !(hasHu && canZimo)
}

//checkActions 检测进入自询状态下，玩家有哪些可以可行的事件
func (s *ZiXunState) checkActions(flow interfaces.MajongFlow) {
	context := flow.GetMajongContext()
	zixunNtf := &room.RoomZixunNtf{}
	isPengZixun := context.GetZixunType() == majongpb.ZixunType_ZXT_PENG
	playerID := s.getZixunPlayer(flow)
	player := utils.GetPlayerByID(context.Players, playerID)

	if !isPengZixun {
		zixunNtf.EnableAngangCards = s.checkAnGang(flow)
		zixunNtf.EnableBugangCards = s.checkBuGang(flow)
		canZimo := s.checkZiMo(flow)
		zixunNtf.EnableZimo = proto.Bool(canZimo)

		canAngang := len(zixunNtf.GetEnableAngangCards()) > 0
		canBugang := len(zixunNtf.GetEnableBugangCards()) > 0
		hasHu := len(player.GetHuCards()) > 0
		canQi := s.canQi(canAngang, canBugang, canZimo, hasHu)
		zixunNtf.EnableQi = proto.Bool(canQi)
	}
	//TODO:暂时将所有的手牌都设置为可以出的牌
	zixunNtf.EnableChupaiCards = utils.ServerCards2Uint32(player.GetHandCards())
	playerIDs := make([]uint64, 0, 0)
	playerIDs = append(playerIDs, playerID)
	toClient := interfaces.ToClientMessage{
		MsgID: int(msgid.MsgID_ROOM_ZIXUN_NTF),
		Msg:   zixunNtf,
	}
	flow.PushMessages(playerIDs, toClient)
	logrus.WithFields(logrus.Fields{
		"ntf":       zixunNtf.String(),
		"player_id": playerID,
	}).Infoln("自询通知")
}

//checkZiMo 查自摸
func (s *ZiXunState) checkZiMo(flow interfaces.MajongFlow) bool {
	context := flow.GetMajongContext()
	activePlayerID := s.getZixunPlayer(flow)
	activePlayer := utils.GetPlayerByID(context.Players, activePlayerID)
	handCard := activePlayer.GetHandCards()
	if utils.CheckHasDingQueCard(handCard, activePlayer.GetDingqueColor()) {
		return false
	}
	l := len(handCard)
	if l%3 != 2 {
		return false
	}
	flag := utils.CheckHu(handCard, 0)
	return flag
}

// checkAnGang 查暗杠
func (s *ZiXunState) checkAnGang(flow interfaces.MajongFlow) (enableAngangCards []uint32) {
	enableAngangCards = make([]uint32, 0, 0)
	context := flow.GetMajongContext()
	if len(context.WallCards) == 0 {
		return
	}
	activePlayerID := s.getZixunPlayer(flow)
	activePlayer := utils.GetPlayerByID(context.Players, activePlayerID)
	//分两种情况查暗杠，一种是胡牌前，一种胡牌后
	hasHu := len(activePlayer.GetHuCards()) > 0
	handCard := activePlayer.GetHandCards()

	cardNum := make(map[*majongpb.Card]int)
	for i := 0; i < len(handCard); i++ {
		num := cardNum[handCard[i]]
		num++
		cardNum[handCard[i]] = num
	}
	color := activePlayer.GetDingqueColor()
	for k, num := range cardNum {
		if k.Color != color && num == 4 {
			if hasHu {
				newCards := []*majongpb.Card{}
				newCards = append(newCards, handCard...)
				newCards, _ = utils.RemoveCards(newCards, k, 4)
				utilCards := utils.CardsToUtilCards(newCards)
				huCards := utils.FastCheckTingV2(utilCards, map[utils.Card]bool{})
				if utils.ContainHuCards(huCards, utils.HuCardsToUtilCards(activePlayer.HuCards)) {
					enableAngangCards = append(enableAngangCards, utils.ServerCard2Uint32(k))
				}
			} else {
				enableAngangCards = append(enableAngangCards, utils.ServerCard2Uint32(k))
			}
		}
	}
	return
}

//checkBuGang 查补杠
func (s *ZiXunState) checkBuGang(flow interfaces.MajongFlow) []uint32 {
	enableBugangCards := []uint32{}
	context := flow.GetMajongContext()
	// 没有墙牌不能杠
	if len(context.WallCards) == 0 {
		return enableBugangCards
	}
	activePlayerID := s.getZixunPlayer(flow)
	activePlayer := utils.GetPlayerByID(context.Players, activePlayerID)
	//分两种情况查暗杠，一种是胡牌前，一种胡牌后
	hasHu := len(activePlayer.GetHuCards()) > 0
	pengCards := activePlayer.GetPengCards()
	// actioninfos := []*clientpb.ActionInfo{}
	for _, touchCard := range activePlayer.HandCards {
		for _, pengCard := range pengCards {
			if *pengCard.Card == *touchCard {
				removeCard, _ := utils.CardToInt(*touchCard)
				if hasHu {
					//创建副本，移除相应的杠牌进行查胡
					cardsI, _ := utils.CardsToInt(activePlayer.HandCards)
					newcardsI := make([]int32, 0, len(cardsI))
					newcardsI = append(newcardsI, cardsI...)
					newcardsI, _ = utils.DeleteIntCardFromLast(newcardsI, *removeCard)
					utilCards := utils.IntToUtilCard(newcardsI)
					laizi := make(map[utils.Card]bool)
					huCards := utils.FastCheckTingV2(utilCards, laizi)
					if utils.ContainHuCards(huCards, utils.HuCardsToUtilCards(activePlayer.HuCards)) {
						enableBugangCards = append(enableBugangCards, uint32(*removeCard))
					}
				} else {
					enableBugangCards = append(enableBugangCards, uint32(*removeCard))
				}
			}
		}
	}
	return enableBugangCards
}

// OnEntry 进入状态
func (s *ZiXunState) OnEntry(flow interfaces.MajongFlow) {
	s.checkActions(flow)
}

// OnExit 退出状态
func (s *ZiXunState) OnExit(flow interfaces.MajongFlow) {

}
