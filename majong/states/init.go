package states

import (
	"steve/majong/interfaces"
	majongpb "steve/server_pb/majong"
)

// InitState 初始化状态
type InitState struct {
}

var _ interfaces.MajongState = new(InitState)

// ProcessEvent 处理事件
func (s *InitState) ProcessEvent(eventID majongpb.EventID, eventContext []byte, flow interfaces.MajongFlow) (newState majongpb.StateID, err error) {
	if eventID == majongpb.EventID_event_start_game {
		return majongpb.StateID_state_xipai, nil
	}
	return majongpb.StateID_state_init, errInvalidEvent
}

// OnEntry 进入状态
func (s *InitState) OnEntry(flow interfaces.MajongFlow) {

}

// OnExit 退出状态
func (s *InitState) OnExit(flow interfaces.MajongFlow) {

}
