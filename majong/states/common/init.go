package common

import (
	msgid "steve/client_pb/msgId"
	"steve/client_pb/room"
	"steve/gutils"
	"steve/majong/global"
	"steve/majong/interfaces"
	"steve/majong/interfaces/facade"
	majongpb "steve/server_pb/majong"

	"github.com/golang/protobuf/proto"
)

// InitState 初始化状态
type InitState struct {
}

var _ interfaces.MajongState = new(InitState)

// ProcessEvent 处理事件
func (s *InitState) ProcessEvent(eventID majongpb.EventID, eventContext []byte, flow interfaces.MajongFlow) (newState majongpb.StateID, err error) {
	if eventID == majongpb.EventID_event_start_game {
		s.notifyPlayers(flow)
		return majongpb.StateID_state_xipai, nil
	}
	return majongpb.StateID_state_init, global.ErrInvalidEvent
}

// notifyPlayers 通知玩家游戏开始
func (s *InitState) notifyPlayers(flow interfaces.MajongFlow) {
	//通知前端是否要进行换三张，默认通知
	isHsz := true
	if flow.GetMajongContext().GetGameId() == gutils.SCXZGameID {
		isHsz = flow.GetMajongContext().GetOption().GetHasHuansanzhang()
	}
	facade.BroadcaseMessage(flow, msgid.MsgID_ROOM_START_GAME_NTF, &room.RoomStartGameNtf{
		NeedHsz: proto.Bool(isHsz),
	})
}

// OnEntry 进入状态
func (s *InitState) OnEntry(flow interfaces.MajongFlow) {

}

// OnExit 退出状态
func (s *InitState) OnExit(flow interfaces.MajongFlow) {

}
