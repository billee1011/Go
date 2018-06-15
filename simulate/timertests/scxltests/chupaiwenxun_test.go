package timertests

import (
	"steve/simulate/global"
	"steve/simulate/utils"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Test_ChuPaiWenXubTimeOut 测试出牌问询超时
// 步骤：
//	1. 登录4个用户，并且申请开局, 执行换三张，定缺
//  2. 用户0-2在收到换三张完成通知后，请求定缺，花色为万. 用户3不请求定缺
// 期望：
// 	1. 16秒后，所有用户收到定缺完成通知
func Test_ChuPaiWenXubTimeOut(t *testing.T) {
	params := global.NewCommonStartGameParams()
	params.WallCards = []uint32{31}
	params.Cards[0] = []uint32{11, 11, 11, 31, 12, 12, 12, 32, 13, 13, 13, 33, 14, 18}
	deskData, err := utils.StartGame(params)
	assert.Nil(t, err)
	assert.NotNil(t, deskData)

	assert.Nil(t, utils.SendChupaiReq(deskData, 0, 18))
	assert.Nil(t, utils.WaitMoPaiNtf(deskData, time.Second*16, []int{0, 1, 2, 3}, 31, 1))
}
