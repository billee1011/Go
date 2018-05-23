package timertests

import (
	"steve/client_pb/room"
	"steve/simulate/global"
	"steve/simulate/utils"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Test_DingqueTimeOut 测试定缺超时
// 步骤：
//	1. 登录4个用户，并且申请开局, 执行换三张
//  2. 用户0-2在收到换三张完成通知后，请求定缺，花色为万. 用户3不请求定缺
// 期望：
// 	1. 11秒后，所有用户收到定缺完成通知
func Test_DingqueTimeOut(t *testing.T) {
	params := global.NewCommonStartGameParams()
	params.DingqueColor = nil
	deskData, err := utils.StartGame(params)
	assert.Nil(t, err)
	assert.NotNil(t, deskData)

	assert.Nil(t, utils.SendDingqueReq(0, deskData, room.CardColor_CC_WAN))
	assert.Nil(t, utils.SendDingqueReq(1, deskData, room.CardColor_CC_WAN))
	assert.Nil(t, utils.SendDingqueReq(2, deskData, room.CardColor_CC_WAN))
	assert.Nil(t, utils.WaitDingqueFinish(deskData, time.Second*11, nil))
}
