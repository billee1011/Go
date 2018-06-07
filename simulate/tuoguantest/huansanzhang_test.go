package tuoguantest

import (
	"steve/client_pb/room"
	"steve/simulate/global"
	"steve/simulate/utils"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestHuansanzhangTuoguan 测试换三张时，退出房间托管
// 步骤：
//	1. 登录4个用户，并且申请开局, 执行发牌
//  2. 用户0-2在收到发牌完成通知后，请求换三张
//  3. 用户 3 请求退出游戏，
// 期望：
// 	1. 最迟1秒后，用户0-2收到换三张完成通知， 用户3不会收到换三张完成通知
func TestHuansanzhangTuoguan(t *testing.T) {
	params := global.NewCommonStartGameParams()
	params.HszCards = nil
	params.DingqueColor = nil
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
	assert.Nil(t, utils.SendQuitReq(deskData, 3))
	assert.Nil(t, utils.WaitHuansanzhangFinish(deskData, time.Second*2, []int{0, 1, 2}, nil, 3))
}

// TestHuansanzhangTuoguan 测试换三张时，退出房间托管
// 步骤：
//	1. 登录4个用户，并且申请开局, 执行发牌
//  2. 用户0-2在收到发牌完成通知后，请求换三张
//  3. 用户 3 超时，自动执行换三张
//  4. 用户0-2在收到换三张完成通知后，请求定缺，花色为万
//  5. 用户 3 超时，自动执行定缺，受到托管请求
// 期望：
// 	1.用户3超时2次后接受到托管请求
func Test_HuansanzhangTuoguan2(t *testing.T) {
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

	assert.Nil(t, utils.SendDingqueReq(0, deskData, room.CardColor_CC_WAN))
	assert.Nil(t, utils.SendDingqueReq(1, deskData, room.CardColor_CC_WAN))
	assert.Nil(t, utils.SendDingqueReq(2, deskData, room.CardColor_CC_WAN))
	assert.Nil(t, utils.WaitDingqueFinish(deskData, time.Second*16, nil, []int{0, 1, 2, 3}))

	assert.Nil(t, utils.WaitTuoGuanNtf(deskData, time.Second*6, []int{3}))

}
