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
	if !utils.IsPlayerContinue(mopaiPlayer.GetXpState(), mjContext.GetOption()) {
		mopaiPlayer = utils.GetNextXpPlayerByID(mopaiPlayerID, players, mjContext.GetOption())
	}
	mjContext.MopaiPlayer = mopaiPlayer.GetPalyerId()
	mjContext.MopaiType = majongpb.MopaiType_MT_NORMAL
}

// doQiangGangHuSettle 抢杠胡结算
func (s *QiangGangHuSettleState) doQiangGangHuSettle(flow interfaces.MajongFlow) {
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
		HuType:     majongpb.HuType_hu_qiangganghu,
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

func (s *QiangGangHuSettleState) settleFinishEvent(eventContext []byte, flow interfaces.MajongFlow) (majongpb.StateID, error) {
	message := &majongpb.SettleFinishEvent{}
	err := proto.Unmarshal(eventContext, message)
	if err != nil {
		return majongpb.StateID_state_qiangganghu_settle, global.ErrInvalidEvent
	}
	utils.SettleOver(flow, message)

	nextState := utils.GetNextState(flow.GetMajongContext())
	if nextState == majongpb.StateID_state_mopai {
		s.setMopaiPlayer(flow)
	}

	return nextState, nil
}
