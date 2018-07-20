package common

import (
	"fmt"
	"steve/client_pb/msgid"
	"steve/client_pb/room"
	"steve/common/mjoption"
	"steve/gutils"
	"steve/majong/fantype"
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
		return s.doHu(eventContext, flow)
	case majongpb.EventID_event_chupai_request:
		return s.doChupai(eventContext, flow)
	case majongpb.EventID_event_gang_request:
		return s.doGang(eventContext, flow)
	}
	return majongpb.StateID_state_zixun, nil
}

func (s *ZiXunState) doGang(eventContext []byte, flow interfaces.MajongFlow) (majongpb.StateID, error) {
	message := &majongpb.GangRequestEvent{}
	err := proto.Unmarshal(eventContext, message)
	if err != nil {
		return majongpb.StateID_state_zixun, global.ErrInvalidEvent
	}
	return s.gang(flow, message)
}

func (s *ZiXunState) doHu(eventContext []byte, flow interfaces.MajongFlow) (majongpb.StateID, error) {
	message := &majongpb.HuRequestEvent{}
	err := proto.Unmarshal(eventContext, message)
	if err != nil {
		return majongpb.StateID_state_zixun, global.ErrInvalidEvent
	}
	return s.zimo(flow, message)
}

func (s *ZiXunState) doChupai(eventContext []byte, flow interfaces.MajongFlow) (majongpb.StateID, error) {
	message := &majongpb.ChupaiRequestEvent{}
	err := proto.Unmarshal(eventContext, message)
	if err != nil {
		return majongpb.StateID_state_zixun, global.ErrInvalidEvent
	}
	return s.chupai(flow, message)
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
	mjContext := flow.GetMajongContext()
	pid := message.GetHead().GetPlayerId()
	if gutils.GetZixunPlayer(mjContext) != pid {
		return majongpb.StateID_state_zixun, fmt.Errorf("未到玩家：%v 出牌，当前应该出牌的玩家是：%v", pid, gutils.GetZixunPlayer(mjContext))
	}

	//检查玩家是否胡过牌,胡过牌的话,摸啥打啥,能胡不让打
	card := message.GetCards()
	activePlayer := utils.GetPlayerByID(mjContext.GetPlayers(), pid)
	if len(activePlayer.GetHuCards()) > 0 {
		if !utils.CardEqual(card, mjContext.GetLastMopaiCard()) {
			return majongpb.StateID_state_zixun, nil
		}
		if s.checkZiMo(flow) {
			return majongpb.StateID_state_zixun, fmt.Errorf("玩家当前只能选择胡牌,不能进行出牌操作")
		}
		mjContext.LastOutCard = card
		mjContext.LastChupaiPlayer = pid
		return majongpb.StateID_state_chupai, nil
	}
	logEntry := logrus.WithFields(logrus.Fields{
		"player":   activePlayer.GetPalyerId(),
		"reqCard":  card,
		"hanCards": gutils.FmtMajongpbCards(activePlayer.GetHandCards()),
	})
	for _, c := range activePlayer.GetHandCards() {
		if !utils.CardEqual(c, card) {
			continue
		}
		mjContext.LastOutCard = card
		mjContext.LastChupaiPlayer = pid
		//出牌后标志听的状态
		if !gutils.IsTing(activePlayer) {
			TingAction := message.GetTingAction()
			if TingAction.GetEnableTing() {
				activePlayer.GetTingStateInfo().IsTing = true
				if TingAction.GetTingType() == majongpb.TingType_TT_TIAN_TING {
					activePlayer.GetTingStateInfo().IsTianting = true
				}
				activePlayer.GetTingStateInfo().BaotingyifaCount = activePlayer.GetChupaiCount() +
					int32(len(activePlayer.GetGangCards()))
			}
		}
		logEntry = logEntry.WithFields(logrus.Fields{
			"selectTing": gutils.IsTing(activePlayer),
			"tingType":   gutils.GetTingType(activePlayer),
		})
		logEntry.Infoln("玩家出牌请求")
		return majongpb.StateID_state_chupai, nil
	}
	return majongpb.StateID_state_zixun, nil
}

