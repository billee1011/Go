package desk

type DeskEvent struct {
	EventID  int
	ParamLen int
	Params   []interface{}
	Desk     *Desk
}

func NewDeskEvent(id int, len int, desk *Desk, params ...interface{}) DeskEvent {
	return DeskEvent{EventID: id, ParamLen: len, Params: params, Desk: desk}
}
