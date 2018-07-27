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
//  2. 用户0-2在收到换三张完成通知后，请求定缺，花色为万
//  3. 用户 3 请求退出游戏，
// 期望：
// 	1. 最迟1秒后，用户0-2收到定缺完成通知， 用户3不会收到定缺完成通知
func Test_DingqueTuoguan(t *testing.T) {
	params := global.NewCommonStartGameParams()
	params.DingqueColor = nil
	params.GameID = common.GameId_GAMEID_XUEZHAN // 血战
	params.PeiPaiGame = "scxz"
	deskData, err := utils.StartGame(params)
	assert.Nil(t, err)
	assert.NotNil(t, deskData)

	assert.Nil(t, utils.SendDingqueReq(0, deskData, room.CardColor_CC_WAN))
	assert.Nil(t, utils.SendDingqueReq(1, deskData, room.CardColor_CC_WAN))
	assert.Nil(t, utils.SendDingqueReq(2, deskData, room.CardColor_CC_WAN))
	assert.Nil(t, utils.SendQuitReq(deskData, 3))
	assert.Nil(t, utils.WaitDingqueFinish(deskData, time.Second*2, nil, []int{0, 1, 2}))
}