//checkAnGang 检查暗杠 (判断当前事件是否可行)
func (s *ZiXunState) canAnGang(flow interfaces.MajongFlow, message *majongpb.GangRequestEvent) (bool, error) {
	angangCard := message.GetCard()
	mjContext := flow.GetMajongContext()
	wallCards := mjContext.GetWallCards()
	playerID := gutils.GetZixunPlayer(mjContext)
	activePlayer := utils.GetPlayerByID(mjContext.Players, playerID)
	if len(wallCards) == 0 {
		return false, fmt.Errorf("墙牌为0，不允许暗杠")
	}
	if message.GetHead().GetPlayerId() != playerID {
		return false, fmt.Errorf("当前玩家不是可执行玩家，不予操作")
	}
	//检查手牌中是否有足够的暗杠牌
	gangCardsNum := utils.GetCardNum(angangCard, activePlayer.GetHandCards())
	if gangCardsNum != 4 {
		return false, fmt.Errorf("暗杠的牌不足4张")
	}
	//判断当前玩家是否胡过牌，胡过牌了，当前玩家需要移除杠牌进行查胡，判断移除后是否会影响胡牌
	if (gutils.IsHu(activePlayer) || gutils.IsTing(activePlayer)) && !utils.CheckHuByRemoveGangCards(activePlayer, angangCard, 4) {
		return false, fmt.Errorf("当前的明杠操作会影响胡牌后的胡牌牌型，不允许暗杠")
	}
	return true, nil
}

//checkBuGang 检查补杠 (判断当前事件是否可行)
func (s *ZiXunState) canBuGang(flow interfaces.MajongFlow, message *majongpb.GangRequestEvent) (bool, error) {
	mjContext := flow.GetMajongContext()
	activePlayer := utils.GetPlayerByID(mjContext.Players, gutils.GetZixunPlayer(mjContext))
	if len(mjContext.WallCards) == 0 {
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
		laizi := make(map[utils.Card]bool)
		huCards, _ := utils.GetTingCards(newcards, laizi)
		if len(huCards) == 0 || !utils.ContainHuCards(huCards, utils.HuCardsToUtilCards(activePlayer.HuCards)) {
			return false, fmt.Errorf("当前的补杠杠操作会影响胡牌后的胡牌牌型，不允许补杠")
		}
	}
	return true, nil
}

func (s *ZiXunState) canPlayerZimo(flow interfaces.MajongFlow) bool {
	mjContext := flow.GetMajongContext()
	playerID := gutils.GetZixunPlayer(mjContext)
	player := utils.GetPlayerByID(mjContext.GetPlayers(), playerID)
	handCard := player.GetHandCards()
	if gutils.CheckHasDingQueCard(mjContext, player) {
		return false
	}
	l := len(handCard)
	if l%3 != 2 {
		return false
	}
	result := utils.CheckHu(handCard, 0, false)
	return result.Can
}

// canZiMo 检查自摸 (判断当前事件是否可行)
func (s *ZiXunState) canZiMo(flow interfaces.MajongFlow, message *majongpb.HuRequestEvent) (bool, error) {
	if gutils.GetZixunPlayer(flow.GetMajongContext()) != message.GetHead().GetPlayerId() {
		return false, fmt.Errorf("当前玩家不允许操作")
	}
	return s.canPlayerZimo(flow), nil
}

