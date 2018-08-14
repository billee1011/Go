package common

import (
	"steve/client_pb/msgid"
	"steve/client_pb/room"
	"steve/common/mjoption"
	majongpb "steve/entity/majong"
	"steve/room/majong/global"
	"steve/room/majong/interfaces"
	"steve/room/majong/interfaces/facade"
	"time"

	"github.com/golang/protobuf/proto"
)

// InitState 初始化状态
type InitState struct {
}

var _ interfaces.MajongState = new(InitState)

// ProcessEvent 处理事件
func (s *InitState) ProcessEvent(eventID majongpb.EventID, eventContext interface{}, flow interfaces.MajongFlow) (newState majongpb.StateID, err error) {
	if eventID == majongpb.EventID_event_start_game {
		s.notifyPlayers(flow)
		s.setStartTime(flow)
		return majongpb.StateID_state_xipai, nil
	}
	return majongpb.StateID_state_init, global.ErrInvalidEvent
}

// notifyPlayers 通知玩家游戏开始
func (s *InitState) notifyPlayers(flow interfaces.MajongFlow) {
	//先要判断游戏有没有换三张的玩法，有换三张的玩法，再判断需不需要配置换三张
	mjContext := flow.GetMajongContext()
	mjContext.NextBankerSeat = uint32(len(mjContext.GetPlayers()))
	isHsz := mjoption.GetXingpaiOption(int(mjContext.GetXingpaiOptionId())).Hnz.Enable
	if isHsz {
		isHsz = mjContext.GetOption().GetHasHuansanzhang()
	}
	facade.BroadcaseMessage(flow, msgid.MsgID_ROOM_START_GAME_NTF, &room.RoomStartGameNtf{
		NeedHsz: proto.Bool(isHsz),
	})
}

func (s *InitState) setStartTime(flow interfaces.MajongFlow) {
	mjContext := flow.GetMajongContext()
	mjContext.GameStartTime = time.Now()
}

// OnEntry 进入状态
func (s *InitState) OnEntry(flow interfaces.MajongFlow) {

}

// OnExit 退出状态
func (s *InitState) OnExit(flow interfaces.MajongFlow) {

}
