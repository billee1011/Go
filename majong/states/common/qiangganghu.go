package common

import (
	msgid "steve/client_pb/msgId"
	"steve/client_pb/room"
	"steve/majong/global"
	"steve/majong/interfaces"
	"steve/majong/interfaces/facade"
	"steve/majong/utils"
	majongpb "steve/server_pb/majong"

	"github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
)

// QiangganghuState 抢杠胡状态
// 执行抢杠胡操作，并广播
// 从上个摸牌的玩家算起，最后胡的玩家的下家摸牌
type QiangganghuState struct {
}

var _ interfaces.MajongState = new(QiangganghuState)

// ProcessEvent 处理事件
func (s *QiangganghuState) ProcessEvent(eventID majongpb.EventID, eventContext []byte, flow interfaces.MajongFlow) (newState majongpb.StateID, err error) {
	if eventID == majongpb.EventID_event_qiangganghu_finish {
		s.setMopaiPlayer(flow)
		return majongpb.StateID_state_mopai, nil
	}
	return majongpb.StateID_state_qiangganghu, global.ErrInvalidEvent
}

// OnEntry 进入状态
func (s *QiangganghuState) OnEntry(flow interfaces.MajongFlow) {
	s.doHu(flow)
	s.doQiangGangHuSettle(flow)
	flow.SetAutoEvent(majongpb.AutoEvent{
		EventId:      majongpb.EventID_event_qiangganghu_finish,
		EventContext: nil,
	})
}

// OnExit 退出状态
func (s *QiangganghuState) OnExit(flow interfaces.MajongFlow) {

}

// setMopaiPlayer 设置摸牌玩家
func (s *QiangganghuState) setMopaiPlayer(flow interfaces.MajongFlow) {
	mjContext := flow.GetMajongContext()
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "QiangganghuState.setMopaiPlayer",
	})
	logEntry = utils.WithMajongContext(logEntry, mjContext)
	huPlayers := mjContext.GetLastHuPlayers()
	srcPlayer := mjContext.GetLastMopaiPlayer()
	players := mjContext.GetPlayers()

	mjContext.MopaiPlayer = calcMopaiPlayer(logEntry, huPlayers, srcPlayer, players)
	mjContext.MopaiType = majongpb.MopaiType_MT_NORMAL
}

// addHuCard 添加胡的牌
func (s *QiangganghuState) addHuCard(card *majongpb.Card, player *majongpb.Player, srcPlayerID uint64, isReal bool) {
	addHuCard(card, player, srcPlayerID, majongpb.HuType_hu_dianpao, isReal)
}

func (s *QiangganghuState) removeSrcCard(card *majongpb.Card, srcPlayer *majongpb.Player) {
	var succ bool
	srcPlayer.HandCards, succ = utils.RemoveCards(srcPlayer.GetHandCards(), card, 1)
	if !succ {
		logrus.WithFields(logrus.Fields{
			"func_name":      "QiangganghuState.removeSrcCard",
			"hand_cards":     srcPlayer.GetHandCards(),
			"gang_player_id": srcPlayer.GetPalyerId(),
		}).Errorln("移除杠者的杠牌失败")
	}
}

// doHu 执行胡操作
func (s *QiangganghuState) doHu(flow interfaces.MajongFlow) {
	mjContext := flow.GetMajongContext()

	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "QiangganghuState.doHu",
	})
	logEntry = utils.WithMajongContext(logEntry, mjContext)
	players := mjContext.GetLastHuPlayers()
	srcPlayerID := mjContext.GetLastGangPlayer()
	srcPlayer := utils.GetPlayerByID(mjContext.GetPlayers(), srcPlayerID)
	card := mjContext.GetGangCard() // 杠的牌为抢杠胡的牌

	isReal := true
	for _, playerID := range players {
		player := utils.GetMajongPlayer(playerID, mjContext)
		s.addHuCard(card, player, srcPlayerID, isReal)
		isReal = false
	}
	s.removeSrcCard(card, srcPlayer)
	s.notifyHu(flow)
	return
}

// QiangganghuState 广播胡
func (s *QiangganghuState) notifyHu(flow interfaces.MajongFlow) {
	mjContext := flow.GetMajongContext()
	card := mjContext.GetGangCard()
	body := room.RoomHuNtf{
		Players:      mjContext.GetLastHuPlayers(),
		FromPlayerId: proto.Uint64(mjContext.GetLastMopaiPlayer()),
		Card:         proto.Uint32(uint32(utils.ServerCard2Number(card))),
		HuType:       room.HuType_HT_QIANGGANGHU.Enum(),
	}
	facade.BroadcaseMessage(flow, msgid.MsgID_ROOM_HU_NTF, &body)
}

// doQiangGangHuSettle 抢杠胡结算
func (s *QiangganghuState) doQiangGangHuSettle(flow interfaces.MajongFlow) {
	mjContext := flow.GetMajongContext()

	allPlayers := make([]uint64, 0)
	for _, player := range mjContext.Players {
		allPlayers = append(allPlayers, player.GetPalyerId())
	}

	cardValues := make(map[uint64]uint32, 0)
	cardTypes := make(map[uint64][]majongpb.CardType, 0)
	genCount := make(map[uint64]uint32, 0)
	gameID := int(mjContext.GetGameId())

	huPlayers := mjContext.GetLastHuPlayers()
	for _, huPlayerID := range huPlayers {
		huPlayer := utils.GetPlayerByID(mjContext.Players, huPlayerID)
		cardParams := interfaces.CardCalcParams{
			HandCard: huPlayer.HandCards,
			PengCard: utils.TransPengCard(huPlayer.PengCards),
			GangCard: utils.TransGangCard(huPlayer.GangCards),
			HuCard:   mjContext.GetGangCard(),
			GameID:   gameID,
		}
		calculator := global.GetCardTypeCalculator()
		cardType, gen := calculator.Calculate(cardParams)
		cardValue, _ := calculator.CardTypeValue(gameID, cardType, gen)

		cardTypes[huPlayerID] = cardType
		cardValues[huPlayerID] = cardValue
		genCount[huPlayerID] = gen
	}

	params := interfaces.HuSettleParams{
		HuPlayers:  huPlayers,
		SrcPlayer:  mjContext.GetLastGangPlayer(),
		AllPlayers: allPlayers,
		SettleType: majongpb.SettleType_settle_dianpao,
		HuType:     majongpb.SettleHuType_settle_hu_qiangganghu,
		CardTypes:  cardTypes,
		CardValues: cardValues,
		GenCount:   genCount,
		SettleID:   mjContext.CurrentSettleId,
	}
	settleInfos := facade.SettleHu(global.GetGameSettlerFactory(), int(mjContext.GetGameId()), params)
	maxSID := uint64(0)
	for _, settleInfo := range settleInfos {
		mjContext.SettleInfos = append(mjContext.SettleInfos, settleInfo)
		if settleInfo.Id > maxSID {
			maxSID = settleInfo.Id
		}
	}
	mjContext.CurrentSettleId = maxSID
}