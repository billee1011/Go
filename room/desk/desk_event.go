package desk


type DeskEvent struct {
	EventID  int
	EventType int      // 事件类型
	Params   EventParams
	Desk     *Desk
}

func NewDeskEvent(id int,eventType int, desk *Desk, params EventParams) DeskEvent {
	return DeskEvent{EventID: id,EventType:eventType, Params: params, Desk: desk}
}
