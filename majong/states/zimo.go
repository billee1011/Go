package states

import (
	"steve/majong/interfaces"
	"steve/majong/utils"
	majongpb "steve/server_pb/majong"
)

// ZimoState 自摸状态
type ZimoState struct {
}

var _ interfaces.MajongState = new(ZimoState)

// ProcessEvent 处理事件
func (s *ZimoState) ProcessEvent(eventID majongpb.EventID, eventContext []byte, flow interfaces.MajongFlow) (newState majongpb.StateID, err error) {
	if eventID == majongpb.EventID_event_zimo_finish {
		s.mopai(flow)
		return majongpb.StateID_state_mopai, nil
	}
	return majongpb.StateID_state_zimo, errInvalidEvent
}

//mopai 摸牌处理
func (s *ZimoState) mopai(flow interfaces.MajongFlow) (majongpb.StateID, error) {
	context := flow.GetMajongContext()
	players := context.GetPlayers()
	activePlayer := utils.GetNextPlayerByID(players, context.ActivePlayer)
	//TODO：目前只在这个地方改变操作玩家（感觉碰，明杠，点炮这三种情况也需要改变activePlayer）
	context.ActivePlayer = activePlayer.GetPalyerId()
	//从墙牌中移除一张牌
	drowCard := context.WallCards[0]
	context.WallCards = context.WallCards[1:]
	//将这张牌添加到手牌中
	activePlayer.HandCards = append(activePlayer.HandCards, drowCard)
	return majongpb.StateID_state_zixun, nil
}

// OnEntry 进入状态
func (s *ZimoState) OnEntry(flow interfaces.MajongFlow) {
	flow.SetAutoEvent(majongpb.AutoEvent{
		EventId:      majongpb.EventID_event_zimo_finish,
		EventContext: nil,
	})
}

// OnExit 退出状态
func (s *ZimoState) OnExit(flow interfaces.MajongFlow) {

}
