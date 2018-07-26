package common

import (
	"steve/common/mjoption"
	"steve/gutils"
	"steve/majong/fantype"
	"steve/majong/global"
	"steve/majong/interfaces"
	"steve/majong/utils"
	majongpb "steve/server_pb/majong"

	"steve/majong/settle"

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
	if !gutils.IsPlayerContinue(mopaiPlayer.GetXpState(), mjContext) {
		mopaiPlayer = utils.GetNextXpPlayerByID(mopaiPlayerID, players, mjContext)
	}
	mjContext.MopaiPlayer = mopaiPlayer.GetPalyerId()
	mjContext.MopaiType = majongpb.MopaiType_MT_NORMAL
}

// doHuSettle 胡的结算
func (s *HuSettleState) doHuSettle(flow interfaces.MajongFlow) {
	mjContext := flow.GetMajongContext()

	cardValues := make(map[uint64]uint64, 0)
	cardTypes := make(map[uint64][]int64, 0)
	genCount := make(map[uint64]uint64, 0)
	huaCount := make(map[uint64]uint64, 0)
	cardsGroup := make(map[uint64][]*majongpb.CardsGroup, 0)
	huPlayers := mjContext.GetLastHuPlayers()
	for _, huPlayerID := range huPlayers {
		huPlayer := utils.GetPlayerByID(mjContext.Players, huPlayerID)
		huCard := huPlayer.GetHuCards()[len(huPlayer.GetHuCards())-1]

		fanTypes, genSum, huaSum := fantype.CalculateFanTypes(mjContext, huPlayerID, huPlayer.GetHandCards(), huCard)
		totalValue := fantype.CalculateScore(mjContext, fanTypes, genSum, huaSum)
		cardsGroup[huPlayerID] = utils.GetCardsGroup(huPlayer, mjContext.GetLastOutCard())

		cardOptionID := int(mjContext.GetCardtypeOptionId())
		HfanTypes := gutils.GetShowFan(cardOptionID, fanTypes)
		cardTypes[huPlayerID] = HfanTypes
		cardValues[huPlayerID] = totalValue
		genCount[huPlayerID] = uint64(genSum)
		huaCount[huPlayerID] = uint64(huaSum)
	}

	params := interfaces.HuSettleParams{
		SettleOptionID: int(mjContext.GetSettleOptionId()),
		HuPlayers:      huPlayers,
		SrcPlayer:      mjContext.GetLastChupaiPlayer(),
		AllPlayers:     utils.GetAllPlayers(mjContext),
		HasHuPlayers:   utils.GetHuPlayers(mjContext, append([]uint64{}, huPlayers...)),
		QuitPlayers:    utils.GetQuitPlayers(mjContext),
		GiveupPlayers:  utils.GetGiveupPlayers(mjContext),
		SettleType:     majongpb.SettleType_settle_dianpao,
		HuType:         majongpb.HuType_hu_dianpao,
		CardTypes:      cardTypes,
		CardValues:     cardValues,
		GenCount:       genCount,
		HuaCount:       huaCount,
		SettleID:       mjContext.CurrentSettleId,
	}
	if s.isAfterGang(mjContext) {
		GangCards := utils.GetMajongPlayer(mjContext.GetLastChupaiPlayer(), mjContext).GangCards
		params.HuType = majongpb.HuType_hu_ganghoupao
		params.GangCard = *GangCards[len(GangCards)-1]
	}
	settlerFactory := settle.SettlerFactory{}
	settleInfos := settlerFactory.CreateHuSettler(mjContext.GameId).Settle(params)
	if s.isAfterGang(mjContext) {
		settleOption := mjoption.GetSettleOption(int(mjContext.GetSettleOptionId()))
		if settleOption.GangInstantSettle {
			lastSettleInfo := mjContext.SettleInfos[len(mjContext.SettleInfos)-1]
			if lastSettleInfo.SettleType == majongpb.SettleType_settle_angang || lastSettleInfo.SettleType == majongpb.SettleType_settle_minggang || lastSettleInfo.SettleType == majongpb.SettleType_settle_bugang {
				lastSettleInfo.CallTransfer = true
			}
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
			huPlayer.MaxCardValue = totalValue
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
	return utils.IsGameOverReturnState(mjcontext)
}