func (s *ZiXunState) hasQiangGangHu(flow interfaces.MajongFlow) bool {
	ctx := flow.GetMajongContext()
	card := ctx.GetGangCard()
	cardI := utils.ServerCard2Uint32(card)
	var hasQGanghu bool
	for _, player := range utils.GetCanXpPlayers(ctx.GetPlayers(), ctx) {
		player.PossibleActions = []majongpb.Action{}
		if player.GetPalyerId() != ctx.GetLastGangPlayer() {
			if gutils.CheckHasDingQueCard(flow.GetMajongContext(), player) {
				continue
			}
			result := utils.CheckHu(player.HandCards, cardI, false)
			if result.Can {
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
	mjContext := flow.GetMajongContext()
	zixunNtf := &room.RoomZixunNtf{}
	isNomarlZixun := mjContext.GetZixunType() != majongpb.ZixunType_ZXT_PENG && mjContext.GetZixunType() != majongpb.ZixunType_ZXT_CHI
	playerID := gutils.GetZixunPlayer(mjContext)
	player := utils.GetPlayerByID(mjContext.Players, playerID)
	player.ZixunRecord = &majongpb.ZiXunRecord{}
	record := player.GetZixunRecord()
	if isNomarlZixun {
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
			zixunNtf.EnableChupaiCards = []uint32{utils.ServerCard2Uint32(mjContext.GetLastMopaiCard())}
		}
	} else {
		//3没胡过,所有手牌都可以打
		zixunNtf.EnableChupaiCards = utils.ServerCards2Uint32(player.GetHandCards())
	}
	//查听,打什么,听什么
	s.checkTing(zixunNtf, player, mjContext)
	if zixunNtf.GetEnableZimo() == true {
		//roomHuType, majongHuType := s.getHuType(playerID, mjContext)
		s.checkFanType(record, mjContext, playerID, player.GetHandCards(), s.getHuCard(mjContext, player))
		zixunNtf.HuType = room.HuType(int32(record.HuType)).Enum()
	}
	xpOption := mjoption.GetXingpaiOption(int(mjContext.GetXingpaiOptionId()))
	if len(zixunNtf.GetCanTingCardInfo()) != 0 && !gutils.IsTing(player) && xpOption.EnableTingButton {
		s.addTingButton(mjContext, player, zixunNtf, record)
	}
	s.recordZixunMsg(record, zixunNtf)
	logrus.WithFields(logrus.Fields{
		"EnableAngangCards": record.GetEnableAngangCards(),
		"EnableBugangCards": record.GetEnableBugangCards(),
		"EnableTing":        record.GetEnableTing(),
		"EnableZimo":        record.GetEnableZimo(),
		"EnableQi":          record.GetEnableQi(),
		"EnableChupaiCards": record.GetEnableChupaiCards(),
		"HuType":            record.GetHuType(),
		"TingType":          record.GetTingType(),
		"CanTingCardInfo":   record.GetCanTingCardInfo(),
	}).Infoln("自询记录")
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
		"pengCards": gutils.FmtPengCards(player.GetPengCards()),
		"gangCards": gutils.FmtGangCards(player.GetGangCards()),
		"handCards": gutils.FmtMajongpbCards(player.GetHandCards()),
		"wallCards": gutils.FmtMajongpbCards(mjContext.GetWallCards()),
		"huCards":   gutils.FmtHuCards(player.GetHuCards()),
	}).Infoln("自询通知")
}

func (s *ZiXunState) addTingButton(mjContext *majongpb.MajongContext, player *majongpb.Player, zixunNtf *room.RoomZixunNtf, record *majongpb.ZiXunRecord) {
	zixunNtf.EnableTing = proto.Bool(true)
	if *zixunNtf.EnableTing {
		if int(player.GetZixunCount()) == 1 && s.getOperateCount(player) == 0 {
			zixunNtf.TingType = room.TingType_TT_TIAN_TING.Enum()
			record.TingType = majongpb.TingType_TT_TIAN_TING
		} else {
			zixunNtf.TingType = room.TingType_TT_NORMAL_TING.Enum()
			record.TingType = majongpb.TingType_TT_NORMAL_TING
		}
		zixunNtf.EnableQi = proto.Bool(true)
	}
}

func (s *ZiXunState) getOperateCount(player *majongpb.Player) int {
	return len(player.GetChiCards()) + len(player.GetPengCards()) + len(player.GetGangCards())
}

func (s *ZiXunState) checkFanType(record *majongpb.ZiXunRecord, context *majongpb.MajongContext, huPlayerID uint64, handCards []*majongpb.Card, huCard *majongpb.Card) {
	calcHandCard, _ := utils.RemoveCards(handCards, huCard, 1)
	calcHuCard := &majongpb.HuCard{
		Card:      huCard,
		SrcPlayer: huPlayerID,
		Type:      majongpb.HuType_hu_zimo,
	}
	cardOptionID := int(context.GetCardtypeOptionId())
	fanTypes, genCount, huaCount := fantype.CalculateFanTypes(context, huPlayerID, calcHandCard, calcHuCard)
	record.HuFanType = new(majongpb.HuFanType)
	record.HuFanType.GenCount = uint64(genCount)
	record.HuFanType.HuaCount = uint64(huaCount)
	record.HuType = majongpb.HuType_hu_zimo
	huType := gutils.ServerFanType2ClientHuType(cardOptionID, fanTypes)
	if huType != -1 {
		record.HuType = majongpb.HuType(huType)
	}
	for _, fanType := range fanTypes {
		record.HuFanType.FanTypes = append(record.HuFanType.FanTypes, int64(fanType))
	}
	logrus.WithFields(logrus.Fields{
		"calcHandCard": calcHandCard,
		"calcHuCard":   calcHuCard,
		"huType":       huType,
		"fanTypes":     fanTypes,
		"HfanTypes":    record.HuFanType.FanTypes,
		"genCount":     genCount,
		"huaCount":     huaCount,
	}).Infoln("自询查番")

}

