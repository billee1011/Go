package states

import (
	"steve/majong/interfaces"
	"steve/majong/utils"
	majongpb "steve/server_pb/majong"
)

// MoPaiState 摸牌状态
type MoPaiState struct {
}

var _ interfaces.MajongState = new(MoPaiState)

// ProcessEvent 处理事件
func (s *MoPaiState) ProcessEvent(eventID majongpb.EventID, eventContext []byte, flow interfaces.MajongFlow) (newState majongpb.StateID, err error) {
	if eventID == majongpb.EventID_event_mopai_finish {
		mjContext := flow.GetMajongContext()
		wallCards := mjContext.GetWallCards()
		if len(wallCards) == 0 {
			return majongpb.StateID_state_gameover, nil
		}
		return s.mopai(flow)
	}
	return majongpb.StateID_state_mopai, errInvalidEvent
}

func (s *MoPaiState) mopai(flow interfaces.MajongFlow) (majongpb.StateID, error) {
	context := flow.GetMajongContext()
	players := context.GetPlayers()
	activePlayer := utils.GetPlayerByID(players, context.ActivePlayer)
	//从墙牌中移除一张牌
	drowCard := context.WallCards[0]
	context.WallCards = context.WallCards[1:]
	//将这张牌添加到手牌中
	activePlayer.HandCards = append(activePlayer.HandCards, drowCard)
	return majongpb.StateID_state_zixun, nil
}

// OnEntry 进入状态
func (s *MoPaiState) OnEntry(flow interfaces.MajongFlow) {
	flow.SetAutoEvent(majongpb.AutoEvent{
		EventId:      majongpb.EventID_event_mopai_finish,
		EventContext: nil,
	})
}

// OnExit 退出状态
func (s *MoPaiState) OnExit(flow interfaces.MajongFlow) {

}
