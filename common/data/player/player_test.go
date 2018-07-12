package player

import (
	"testing"

	"github.com/spf13/viper"

	"github.com/stretchr/testify/assert"
)

func init() {
	viper.SetDefault("redis_addr", "127.0.0.1:6379")
	viper.SetDefault("redis_passwd", "")
}

func TestGetPlayerPlayStates(t *testing.T) {
	var testPlayerID uint64 = 100000000

	defState := PlayStates{
		State:    99,
		GameID:   100,
		RoomAddr: "abc",
	}
	state, err := GetPlayerPlayStates(testPlayerID, defState)
	assert.Nil(t, err)
	assert.Equal(t, state.State, 99)
	assert.Equal(t, state.GameID, 100)
	assert.Equal(t, state.RoomAddr, "abc")
}

func TestSetPlayerPlayStates(t *testing.T) {
	var testPlayerID uint64 = 100000001
	SetPlayerPlayStates(testPlayerID, PlayStates{
		State:    1,
		GameID:   2,
		RoomAddr: "cde",
	})
	states, err := GetPlayerPlayStates(testPlayerID, PlayStates{})
	assert.Nil(t, err)
	assert.Equal(t, states.State, 1)
	assert.Equal(t, states.GameID, 2)
	assert.Equal(t, states.RoomAddr, "cde")
}
