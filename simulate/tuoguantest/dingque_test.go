package tuoguantest

import (
	"steve/client_pb/room"
	"steve/simulate/global"
	"steve/simulate/utils"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Test_DingqueTuoguan 测试定缺时，退出房间托管
// 步骤：
//	1. 登录4个用户，并且申请开局, 执行换三张
//  2. 用户0打出8Ｗ
//  3. 用户1请求退出游戏
// 期望：
// 	1. 最迟1秒后，用户0,2，3收到用户1弃通知， 用户1不会收到弃通知
func Test_DingqueTuoguan(t *testing.T) {
	params := global.NewCommonStartGameParams()
	params.DingqueColor = nil
	deskData, err := utils.StartGame(params)
	assert.Nil(t, err)
	assert.NotNil(t, deskData)

	assert.Nil(t, utils.SendDingqueReq(0, deskData, room.CardColor_CC_WAN))
	assert.Nil(t, utils.SendDingqueReq(1, deskData, room.CardColor_CC_WAN))
	assert.Nil(t, utils.SendDingqueReq(2, deskData, room.CardColor_CC_WAN))
	assert.Nil(t, utils.SendQuitReq(deskData, 3))
	assert.Nil(t, utils.WaitDingqueFinish(deskData, time.Second*2, nil, []int{0, 1, 2}))
}
