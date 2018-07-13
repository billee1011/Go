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
			return majongpb.StateID_state_zimo_settle, global.ErrInvalidEvent
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
	return majongpb.StateID(majongpb.StateID_state_zimo_settle), global.ErrInvalidEvent
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
	if !gutils.IsPlayerContinue(mopaiPlayer.GetXpState(), mjContext) {
		mopaiPlayer = utils.GetNextXpPlayerByID(mopaiPlayerID, players, mjContext)
	}
	mjContext.MopaiPlayer = mopaiPlayer.GetPalyerId()
	mjContext.MopaiType = majongpb.MopaiType_MT_NORMAL
}

// doZiMoSettle 自摸的结算
func (s *ZiMoSettleState) doZiMoSettle(flow interfaces.MajongFlow) {
	mjContext := flow.GetMajongContext()

	huPlayerID := mjContext.GetLastMopaiPlayer()
	huPlayer := utils.GetPlayerByID(mjContext.Players, huPlayerID)
	huCard := huPlayer.HuCards[len(huPlayer.HuCards)-1]

	cardValues := make(map[uint64]uint64, 0)
	cardTypes := make(map[uint64][]int64, 0)
	genCount := make(map[uint64]uint64, 0)
	huaCount := make(map[uint64]uint64, 0)

	record := huPlayer.GetZixunRecord()
	cardTypes[huPlayerID] = record.GetHuFanType().GetFanTypes()
	cardValues[huPlayerID] = s.calculateScore(mjContext, record)
	genCount[huPlayerID] = record.GetHuFanType().GetGenCount()
	huaCount[huPlayerID] = record.GetHuFanType().GetHuaCount()

	params := interfaces.HuSettleParams{
		SettleOptionID: int(mjContext.GetSettleOptionId()),
		HuPlayers:      []uint64{huPlayerID},
		SrcPlayer:      huPlayerID,
		AllPlayers:     utils.GetAllPlayers(mjContext),
		HasHuPlayers:   utils.GetHuPlayers(mjContext),
		QuitPlayers:    utils.GetQuitPlayers(mjContext),
		GiveupPlayers:  utils.GetGiveupPlayers(mjContext),
		SettleType:     majongpb.SettleType_settle_zimo,
		HuType:         huCard.GetType(),
		CardTypes:      cardTypes,
		CardValues:     cardValues,
		GenCount:       genCount,
		HuaCount:       huaCount,
		SettleID:       mjContext.CurrentSettleId,
	}
	totalValue := uint32(0)
	settlerFactory := majong.SettlerFactory{}
	settleInfos := settlerFactory.CreateHuSettler().Settle(params)
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

func (s *ZiMoSettleState) calculateScore(mjcontext *majongpb.MajongContext, record *majongpb.ZiXunRecord) uint64 {
	hufanType := record.GetHuFanType()
	fanTypes := make([]int, 0)
	for _, fType := range hufanType.GetFanTypes() {
		fanTypes = append(fanTypes, int(fType))
	}
	return fantype.CalculateScore(mjcontext, fanTypes, int(hufanType.GetGenCount()), int(hufanType.GetHuaCount()))
}

// nextState 下个状态
func (s *ZiMoSettleState) nextState(mjcontext *majongpb.MajongContext) majongpb.StateID {
	return utils.IsGameOverReturnState(mjcontext)
}
