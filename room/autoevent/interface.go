package autoevent

import (
	"steve/room/interfaces"
	"steve/server_pb/majong"
	"time"
)

type stateHandler interface {
	Generate(mjContext *majong.MajongContext, stateTime time.Time) []interfaces.Event
}
