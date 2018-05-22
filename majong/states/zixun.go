package states

import (
	"fmt"
	"steve/client_pb/msgId"
	"steve/client_pb/room"
	"steve/majong/cardtype"
	"steve/majong/global"
	"steve/majong/interfaces"
	"steve/majong/interfaces/facade"
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
			return majongpb.StateID_state_zixun, nil
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
	mjContext := flow.GetMajongContext()
	if canAnGang {
		mjContext.GangCard = message.GetCard()
		mjContext.LastGangPlayer = message.GetHead().GetPlayerId()
		return majongpb.StateID_state_angang, nil
	}
	canBuGang, _ := s.canBuGang(flow, message)
	if canBuGang {
		//补杠的时候，可能会有玩家抢杠胡，所以此处也将胡牌玩家列表清空
		flow.GetMajongContext().LastHuPlayers = flow.GetMajongContext().LastHuPlayers[:0]
		mjContext.GangCard = message.GetCard()
		mjContext.LastGangPlayer = message.GetHead().GetPlayerId()
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
	//检查玩家是否胡过牌,胡过牌的话,摸啥打啥,能胡不让打
	card := message.GetCards()
	activePlayer := utils.GetPlayerByID(context.GetPlayers(), pid)
	if len(activePlayer.GetHuCards()) > 0 {
		if !utils.CardEqual(card, context.GetLastMopaiCard()) {
			return majongpb.StateID_state_zixun, nil
		}
		if s.checkZiMo(flow) {
			return majongpb.StateID_state_zixun, fmt.Errorf("玩家当前只能选择胡牌,不能进行出牌操作")
		}
		context.LastOutCard = card
		context.LastChupaiPlayer = pid
		return majongpb.StateID_state_chupai, nil
	}
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
	return true, nil
}

//checkBuGang 检查补杠 (判断当前事件是否可行)
func (s *ZiXunState) canBuGang(flow interfaces.MajongFlow, message *majongpb.GangRequestEvent) (bool, error) {
	context := flow.GetMajongContext()
	activePlayer := utils.GetPlayerByID(context.Players, s.getZixunPlayer(flow))
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
	return true, nil
}

func (s *ZiXunState) canPlayerZimo(flow interfaces.MajongFlow) bool {
	playerID := s.getZixunPlayer(flow)
	mjContext := flow.GetMajongContext()
	player := utils.GetPlayerByID(mjContext.GetPlayers(), playerID)
	handCard := player.GetHandCards()
	if utils.CheckHasDingQueCard(handCard, player.GetDingqueColor()) {
		return false
	}
	l := len(handCard)
	if l%3 != 2 {
		return false
	}
	flag := utils.CheckHu(handCard, 0)
	if !flag {
		return false
	}
	return true
}

// canZiMo 检查自摸 (判断当前事件是否可行)
func (s *ZiXunState) canZiMo(flow interfaces.MajongFlow, message *majongpb.HuRequestEvent) (bool, error) {
	if s.getZixunPlayer(flow) != message.GetHead().GetPlayerId() {
		return false, fmt.Errorf("当前玩家不允许操作")
	}
	return s.canPlayerZimo(flow), nil
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
	//分三种情况,1胡过,且摸到牌能胡,chupaicard字段不给值
	if len(player.GetHuCards()) > 0 {
		if !*zixunNtf.EnableZimo {
			//2胡过,且不能自摸,摸什么打什么
			zixunNtf.EnableChupaiCards = []uint32{utils.ServerCard2Uint32(context.GetLastMopaiCard())}
		}
	} else {
		//3没胡过,所有手牌都可以打
		zixunNtf.EnableChupaiCards = utils.ServerCards2Uint32(player.GetHandCards())
	}
	//查听,打什么,听什么
	s.checkTing(zixunNtf, player, context)
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
		"pengCards": FmtPengCards(player.GetPengCards()),
		"handCards": FmtMajongpbCards(player.GetHandCards()),
		"wallCards": FmtMajongpbCards(context.GetWallCards()),
		"huCards":   FmtHuCards(player.GetHuCards()),
	}).Infoln("自询通知")
}

func (s *ZiXunState) getPengCards(pengCards []*majongpb.PengCard) []*majongpb.Card {
	resultCards := []*majongpb.Card{}
	for _, pengCard := range pengCards {
		resultCards = append(resultCards, pengCard.GetCard())
	}
	return resultCards
}

func (s *ZiXunState) getGangCards(gangCards []*majongpb.GangCard) []*majongpb.Card {
	resultCards := []*majongpb.Card{}
	for _, gangCard := range gangCards {
		resultCards = append(resultCards, gangCard.GetCard())
	}
	return resultCards
}

