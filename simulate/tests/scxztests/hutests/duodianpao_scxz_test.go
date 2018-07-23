package hutests

import (
	"steve/client_pb/common"
	"steve/client_pb/room"
	"steve/simulate/global"
	"steve/simulate/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test_SCXZ_Duo_Dianpao 多家点炮胡测试
// 开始游戏后，庄家出9W，其他玩家都可以胡
// 期望：
// 1. 1，2,3号玩家收到出牌问询通知，且可以胡
// 2. 1,2,3号玩家发送胡请求后，所有玩家收到胡通知， 胡牌者为1,2,3号玩家，胡类型为点炮，胡的牌为9W
func Test_SCXZ_Duo_Dianpao(t *testing.T) {
	var Int9W uint32 = 19
	params := global.NewCommonStartGameParams()
	params.GameID = common.GameId_GAMEID_XUEZHAN // 血战
	params.PeiPaiGame = "scxz"
	params.BankerSeat = 0
	hu1Seat, hu2Seat, hu3Seat := 1, 2, 3
	bankerSeat := params.BankerSeat
	// 修改所有定缺颜色为筒
	params.DingqueColor = []room.CardColor{room.CardColor_CC_TONG, room.CardColor_CC_TONG, room.CardColor_CC_TONG, room.CardColor_CC_TONG}
	// 庄家的最后一张牌改为 9W
	params.Cards[bankerSeat][13] = 19
	// 1 号玩家最后1张牌改为 9W
	params.Cards[hu1Seat][12] = 19
	// 2 号玩家最后1张牌改为 9W
	params.Cards[hu2Seat][12] = 19
	// 3 号玩家最后1张牌改为 9W
	params.Cards[hu3Seat][12] = 19

	deskData, err := utils.StartGame(params)
	assert.Nil(t, err)
	assert.NotNil(t, deskData)
	// 庄家出 9W
	assert.Nil(t, utils.SendChupaiReq(deskData, bankerSeat, Int9W))

	// 1 号玩家收到出牌问询通知， 可以胡
	assert.Nil(t, utils.WaitChupaiWenxunNtf(deskData, hu1Seat, false, true, false))
	// 2 号玩家收到出牌问询通知， 可以胡
	assert.Nil(t, utils.WaitChupaiWenxunNtf(deskData, hu2Seat, false, true, false))
	// 3 号玩家收到出牌问询通知， 可以胡
	assert.Nil(t, utils.WaitChupaiWenxunNtf(deskData, hu3Seat, false, true, false))

	// 1 号玩家发送胡请求
	assert.Nil(t, utils.SendHuReq(deskData, hu1Seat))

	// 2 号玩家发送胡请求
	assert.Nil(t, utils.SendHuReq(deskData, hu2Seat))

	// 3 号玩家发送胡请求
	assert.Nil(t, utils.SendHuReq(deskData, hu3Seat))

	// 检测所有玩家收到点炮通知
	utils.CheckHuNotify(t, deskData, []int{hu1Seat, hu2Seat, hu3Seat}, bankerSeat, Int9W, room.HuType_HT_DIANPAO)

	// 检测0, 2, 3玩家收到点炮结算通知
	utils.CheckDianPaoSettleNotify(t, deskData, []int{hu1Seat, hu2Seat, hu3Seat}, bankerSeat, Int9W, room.HuType_HT_DIANPAO)
}
