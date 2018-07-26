package connecttests

import (
	"steve/client_pb/hall"
	"steve/client_pb/msgid"
	"steve/simulate/config"
	"steve/simulate/connect"
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

// Test_HeartBeat3 测试发送任意消息后，不发送心跳，连接不会断开
// 步骤：
// 1. 登录用户
// 2. 等待 30s，发送任意消息
// 3. 再等待 31s， 连接未关闭
// 4. 再等待 30s， 连接关闭
func Test_HeartBeat3(t *testing.T) {
	const testCount = 2
	test := func() {
		player, err := utils.LoginNewPlayer()
		assert.Nil(t, err)
		client := player.GetClient()
		time.Sleep(time.Second * 30)
		_, err = player.GetClient().SendPackage(utils.CreateMsgHead(msgid.MsgID_HALL_GET_PLAYER_INFO_REQ), &hall.HallGetPlayerInfoReq{})
		assert.Nil(t, err)
		time.Sleep(time.Second * 31)

		assert.False(t, client.Closed())
		time.Sleep(time.Second * 30)
		assert.True(t, client.Closed())
	}
	test()
}

// Test_NotAuth 测试网关超时未认证，连接断开
// 步骤：
//	1. 通过登录服认证，获取到网关服地址
//  2. 等待 61 秒
// 期望：
// 	连接已关闭
func Test_NotAuth(t *testing.T) {
	gateClient := connect.NewTestClient(config.GetGatewayServerAddr(), config.GetClientVersion())
	assert.NotNil(t, gateClient)
	assert.False(t, gateClient.Closed())

	time.Sleep(61 * time.Second)
	assert.True(t, gateClient.Closed())
}
