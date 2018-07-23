package matchtests

import (
	"steve/client_pb/common"
	"steve/client_pb/msgid"
	"steve/simulate/global"
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
	player1, err := utils.LoginNewPlayer()
	utils.ApplyJoinDesk(player1, common.GameId_GAMEID_XUELIU)
	assert.Nil(t, err)
	player1.GetClient().Stop()
	time.Sleep(time.Millisecond * 200) // 等200毫秒，确保连接断开

	params := global.NewCommonStartGameParams()
	deskData, err := utils.StartGame(params)
	assert.NotNil(t, deskData)
	assert.Nil(t, err)
}

// Test_LoginAfterMatch 测试申请匹配后再次登录，不会被匹配到
// 设置机器人匹配时间为100ms, 登录玩家，申请匹配，使用原账号再次登录玩家。
// 预期 200ms 后不会创建房间通知
func Test_LoginAfterMatch(t *testing.T) {

	player, err := utils.LoginNewPlayer()
	_, err = utils.ApplyJoinDesk(player, common.GameId_GAMEID_XUELIU)
	assert.Nil(t, err)

	modifyRobotJoinTime(100 * time.Millisecond)

	player, err = utils.LoginPlayer(player.GetAccountID(), "")
	player.AddExpectors(msgid.MsgID_ROOM_DESK_CREATED_NTF)
	createExpector := player.GetExpector(msgid.MsgID_ROOM_DESK_CREATED_NTF)
	assert.NotNil(t, createExpector.Recv(time.Second*1, nil))
}
