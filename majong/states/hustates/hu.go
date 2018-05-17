package hustates

import (
	"steve/client_pb/msgId"
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

// HuState 胡状态
// 进入胡状态时， 执行胡操作。设置胡完成事件
// 收到胡完成事件时，设置摸牌玩家，返回摸牌状态
type HuState struct {
}

var _ interfaces.MajongState = new(HuState)

// ProcessEvent 处理事件
func (s *HuState) ProcessEvent(eventID majongpb.EventID, eventContext []byte, flow interfaces.MajongFlow) (newState majongpb.StateID, err error) {
	if eventID == majongpb.EventID_event_hu_finish {
		s.setMopaiPlayer(flow)
		return majongpb.StateID_state_mopai, nil
	}
	return majongpb.StateID_state_hu, global.ErrInvalidEvent
}

// OnEntry 进入状态
func (s *HuState) OnEntry(flow interfaces.MajongFlow) {
	s.doHu(flow)
	flow.SetAutoEvent(majongpb.AutoEvent{
		EventId:      majongpb.EventID_event_hu_finish,
		EventContext: nil,
	})
}

// OnExit 退出状态
func (s *HuState) OnExit(flow interfaces.MajongFlow) {

}

// addHuCard 添加胡的牌
func (s *HuState) addHuCard(card *majongpb.Card, player *majongpb.Player, srcPlayerID uint64) {
	addHuCard(card, player, srcPlayerID, majongpb.HuType_hu_dianpao)
}

// doHu 执行胡操作
func (s *HuState) doHu(flow interfaces.MajongFlow) {
	mjContext := flow.GetMajongContext()

	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "HuState.doHu",
	})
	logEntry = utils.WithMajongContext(logEntry, mjContext)
	players := mjContext.GetLastHuPlayers()

	for _, playerID := range players {
		player := utils.GetMajongPlayer(playerID, mjContext)
		card := mjContext.GetLastOutCard()
		s.addHuCard(card, player, playerID)
	}
	s.notifyHu(flow)
	s.doHuSettle(flow)
	return
}

// isAfterGang 是否为杠后炮
// 杠后摸牌、自询出牌则为杠后炮
func (s *HuState) isAfterGang(mjContext *majongpb.MajongContext) bool {
	cpType := mjContext.GetChupaiType()
	mpType := mjContext.GetMopaiType()
	return mpType == majongpb.MopaiType_MT_GANG && cpType == majongpb.ChupaiType_CT_ZIXUN
}

// HuState 广播胡
func (s *HuState) notifyHu(flow interfaces.MajongFlow) {
	mjContext := flow.GetMajongContext()
	huType := room.HuType_DianPao.Enum()
	if s.isAfterGang(mjContext) {
		huType = room.HuType_GangouPao.Enum()
	}
	body := room.RoomHuNtf{
		Players:      mjContext.GetLastHuPlayers(),
		FromPlayerId: proto.Uint64(mjContext.GetLastChupaiPlayer()),
		Card:         proto.Uint32(uint32(utils.ServerCard2Number(mjContext.GetLastOutCard()))),
		HuType:       huType,
	}
	facade.BroadcaseMessage(flow, msgid.MsgID_ROOM_HU_NTF, &body)
}

// setMopaiPlayer 设置摸牌玩家
func (s *HuState) setMopaiPlayer(flow interfaces.MajongFlow) {
	mjContext := flow.GetMajongContext()
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "QiangganghuState.setMopaiPlayer",
	})
	logEntry = utils.WithMajongContext(logEntry, mjContext)
	huPlayers := mjContext.GetLastHuPlayers()
	srcPlayer := mjContext.GetLastChupaiPlayer()
	players := mjContext.GetPlayers()

	mjContext.MopaiPlayer = calcMopaiPlayer(logEntry, huPlayers, srcPlayer, players)
	mjContext.MopaiType = majongpb.MopaiType_MT_NORMAL
}

// doHuSettle 胡的结算
func (s *HuState) doHuSettle(flow interfaces.MajongFlow) {
	mjContext := flow.GetMajongContext()

	allPlayers := make([]uint64, 0)
	for _, player := range mjContext.Players {
		allPlayers = append(allPlayers, player.GetPalyerId())
	}

	cardValues := make(map[uint64]uint32, 0)
	cardTypes := make(map[uint64][]majongpb.CardType, 0)
	genCount := make(map[uint64]uint32, 0)

	huPlayers := mjContext.GetLastHuPlayers()
	for _, huPlayerID := range huPlayers {
		huPlayer := utils.GetPlayerByID(mjContext.Players, huPlayerID)
		cardParams := interfaces.CardCalcParams{
			HandCard: huPlayer.HandCards,
			PengCard: utils.TransPengCard(huPlayer.PengCards),
			GangCard: utils.TransGangCard(huPlayer.GangCards),
			HuCard:   mjContext.GetLastMopaiCard(),
		}
		calculator := new(cardtype.ScxlCardTypeCalculator)
		cardType, gen := calculator.Calculate(cardParams)
		cardValue, _ := calculator.CardTypeValue(cardType, gen)

		cardTypes[huPlayerID] = cardType
		cardValues[huPlayerID] = cardValue
		genCount[huPlayerID] = gen
	}

	huType := majongpb.SettleHuType_settle_hu_noramaldianpao

	params := interfaces.HuSettleParams{
		HuPlayers:  huPlayers,
		SrcPlayer:  mjContext.GetLastChupaiPlayer(),
		AllPlayers: allPlayers,
		SettleType: majongpb.SettleType_settle_dianpao,
		HuType:     huType,
		CardTypes:  cardTypes,
		CardValues: cardValues,
		GenCount:   genCount,
		SettleID:   mjContext.CurrentSettleId,
	}
	if s.isAfterGang(mjContext) {
		huType = majongpb.SettleHuType_settle_hu_ganghoupao
		GangCards := utils.GetMajongPlayer(mjContext.GetLastChupaiPlayer(), mjContext).GangCards
		params.HuType = huType
		params.GangCard = *GangCards[len(GangCards)-1]
	}
	huSettle := new(settle.HuSettle)
	settleInfos := huSettle.Settle(params)
	for _, settleInfo := range settleInfos {
		mjContext.SettleInfos = append(mjContext.SettleInfos, settleInfo)
		mjContext.CurrentSettleId++
	}
}
