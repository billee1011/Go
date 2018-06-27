package connecttests

import (
	"steve/simulate/config"
	"steve/simulate/connect"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Test_LoginDisconnect 测试登录服定时断开连接的功能
// 步骤：
//	1. 连接登录服
//	2. 等待 66 秒
func Test_LoginDisconnect(t *testing.T) {
	loginClient := connect.NewTestClient(config.GetLoginServerAddr(), config.GetClientVersion())
	assert.NotNil(t, loginClient)
	assert.False(t, loginClient.Closed())

	time.Sleep(66 * time.Second)
	assert.True(t, loginClient.Closed())
}

// Test_LoginDisconnectN like Test_LoginDisconnect , but test multi connects
func Test_LoginDisconnectN(t *testing.T) {
	var testCount = 10

	wg := sync.WaitGroup{}
	wg.Add(testCount)

	for i := 0; i < testCount; i++ {
		go func() {
			defer wg.Done()
			Test_LoginDisconnect(t)
		}()
	}
	wg.Wait()
}
