package tuoguantest

import (
	"steve/client_pb/room"
	"steve/simulate/global"
	"steve/simulate/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test_ZixunTimeOut01 测试定缺超时
// 步骤：
//	1. 登录4个用户，并且申请开局, 执行换三张和定缺
//  2. 用户0也就是庄家起手托管
// 期望：
// 	1. 托管后，所有用户收到玩家0的出牌通知
func Test_ZixunTuoguan01(t *testing.T) {
	params := global.NewCommonStartGameParams()
	params.GameID = room.GameId_GAMEID_XUEZHAN // 血战
	params.PeiPaiGame = "scxz"
	deskData, err := utils.StartGame(params)
	assert.Nil(t, err)
	assert.NotNil(t, deskData)
	assert.Nil(t, utils.WaitZixunNtf(deskData, deskData.BankerSeat))
	//庄家退出桌子自动托管
	assert.Nil(t, utils.SendQuitReq(deskData, 0))
	utils.CheckChuPaiNotifyWithSeats(t, deskData, uint32(25), deskData.BankerSeat, []int{1, 2, 3})
}

// Test_DingqueTimeOut02 测试定缺超时
// 步骤：
//	1. 登录4个用户，并且申请开局, 执行换三张和定缺
//  2. 用户0也就是庄家选择天胡,下家对家尾家摸牌打1b,再次到庄家摸1w
//  3. 庄家可胡,但退出桌子进行托管
// 期望：
// 	1. 庄家之前开过胡了,因此托管后庄家自动胡,且墙牌没了,所以胡牌类型时海底捞
func Test_ZixunTuoguan02(t *testing.T) {
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
	//庄家退出桌子自动托管
	assert.Nil(t, utils.SendQuitReq(deskData, 0))
	utils.CheckHuNotifyBySeats(t, deskData, []int{0}, 0, uint32(11), room.HuType_HT_HAIDILAO, []int{1, 2, 3})
}
