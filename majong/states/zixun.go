package states

import (
	"steve/majong/interfaces"
	majongpb "steve/server_pb/majong"
)

// ZiXunState 摸牌状态
type ZiXunState struct {
}

var _ interfaces.MajongState = new(ZiXunState)

// ProcessEvent 处理事件
func (s *ZiXunState) ProcessEvent(eventID majongpb.EventID, eventContext []byte, flow interfaces.MajongFlow) (newState majongpb.StateID, err error) {
	switch eventID {
	case majongpb.EventID_event_angang_request:
		{
			mjContext := flow.GetMajongContext()
			wallCards := mjContext.GetWallCards()
			if len(wallCards) == 0 {
				return majongpb.StateID_state_zixun, errInvalidEvent
			}
			return s.angang(flow)
		}
	case majongpb.EventID_event_zimo_request:
		{
			return s.zimo(flow)

		}
	case majongpb.EventID_event_chupai_request:
		{
			return s.chupai(flow)
		}
	default:
		{
			return majongpb.StateID_state_zixun, errInvalidEvent
		}
	}
}

func (s *ZiXunState) angang(flow interfaces.MajongFlow) (majongpb.StateID, error) {
	s.checkAnGang(flow)
	return majongpb.StateID_state_angang, nil
}

func (s *ZiXunState) zimo(flow interfaces.MajongFlow) (majongpb.StateID, error) {

	return majongpb.StateID_state_zimo, nil
}
func (s *ZiXunState) chupai(flow interfaces.MajongFlow) (majongpb.StateID, error) {

	return majongpb.StateID_state_chupai, nil
}

func (s *ZiXunState) checkAnGang(flow interfaces.MajongFlow) bool {
	return true
}
func (s *ZiXunState) checkBuGang(flow interfaces.MajongFlow) bool {
	return true
}

func (s *ZiXunState) checkZiMo(flow interfaces.MajongFlow) bool {
	return true
}

// OnEntry 进入状态
func (s *ZiXunState) OnEntry(flow interfaces.MajongFlow) {
	//进行相关的check
}

// OnExit 退出状态
func (s *ZiXunState) OnExit(flow interfaces.MajongFlow) {

}
