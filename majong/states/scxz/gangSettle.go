package scxz

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
	"steve/majong/global"
	"steve/majong/interfaces"
	"steve/majong/interfaces/facade"
	"steve/majong/utils"
	majongpb "steve/server_pb/majong"

	"github.com/golang/protobuf/proto"
)

// GangSettleState 杠结算状态
type GangSettleState struct {
}

var _ interfaces.MajongState = new(GangSettleState)

// ProcessEvent 处理事件
// 杠逻辑执行完后，进入杠结算状态
// 1.处理结算完成事件，返回摸牌状态
// 2.处理玩家认输事件，返回游戏结束状态
func (s *GangSettleState) ProcessEvent(eventID majongpb.EventID, eventContext []byte, flow interfaces.MajongFlow) (newState majongpb.StateID, err error) {
	s.setMopaiPlayer(flow)
	if eventID == majongpb.EventID_event_settle_finish {
		message := &majongpb.SettleFinishEvent{}
		err := proto.Unmarshal(eventContext, message)
		if err != nil {
			return majongpb.StateID_state_gang_settle, global.ErrInvalidEvent
		}
		return s.settleOver(flow, message)
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

	allPlayers := make([]uint64, 0)
	for _, player := range mjContext.Players {
		if player.XpState == majongpb.XingPaiState_normal {
			allPlayers = append(allPlayers, player.GetPalyerId())
		}
	}
	param := interfaces.GangSettleParams{
		GangPlayer: player.GetPalyerId(),
		SrcPlayer:  gangCard.GetSrcPlayer(),
		AllPlayers: allPlayers,
		GangType:   gangCard.GetType(),
		SettleID:   mjContext.CurrentSettleId,
	}

	f := global.GetGameSettlerFactory()
	gameID := int(mjContext.GetGameId())
	settleInfo := facade.SettleGang(f, gameID, param)
	if settleInfo != nil {
		mjContext.SettleInfos = append(mjContext.SettleInfos, settleInfo)
		mjContext.CurrentSettleId++
	}
}

//settleOver 结算完成
func (s *GangSettleState) settleOver(flow interfaces.MajongFlow, message *majongpb.SettleFinishEvent) (majongpb.StateID, error) {
	mjContext := flow.GetMajongContext()
	playerIds := message.GetPlayerId()
	if len(playerIds) != 0 {
		for _, pid := range playerIds {
			player := utils.GetMajongPlayer(pid, mjContext)
			if player == nil {
				return majongpb.StateID_state_gang_settle, global.ErrInvalidEvent
			}
			player.XpState = majongpb.XingPaiState_give_up
		}
		return majongpb.StateID_state_gameover, nil
	}
	return majongpb.StateID(majongpb.StateID_state_mopai), nil
}
