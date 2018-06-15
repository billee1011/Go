package tests

import (
	"steve/client_pb/room"
	"steve/simulate/config"
	"steve/simulate/global"
	"steve/simulate/utils"

	"github.com/golang/protobuf/proto"

	"steve/simulate/connect"

	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogin(t *testing.T) {
	client := connect.NewTestClient(config.ServerAddr, config.ClientVersion)
	assert.NotNil(t, client)
	userName := global.AllocUserName()
	player, err := utils.LoginUser(client, userName)
	assert.Nil(t, err)
	assert.NotNil(t, player)
	assert.NotEqual(t, 0, player.GetID())
	player2, err := utils.LoginUser(client, userName)
	assert.Nil(t, err)
	assert.NotNil(t, player2)
	assert.Equal(t, player2.GetID(), player.GetID())
}

func TestVisitorLogin(t *testing.T) {
	var playerID uint64
	loginCount := 5
	for i := loginCount; i > 0; i-- {
		client := connect.NewTestClient(config.ServerAddr, config.ClientVersion)
		assert.NotNil(t, client)
		visitorLoginReq := &room.RoomVisitorLoginReq{
			DeviceInfo: &room.DeviceInfo{
				DeviceType: room.DeviceType_DT_ANDROID.Enum(),
				Uuid:       proto.String("1013210cc"),
			},
		}
		player, err := utils.LoginVisitor(client, visitorLoginReq)
		if i == loginCount {
			playerID = player.GetID()
		}
		assert.Nil(t, err)
		assert.NotNil(t, player)
		assert.NotEqual(t, 0, player.GetID())
		if i != loginCount {
			assert.Equal(t, playerID, player.GetID())
		}
		client.Stop()
	}
}
