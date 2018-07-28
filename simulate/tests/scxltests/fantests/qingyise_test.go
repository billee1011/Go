package fantests

import (
	"steve/client_pb/room"
	"steve/simulate/global"
	"steve/simulate/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

//qingyise 共同步骤
// 玩家换三张后的牌
//庄家0手牌 12, 13, 14,15, 16, 17, 14, 15, 16, 17, 18, 19, 11, 31
//1玩家手牌 12, 13, 14,15, 16, 17, 14, 15, 16, 17, 18, 19, 11
//2玩家手牌 22, 23, 24, 25, 26, 27, 24, 25, 26, 27, 28, 29, 21
//3玩家手牌 22, 23, 24, 25, 26, 27, 24, 25, 26, 27, 28, 29, 21
// 庄家出31,1玩家摸牌31,出31,2玩家摸牌31,出31,3玩家摸牌31,出31,
// 最后庄家摸牌11 结果庄家和1牌型
// 12, 13, 14,15, 16, 17, 14, 15, 16, 17, 18, 19, 11,11
func qingyise(t *testing.T) *utils.DeskData {
	params := global.NewCommonStartGameParams()
	params.BankerSeat = 0
	params.Cards = [][]uint32{
		{22, 23, 24, 15, 16, 17, 14, 15, 16, 17, 18, 19, 11, 31},
		{22, 23, 24, 15, 16, 17, 14, 15, 16, 17, 18, 19, 11},
		{12, 13, 14, 25, 26, 27, 24, 25, 26, 27, 28, 29, 21},
		{12, 13, 14, 25, 26, 27, 24, 25, 26, 27, 28, 29, 21},
	}
	// 对家换牌
	params.HszDir = room.Direction_Opposite
	params.HszCards = [][]uint32{
		{22, 23, 24},
		{22, 23, 24},
		{14, 12, 13},
		{14, 12, 13},
	}
	params.DingqueColor = []room.CardColor{room.CardColor_CC_TONG, room.CardColor_CC_TONG, room.CardColor_CC_TONG, room.CardColor_CC_TONG}
	params.WallCards = []uint32{31, 31, 31, 11, 33}
	deskData, err2 := utils.StartGame(params)
	assert.NotNil(t, deskData)
	assert.Nil(t, err2)

	//开局 0 自询
	assert.Nil(t, utils.WaitZixunNtf(deskData, 0))
	// 0 出 31
	assert.Nil(t, utils.SendChupaiReq(deskData, 0, 31))
	//1 自询
	assert.Nil(t, utils.WaitZixunNtf(deskData, 1))
	// 1 出 31
	assert.Nil(t, utils.SendChupaiReq(deskData, 1, 31))
	//2 自询
	assert.Nil(t, utils.WaitZixunNtf(deskData, 2))
	// 2 出 31
	assert.Nil(t, utils.SendChupaiReq(deskData, 2, 31))
	//3 自询
	assert.Nil(t, utils.WaitZixunNtf(deskData, 3))
	// 3 出 31
	assert.Nil(t, utils.SendChupaiReq(deskData, 3, 31))
	return deskData
}

//TestFan_Qingyise_Zimo 清一色立即结算自摸测试
// 庄摸牌11,自摸11
//期望赢分：24 = 4 * 2 *3
func TestFan_Qingyise_Zimo(t *testing.T) {
	deskData := qingyise(t)
	//0 号玩家摸牌11后 检测自询,能自摸
	utils.CheckZixunNtf(t, deskData, 0, false, false, true)
	// 0 号玩家发送胡请求
	assert.Nil(t, utils.SendHuReq(deskData, 0))

	// 检测所有玩家收到自摸通知
	utils.CheckHuNotify(t, deskData, []int{0}, 0, 11, room.HuType_HT_ZIMO)

	// 检测清一色自摸分数
	winScro := 4 * 2 * (len(deskData.Players) - 1)
	utils.CheckInstantSettleScoreNotify(t, deskData, 0, int64(winScro), deskData.DiFen)
}

//TestFan_Qingyise_Dianpao 清一色立即点炮自摸测试
// 庄摸牌11,不胡，出11,1玩家点炮11
//期望赢分：4
func TestFan_Qingyise_Dianpao(t *testing.T) {
	deskData := qingyise(t)
	//0 号玩家摸牌后 检测自询
	assert.Nil(t, utils.WaitZixunNtf(deskData, 0))
	//0 号玩家出11
	assert.Nil(t, utils.SendChupaiReq(deskData, 0, 11))
	//1 号玩家收到到出牌问询能点炮
	assert.Nil(t, utils.WaitChupaiWenxunNtf(deskData, 1, false, true, false))
	// 1 号玩家发送胡请求
	assert.Nil(t, utils.SendHuReq(deskData, 1))

	// 检测所有玩家收到点炮通知x
	utils.CheckHuNotify(t, deskData, []int{1}, 0, 11, room.HuType_HT_DIANPAO)
	// 检测清一色点炮分数
	utils.CheckInstantSettleScoreNotify(t, deskData, 1, 4, deskData.DiFen)
}
