package fantests

import (
	"steve/client_pb/common"
	"steve/client_pb/room"
	"steve/simulate/global"
	"steve/simulate/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

//qidui 共同步骤
// 玩家换三张后的牌
//庄家0手牌 11, 11, 12,12, 21, 21, 22, 22, 13, 13, 23, 23, 14, 15
//1玩家手牌 11, 11, 12, 12, 21, 21, 22, 22, 13, 13, 23, 23, 14
//2玩家手牌 16, 16, 17 17, 18, 18, 31, 31, 32, 32, 33, 33, 34
//3玩家手牌 16, 16, 17 17, 18, 18, 31, 31, 32, 32, 33, 33, 34
// 庄家出15,1玩家摸牌36,出36,2玩家摸牌36,出36,3玩家摸牌36,出36
// 最后庄家摸牌14 结果庄家和1牌型
// 11, 11, 12,12, 21, 21, 22, 22, 13, 13, 23, 23, 14, 14
func qidui(t *testing.T) *utils.DeskData {
	params := global.NewCommonStartGameParams()
	params.GameID = common.GameId_GAMEID_XUEZHAN // 血战
	params.PeiPaiGame = "scxz"
	params.BankerSeat = 0
	params.Cards = [][]uint32{
		{16, 16, 17, 12, 21, 21, 22, 22, 13, 13, 23, 23, 14, 15},
		{16, 16, 17, 12, 21, 21, 22, 22, 13, 13, 23, 23, 14},
		{11, 11, 12, 17, 18, 18, 31, 31, 32, 32, 33, 33, 34},
		{11, 11, 12, 17, 18, 18, 31, 31, 32, 32, 33, 33, 34},
	}
	params.WallCards = []uint32{36, 36, 36, 14, 14, 37}
	// 对家换牌
	params.HszDir = room.Direction_Opposite
	params.HszCards = [][]uint32{
		{16, 16, 17},
		{16, 16, 17},
		{11, 11, 12},
		{11, 11, 12},
	}
	params.DingqueColor = []room.CardColor{room.CardColor_CC_TONG, room.CardColor_CC_TONG, room.CardColor_CC_TONG, room.CardColor_CC_TONG}
	deskData, err2 := utils.StartGame(params)
	assert.NotNil(t, deskData)
	assert.Nil(t, err2)

	//开局 0 自询
	assert.Nil(t, utils.WaitZixunNtf(deskData, 0))
	// 0 出 15
	assert.Nil(t, utils.SendChupaiReq(deskData, 0, 15))
	//1 自询
	assert.Nil(t, utils.WaitZixunNtf(deskData, 1))
	// 1 出 36
	assert.Nil(t, utils.SendChupaiReq(deskData, 1, 36))
	//2 自询
	assert.Nil(t, utils.WaitZixunNtf(deskData, 2))
	// 2 出 36
	assert.Nil(t, utils.SendChupaiReq(deskData, 2, 36))
	//3 自询
	assert.Nil(t, utils.WaitZixunNtf(deskData, 3))
	// 3 出 36
	assert.Nil(t, utils.SendChupaiReq(deskData, 3, 36))
	return deskData
}

//TestFan_Qidui_Zimo_SCXZ 七对立即结算自摸测试
// 庄摸14,自摸14
//期望赢分：4*2*3=24
func TestFan_Qidui_Zimo_SCXZ(t *testing.T) {
	deskData := qidui(t)
	//0 号玩家摸牌14后 检测自询,能自摸
	utils.CheckZixunNtf(t, deskData, 0, false, false, true)
	// 0 号玩家发送胡请求
	assert.Nil(t, utils.SendHuReq(deskData, 0))

	// 检测所有玩家收到自摸通知
	utils.CheckHuNotify(t, deskData, []int{0}, 0, 14, room.HuType_HT_ZIMO)

	// 检测七对自摸分数
	winScro := 4 * 2 * (len(deskData.Players) - 1)
	utils.CheckInstantSettleScoreNotify(t, deskData, 0, int64(winScro))
}

//TestFan_Qidui_Dianpao 七对立即点炮自摸测试
// 庄摸14,能自摸14，不自摸，出14,1玩家点炮14
//期望赢分：4
func TestFan_Qidui_Dianpao_SCXZ(t *testing.T) {
	deskData := qidui(t)
	//0 号玩家摸牌后 检测自询
	assert.Nil(t, utils.WaitZixunNtf(deskData, 0))
	//0 号玩家出14
	assert.Nil(t, utils.SendChupaiReq(deskData, 0, 14))
	//1 号玩家收到到出牌问询能点炮
	assert.Nil(t, utils.WaitChupaiWenxunNtf(deskData, 1, false, true, false))
	// 1 号玩家发送胡请求
	assert.Nil(t, utils.SendHuReq(deskData, 1))

	// 检测所有玩家收到点炮通知x
	utils.CheckHuNotify(t, deskData, []int{1}, 0, 14, room.HuType_HT_DIANPAO)
	// 检测七对点炮分数
	utils.CheckInstantSettleScoreNotify(t, deskData, 1, 4)
}
