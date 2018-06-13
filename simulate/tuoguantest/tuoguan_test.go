package tuoguantest

import (
	"steve/simulate/global"
	"steve/simulate/utils"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// 目标： Test_OverTimeTuoguan 测试玩家因为超时两次进入托管状态
// 步骤：
//	step1. 开始一局游戏
//  step2. 玩家 0 出 1W, 玩家 1 可以碰。 玩家 1 收到出牌问询通知，但不操作直到超时
//  step3. 玩家 1 摸 1T，收到自询通知后不操作。
// 期望：
//  1. 玩家 1 收到托管通知进入托管状态
func Test_OverTimeTuoguan(t *testing.T) {
	params := global.NewCommonStartGameParams()

	deskData, err := utils.StartGame(params)
	assert.Nil(t, err)

	// 庄家出 1W
	assert.Nil(t, utils.WaitZixunNtf(deskData, 0))
	assert.Nil(t, utils.SendChupaiReq(deskData, 0, 11))

	// 1 号玩家收到出牌问询通知，可以碰
	assert.Nil(t, utils.WaitChupaiWenxunNtf(deskData, 1, true, false, true))

	// 等待 超时
	time.Sleep(global.OverTimeInterval)

	// 1 号玩家收到自询通知
	assert.Nil(t, utils.WaitZixunNtf(deskData, 1))
	// 等待 超时
	time.Sleep(global.OverTimeInterval)

	assert.Nil(t, utils.WaitTuoGuanNtf(deskData, global.DefaultWaitMessageTime, []int{1}))
}
