package states

import (
	"steve/majong/interfaces"
	majongpb "steve/server_pb/majong"
)

//BuGangState 补杠状态 @Author:wuhongwei
type BuGangState struct {
}

var _ interfaces.MajongState = new(BuGangState)

// ProcessEvent 处理事件
func (s *BuGangState) ProcessEvent(eventID majongpb.EventID, eventContext []byte, flow interfaces.MajongFlow) (newState majongpb.StateID, err error) {
	if eventID == majongpb.EventID_event_bugang_finish {
		return majongpb.StateID(majongpb.StateID_state_mopai), nil
	}
	return majongpb.StateID(majongpb.StateID_state_bugang), nil
}

// OnEntry 进入状态
func (s *BuGangState) OnEntry(flow interfaces.MajongFlow) {
	flow.SetAutoEvent(majongpb.AutoEvent{
		EventId:      majongpb.EventID_event_bugang_finish,
		EventContext: nil,
	})
}

// OnExit 退出状态
func (s *BuGangState) OnExit(flow interfaces.MajongFlow) {

}
