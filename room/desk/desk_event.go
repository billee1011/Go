package desk

type DeskEvent struct {
	EventID     int
	EventType   int // 事件类型
	Context     interface{}
	PlayerID    uint64
	StateNumber int
	Desk        *Desk
}
