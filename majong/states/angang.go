package states

import (
	"steve/majong/interfaces"
	majongpb "steve/server_pb/majong"
)

//AnGangState 暗杠状态 @Author:wuhongwei
type AnGangState struct {
}

var _ interfaces.MajongState = new(AnGangState)

// ProcessEvent 处理事件
// 暗杠逻辑执行完后，进入暗杠状态，确认接收到暗杠完成请求，返回摸牌状态
func (s *AnGangState) ProcessEvent(eventID majongpb.EventID, eventContext []byte, flow interfaces.MajongFlow) (newState majongpb.StateID, err error) {
	if eventID == majongpb.EventID_event_angang_finish {
		return majongpb.StateID(majongpb.StateID_state_mopai), nil
	}
	return majongpb.StateID(majongpb.StateID_state_angang), nil
}

// OnEntry 进入状态
func (s *AnGangState) OnEntry(flow interfaces.MajongFlow) {
	flow.SetAutoEvent(majongpb.AutoEvent{
		EventId:      majongpb.EventID_event_angang_finish,
		EventContext: nil,
	})
}

// OnExit 退出状态
func (s *AnGangState) OnExit(flow interfaces.MajongFlow) {
}
