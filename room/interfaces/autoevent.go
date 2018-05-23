package interfaces

import (
	"steve/server_pb/majong"
	"time"
)

// Event 事件
type Event struct {
	ID      majong.EventID
	Context []byte
}

// DeskAutoEventGenerator 牌桌自动事件产生器
type DeskAutoEventGenerator interface {
	Generate(mjContext *majong.MajongContext, stateTime time.Time) []Event
}

// DeskAutoEventGeneratorFactory 工厂
type DeskAutoEventGeneratorFactory interface {
	CreateGenerator() DeskAutoEventGenerator
}
