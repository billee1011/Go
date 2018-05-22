package matchtests

import (
	msgid "steve/client_pb/msgId"
	"steve/client_pb/room"
	"steve/simulate/config"
	"steve/simulate/connect"
	"steve/simulate/global"
	"steve/simulate/interfaces"
	"steve/simulate/utils"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Test_OfflineMatch 测试离线时不会被匹配到
// 1. 登录玩家1，发送加入房间请求
// 2. 玩家1断开连接
// 3. 登录4个玩家，分别发送加入房间请求
// 预期：
//  后4个玩家都收到了创建房间通知和游戏开始通知
func Test_OfflineMatch(t *testing.T) {
	client1 := connect.NewTestClient(config.ServerAddr, config.ClientVersion)
	assert.NotNil(t, client1)
	_, err := utils.LoginUser(client1, "test_user")
	assert.Nil(t, err)
	client1.Stop()
	time.Sleep(time.Millisecond * 200) // 等200毫秒，确保连接断开

	createNtfExpectors := map[int]interfaces.MessageExpector{}
	gameStartNtfExpectors := map[int]interfaces.MessageExpector{}

	for i := 0; i < 4; i++ {
		// 创建客户端连接
		client := connect.NewTestClient(config.ServerAddr, config.ClientVersion)
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
		assert.Nil(t, e.Recv(global.DefaultWaitMessageTime, ntf))
		assert.Equal(t, 4, len(ntf.GetPlayers()))
	}
	for _, e := range gameStartNtfExpectors {
		ntf := &room.RoomStartGameNtf{}
		assert.Nil(t, e.Recv(global.DefaultWaitMessageTime, ntf))
	}
}
