package tests

import (
	"steve/simulate/utils"

	"steve/simulate/connect"

	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogin(t *testing.T) {
	client := connect.NewTestClient(ServerAddr, ClientVersion)
	assert.NotNil(t, client)
	player, err := utils.LoginUser(client, "test_user")
	assert.Nil(t, err)
	assert.NotNil(t, player)
	assert.NotEqual(t, 0, player.GetID())
}
