package interfaces

// EventType 事件类型
type EventType int

const (
	// NormalEvent 普通事件
	NormalEvent EventType = iota
	// OverTimeEvent 超时事件
	OverTimeEvent
	// TuoGuanEvent 托管事件
	TuoGuanEvent
	// RobotEvent 机器人事件
	RobotEvent
	// HuStateEvent 胡状态的事件
	HuStateEvent
	// TingStateEvent 听状态的事件
	TingStateEvent
)

// Event 事件
type Event struct {
	ID        int32       // 事件 ID
	Context   interface{} // 事件现场
	EventType EventType   // 事件类型
	PlayerID  uint64      // 针对哪个玩家的事件
}
