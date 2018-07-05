package common

import (
	"steve/majong/global"
	"steve/majong/interfaces"
	"steve/majong/interfaces/facade"
	"steve/majong/utils"
	majongpb "steve/server_pb/majong"

	"github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
)

// HuSettleState 杠结算状态
type HuSettleState struct {
}

var _ interfaces.MajongState = new(HuSettleState)

// ProcessEvent 处理事件
// 点炮逻辑执行完后，进入点炮结算状态
// 1.处理结算完成事件，返回摸牌状态
// 2.处理玩家认输事件，返回游戏结束状态
func (s *HuSettleState) ProcessEvent(eventID majongpb.EventID, eventContext []byte, flow interfaces.MajongFlow) (newState majongpb.StateID, err error) {
	if eventID == majongpb.EventID_event_settle_finish {
		message := &majongpb.SettleFinishEvent{}
		err := proto.Unmarshal(eventContext, message)
		if err != nil {
			return majongpb.StateID_state_hu_settle, global.ErrInvalidEvent
		}
		utils.SettleOver(flow, message)
		nextState := s.nextState(flow.GetMajongContext())
		if nextState == majongpb.StateID_state_mopai {
			s.setMopaiPlayer(flow)
		}
		logrus.WithFields(logrus.Fields{
			"func_name": "HuSettleState.ProcessEvent",
			"nextState": nextState,
		}).Infoln("点炮结算下个状态")
		return nextState, nil
	}
	return majongpb.StateID(majongpb.StateID_state_hu_settle), global.ErrInvalidEvent
}

// OnEntry 进入状态
func (s *HuSettleState) OnEntry(flow interfaces.MajongFlow) {
	s.doHuSettle(flow)
}

// OnExit 退出状态
func (s *HuSettleState) OnExit(flow interfaces.MajongFlow) {
}

// setMopaiPlayer 设置摸牌玩家
func (s *HuSettleState) setMopaiPlayer(flow interfaces.MajongFlow) {
	mjContext := flow.GetMajongContext()
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "huSettleState.setMopaiPlayer",
	})
	logEntry = utils.WithMajongContext(logEntry, mjContext)
	huPlayers := mjContext.GetLastHuPlayers()
	srcPlayer := mjContext.GetLastChupaiPlayer()
	players := mjContext.GetPlayers()

	mopaiPlayerID := CalcMopaiPlayer(logEntry, huPlayers, srcPlayer, players)
	// 摸牌玩家不能是非正常状态玩家
	mopaiPlayer := utils.GetPlayerByID(players, mopaiPlayerID)
	if !utils.IsPlayerContinue(mopaiPlayer.GetXpState(), mjContext.GetOption()) {
		mopaiPlayer = utils.GetNextXpPlayerByID(mopaiPlayerID, players, mjContext.GetOption())
	}
	mjContext.MopaiPlayer = mopaiPlayer.GetPalyerId()
	mjContext.MopaiType = majongpb.MopaiType_MT_NORMAL
}

// doHuSettle 胡的结算
func (s *HuSettleState) doHuSettle(flow interfaces.MajongFlow) {
	mjContext := flow.GetMajongContext()

	allPlayers := make([]uint64, 0)
	for _, player := range mjContext.Players {
		allPlayers = append(allPlayers, player.GetPalyerId())
	}

	cardValues := make(map[uint64]uint32, 0)
	cardTypes := make(map[uint64][]majongpb.CardType, 0)
	genCount := make(map[uint64]uint32, 0)
	cardsGroup := make(map[uint64][]*majongpb.CardsGroup, 0)

	huPlayers := mjContext.GetLastHuPlayers()
	gameID := int(mjContext.GetGameId())
	for _, huPlayerID := range huPlayers {
		huPlayer := utils.GetPlayerByID(mjContext.Players, huPlayerID)
		cardParams := interfaces.CardCalcParams{
			HandCard: huPlayer.HandCards,
			PengCard: utils.TransPengCard(huPlayer.PengCards),
			GangCard: utils.TransGangCard(huPlayer.GangCards),
			HuCard:   mjContext.GetLastOutCard(),
			GameID:   gameID,
		}
		calculator := global.GetCardTypeCalculator()
		cardType, gen := calculator.Calculate(cardParams)
		cardValue, _ := calculator.CardTypeValue(gameID, cardType, gen)

		cardTypes[huPlayerID] = cardType
		cardValues[huPlayerID] = cardValue
		genCount[huPlayerID] = gen
		cardsGroup[huPlayerID] = utils.GetCardsGroup(huPlayer, mjContext.GetLastOutCard())
	}

	huType := majongpb.HuType_hu_dianpao

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
		huType = majongpb.HuType_hu_ganghoupao
		GangCards := utils.GetMajongPlayer(mjContext.GetLastChupaiPlayer(), mjContext).GangCards
		params.HuType = huType
		params.GangCard = *GangCards[len(GangCards)-1]
	}
	settleInfos := facade.SettleHu(global.GetGameSettlerFactory(), int(mjContext.GetGameId()), params)
	if s.isAfterGang(mjContext) {
		lastSettleInfo := mjContext.SettleInfos[len(mjContext.SettleInfos)-1]
		if lastSettleInfo.SettleType == majongpb.SettleType_settle_angang || lastSettleInfo.SettleType == majongpb.SettleType_settle_minggang || lastSettleInfo.SettleType == majongpb.SettleType_settle_bugang {
			lastSettleInfo.CallTransfer = true
		}
	}
	maxSID := uint64(0)
	totalValue := uint32(0)
	for _, settleInfo := range settleInfos {
		mjContext.SettleInfos = append(mjContext.SettleInfos, settleInfo)
		if settleInfo.Id > maxSID {
			maxSID = settleInfo.Id
		}
		totalValue = settleInfo.CardValue
	}
	for _, huPlayerID := range huPlayers {
		huPlayer := utils.GetPlayerByID(mjContext.Players, huPlayerID)
		if totalValue > huPlayer.MaxCardValue {
			huPlayer.CardsGroup = cardsGroup[huPlayerID]
		}
	}
	mjContext.CurrentSettleId = maxSID
}

// isAfterGang 是否为杠后炮
// 杠后摸牌、自询出牌则为杠后炮
func (s *HuSettleState) isAfterGang(mjContext *majongpb.MajongContext) bool {
	zxType := mjContext.GetZixunType()
	mpType := mjContext.GetMopaiType()
	return mpType == majongpb.MopaiType_MT_GANG && zxType == majongpb.ZixunType_ZXT_NORMAL
}

// nextState 下个状态
func (s *HuSettleState) nextState(mjcontext *majongpb.MajongContext) majongpb.StateID {
	return utils.GetNextState(mjcontext)
}
