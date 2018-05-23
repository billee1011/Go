package autoevent

import (
	"steve/room/interfaces"
	"steve/server_pb/majong"
	"time"
)

type autoEventGenerator struct {
	stateHandlers map[majong.StateID]stateHandler
}

func (aeg *autoEventGenerator) Generate(mjContext *majong.MajongContext, stateTime time.Time) []interfaces.Event {
	state := mjContext.GetCurState()
	if handler, ok := aeg.stateHandlers[state]; ok {
		return handler.Generate(mjContext, stateTime)
	}
	return []interfaces.Event{}
}

func newEventGenerator() interfaces.DeskAutoEventGenerator {
	return &autoEventGenerator{
		stateHandlers: map[majong.StateID]stateHandler{
			majong.StateID_state_dingque: newDingqueHandler(),
		},
	}
}
