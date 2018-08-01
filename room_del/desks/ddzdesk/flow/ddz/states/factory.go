package states

import (
	"steve/entity/poker/ddz"
	"steve/room/desks/ddzdesk/flow/machine"
)

type stateFactory struct {
}

func (f *stateFactory) NewState(stateID int) machine.State {
	switch ddz.StateID(stateID) {
	case ddz.StateID_state_init:
		{
			return new(initState)
		}
	case ddz.StateID_state_deal:
		{
			return new(dealState)
		}
	case ddz.StateID_state_grab:
		{
			return new(grabState)
		}
	case ddz.StateID_state_double:
		{
			return new(doubleState)
		}
	case ddz.StateID_state_playing:
		{
			return new(playState)
		}
	case ddz.StateID_state_settle:
		{
			return new(settleState)
		}
	case ddz.StateID_state_over:
		{
			return new(overState)
		}
	}
	return nil
}

// NewFactory 创建工厂
func NewFactory() machine.StateFactory {
	return new(stateFactory)
}
