package tests

import (
	"steve/entity/poker/ddz"
	"steve/room/desks/ddzdesk/flow/machine"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_StartGame(t *testing.T) {
	m := createMachine(ddz.StateID_state_init)
	err := m.ProcessEvent(machine.Event{
		EventID: int(ddz.EventID_event_start_game),
	})
	assert.Nil(t, err)
	context := m.GetDDZContext()
	assert.Equal(t, ddz.StateID_state_deal, context.GetCurState())
}
