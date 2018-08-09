package common

//适用麻将：四川血流
//前置条件：取麻将现场最后杠玩家和玩家杠牌
//处理的事件请求：杠结算扣费失败/成功请求
//处理请求的过程：设置麻将现场的摸牌玩家，生成结算信息
//处理请求的结果：1.杠结算扣费失败，返回游戏结束状态
//			   2.杠结算扣费成功，返回摸牌状态
//状态退出行为：无
//状态进入行为：触发生成杠结算信息
//约束条件：无
import (
	majongpb "steve/entity/majong"
	"steve/room/majong/global"
	"steve/room/majong/interfaces"
	"steve/room/majong/settle"
	"steve/room/majong/utils"

	"github.com/Sirupsen/logrus"
)

// GangSettleState 杠结算状态
type GangSettleState struct {
}

var _ interfaces.MajongState = new(GangSettleState)

// ProcessEvent 处理事件
// 杠逻辑执行完后，进入杠结算状态
// 1.处理结算完成事件，返回摸牌状态
// 2.处理玩家认输事件，返回游戏结束状态
func (s *GangSettleState) ProcessEvent(eventID majongpb.EventID, eventContext interface{}, flow interfaces.MajongFlow) (newState majongpb.StateID, err error) {
	s.setMopaiPlayer(flow)
	if eventID == majongpb.EventID_event_settle_finish {
		message := eventContext.(*majongpb.SettleFinishEvent)
		utils.SettleOver(flow, message)
		return s.nextState(flow.GetMajongContext()), nil
	}
	return majongpb.StateID(majongpb.StateID_state_gang_settle), global.ErrInvalidEvent
}

// OnEntry 进入状态
func (s *GangSettleState) OnEntry(flow interfaces.MajongFlow) {
	s.doGangSettle(flow)
}

// OnExit 退出状态
func (s *GangSettleState) OnExit(flow interfaces.MajongFlow) {
}

// setMopaiPlayer 设置摸牌玩家
func (s *GangSettleState) setMopaiPlayer(flow interfaces.MajongFlow) {
	mjContext := flow.GetMajongContext()
	mjContext.MopaiPlayer = mjContext.GetLastGangPlayer()
	mjContext.MopaiType = majongpb.MopaiType_MT_GANG
}

// doGangSettle 杠结算
func (s *GangSettleState) doGangSettle(flow interfaces.MajongFlow) {
	mjContext := flow.GetMajongContext()
	playerID := mjContext.GetLastGangPlayer()
	player := utils.GetMajongPlayer(playerID, mjContext)

	gangCard := player.GetGangCards()[len(player.GetGangCards())-1]

	param := interfaces.GangSettleParams{
		SettleOptionID: int(mjContext.GetSettleOptionId()),
		GangPlayer:     player.GetPlayerId(),
		SrcPlayer:      gangCard.GetSrcPlayer(),
		AllPlayers:     utils.GetAllPlayers(mjContext),
		HasHuPlayers:   utils.GetHuPlayers(mjContext, []uint64{}),
		QuitPlayers:    utils.GetQuitPlayers(mjContext),
		GiveupPlayers:  utils.GetGiveupPlayers(mjContext),
		GangType:       gangCard.GetType(),
		SettleID:       mjContext.CurrentSettleId,
	}

	settlerFactory := settle.SettlerFactory{}
	settleInfo := settlerFactory.CreateGangSettler(mjContext.GameId).Settle(param)
	if settleInfo != nil {
		mjContext.SettleInfos = append(mjContext.SettleInfos, settleInfo)
		mjContext.CurrentSettleId++
	}
}

// nextState 下个状态
func (s *GangSettleState) nextState(mjcontext *majongpb.MajongContext) majongpb.StateID {
	nextState := utils.IsGameOverReturnState(mjcontext)
	logrus.WithFields(logrus.Fields{
		"func_name": "GangSettleState.nextState",
		"newState":  nextState,
	}).Infoln("杠结算下个状态")
	return nextState
}
