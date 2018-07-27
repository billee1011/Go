package matchtests

import (
	"steve/client_pb/common"
	msgid "steve/client_pb/msgid"
	"steve/client_pb/room"
	"steve/simulate/cheater"
	"steve/simulate/global"
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
		// 登录用户
		player, err := utils.LoginNewPlayer()
		assert.Nil(t, err)
		assert.NotNil(t, player)
		client := player.GetClient()

		createNtfExpector, err := client.ExpectMessage(msgid.MsgID_ROOM_DESK_CREATED_NTF)
		assert.Nil(t, err)
		createNtfExpectors[i] = createNtfExpector

		gameStartNtfExpector, err := client.ExpectMessage(msgid.MsgID_ROOM_START_GAME_NTF)
		assert.Nil(t, err)
		gameStartNtfExpectors[i] = gameStartNtfExpector

		_, err = utils.ApplyJoinDesk(player, common.GameId_GAMEID_XUELIU)
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

// TestNoMoneyMatch 测试金币数为0时参与匹配
// 登录玩家，设置其金币数为0，发起匹配请求
// 期望收到回复，匹配失败
func TestNoMoneyMatch(t *testing.T) {
	player, err := utils.LoginNewPlayer()
	assert.Nil(t, err)
	assert.NotNil(t, player)
	playerID := player.GetID()
	cheater.SetPlayerCoin(playerID, 0)
	rsp, err := utils.ApplyJoinDesk(player, common.GameId_GAMEID_XUELIU)
	assert.Nil(t, err)
	assert.NotEqual(t, int32(0), rsp.GetErrCode())
}

/// 等待时间太久，先注释
// TestRobotMatch 测试机器人匹配
// 步骤：
//  1. 登录 1 个玩家，申请匹配四川血流
//  2. 等待 5s
// 期望：
//  1. 玩家收到创建房间通知和开始游戏通知
func TestRobotMatch(t *testing.T) {
	// 修改机器人加入匹配的时间为 100ms
	modifyRobotJoinTime(100 * time.Millisecond)
	player, err := utils.LoginNewPlayer()
	assert.Nil(t, err)
	assert.NotNil(t, player)
	player.AddExpectors(msgid.MsgID_ROOM_DESK_CREATED_NTF, msgid.MsgID_ROOM_START_GAME_NTF)

	utils.ApplyJoinDesk(player, common.GameId_GAMEID_XUELIU)
	createExpector := player.GetExpector(msgid.MsgID_ROOM_DESK_CREATED_NTF)
	assert.Nil(t, createExpector.Recv(global.DefaultWaitMessageTime, nil))

	startExpector := player.GetExpector(msgid.MsgID_ROOM_START_GAME_NTF)
	assert.Nil(t, startExpector.Recv(global.DefaultWaitMessageTime, nil))
}
