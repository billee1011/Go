package connecttests

import (
	"fmt"
	"steve/simulate/config"
	"steve/simulate/connect"
	"steve/simulate/global"
	"steve/simulate/utils"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Test_HeartBeat 测试一段时间没有发心跳，连接断开
// 步骤：
// 	1. 登录用户
//  2. 等待 61 秒
// 期望：
//  连接已关闭
func Test_HeartBeat(t *testing.T) {
	player, err := utils.LoginNewPlayer()
	assert.Nil(t, err)
	client := player.GetClient()
	time.Sleep(time.Second * 61)
	assert.True(t, client.Closed())
}

// Test_HeartBeat2 测试多个连接一段时间没有发心跳，连接断开
// 步骤：
// 	1. 登录用户
//  2. 等待 61 秒
// 期望：
//  连接已关闭
func Test_HeartBeat2(t *testing.T) {
	const testCount = 2
	test := func() {
		player, err := utils.LoginNewPlayer()
		assert.Nil(t, err)
		client := player.GetClient()
		time.Sleep(time.Second * 61)
		assert.True(t, client.Closed())
	}

	wg := sync.WaitGroup{}
	wg.Add(testCount)
	for i := 0; i < testCount; i++ {
		go func() {
			test()
			wg.Done()
		}()
	}
	wg.Wait()
}

// Test_NotAuth 测试网关超时未认证，连接断开
// 步骤：
//	1. 通过登录服认证，获取到网关服地址
//  2. 等待 61 秒
// 期望：
// 	连接已关闭
func Test_NotAuth(t *testing.T) {
	loginClient := connect.NewTestClient(config.GetLoginServerAddr(), config.GetClientVersion())
	assert.NotNil(t, loginClient)
	accountID := global.AllocAccountID()
	accountName := utils.GenerateAccountName(accountID)
	loginResp, err := utils.RequestAuth(loginClient, accountID, accountName, global.DefaultWaitMessageTime)
	assert.Nil(t, err)

	gateIP := loginResp.GetGateIp()
	gatePort := loginResp.GetGatePort()
	gateAddr := fmt.Sprintf("%s:%d", gateIP, gatePort)
	gateClient := connect.NewTestClient(gateAddr, config.GetClientVersion())
	assert.NotNil(t, gateClient)
	assert.False(t, gateClient.Closed())

	time.Sleep(61 * time.Second)
	assert.True(t, gateClient.Closed())
}
