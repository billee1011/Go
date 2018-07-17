package desk

import "steve/room2/common"

type DeskEvent struct {
	EventID  int
	EventType int      // 事件类型
	Params   common.EventParams
	Desk     *Desk
}

func NewDeskEvent(id int,eventType int, desk *Desk, params common.EventParams) DeskEvent {
	return DeskEvent{EventID: id,EventType:eventType, Params: params, Desk: desk}
}