// getHuCard 获取胡牌
func (s *ZiXunState) getHuCard(mjContext *majongpb.MajongContext, player *majongpb.Player) (card *majongpb.Card) {
	// 没有上个摸牌的玩家，是为天胡， 取庄家作为胡牌玩家
	if player.GetZixunCount() == 1 && player.GetPalyerId() == mjContext.Players[int(mjContext.GetZhuangjiaIndex())].GetPalyerId() {
		xpOption := mjoption.GetXingpaiOption(int(mjContext.GetXingpaiOptionId()))
		switch xpOption.TianhuCardType {
		case mjoption.MostTingsCard:
			_, card = utils.CalcTianHuCardNum(mjContext, player.GetPalyerId())
		case mjoption.RightCard:
			card = player.HandCards[len(player.GetHandCards())-1]
		case mjoption.MoCard:
			card = mjContext.GetLastMopaiCard()
		}
	} else {
		card = mjContext.GetLastMopaiCard()
	}
	return
}

func (s *ZiXunState) recordZixunMsg(record *majongpb.ZiXunRecord, ntf *room.RoomZixunNtf) {
	record.EnableAngangCards = ntf.GetEnableAngangCards()
	record.EnableBugangCards = ntf.GetEnableBugangCards()
	record.EnableChupaiCards = ntf.GetEnableChupaiCards()
	record.EnableQi = ntf.GetEnableQi()
	record.EnableZimo = ntf.GetEnableZimo()
	record.EnableTing = ntf.GetEnableTing()
}

// getHuType 计算胡牌类型
func (s *ZiXunState) getHuType(huPlayerID uint64, mjContext *majongpb.MajongContext) (room.HuType, majongpb.HuType) {
	huPlayer := utils.GetMajongPlayer(huPlayerID, mjContext)
	if len(huPlayer.PengCards) == 0 && len(huPlayer.GangCards) == 0 && len(huPlayer.HuCards) == 0 {
		if huPlayer.ZixunCount == 1 && huPlayerID == mjContext.Players[mjContext.ZhuangjiaIndex].GetPalyerId() {
			return room.HuType_HT_TIANHU, majongpb.HuType_hu_tianhu
		}
		if huPlayer.MopaiCount == 1 && huPlayerID != mjContext.Players[mjContext.ZhuangjiaIndex].GetPalyerId() {
			return room.HuType_HT_DIHU, majongpb.HuType_hu_dihu
		}
	}
	return room.HuType_HT_ZIMO, majongpb.HuType_hu_zimo
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
	xpOption := mjoption.GetXingpaiOption(int(context.GetXingpaiOptionId()))

	dingqueCards, checkCards := s.getDingqueCardNum(player, xpOption.EnableDingque)
	dqnum := len(dingqueCards)
	if gutils.IsTing(player) {
		// 当玩家已经是听牌状态，只有摸上来的牌才有听牌提示
		moCard := context.GetLastMopaiCard()
		newTingInfos := map[utils.Card][]utils.Card{}
		tingInfos := utils.GetPlayCardCheckTing(player.GetHandCards(), nil)
		for outCard, tingCard := range tingInfos {
			card, _ := utils.IntToCard(int32(outCard))
			if card.GetColor() == moCard.GetColor() && card.GetPoint() == moCard.GetPoint() {
				newTingInfos[outCard] = tingCard
			}
		}
		//满足条件说明,打出这张定缺牌可以进入听牌状态
		s.addTingInfo(zixunNtf, player, context, newTingInfos)
	} else if dqnum == 0 {
		//没有定缺牌的时候正常查听
		tingInfos := utils.GetPlayCardCheckTing(player.GetHandCards(), nil)
		s.addTingInfo(zixunNtf, player, context, tingInfos)
	} else if dqnum == 1 {
		tingInfos, err := utils.GetTingCards(checkCards, nil)
		if err == nil && len(tingInfos) > 0 {
			dingqueCard := utils.Card(utils.ServerCard2Number(dingqueCards[0]))
			s.addTingInfo(zixunNtf, player, context, map[utils.Card][]utils.Card{dingqueCard: tingInfos})
		}
	}
	// 玩家的定缺牌数量超过1张的时候,不查听
}