// checkTing 查听
func (s *ZiXunState) checkTing(zixunNtf *room.RoomZixunNtf, player *majongpb.Player, context *majongpb.MajongContext) {
	dqnum := 0
	for _, card := range player.GetHandCards() {
		if card.Color == player.DingqueColor {
			dqnum++
		}
	}
	logrus.WithFields(logrus.Fields{
		"手中定缺牌的个数": dqnum,
	}).Info("查听")
	switch dqnum {
	//没有定缺牌的时候正常查听
	case 0:
		{
			tingInfos := utils.GetPlayCardCheckTing(player.GetHandCards(), nil)
			s.addTingInfo(zixunNtf, player, context, tingInfos)
		}
	// 当定缺牌为1的时候,只有定缺牌才有听牌提示
	case 1:
		{
			newTingInfos := map[utils.Card][]utils.Card{}
			tingInfos := utils.GetPlayCardCheckTing(player.GetHandCards(), nil)
			for outCard, tingCard := range tingInfos {
				card, _ := utils.IntToCard(int32(outCard))
				if card.GetColor() == player.DingqueColor {
					newTingInfos[outCard] = tingCard
				}
			}
			//满足条件说明,打出这张定缺牌可以进入听牌状态
			s.addTingInfo(zixunNtf, player, context, newTingInfos)
		}
	// 玩家的定缺牌数量超过1张的时候,不查听
	default:
		{
			return
		}
	}
}

// addTingInfo 自询通知添加听牌信息
func (s *ZiXunState) addTingInfo(zixunNtf *room.RoomZixunNtf, player *majongpb.Player, context *majongpb.MajongContext, tingInfos map[utils.Card][]utils.Card) {
	canTingInfos := []*room.CanTingCardInfo{}
	for outCard, tingInfo := range tingInfos {
		tingCardInfo := []*room.TingCardInfo{}
		outCard0, err := utils.IntToCard(int32(outCard))
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"func_name": "utils.IntToCard",
			}).Error("牌型转换失败")
		}
		newHand, success := utils.RemoveCards(player.GetHandCards(), outCard0, 1)
		if !success {
			logrus.WithFields(logrus.Fields{
				"func_name": "addTingInfo",
			}).Error("牌型移除失败")
		}
		for _, tt := range tingInfo {
			huCard, _ := utils.IntToCard(int32(tt))
			times, _ := facade.CalculateCardValue(&cardtype.ScxlCardTypeCalculator{}, interfaces.CardCalcParams{
				HandCard: newHand,
				PengCard: s.getPengCards(player.GetPengCards()),
				GangCard: s.getGangCards(player.GetGangCards()),
				HuCard:   huCard,
				GameID:   int(context.GetGameId()),
			})
			tingCardInfo = append(tingCardInfo, &room.TingCardInfo{
				TingCard: proto.Uint32(uint32(tt)),
				Times:    proto.Uint32(times),
			})
		}
		canTingInfos = append(canTingInfos, &room.CanTingCardInfo{
			OutCard:      proto.Uint32(uint32(outCard)),
			TingCardInfo: tingCardInfo,
		})
	}
	zixunNtf.CanTingCardInfo = canTingInfos
}

// checkZiMo 查自摸
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

func (s *ZiXunState) checkPlayerAngang(player *majongpb.Player) []uint32 {
	result := make([]uint32, 0, 0)

	//分两种情况查暗杠，一种是胡牌前，一种胡牌后
	huCards := player.GetHuCards()
	handCard := player.GetHandCards()

	cardNum := make(map[majongpb.Card]int)
	for i := 0; i < len(handCard); i++ {
		num := cardNum[*handCard[i]]
		num++
		cardNum[*handCard[i]] = num
	}
	color := player.GetDingqueColor()
	for k, num := range cardNum {
		if k.Color != color && num == 4 {
			if len(huCards) > 0 {
				newCards := []*majongpb.Card{}
				newCards = append(newCards, handCard...)
				newCards, _ = utils.RemoveCards(newCards, &k, 4)
				utilCards := utils.CardsToUtilCards(newCards)
				tingCards := utils.FastCheckTingV2(utilCards, map[utils.Card]bool{})
				if utils.ContainHuCards(tingCards, utils.HuCardsToUtilCards(huCards)) {
					result = append(result, utils.ServerCard2Uint32(&k))
				}
			} else {
				result = append(result, utils.ServerCard2Uint32(&k))
			}
		}
	}
	logrus.WithFields(logrus.Fields{
		"func_name": "ZiXunState.checkPlayerAngang",
		"hand_card": FmtMajongpbCards(handCard),
		"hu_cards":  FmtHuCards(huCards),
		"card_num":  cardNum,
		"result":    result,
	}).Debugln("查暗杠")
	return result
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
	return s.checkPlayerAngang(activePlayer)
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
