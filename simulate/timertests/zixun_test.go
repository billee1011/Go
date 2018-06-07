package timertests

import (
	"steve/client_pb/room"
	"steve/simulate/global"
	"steve/simulate/utils"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Test_ZixunTimeOut01 测试定缺超时
// 步骤：
//	1. 登录4个用户，并且申请开局, 执行换三张和定缺
//  2. 用户0也就是庄家挂机,时间到后进入托管打牌
// 期望：
// 	1. 16秒后，所有用户收到玩家0的出牌通知
func Test_ZixunTimeOut01(t *testing.T) {
	params := global.NewCommonStartGameParams()
	deskData, err := utils.StartGame(params)
	assert.Nil(t, err)
	assert.NotNil(t, deskData)
	//庄家不出牌,倒计时结束后自动
	assert.Nil(t, utils.WaitZixunNtf(deskData, deskData.BankerSeat))
	utils.CheckChuPaiNotify0(t, deskData, time.Second*16, uint32(25), deskData.BankerSeat)
}

// Test_DingqueTimeOut02 测试定缺超时
// 步骤：
//	1. 登录4个用户，并且申请开局, 执行换三张和定缺
//  2. 用户0也就是庄家选择天胡,下家对家尾家摸牌打1b,再次到庄家摸1w
//  3. 庄家可胡,但不进行操作,等待倒计时结束,自动处理
// 期望：
// 	1. 庄家之前开过胡了,因此16秒后，庄家自动胡,且墙牌没了,所以胡牌类型时海底捞
func Test_ZixunTimeOut02(t *testing.T) {
	params := global.NewCommonStartGameParams()
	params.Cards[0][1] = 19
	params.Cards[1][1] = 33
	params.Cards[3][1] = 15
	params.DingqueColor[0] = room.CardColor_CC_TONG
	params.WallCards = []uint32{31, 31, 31, 11}
	params.HszCards[0] = []uint32{11, 19, 11}
	deskData, err := utils.StartGame(params)
	assert.Nil(t, err)
	assert.NotNil(t, deskData)
	//庄家不出牌,倒计时结束后自动
	assert.Nil(t, utils.WaitZixunNtf(deskData, deskData.BankerSeat))
	assert.Nil(t, utils.SendHuReq(deskData, 0))
	utils.CheckHuNotify(t, deskData, []int{0}, 0, uint32(11), room.HuType_HT_TIANHU)
	assert.Nil(t, utils.WaitZixunNtf(deskData, 1))
	assert.Nil(t, utils.SendChupaiReq(deskData, 1, uint32(31)))
	utils.CheckChuPaiNotify(t, deskData, uint32(31), 1)
	assert.Nil(t, utils.WaitZixunNtf(deskData, 2))
	assert.Nil(t, utils.SendChupaiReq(deskData, 2, uint32(31)))
	utils.CheckChuPaiNotify(t, deskData, uint32(31), 2)
	assert.Nil(t, utils.WaitZixunNtf(deskData, 3))
	assert.Nil(t, utils.SendChupaiReq(deskData, 3, uint32(31)))
	utils.CheckChuPaiNotify(t, deskData, uint32(31), 3)
	assert.Nil(t, utils.WaitZixunNtf(deskData, 0))
	utils.CheckHuNotify0(t, deskData, time.Second*16, []int{0}, 0, uint32(11), room.HuType_HT_HAIDILAO)
}
