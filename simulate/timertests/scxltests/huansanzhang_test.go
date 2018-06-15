package timertests

import (
	"steve/client_pb/room"
	"steve/simulate/global"
	"steve/simulate/utils"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Test_HuansanzhangTimeOut 测试换三张超时
// 步骤：
//	1. 登录4个用户，并且申请开局
//  2. 用户0-2在收到发牌通知后，请求换三张，用户3不请求换三张
// 期望：
// 	1. 16秒后，所有用户收到定缺完成通知
func Test_HuansanzhangTimeOut(t *testing.T) {
	params := global.NewCommonStartGameParams()
	params.HszCards = nil
	params.DingqueColor = nil
	params.HszDir = room.Direction_Opposite
	deskData, err := utils.StartGame(params)
	assert.Nil(t, err)
	assert.NotNil(t, deskData)

	hszCards := [][]uint32{
		{11, 11, 11},
		{15, 15, 15},
		{21, 21, 21},
		{25, 25, 25},
	}
	assert.Nil(t, utils.SendHuansanzhangReq(0, deskData, hszCards[0], true))
	assert.Nil(t, utils.SendHuansanzhangReq(1, deskData, hszCards[1], true))
	assert.Nil(t, utils.SendHuansanzhangReq(2, deskData, hszCards[2], true))
	assert.Nil(t, utils.WaitHuansanzhangFinish(deskData, time.Second*16, []int{0, 1, 2, 3}, []uint32{15, 15, 15}, 3))
}
