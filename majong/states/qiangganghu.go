package states

import (
	"steve/majong/interfaces"
	majongpb "steve/server_pb/majong"
)

// QiangganghuState 抢杠胡状态
type QiangganghuState struct {
}

var _ interfaces.MajongState = new(QiangganghuState)

// ProcessEvent 处理事件
func (s *QiangganghuState) ProcessEvent(eventID majongpb.EventID, eventContext []byte, flow interfaces.MajongFlow) (newState majongpb.StateID, err error) {
	if eventID == majongpb.EventID_event_qiangganghu_finish {
		return majongpb.StateID_state_mopai, nil
	}
	return majongpb.StateID_state_qiangganghu, errInvalidEvent
}

// OnEntry 进入状态
func (s *QiangganghuState) OnEntry(flow interfaces.MajongFlow) {
	flow.SetAutoEvent(majongpb.AutoEvent{
		EventId:      majongpb.EventID_event_qiangganghu_finish,
		EventContext: nil,
	})
}

// OnExit 退出状态
func (s *QiangganghuState) OnExit(flow interfaces.MajongFlow) {

}
