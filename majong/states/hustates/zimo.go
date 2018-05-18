package hustates

import (
	msgid "steve/client_pb/msgId"
	"steve/client_pb/room"
	"steve/majong/cardtype"
	"steve/majong/global"
	"steve/majong/interfaces"
	"steve/majong/interfaces/facade"
	"steve/majong/settle"
	"steve/majong/utils"
	majongpb "steve/server_pb/majong"

	"github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
)

// ZimoState 自摸状态
// 进入状态时，执行自摸动作，并广播给玩家
// 自摸完成事件，进入下家摸牌状态
type ZimoState struct {
}

var _ interfaces.MajongState = new(ZimoState)

// ProcessEvent 处理事件
func (s *ZimoState) ProcessEvent(eventID majongpb.EventID, eventContext []byte, flow interfaces.MajongFlow) (newState majongpb.StateID, err error) {
	if eventID == majongpb.EventID_event_zimo_finish {
		s.setMopaiPlayer(flow)
		return majongpb.StateID_state_mopai, nil
	}
	return majongpb.StateID_state_zimo, global.ErrInvalidEvent
}

// OnEntry 进入状态
func (s *ZimoState) OnEntry(flow interfaces.MajongFlow) {
	s.doZimo(flow)
	flow.SetAutoEvent(majongpb.AutoEvent{
		EventId:      majongpb.EventID_event_zimo_finish,
		EventContext: nil,
	})
}

// OnExit 退出状态
func (s *ZimoState) OnExit(flow interfaces.MajongFlow) {

}

// doZimo 执行自摸操作
func (s *ZimoState) doZimo(flow interfaces.MajongFlow) {
	mjContext := flow.GetMajongContext()

	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "ZimoState.doZimo",
	})
	logEntry = utils.WithMajongContext(logEntry, mjContext)
	player, card, err := s.getZimoInfo(mjContext)
	if err != nil {
		logEntry.Errorln(err)
		return
	}
	mjContext.LastHuPlayers = []uint64{player.GetPalyerId()}

	player.HandCards, _ = utils.RemoveCards(player.GetHandCards(), card, 1)
	addHuCard(card, player, player.GetPalyerId(), majongpb.HuType_hu_zimo)
	s.notifyHu(card, player.GetPalyerId(), flow)
	s.doZiMoSettle(card, player.GetPalyerId(), flow)
}

// isAfterGang 判断是否为杠开
// 杠后摸牌则为杠开
func (s *ZimoState) isAfterGang(mjContext *majongpb.MajongContext) bool {
	return mjContext.GetMopaiType() == majongpb.MopaiType_MT_GANG
}

// calcHuType 计算胡牌类型
func (s *ZimoState) calcHuType(huPlayerID uint64, mjContext *majongpb.MajongContext) room.HuType {
	afterGang := s.isAfterGang(mjContext)
	isLast := (len(mjContext.WallCards) == 0)
	if afterGang && isLast {
		return room.HuType_HT_GANGSHANGHAIDILAO
	} else if afterGang {
		return room.HuType_HT_GANGKAI
	} else if isLast {
		return room.HuType_HT_HAIDILAO
	}
	huPlayer := utils.GetMajongPlayer(huPlayerID, mjContext)
	if huPlayer.MopaiCount == 0 && len(huPlayer.PengCards) == 0 && len(huPlayer.GangCards) == 0 {
		if huPlayerID == mjContext.Players[mjContext.ZhuangjiaIndex].GetPalyerId() {
			return room.HuType_HT_TIANHU
		}
		return room.HuType_HT_DIHU
	}
	return room.HuType_HT_ZIMO
}

// notifyHu 广播胡
func (s *ZimoState) notifyHu(card *majongpb.Card, playerID uint64, flow interfaces.MajongFlow) {
	mjContext := flow.GetMajongContext()
	huType := s.calcHuType(playerID, mjContext)
	body := room.RoomHuNtf{
		Players:      []uint64{playerID},
		FromPlayerId: proto.Uint64(playerID),
		Card:         proto.Uint32(uint32(utils.ServerCard2Number(card))),
		HuType:       huType.Enum(),
	}
	facade.BroadcaseMessage(flow, msgid.MsgID_ROOM_HU_NTF, &body)
}

