package tuoguantest

import (
	"steve/client_pb/room"
	"steve/simulate/global"
	"steve/simulate/utils"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Test_ChupauiwenxunTuoguan 测试出牌问询时，退出房间托管
// 步骤：
//	1. 登录4个用户，并且申请开局, 执行换三张,定缺
//  2. 用户0-2在收到换三张完成通知后，请求定缺，花色为万
//  3. 用户1请求退出游戏，用户1执行摸牌
// 期望：
// 	1. 最迟1秒后，用户0-2收到用户1摸牌通知， 用户3不会收到定缺完成通知
func Test_ChupauiwenxunTuoguan(t *testing.T) {
	params := global.NewCommonStartGameParams()
	params.WallCards = []uint32{31}
	params.Cards[0] = []uint32{11, 11, 11, 31, 12, 12, 12, 32, 13, 13, 13, 33, 14, 18}
	params.GameID = room.GameId_GAMEID_XUEZHAN // 血战
	params.PeiPaiGame = "scxz"
	deskData, err := utils.StartGame(params)
	assert.Nil(t, err)
	assert.NotNil(t, deskData)

	assert.Nil(t, utils.SendChupaiReq(deskData, 0, 18))

	assert.Nil(t, utils.SendQuitReq(deskData, 1))
	assert.Nil(t, utils.WaitMoPaiNtf(deskData, time.Second*2, []int{0, 2, 3}, 31, -1))
}

// Test_ChupaiTuoguan 庄家定缺万，这时候手牌为万条筒的组合，因为定缺牌为万，所以优先决策出最大的万字牌
func Test_ChupaiTuoguan(t *testing.T) {
	params := global.NewCommonStartGameParams()
	params.WallCards = []uint32{32, 32, 32, 32}
	params.Cards[0][12] = 31
	params.Cards[0][13] = 31
	params.DingqueColor[0] = room.CardColor_CC_WAN
	params.GameID = room.GameId_GAMEID_XUEZHAN // 血战
	params.PeiPaiGame = "scxz"
	deskData, err := utils.StartGame(params)
	assert.Nil(t, err)
	assert.NotNil(t, deskData)
	time.Sleep(time.Second * 10)
	utils.CheckChuPaiNotify(t, deskData, uint32(13), 0)
}
