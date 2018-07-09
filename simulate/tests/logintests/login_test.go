package logintests

import (
	"fmt"
	"steve/client_pb/gate"
	"steve/client_pb/msgId"
	"steve/simulate/config"
	"steve/simulate/connect"
	"steve/simulate/global"
	"steve/simulate/utils"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Test_Login 测试登录
// 步骤
// step 1. 向登录服请求认证， 获取到玩家 ID 和网关地址等信息
// step 2. 连接网关服，并且向网关服认证
// 期望
// 连接网关服成功，认证成功

func Test_Login(t *testing.T) {
	loginClient := connect.NewTestClient(config.GetLoginServerAddr(), config.GetClientVersion())
	accountID := global.AllocAccountID()
	accountName := utils.GenerateAccountName(accountID)
	loginResponse, err := utils.RequestAuth(loginClient, accountID, accountName, time.Second*5)
	assert.Nil(t, err)
	playerID := loginResponse.GetPlayerId()
	expire := loginResponse.GetExpire()
	token := loginResponse.GetGateToken()

	gateIP := loginResponse.GetGateIp()
	gatePort := loginResponse.GetGatePort()

	gateClient := connect.NewTestClient(fmt.Sprintf("%s:%d", gateIP, gatePort), config.GetClientVersion())
	assert.NotNil(t, gateClient)
	assert.Nil(t, utils.RequestGateAuth(gateClient, playerID, expire, token))
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
