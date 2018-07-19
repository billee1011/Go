package common

import (
	"steve/gutils"
	"steve/majong/fantype"
	"steve/majong/global"
	"steve/majong/interfaces"
	"steve/majong/settle/majong"
	"steve/majong/utils"
	majongpb "steve/server_pb/majong"

	"github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
)

// QiangGangHuSettleState 枪杠胡结算状态
type QiangGangHuSettleState struct {
}

var _ interfaces.MajongState = new(GangSettleState)

// ProcessEvent 处理事件
// 枪杠胡逻辑执行完后，进入枪杠胡结算状态
// 1.处理结算完成事件，返回摸牌状态
// 2.处理玩家认输事件，返回游戏结束状态
func (s *QiangGangHuSettleState) ProcessEvent(eventID majongpb.EventID, eventContext []byte, flow interfaces.MajongFlow) (newState majongpb.StateID, err error) {
	if eventID == majongpb.EventID_event_settle_finish {
		return s.settleFinishEvent(eventContext, flow)
	}
	return majongpb.StateID(majongpb.StateID_state_qiangganghu_settle), global.ErrInvalidEvent
}

// OnEntry 进入状态
func (s *QiangGangHuSettleState) OnEntry(flow interfaces.MajongFlow) {
	s.doQiangGangHuSettle(flow)
}

// OnExit 退出状态
func (s *QiangGangHuSettleState) OnExit(flow interfaces.MajongFlow) {
}

// setMopaiPlayer 设置摸牌玩家
func (s *QiangGangHuSettleState) setMopaiPlayer(flow interfaces.MajongFlow) {
	mjContext := flow.GetMajongContext()
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "QiangGangHuSettleState.setMopaiPlayer",
	})
	logEntry = utils.WithMajongContext(logEntry, mjContext)
	huPlayers := mjContext.GetLastHuPlayers()
	srcPlayer := mjContext.GetLastMopaiPlayer()
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

// doQiangGangHuSettle 抢杠胡结算
func (s *QiangGangHuSettleState) doQiangGangHuSettle(flow interfaces.MajongFlow) {
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

		cardOptionID := int(mjContext.GetCardtypeOptionId())
		HfanTypes := gutils.GetShowFan(cardOptionID, fanTypes)
		cardTypes[huPlayerID] = HfanTypes
		cardValues[huPlayerID] = totalValue
		genCount[huPlayerID] = uint64(genSum)
		huaCount[huPlayerID] = uint64(huaSum)
		cardsGroup[huPlayerID] = utils.GetCardsGroup(huPlayer, mjContext.GetGangCard())
	}

	params := interfaces.HuSettleParams{
		SettleOptionID: int(mjContext.GetSettleOptionId()),
		HuPlayers:      huPlayers,
		SrcPlayer:      mjContext.GetLastGangPlayer(),
		AllPlayers:     utils.GetAllPlayers(mjContext),
		HasHuPlayers:   utils.GetHuPlayers(mjContext),
		QuitPlayers:    utils.GetQuitPlayers(mjContext),
		GiveupPlayers:  utils.GetGiveupPlayers(mjContext),
		SettleType:     majongpb.SettleType_settle_dianpao,
		HuType:         majongpb.HuType_hu_qiangganghu,
		CardTypes:      cardTypes,
		CardValues:     cardValues,
		GenCount:       genCount,
		HuaCount:       huaCount,
		SettleID:       mjContext.CurrentSettleId,
	}
	settlerFactory := majong.SettlerFactory{}
	settleInfos := settlerFactory.CreateHuSettler().Settle(params)
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

func (s *QiangGangHuSettleState) settleFinishEvent(eventContext []byte, flow interfaces.MajongFlow) (majongpb.StateID, error) {
	message := &majongpb.SettleFinishEvent{}
	err := proto.Unmarshal(eventContext, message)
	if err != nil {
		return majongpb.StateID_state_qiangganghu_settle, global.ErrInvalidEvent
	}
	utils.SettleOver(flow, message)

	nextState := utils.IsGameOverReturnState(flow.GetMajongContext())
	if nextState == majongpb.StateID_state_mopai {
		s.setMopaiPlayer(flow)
	}

	return nextState, nil
}
