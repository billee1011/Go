package qitests

import (
	"steve/simulate/global"
	"steve/simulate/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test_AnGang_Qi 测试自询暗杠弃
// 庄家出3W,所有客户端接受出牌通知,其他玩家无操作,下家摸牌1B
// 期望：
// 下家摸牌1B后，将收到自询通知能暗杠6W,7W
// 下家不杠选择弃，等待下家出牌
func Test_AnGang_Qi(t *testing.T) {
	// utils.StartGameParams
	thisParams := global.NewCommonStartGameParams()
	thisParams.WallCards = []uint32{31, 39}
	deskData, err := utils.StartGame(thisParams)
	assert.Nil(t, err)
	assert.NotNil(t, deskData)
	//庄家出牌
	assert.Nil(t, utils.SendChupaiReq(deskData, deskData.BankerSeat, uint32(13)))
	//所有客户端接受出牌通知
	utils.CheckChuPaiNotify(t, deskData, uint32(13), deskData.BankerSeat)
	//下家这时候摸到牌后，进入自询状态，自询状态下可以暗杠
	xjPlayer := utils.CheckMoPaiNotify(t, deskData, (deskData.BankerSeat+1)%len(deskData.Players), 31)

	// 发送弃请求
	assert.Nil(t, utils.SendQiReq(deskData, xjPlayer.Seat))
}
