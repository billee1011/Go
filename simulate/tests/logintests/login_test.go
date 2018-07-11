package logintests

import (
	"steve/client_pb/gate"
	"steve/client_pb/msgid"
	"steve/simulate/global"
	"steve/simulate/utils"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Test_Login 测试登录

func Test_Login(t *testing.T) {
	player, err := utils.LoginNewPlayer()
	assert.Nil(t, err)
	assert.NotNil(t, player)
}

// Test_AnotherLogin 顶号测试
// step 1. 登录新玩家
// step 2. 创建新的连接，向网关服认证同一个用户
// 期望：
// 原玩家收到顶号通知
func Test_AnotherLogin(t *testing.T) {
	player, err := utils.LoginNewPlayer()
	assert.Nil(t, err)
	assert.NotNil(t, player)
	player.AddExpectors(msgid.MsgID_GATE_ANOTHER_LOGIN_NTF)

	accountID := player.GetAccountID()
	accountName := utils.GenerateAccountName(accountID)

	newPlayer, err := utils.LoginPlayer(accountID, accountName)
	assert.Nil(t, err)
	assert.NotNil(t, newPlayer)

	expector := player.GetExpector(msgid.MsgID_GATE_ANOTHER_LOGIN_NTF)
	notify := gate.GateAnotherLoginNtf{}
	assert.Nil(t, expector.Recv(global.DefaultWaitMessageTime, &notify))

	time.Sleep(time.Millisecond * 200) // 确保连接断开
	assert.True(t, player.GetClient().Closed())
}
