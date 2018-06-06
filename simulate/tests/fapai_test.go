package tests

import (
	msgid "steve/client_pb/msgId"
	"steve/client_pb/room"
	"steve/simulate/connect"
	"steve/simulate/global"
	"steve/simulate/interfaces"
	"steve/simulate/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Fapai(t *testing.T) {

	fapaiNtfExpectors := map[int]interfaces.MessageExpector{}

	for i := 0; i < 4; i++ {
		// 登录玩家
		client := connect.NewTestClient(ServerAddr, ClientVersion)
		assert.NotNil(t, client)
		player, err := utils.LoginUser(client, global.AllocUserName())
		assert.Nil(t, err)
		assert.NotNil(t, player)

		// 创建消息期望： 期望收到发牌通知消息
		expector, err := client.ExpectMessage(msgid.MsgID_ROOM_FAPAI_NTF)
		assert.Nil(t, err)
		fapaiNtfExpectors[i] = expector

		assert.Nil(t, utils.ApplyJoinDesk(player))
	}

	cardCountMap := map[int]int{}

	for _, expector := range fapaiNtfExpectors {
		ntf := room.RoomFapaiNtf{}
		assert.Nil(t, expector.Recv(global.DefaultWaitMessageTime, &ntf))
		// 收到消息后，有4个玩家的卡牌数量
		assert.Equal(t, 4, len(ntf.GetPlayerCardCounts()))
		cardCount := len(ntf.GetCards())
		// 自己的卡牌数量为 13 或者 14
		assert.True(t, cardCount == 13 || cardCount == 14)
		cardCountMap[cardCount] = cardCountMap[cardCount] + 1
	}

	// 有 3 个玩家的卡牌数量时 13， 有 1 个玩家的卡牌数量是 14
	assert.Equal(t, 3, cardCountMap[13])
	assert.Equal(t, 1, cardCountMap[14])
}
