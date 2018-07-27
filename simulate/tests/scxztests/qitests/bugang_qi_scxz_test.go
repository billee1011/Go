package qitests

import (
	"steve/client_pb/common"
	"steve/simulate/global"
	"steve/simulate/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test_SCXZ_Bugang_qi 补杠弃测试
// 开始游戏后，庄家出9W，1 号玩家请求碰9W再打出6W，继续行牌，每个玩家出9B，直到1号玩家摸到9W请求补杠
// 期望：
// 1号玩家点弃后，继续在自询状态，等待1号玩家出牌
func Test_SCXZ_Bugang_qi(t *testing.T) {
	params := global.NewCommonStartGameParams()
	params.GameID = common.GameId_GAMEID_XUEZHAN // 血战
	params.PeiPaiGame = "scxz"
	params.BankerSeat = 0
	// 庄家的初始手牌： 11,11,11,11,12,12,12,12,13,13,13,39,31,19
	params.Cards[0][13] = 39
	params.Cards[0][12] = 31
	params.Cards[0][11] = 19
	// 1 号玩家初始手牌： 15,15,15,15,16,16,16,16,17,17,29,19,19
	params.Cards[1][12] = 19
	params.Cards[1][11] = 19
	params.Cards[1][10] = 29
	// 2 号玩家初始手牌： 21,21,21,21,22,22,22,22,23,23,23,17,39
	params.Cards[2][12] = 17
	params.Cards[2][11] = 39
	// 3 号玩家初始手牌： 25,25,25,25,26,26,26,26,27,27,27,27,39
	params.Cards[3][12] = 39
	// 墙牌改为 8W, 8W, 8W, 9W， 3B
	params.WallCards = []uint32{18, 18, 18, 19, 33}

	deskData, err := utils.StartGame(params)
	assert.Nil(t, err)
	assert.NotNil(t, deskData)
	// 庄家出 9W
	assert.Nil(t, utils.WaitZixunNtf(deskData, 0))
	assert.Nil(t, utils.SendChupaiReq(deskData, 0, 19))
	// 1 号玩家等可碰通知， 然后请求碰， 再打出6W
	assert.Nil(t, utils.WaitChupaiWenxunNtf(deskData, 1, true, false, false))
	assert.Nil(t, utils.SendPengReq(deskData, 1))
	assert.Nil(t, utils.WaitZixunNtf(deskData, 1))
	assert.Nil(t, utils.SendChupaiReq(deskData, 1, 16))
	// 2 号玩家等待自询通知， 然后打出9筒
	assert.Nil(t, utils.WaitZixunNtf(deskData, 2))
	assert.Nil(t, utils.SendChupaiReq(deskData, 2, 39))
	// 3 号玩家等待自询通知， 然后打出9筒
	assert.Nil(t, utils.WaitZixunNtf(deskData, 3))
	assert.Nil(t, utils.SendChupaiReq(deskData, 3, 39))
	// 0 号玩家等待自询通知， 然后打出9筒
	assert.Nil(t, utils.WaitZixunNtf(deskData, 0))
	assert.Nil(t, utils.SendChupaiReq(deskData, 0, 39))
	// 1 号玩家等待自询通知， 然后请求杠 9万
	assert.Nil(t, utils.WaitZixunNtf(deskData, 1))
	// 发送弃请求
	assert.Nil(t, utils.SendQiReq(deskData, 1))
}