func (s *ZiXunState) getDingqueCardNum(player *majongpb.Player, hasDingqueOption bool) (dingqueCards []*majongpb.Card, checkCards []*majongpb.Card) {
	dingqueCards = []*majongpb.Card{}
	checkCards = []*majongpb.Card{}
	if hasDingqueOption {
		for _, card := range player.GetHandCards() {
			if card.GetColor() == player.DingqueColor {
				dingqueCards = append(dingqueCards, card)
			} else {
				checkCards = append(checkCards, card)
			}
		}
	} else {
		checkCards = player.GetHandCards()
	}
	return
}

// addTingInfo 自询通知添加听牌信息
func (s *ZiXunState) addTingInfo(zixunNtf *room.RoomZixunNtf, player *majongpb.Player, context *majongpb.MajongContext, tingInfos map[utils.Card][]utils.Card) {
	canTingInfos := []*room.CanTingCardInfo{}
	recordCanTingInfos := []*majongpb.CanTingCardInfo{}
	for outCard, tingInfo := range tingInfos {
		tingCardInfo := []*room.TingCardInfo{}
		recordTingCardInfo := []*majongpb.TingCardInfo{}
		outCard0, err := utils.IntToCard(int32(outCard))
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"func_name":   "utils.IntToCard",
				"failed_card": outCard,
				"tingInfo":    tingInfo,
			}).Error("牌型转换失败")
		}
		newHand, success := utils.RemoveCards(player.GetHandCards(), outCard0, 1)
		if !success {
			logrus.WithFields(logrus.Fields{
				"func_name": "addTingInfo",
			}).Error("牌型移除失败")
		}
		for _, tt := range tingInfo {
			hCard, _ := utils.IntToCard(int32(tt))
			times, _, _ := facade.CalculateCardValue(global.GetFanTypeCalculator(), context, interfaces.FantypeParams{
				PlayerID:  player.GetPalyerId(),
				MjContext: context,
				HandCard:  newHand,
				PengCard:  s.getPengCards(player.GetPengCards()),
				GangCard:  player.GetGangCards(),
				HuCard: &majongpb.HuCard{
					Card: hCard,
					Type: majongpb.HuType_hu_dianpao,
				},
			})
			tingCardInfo = append(tingCardInfo, &room.TingCardInfo{
				TingCard: proto.Uint32(uint32(tt)),
				Times:    proto.Uint32(uint32(times)),
			})
			recordTingCardInfo = append(recordTingCardInfo, &majongpb.TingCardInfo{
				TingCard: uint32(tt),
				Times:    uint32(times),
			})
		}
		canTingInfos = append(canTingInfos, &room.CanTingCardInfo{
			OutCard:      proto.Uint32(uint32(outCard)),
			TingCardInfo: tingCardInfo,
		})
		recordCanTingInfos = append(recordCanTingInfos, &majongpb.CanTingCardInfo{
			OutCard:      uint32(outCard),
			TingCardInfo: recordTingCardInfo,
		})
	}
	zixunNtf.CanTingCardInfo = canTingInfos
	player.GetZixunRecord().CanTingCardInfo = recordCanTingInfos
}

// checkZiMo 查自摸
func (s *ZiXunState) checkZiMo(flow interfaces.MajongFlow) bool {
	mjContext := flow.GetMajongContext()
	activePlayerID := gutils.GetZixunPlayer(mjContext)
	activePlayer := utils.GetPlayerByID(mjContext.Players, activePlayerID)
	handCard := activePlayer.GetHandCards()
	if gutils.CheckHasDingQueCard(mjContext, activePlayer) {
		return false
	}
	l := len(handCard)
	if l%3 != 2 {
		return false
	}
	result := utils.CheckHu(handCard, 0, false)
	return result.Can
}

