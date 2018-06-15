package tests

import (
	"steve/client_pb/msgId"
	"steve/client_pb/room"
	"steve/simulate/config"
	"steve/simulate/connect"
	"steve/simulate/global"
	"steve/simulate/interfaces"
	"steve/simulate/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestApplyJoinDesk(t *testing.T) {

	createNtfExpectors := map[int]interfaces.MessageExpector{}
	gameStartNtfExpectors := map[int]interfaces.MessageExpector{}

	for i := 0; i < 4; i++ {
		// 创建客户端连接
		client := connect.NewTestClient(config.ServerAddr, config.ClientVersion)
		assert.NotNil(t, client)

		// 登录用户
		player, err := utils.LoginUser(client, global.AllocUserName())
		assert.Nil(t, err)
		assert.NotNil(t, player)

		createNtfExpector, err := client.ExpectMessage(msgid.MsgID_ROOM_DESK_CREATED_NTF)
		assert.Nil(t, err)
		createNtfExpectors[i] = createNtfExpector

		gameStartNtfExpector, err := client.ExpectMessage(msgid.MsgID_ROOM_START_GAME_NTF)
		assert.Nil(t, err)
		gameStartNtfExpectors[i] = gameStartNtfExpector

		_, err = utils.ApplyJoinDesk(player)
		assert.Nil(t, err)
	}

	for _, e := range createNtfExpectors {
		ntf := &room.RoomDeskCreatedNtf{}
		assert.Nil(t, e.Recv(global.DefaultWaitMessageTime, ntf))
		assert.Equal(t, 4, len(ntf.GetPlayers()))
	}
	for _, e := range gameStartNtfExpectors {
		ntf := &room.RoomStartGameNtf{}
		assert.Nil(t, e.Recv(global.DefaultWaitMessageTime, ntf))
	}
}
