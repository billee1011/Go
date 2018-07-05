package common

import (
	"steve/gutils"
	"steve/majong/global"
	"steve/majong/interfaces"
	"steve/majong/interfaces/facade"
	"steve/majong/utils"
	majongpb "steve/server_pb/majong"

	"github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
)

// ZiMoSettleState 自摸结算状态
type ZiMoSettleState struct {
}

var _ interfaces.MajongState = new(ZiMoSettleState)

// ProcessEvent 处理事件
// 自摸逻辑执行完后，进入自摸结算状态
func (s *ZiMoSettleState) ProcessEvent(eventID majongpb.EventID, eventContext []byte, flow interfaces.MajongFlow) (newState majongpb.StateID, err error) {
	if eventID == majongpb.EventID_event_settle_finish {
		message := &majongpb.SettleFinishEvent{}
		err := proto.Unmarshal(eventContext, message)
		if err != nil {
			return majongpb.StateID_state_gang_settle, global.ErrInvalidEvent
		}
		utils.SettleOver(flow, message)
		nextState := s.nextState(flow.GetMajongContext())
		if nextState == majongpb.StateID_state_mopai {
			s.setMopaiPlayer(flow)
		}
		logrus.WithFields(logrus.Fields{
			"func_name": "ZiMoSettleState.ProcessEvent",
			"nextState": nextState,
		}).Infoln("自摸结算下个状态")
		return nextState, nil
	}
	return majongpb.StateID(majongpb.StateID_state_gang_settle), global.ErrInvalidEvent
}

// OnEntry 进入状态
func (s *ZiMoSettleState) OnEntry(flow interfaces.MajongFlow) {
	s.doZiMoSettle(flow)
}

// OnExit 退出状态
func (s *ZiMoSettleState) OnExit(flow interfaces.MajongFlow) {
}

// setMopaiPlayer 设置摸牌玩家
func (s *ZiMoSettleState) setMopaiPlayer(flow interfaces.MajongFlow) {
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
	mopaiPlayerID := CalcMopaiPlayer(logEntry, huPlayers, huPlayers[0], players)
	// 摸牌玩家不能是非正常状态玩家
	mopaiPlayer := utils.GetPlayerByID(players, mopaiPlayerID)
	if !utils.IsPlayerContinue(mopaiPlayer.GetXpState(), mjContext.GetOption()) {
		mopaiPlayer = utils.GetNextXpPlayerByID(mopaiPlayerID, players, mjContext.GetOption())
	}
	mjContext.MopaiPlayer = mopaiPlayer.GetPalyerId()
	mjContext.MopaiType = majongpb.MopaiType_MT_NORMAL
}

// doZiMoSettle 自摸的结算
func (s *ZiMoSettleState) doZiMoSettle(flow interfaces.MajongFlow) {
	mjContext := flow.GetMajongContext()

	huPlayerID := mjContext.GetLastMopaiPlayer()

	allPlayers := make([]uint64, 0)
	for _, player := range mjContext.Players {
		if mjContext.GetGameId() == gutils.SCXZGameID {
			if len(player.HuCards) == 0 || huPlayerID == player.GetPalyerId() {
				allPlayers = append(allPlayers, player.GetPalyerId())
			}
		} else {
			allPlayers = append(allPlayers, player.GetPalyerId())
		}
	}

	cardValues := make(map[uint64]uint32, 0)
	cardTypes := make(map[uint64][]majongpb.CardType, 0)
	genCount := make(map[uint64]uint32, 0)
	gameID := int(mjContext.GetGameId())
	huPlayer := utils.GetPlayerByID(mjContext.Players, huPlayerID)
	huCard := huPlayer.HuCards[len(huPlayer.HuCards)-1]
	cardParams := interfaces.CardCalcParams{
		HandCard: append(huPlayer.HandCards, huCard.GetCard()),
		PengCard: utils.TransPengCard(huPlayer.PengCards),
		GangCard: utils.TransGangCard(huPlayer.GangCards),
		HuCard:   nil,
		GameID:   gameID,
	}
	calculator := global.GetCardTypeCalculator()
	cardType, gen := calculator.Calculate(cardParams)
	cardValue, _ := calculator.CardTypeValue(gameID, cardType, gen)

	cardTypes[huPlayerID] = cardType
	cardValues[huPlayerID] = cardValue
	genCount[huPlayerID] = gen

	params := interfaces.HuSettleParams{
		HuPlayers:  []uint64{huPlayerID},
		SrcPlayer:  huPlayerID,
		AllPlayers: allPlayers,
		SettleType: majongpb.SettleType_settle_zimo,
		HuType:     huCard.GetType(),
		CardTypes:  cardTypes,
		CardValues: cardValues,
		GenCount:   genCount,
		SettleID:   mjContext.CurrentSettleId,
	}
	settleInfos := facade.SettleHu(global.GetGameSettlerFactory(), int(mjContext.GetGameId()), params)
	totalValue := uint32(0)

	for _, settleInfo := range settleInfos {
		mjContext.SettleInfos = append(mjContext.SettleInfos, settleInfo)
		mjContext.CurrentSettleId++
		totalValue = settleInfo.CardValue
	}
	if totalValue > huPlayer.MaxCardValue {
		huPlayer.CardsGroup = utils.GetCardsGroup(huPlayer, huCard.Card)
		huPlayer.MaxCardValue = totalValue
	}
}

// nextState 下个状态
func (s *ZiMoSettleState) nextState(mjcontext *majongpb.MajongContext) majongpb.StateID {
	return utils.GetNextState(mjcontext)
}
