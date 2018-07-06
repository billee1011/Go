package logintests

import (
	"fmt"
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
