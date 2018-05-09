package states

import (
	"steve/majong/interfaces"
	majongpb "steve/server_pb/majong"
)

//MingGangState 明杠状态 @Author:wuhongwei
type MingGangState struct {
}

var _ interfaces.MajongState = new(MingGangState)

// ProcessEvent 处理事件
func (s *MingGangState) ProcessEvent(eventID majongpb.EventID, eventContext []byte, flow interfaces.MajongFlow) (newState majongpb.StateID, err error) {
	if eventID == majongpb.EventID_event_mopai_finish {
		return majongpb.StateID(majongpb.StateID_state_mopai), nil
	}
	return majongpb.StateID(majongpb.StateID_state_gang), nil
}

// OnEntry 进入状态
func (s *MingGangState) OnEntry(flow interfaces.MajongFlow) {
	flow.SetAutoEvent(majongpb.AutoEvent{
		EventId:      majongpb.EventID_event_mopai_finish,
		EventContext: nil,
	})
}

// OnExit 退出状态
func (s *MingGangState) OnExit(flow interfaces.MajongFlow) {
}
