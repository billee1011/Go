package fantests

import (
	"steve/client_pb/room"
	"steve/simulate/global"
	"steve/simulate/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

//duiduihu 共同步骤
// 玩家换三张后的牌
//庄家0手牌 11, 11, 12,12,13,13,14,21,21,21,14,24,22,23,
//1玩家手牌 26, 26, 26,39,39,39,28,28,28,29,29,36,36,
//2玩家手牌 31, 31, 31,32,32,32,33,33,33,34,34,35,35,
//3玩家手牌 17, 17, 17,11,12,13,14,22,22,23,23,24,24,
// 庄家出22,3碰22,出11,庄家碰11,出23,3碰23,出12,庄家碰12,出24,3碰24,出13,庄碰13,出14,3点炮14
// 最后庄家摸牌11
// 结果牌型：
// 庄家：peng{11,12,13},handcard{21,21,15,15}
// 3玩家：peng{22,23,24},handcard{17,17,17,14}
// 庄家：2碰下，手中2刻子单一张,结算对对胡自摸加倍 = 4
//3玩家: 3碰下，手中1刻子单一张,结算对对胡2
func duiduihu(t *testing.T) *utils.DeskData {
	params := global.NewCommonStartGameParams()
	params.GameID = room.GameId_GAMEID_XUEZHAN // 血战
	params.PeiPaiGame = "scxz"
	params.BankerSeat = 0
	params.Cards = [][]uint32{
		{31, 31, 31, 12, 13, 13, 14, 21, 21, 21, 14, 24, 22, 23},
		{17, 17, 17, 39, 39, 39, 28, 28, 28, 29, 29, 36, 36},
		{11, 11, 12, 32, 32, 32, 33, 33, 33, 34, 34, 35, 35},
		{26, 26, 26, 11, 12, 13, 14, 22, 22, 23, 23, 24, 24},
	}
	params.WallCards = []uint32{13, 35, 36}
	// 对家换牌
	params.HszDir = room.Direction_Opposite
	params.HszCards = [][]uint32{
		{31, 31, 31},
		{17, 17, 17},
		{11, 11, 12},
		{26, 26, 26},
	}
	params.DingqueColor = []room.CardColor{room.CardColor_CC_TONG, room.CardColor_CC_TONG, room.CardColor_CC_TIAO, room.CardColor_CC_TONG}
	deskData, err2 := utils.StartGame(params)
	assert.NotNil(t, deskData)
	assert.Nil(t, err2)

	//开局 0 自询
	assert.Nil(t, utils.WaitZixunNtf(deskData, 0))
	// 0 出 22
	assert.Nil(t, utils.SendChupaiReq(deskData, 0, 22))

	// 3玩家能碰22
	assert.Nil(t, utils.WaitChupaiWenxunNtf(deskData, 3, true, false, false))
	// 3玩家发送碰请求22
	assert.Nil(t, utils.SendPengReq(deskData, 3))
	// 3玩家碰牌成功进入自询
	assert.Nil(t, utils.WaitZixunNtf(deskData, 3))
	// 3 出 11
	assert.Nil(t, utils.SendChupaiReq(deskData, 3, 11))

	// 0玩家能碰11
	assert.Nil(t, utils.WaitChupaiWenxunNtf(deskData, 0, true, false, false))
	// 0玩家发送碰请求11
	assert.Nil(t, utils.SendPengReq(deskData, 0))
	// 0玩家碰牌成功进入自询
	assert.Nil(t, utils.WaitZixunNtf(deskData, 0))
	// 0 出 23
	assert.Nil(t, utils.SendChupaiReq(deskData, 0, 23))

	// 3玩家能碰23
	assert.Nil(t, utils.WaitChupaiWenxunNtf(deskData, 3, true, true, false))
	// 3玩家发送碰请求23
	assert.Nil(t, utils.SendPengReq(deskData, 3))
	// 3玩家碰牌成功进入自询
	assert.Nil(t, utils.WaitZixunNtf(deskData, 3))
	// 3 出 12
	assert.Nil(t, utils.SendChupaiReq(deskData, 3, 12))

	// 0玩家能碰12
	assert.Nil(t, utils.WaitChupaiWenxunNtf(deskData, 0, true, false, false))
	// 0玩家发送碰请求12
	assert.Nil(t, utils.SendPengReq(deskData, 0))
	// 0玩家碰牌成功进入自询
	assert.Nil(t, utils.WaitZixunNtf(deskData, 0))
	// 0 出 24
	assert.Nil(t, utils.SendChupaiReq(deskData, 0, 24))

	// 3玩家能点炮 24
	assert.Nil(t, utils.WaitChupaiWenxunNtf(deskData, 3, true, false, false))
	// 3玩家发送碰请求24
	assert.Nil(t, utils.SendPengReq(deskData, 3))
	// 3玩家碰牌成功进入自询
	assert.Nil(t, utils.WaitZixunNtf(deskData, 3))
	// 3 出 13
	assert.Nil(t, utils.SendChupaiReq(deskData, 3, 13))

	return deskData
}

//TestFan_Duiduihu_Zimo_SCXZ 对对胡立即结算自摸测试
// 庄放弃动作点炮和碰3玩家的13后，庄摸牌13,自摸13
//期望赢分：12 = 2 * 2 *3
func TestFan_Duiduihu_Zimo_SCXZ(t *testing.T) {
	deskData := duiduihu(t)
	// 0玩家能碰,能胡13,
	assert.Nil(t, utils.WaitChupaiWenxunNtf(deskData, 0, true, true, false))
	// 0玩家放弃点炮13
	utils.SendQiReq(deskData, 0)
	//0 号玩家摸牌13后 检测自询,能自摸
	utils.CheckZixunNtf(t, deskData, 0, false, false, true)
	// 0 号玩家发送胡请求
	assert.Nil(t, utils.SendHuReq(deskData, 0))

	// 检测所有玩家收到自摸通知
	utils.CheckHuNotify(t, deskData, []int{0}, 0, 13, room.HuType_HT_ZIMO)

	// 检测对对胡自摸分数
	winScro := 2 * 2 * (len(deskData.Players) - 1)
	utils.CheckInstantSettleScoreNotify(t, deskData, 0, int64(winScro))
}

//TestFan_Duiduihu_Dianpao_SCXZ 对对胡立即点炮自摸测试
// 庄放弃点炮，去碰3玩家的13后，出14，3玩家点炮14
//期望赢分：2
func TestFan_Duiduihu_Dianpao_SCXZ(t *testing.T) {
	deskData := duiduihu(t)
	// 0玩家能碰,能胡13
	assert.Nil(t, utils.WaitChupaiWenxunNtf(deskData, 0, true, true, false))
	// 0玩家发送碰请求13
	assert.Nil(t, utils.SendPengReq(deskData, 0))
	// 0玩家碰牌成功进入自询
	assert.Nil(t, utils.WaitZixunNtf(deskData, 0))
	// 0 出 14
	assert.Nil(t, utils.SendChupaiReq(deskData, 0, 14))
	// 3玩家能点炮 14
	assert.Nil(t, utils.WaitChupaiWenxunNtf(deskData, 3, false, true, false))
	// 3号玩家发送胡请求
	assert.Nil(t, utils.SendHuReq(deskData, 3))
	// 检测所有玩家收到点炮通知x
	utils.CheckHuNotify(t, deskData, []int{3}, 0, 14, room.HuType_HT_DIANPAO)
	// 检测对对胡点炮分数
	utils.CheckInstantSettleScoreNotify(t, deskData, 3, 2)
}
