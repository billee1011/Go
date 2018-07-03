package ddztests

import (
	"steve/simulate/tests/ddztests/game"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_StartGame(t *testing.T) {
	ddzGame, err := game.StartGame(game.StartGameParams{})
	assert.NotNil(t, ddzGame)
	assert.Nil(t, err)
}
