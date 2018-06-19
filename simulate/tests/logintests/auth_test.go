package logintests

import (
	"steve/client_pb/login"
	"steve/simulate/config"
	"steve/simulate/connect"
	"steve/simulate/global"
	"steve/simulate/utils"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// 目的：  验证认证是否成功
// 步骤：
// 	1. 创建客户端并连接登录服
//  2. 分配账号 ID， 生成账号名字， 发起认证请求
// 期望：
//	1. 收到服务器响应
//	2. 错误码为成功， 玩家 ID 不为0， 网关 IP 和网关端口合法，到期时间大于当前时间， token 不为空
func Test_Auth(t *testing.T) {
	client := connect.NewTestClient(config.ServerAddr, config.ClientVersion)
	assert.NotNil(t, client)

	accountID := global.AllocAccountID()
	accountName := utils.GenerateAccountName(accountID)

	response, err := utils.RequestAuth(client, accountID, accountName, time.Second*5)
	assert.Nil(t, err)
	assert.Equal(t, login.ErrorCode_SUCCESS, response.GetErrCode())
	assert.NotEqual(t, 0, response.GetPlayerId())
	assert.NotEmpty(t, response.GetGateIp())
	assert.NotEqual(t, 0, response.GetGatePort())

	tokenExpire := time.Unix(response.GetExpire(), 0)
	assert.True(t, tokenExpire.After(time.Now()))
	assert.NotEmpty(t, response.GetGateToken())
}

// 目的： 验证同一个账号再次认证， 得到的玩家 ID 相同
// 步骤:
// 	1. 连接登录服，分配账号ID并发起认证请求，并记录玩家 ID
// 	2. 创建另一个客户端，用同一个账号进行认证
// 期望：
//  1. 两次登录得到的玩家 ID 相同
func Test_AuthAgain(t *testing.T) {

	accountID := global.AllocAccountID()
	accountName := utils.GenerateAccountName(accountID)

	client1 := connect.NewTestClient(config.ServerAddr, config.ClientVersion)
	response1, _ := utils.RequestAuth(client1, accountID, accountName, time.Second*5)

	client2 := connect.NewTestClient(config.ServerAddr, config.ClientVersion)
	response2, _ := utils.RequestAuth(client2, accountID, accountName, time.Second*5)

	assert.Equal(t, response1.GetPlayerId(), response2.GetPlayerId())
}