func (s *ZiXunState) checkPlayerAngang(player *majongpb.Player, xpOption *mjoption.XingPaiOption) []uint32 {
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
		if xpOption.EnableDingque && k.Color == color {
			continue
		}
		if num == 4 {
			if gutils.IsHu(player) || gutils.IsTing(player) {
				if utils.CheckHuByRemoveGangCards(player, &k, 4) {
					result = append(result, utils.ServerCard2Uint32(&k))
				}
			} else {
				result = append(result, utils.ServerCard2Uint32(&k))
			}
		}
	}
	logrus.WithFields(logrus.Fields{
		"func_name": "ZiXunState.checkPlayerAngang",
		"hand_card": gutils.FmtMajongpbCards(handCard),
		"hu_cards":  gutils.FmtHuCards(huCards),
		"card_num":  cardNum,
		"result":    result,
	}).Debugln("查暗杠")
	return result
}

// canTing 是否可以听
func (s *ZiXunState) canTing(flow interfaces.MajongFlow) bool {
	mjContext := flow.GetMajongContext()
	activePlayerID := gutils.GetZixunPlayer(mjContext)
	activePlayer := utils.GetPlayerByID(mjContext.Players, activePlayerID)
	tingState := activePlayer.GetTingStateInfo()
	if tingState.GetIsTing() || tingState.GetIsTianting() {
		return false
	}
	return true
}

// checkAnGang 查暗杠
func (s *ZiXunState) checkAnGang(flow interfaces.MajongFlow) (enableAngangCards []uint32) {
	enableAngangCards = make([]uint32, 0, 0)
	mjContext := flow.GetMajongContext()
	if !utils.HasAvailableWallCards(flow) {
		return
	}
	activePlayerID := gutils.GetZixunPlayer(mjContext)
	activePlayer := utils.GetPlayerByID(mjContext.Players, activePlayerID)
	xpOption := mjoption.GetXingpaiOption(int(mjContext.GetXingpaiOptionId()))
	return s.checkPlayerAngang(activePlayer, xpOption)
}

//checkBuGang 查补杠
func (s *ZiXunState) checkBuGang(flow interfaces.MajongFlow) []uint32 {
	enableBugangCards := []uint32{}
	mjContext := flow.GetMajongContext()
	// 没有墙牌不能杠
	if !utils.HasAvailableWallCards(flow) {
		return enableBugangCards
	}
	activePlayerID := gutils.GetZixunPlayer(mjContext)
	activePlayer := utils.GetPlayerByID(mjContext.Players, activePlayerID)
	//分两种情况查暗杠，一种是胡牌前，一种胡牌后
	pengCards := activePlayer.GetPengCards()
	for _, touchCard := range activePlayer.HandCards {
		for _, pengCard := range pengCards {
			if *pengCard.Card == *touchCard {
				removeCard := utils.ServerCard2Uint32(touchCard)
				if gutils.IsTing(activePlayer) || gutils.IsHu(activePlayer) {
					if utils.CheckHuByRemoveGangCards(activePlayer, touchCard, 1) {
						enableBugangCards = append(enableBugangCards, removeCard)
					}
				} else {
					enableBugangCards = append(enableBugangCards, removeCard)
				}
			}
		}
	}
	return enableBugangCards
}

func (s *ZiXunState) sortCards(flow interfaces.MajongFlow) {
	mjContext := flow.GetMajongContext()
	playerID := gutils.GetZixunPlayer(mjContext)
	player := utils.GetPlayerByID(mjContext.GetPlayers(), playerID)
	utils.SortCards(player.GetHandCards())
}

//AddZiXunCount 自询次数递增1
func (s *ZiXunState) AddZiXunCount(flow interfaces.MajongFlow) {
	mjContext := flow.GetMajongContext()
	playerID := gutils.GetZixunPlayer(mjContext)
	player := utils.GetPlayerByID(mjContext.GetPlayers(), playerID)
	player.ZixunCount++
}

// OnEntry 进入状态
func (s *ZiXunState) OnEntry(flow interfaces.MajongFlow) {
	s.AddZiXunCount(flow)
	s.sortCards(flow)
	s.checkActions(flow)

}

// OnExit 退出状态 清除本状态数据
func (s *ZiXunState) OnExit(flow interfaces.MajongFlow) {

}
