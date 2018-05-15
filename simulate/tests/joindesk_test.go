package tests

import (
	"steve/client_pb/msgId"
	"steve/client_pb/room"
	"steve/simulate/connect"
	"steve/simulate/interfaces"
	"steve/simulate/utils"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestApplyJoinDesk(t *testing.T) {

	createNtfExpectors := map[int]interfaces.MessageExpector{}
	gameStartNtfExpectors := map[int]interfaces.MessageExpector{}

	for i := 0; i < 4; i++ {
		// 创建客户端连接
		client := connect.NewTestClient(ServerAddr, ClientVersion)
		assert.NotNil(t, client)

		// 登录用户
		player, err := utils.LoginUser(client, "test_user")
		assert.Nil(t, err)
		assert.NotNil(t, player)

		createNtfExpector, err := client.ExpectMessage(msgid.MsgID_ROOM_DESK_CREATED_NTF)
		assert.Nil(t, err)
		createNtfExpectors[i] = createNtfExpector

		gameStartNtfExpector, err := client.ExpectMessage(msgid.MsgID_ROOM_START_GAME_NTF)
		assert.Nil(t, err)
		gameStartNtfExpectors[i] = gameStartNtfExpector

		assert.Nil(t, utils.ApplyJoinDesk(player))
	}

	for _, e := range createNtfExpectors {
		ntf := &room.RoomDeskCreatedNtf{}
		assert.Nil(t, e.Recv(time.Second*1, ntf))
		assert.Equal(t, 4, len(ntf.GetPlayers()))
	}
	for _, e := range gameStartNtfExpectors {
		ntf := &room.RoomStartGameNtf{}
		assert.Nil(t, e.Recv(time.Second*1, ntf))
	}
}
