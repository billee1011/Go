package states

import (
	"steve/majong/global"
	"steve/majong/interfaces"
	majongpb "steve/server_pb/majong"
)

// GameOverState 游戏结束状态
type GameOverState struct {
}

var _ interfaces.MajongState = new(GameOverState)

// ProcessEvent 处理事件
func (s *GameOverState) ProcessEvent(eventID majongpb.EventID, eventContext []byte, flow interfaces.MajongFlow) (newState majongpb.StateID, err error) {
	return majongpb.StateID_state_gameover, global.ErrInvalidEvent
}

// OnEntry 进入状态
func (s *GameOverState) OnEntry(flow interfaces.MajongFlow) {

}

// OnExit 退出状态
func (s *GameOverState) OnExit(flow interfaces.MajongFlow) {

}
