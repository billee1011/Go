package common

import (
	msgid "steve/client_pb/msgId"
	"steve/client_pb/room"
	"steve/common/mjoption"
	"steve/gutils"
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
		// 通知听牌提示
		utils.NotifyTingCards(flow, context.GetLastChupaiPlayer())
		players := context.GetPlayers()
		card := context.GetLastOutCard()
		var hasChupaiwenxun bool
		//出完牌后，将上轮添加的胡牌玩家列表重置
		context.LastHuPlayers = context.LastHuPlayers[:0]
		for _, player := range utils.GetCanXpPlayers(players, context) { // 能正常行牌的玩家才进行查动作
			//每个玩家的possibleActions都需要清空
			player.PossibleActions = player.PossibleActions[:0]
			player.EnbleChiCards = player.EnbleChiCards[:0]
			logrus.WithFields(logrus.Fields{"playerID": player.GetPalyerId(),
				"xpStates": player.GetXpState()}).Info("出牌：每个玩家的状态")
			if context.GetLastChupaiPlayer() == player.GetPalyerId() {
				continue
			}
			need := s.checkActions(flow, player, card)
			if need {
				hasChupaiwenxun = true
			}
		}
		if hasChupaiwenxun {
			logrus.WithFields(logrus.Fields{
				"player":  context.GetLastChupaiPlayer(),
				"outCard": card,
			}).Info("出牌信息")
			return majongpb.StateID_state_chupaiwenxun, nil
		}
		player := utils.GetNextXpPlayerByID(context.GetLastChupaiPlayer(), context.GetPlayers(), context)
		logrus.WithFields(logrus.Fields{"playerID": player.GetPalyerId(),
			"xpStates": player.GetXpState()}).Info("出牌：下一个摸牌玩家的状态")
		context.MopaiPlayer = player.GetPalyerId()
		context.MopaiType = majongpb.MopaiType_MT_NORMAL
		return majongpb.StateID_state_mopai, nil
	}
	return majongpb.StateID_state_chupai, global.ErrInvalidEvent
}

//checkActions 检查玩家可以有哪些操作
func (s *ChupaiState) checkActions(flow interfaces.MajongFlow, player *majongpb.Player, card *majongpb.Card) bool {
	context := flow.GetMajongContext()
	xpOption := mjoption.GetXingpaiOption(int(context.GetXingpaiOptionId()))
	canMingGang := s.checkMingGang(flow, player, card, xpOption)
	if canMingGang {
		player.PossibleActions = append(player.PossibleActions, majongpb.Action_action_gang)
	}
	canDianPao := s.checkDianPao(context, player, card)
	if canDianPao {
		context.LastHuPlayers = append(context.LastHuPlayers, player.GetPalyerId())
		player.PossibleActions = append(player.PossibleActions, majongpb.Action_action_hu)
	}
	canPeng := s.checkPeng(context, player, card)
	if canPeng {
		player.PossibleActions = append(player.PossibleActions, majongpb.Action_action_peng)
	}
	chiSlice := make([]uint32, 0)
	if xpOption.EnableChi {
		chiSlice = s.checkChi(context, player, card)
		if len(chiSlice) != 0 {
			player.EnbleChiCards = chiSlice
			player.PossibleActions = append(player.PossibleActions, majongpb.Action_action_chi)
		}
	}
	if len(player.PossibleActions) > 0 {
		if len(player.GetHuCards()) == 0 || !canDianPao {
			player.PossibleActions = append(player.PossibleActions, majongpb.Action_action_qi)
		}
	}
	logrus.WithFields(logrus.Fields{
		"func_name":   "checkActions",
		"player":      player.GetPalyerId(),
		"check_card":  card,
		"canPeng":     canPeng,
		"canMingGang": canMingGang,
		"canDianPao":  canDianPao,
		"chiSlice":    chiSlice,
		"handCards":   gutils.FmtMajongpbCards(player.GetHandCards()),
	}).Info("检测玩家是否有特殊操作")
	return canDianPao || canMingGang || canPeng || len(chiSlice) != 0
}

//checkMingGang 查明杠
func (s *ChupaiState) checkMingGang(flow interfaces.MajongFlow, player *majongpb.Player, card *majongpb.Card, xpOption *mjoption.XingPaiOption) bool {
	// 没有墙牌不能明杠
	context := flow.GetMajongContext()
	if !utils.HasAvailableWallCards(flow) {
		// if len(context.WallCards) == 0 {
		return false
	}
	outCard := context.GetLastOutCard()
	color := player.GetDingqueColor()
	//定缺牌不查
	if xpOption.EnableDingque && outCard.Color == color {
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
			// 推倒胡
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
func (s *ChupaiState) checkPeng(context *majongpb.MajongContext, player *majongpb.Player, card *majongpb.Card) bool {
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
func (s *ChupaiState) checkDianPao(context *majongpb.MajongContext, player *majongpb.Player, card *majongpb.Card) bool {
	cpCard := context.GetLastOutCard()
	color := player.GetDingqueColor()
	hasDingQueCard := gutils.CheckHasDingQueCard(player.HandCards, color)
	if hasDingQueCard {
		return false
	}
	handCard := player.GetHandCards() // 当前点炮胡玩家手牌
	cardI, _ := utils.CardToInt(*cpCard)
	result := utils.CheckHu(handCard, uint32(*cardI), false)
	if result.Can {
		return true
	}
	return false
}

// checkChi 查吃
func (s *ChupaiState) checkChi(context *majongpb.MajongContext, player *majongpb.Player, card *majongpb.Card) []uint32 {
	//判断当前玩家是否可以进行吃的出牌问询，二人麻将只能下家吃牌
	chicards := make([]uint32, 0)
	if utils.GetNextXpPlayerByID(context.GetLastChupaiPlayer(), context.GetPlayers(), context).GetPalyerId() != player.GetPalyerId() {
		return chicards
	}
	//只有万条筒可以进行吃的操作
	color := card.GetColor()
	point := card.GetPoint()
	if color == majongpb.CardColor_ColorFeng || card.GetColor() == majongpb.CardColor_ColorHua {
		return chicards
	}
	handCards := player.GetHandCards()
	var A, B, C, D bool
	//将下家手牌拿出来与上家出的牌进行对比
	for _, hc := range handCards {
		//查三种吃的方式，左边吃，中间吃，右边吃
		if hc.GetColor() != color {
			continue
		}
		switch hc.GetPoint() {
		case point - 2:
			A = true
		case point - 1:
			B = true
		case point + 1:
			C = true
		case point + 2:
			D = true
		}
	}
	cardInInt := utils.ServerCard2Uint32(card)
	if A && B {
		chicards = append(chicards, []uint32{cardInInt - 2, cardInInt - 1, cardInInt}...)
	}
	if B && C {
		chicards = append(chicards, []uint32{cardInInt - 1, cardInInt, cardInInt + 1}...)
	}
	if C && D {
		chicards = append(chicards, []uint32{cardInInt, cardInInt + 1, cardInInt + 2}...)
	}
	return chicards
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