// setMopaiPlayer 设置摸牌玩家
func (s *ZimoState) setMopaiPlayer(flow interfaces.MajongFlow) {
	mjContext := flow.GetMajongContext()
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "ZimoState.doZimo",
	})
	logEntry = utils.WithMajongContext(logEntry, mjContext)

	huPlayers := mjContext.GetLastHuPlayers()
	if len(huPlayers) == 0 {
		logEntry.Errorln("胡牌玩家列表为空")
		return
	}
	players := mjContext.GetPlayers()
	mjContext.MopaiPlayer = calcMopaiPlayer(logEntry, huPlayers, huPlayers[0], players)
	mjContext.MopaiType = majongpb.MopaiType_MT_NORMAL
}

// getZimoInfo 获取自摸信息
func (s *ZimoState) getZimoInfo(mjContext *majongpb.MajongContext) (player *majongpb.Player, card *majongpb.Card, err error) {
	playerID := mjContext.GetLastMopaiPlayer()
	players := mjContext.GetPlayers()
	player = utils.GetPlayerByID(players, playerID)

	// 没有上个摸牌的玩家，是为天胡， 取庄家作为胡牌玩家
	if player.GetMopaiCount() == 0 {
		card = s.calcTianhuCard(player.GetHandCards())
	} else {
		card = mjContext.GetLastMopaiCard()
	}
	mjContext.LastHuPlayers = []uint64{playerID}
	return
}

// calcTianhuCard 计算天胡胡的牌
func (s *ZimoState) calcTianhuCard(cards []*majongpb.Card) *majongpb.Card {
	// TODO
	return cards[len(cards)-1]
}

func (s *ZimoState) huType2SettleType(huType room.HuType) majongpb.SettleHuType {
	return map[room.HuType]majongpb.SettleHuType{
		room.HuType_HT_ZIMO:              majongpb.SettleHuType_settle_hu_zimo,
		room.HuType_HT_GANGKAI:           majongpb.SettleHuType_settle_hu_gangkai,
		room.HuType_HT_HAIDILAO:          majongpb.SettleHuType_settle_hu_haidilao,
		room.HuType_HT_GANGSHANGHAIDILAO: majongpb.SettleHuType_settle_hu_gangshanghaidilao,
	}[huType]
}

// doZiMoSettle 自摸的结算
func (s *ZimoState) doZiMoSettle(card *majongpb.Card, huPlayerID uint64, flow interfaces.MajongFlow) {
	mjContext := flow.GetMajongContext()

	allPlayers := make([]uint64, 0)
	for _, player := range mjContext.Players {
		allPlayers = append(allPlayers, player.GetPalyerId())
	}

	cardValues := make(map[uint64]uint32, 0)
	cardTypes := make(map[uint64][]majongpb.CardType, 0)
	genCount := make(map[uint64]uint32, 0)

	huPlayer := utils.GetPlayerByID(mjContext.Players, huPlayerID)
	cardParams := interfaces.CardCalcParams{
		HandCard: huPlayer.HandCards,
		PengCard: utils.TransPengCard(huPlayer.PengCards),
		GangCard: utils.TransGangCard(huPlayer.GangCards),
		HuCard:   nil,
	}
	calculator := new(cardtype.ScxlCardTypeCalculator)
	cardType, gen := calculator.Calculate(cardParams)
	cardValue, _ := calculator.CardTypeValue(cardType, gen)

	cardTypes[huPlayerID] = cardType
	cardValues[huPlayerID] = cardValue
	genCount[huPlayerID] = gen

	huType := s.calcHuType(huPlayerID, mjContext)
	settleType := s.huType2SettleType(huType)
	params := interfaces.HuSettleParams{
		HuPlayers:  []uint64{huPlayerID},
		SrcPlayer:  huPlayerID,
		AllPlayers: allPlayers,
		SettleType: majongpb.SettleType_settle_zimo,
		HuType:     settleType,
		CardTypes:  cardTypes,
		CardValues: cardValues,
		GenCount:   genCount,
		SettleID:   mjContext.CurrentSettleId,
	}
	huSettle := new(settle.HuSettle)
	settleInfos := huSettle.Settle(params)
	for _, settleInfo := range settleInfos {
		mjContext.SettleInfos = append(mjContext.SettleInfos, settleInfo)
		mjContext.CurrentSettleId++
	}
}
