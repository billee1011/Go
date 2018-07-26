package tests

import (
	"steve/client_pb/common"
	"steve/client_pb/room"
	"steve/simulate/global"
	"steve/simulate/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test_SCXZ_Hued_CPWX_NotGang 测试胡过玩家是否还能在杠
// 步骤：庄家天胡自摸，庄下家即1,出牌4万,庄可以杠
// 期望: 庄没有出牌杠问询，而庄对家即2,摸牌
func Test_SCXZ_Hued_CPWX_NotGang(t *testing.T) {
	params := global.NewCommonStartGameParams()
	params.Cards = [][]uint32{
		{11, 11, 11, 12, 12, 12, 13, 13, 13, 14, 14, 14, 31, 31},
		{15, 15, 15, 14, 16, 16, 16, 16, 17, 17, 17, 17, 18},
		{21, 21, 21, 21, 22, 22, 22, 22, 23, 23, 23, 23, 24},
		{25, 25, 25, 25, 26, 26, 26, 26, 27, 27, 27, 27, 28},
	}
	params.HszCards = [][]uint32{}
	params.GameID = common.GameId_GAMEID_XUEZHAN // 血战
	params.PeiPaiGame = "scxz"
	params.IsHsz = false // 不换三张
	params.WallCards = []uint32{31, 31, 32, 33}
	params.DingqueColor = []room.CardColor{room.CardColor_CC_TIAO, room.CardColor_CC_TIAO, room.CardColor_CC_TIAO, room.CardColor_CC_TIAO}
	deskData, err := utils.StartGame(params)
	assert.NotNil(t, deskData)
	assert.Nil(t, err)

	banker := params.BankerSeat
	// 庄家自摸
	assert.Nil(t, utils.WaitZixunNtf(deskData, banker))
	assert.Nil(t, utils.SendHuReq(deskData, banker))
	var Int1B uint32 = 31
	// 检测所有玩家收到天胡通知
	utils.CheckHuNotify(t, deskData, []int{banker}, banker, Int1B, room.HuType_HT_TIANHU)

	// 庄下家即1玩家出牌
	var Int4w uint32 = 14
	assert.Nil(t, utils.WaitZixunNtf(deskData, 1))
	assert.Nil(t, utils.SendChupaiReq(deskData, 1, Int4w))

	// 2号玩家摸牌自询（庄可以杠，但胡过，不提示进行出牌问询）
	assert.Nil(t, utils.WaitZixunNtf(deskData, 2))
	// assert.Nil(t, utils.WaitChupaiWenxunNtf(deskData, banker, false, false, true))
}

// Test_SCXZ_Hued_CPWX_NotHu 测试胡过玩家是否还能在胡
// 步骤：庄家天胡自摸，庄下家即1,出牌2万,庄可以杠
// 期望: 庄没有出牌胡问询，而庄对家即2,摸牌
func Test_SCXZ_Hued_CPWX_NotHu(t *testing.T) {
	params := global.NewCommonStartGameParams()
	params.Cards = [][]uint32{
		{11, 11, 11, 12, 12, 12, 13, 13, 13, 14, 14, 14, 31, 31},
		{15, 15, 15, 12, 16, 16, 16, 16, 17, 17, 17, 17, 18},
		{21, 21, 21, 21, 22, 22, 22, 22, 23, 23, 23, 23, 24},
		{25, 25, 25, 25, 26, 26, 26, 26, 27, 27, 27, 27, 28},
	}
	params.HszCards = [][]uint32{}
	params.GameID = common.GameId_GAMEID_XUEZHAN // 血战
	params.PeiPaiGame = "scxz"
	params.IsHsz = false // 不换三张
	params.WallCards = []uint32{31, 31, 32, 33}
	params.DingqueColor = []room.CardColor{room.CardColor_CC_TIAO, room.CardColor_CC_TIAO, room.CardColor_CC_TIAO, room.CardColor_CC_TIAO}
	deskData, err := utils.StartGame(params)
	assert.NotNil(t, deskData)
	assert.Nil(t, err)

	banker := params.BankerSeat
	// 庄家自摸
	assert.Nil(t, utils.WaitZixunNtf(deskData, banker))
	assert.Nil(t, utils.SendHuReq(deskData, banker))
	var Int1B uint32 = 31
	// 检测所有玩家收到天胡通知
	utils.CheckHuNotify(t, deskData, []int{banker}, banker, Int1B, room.HuType_HT_TIANHU)

	// 庄下家即1玩家出牌
	var Int2w uint32 = 12
	assert.Nil(t, utils.WaitZixunNtf(deskData, 1))
	assert.Nil(t, utils.SendChupaiReq(deskData, 1, Int2w))

	// 2号玩家摸牌自询（庄可以胡，但胡过，不提示进行出牌问询）
	assert.Nil(t, utils.WaitZixunNtf(deskData, 2))
	// assert.Nil(t, utils.WaitChupaiWenxunNtf(deskData, banker, false, true, false))
}
