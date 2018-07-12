package common

import (
	 "steve/client_pb/msgId"
	"steve/client_pb/room"
	"steve/common/mjoption"
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
	//先要判断游戏有没有换三张的玩法，有换三张的玩法，再判断需不需要配置换三张
	mjContext := flow.GetMajongContext()
	isHsz := mjoption.GetXingpaiOption(int(mjContext.GetXingpaiOptionId())).Hnz.Enable
	if isHsz {
		isHsz = mjContext.GetOption().GetHasHuansanzhang()
	}
	facade.BroadcaseMessage(flow, msgId.MsgID_ROOM_START_GAME_NTF, &room.RoomStartGameNtf{
		NeedHsz: proto.Bool(isHsz),
	})
}

// OnEntry 进入状态
func (s *InitState) OnEntry(flow interfaces.MajongFlow) {

}

// OnExit 退出状态
func (s *InitState) OnExit(flow interfaces.MajongFlow) {

}
