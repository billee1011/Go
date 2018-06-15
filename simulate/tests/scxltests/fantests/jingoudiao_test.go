package fantests

import (
	"steve/client_pb/room"
	"steve/simulate/global"
	"steve/simulate/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

//jingoudiao 共同步骤
// 玩家换三张后的牌
//庄家0手牌 11, 11, 13, 22, 22, 13, 24, 24, 15, 25, 21, 12, 23, 14
//1玩家手牌 21, 21, 22, 12, 23, 23, 14, 14, 15, 11, 12, 13, 24
//2玩家手牌 31, 31, 33, 33, 34, 34, 35, 35, 36, 36, 37, 37, 38
//3玩家手牌 31, 31, 33, 33, 34, 34, 35, 35, 36, 36, 37, 37, 38
// 庄家出21,2碰21,出11,庄碰11,出12,2碰12,出22,庄碰22,出23,2碰23,出24,庄碰24,出14,2碰14,出13,庄碰13,出25,
// 2摸牌15,能自摸
// 最后庄家摸牌11
// 结果牌型：
// 庄家：peng{11,22,24,13},handcard{15}
// 2玩家：peng{21,12,23,14},handcard{15}
func jingoudiao(t *testing.T) *utils.DeskData {
	params := global.NewCommonStartGameParams()
	params.BankerSeat = 0
	params.Cards = [][]uint32{
		{31, 31, 33, 22, 22, 13, 24, 24, 15, 25, 21, 12, 23, 14},
		{31, 31, 33, 12, 23, 23, 14, 14, 15, 11, 12, 13, 24},
		{11, 11, 13, 33, 34, 34, 35, 35, 36, 36, 37, 37, 38},
		{21, 21, 22, 33, 34, 34, 35, 35, 36, 36, 37, 37, 38},
	}
	params.WallCards = []uint32{15, 16}
	// 对家换牌
	params.HszDir = room.Direction_Opposite
	params.HszCards = [][]uint32{
		{31, 31, 33},
		{31, 31, 33},
		{11, 11, 13},
		{21, 21, 22},
	}
	params.DingqueColor = []room.CardColor{room.CardColor_CC_TONG, room.CardColor_CC_TONG, room.CardColor_CC_TONG, room.CardColor_CC_TONG}
	deskData, err2 := utils.StartGame(params)
	assert.NotNil(t, deskData)
	assert.Nil(t, err2)

	//开局 0 自询
	assert.Nil(t, utils.WaitZixunNtf(deskData, 0))
	// 0 出 21
	assert.Nil(t, utils.SendChupaiReq(deskData, 0, 21))

	// 1玩家能碰21
	assert.Nil(t, utils.WaitChupaiWenxunNtf(deskData, 1, true, false, false))
	// 1玩家发送碰请求21
	assert.Nil(t, utils.SendPengReq(deskData, 1))
	// 1玩家碰牌成功进入自询
	assert.Nil(t, utils.WaitZixunNtf(deskData, 1))
	// 1 出 11
	assert.Nil(t, utils.SendChupaiReq(deskData, 1, 11))

	// 0玩家能碰11
	assert.Nil(t, utils.WaitChupaiWenxunNtf(deskData, 0, true, false, false))
	// 0玩家发送碰请求11
	assert.Nil(t, utils.SendPengReq(deskData, 0))
	// 0玩家碰牌成功进入自询
	assert.Nil(t, utils.WaitZixunNtf(deskData, 0))
	// 0 出 12
	assert.Nil(t, utils.SendChupaiReq(deskData, 0, 12))

	// 1玩家能碰12
	assert.Nil(t, utils.WaitChupaiWenxunNtf(deskData, 1, true, false, false))
	// 1玩家发送碰请求12
	assert.Nil(t, utils.SendPengReq(deskData, 1))
	// 1玩家碰牌成功进入自询
	assert.Nil(t, utils.WaitZixunNtf(deskData, 1))
	// 1 出 22
	assert.Nil(t, utils.SendChupaiReq(deskData, 1, 22))

	// 0玩家能碰22
	assert.Nil(t, utils.WaitChupaiWenxunNtf(deskData, 0, true, false, false))
	// 0玩家发送碰请求22
	assert.Nil(t, utils.SendPengReq(deskData, 0))
	// 0玩家碰牌成功进入自询
	assert.Nil(t, utils.WaitZixunNtf(deskData, 0))
	// 0 出 23
	assert.Nil(t, utils.SendChupaiReq(deskData, 0, 23))

	// 1玩家能碰23
	assert.Nil(t, utils.WaitChupaiWenxunNtf(deskData, 1, true, false, false))
	// 1玩家发送碰请求23
	assert.Nil(t, utils.SendPengReq(deskData, 1))
	// 1玩家碰牌成功进入自询
	assert.Nil(t, utils.WaitZixunNtf(deskData, 1))
	// 1 出 24
	assert.Nil(t, utils.SendChupaiReq(deskData, 1, 24))

	// 0玩家能碰24
	assert.Nil(t, utils.WaitChupaiWenxunNtf(deskData, 0, true, false, false))
	// 0玩家发送碰请求24
	assert.Nil(t, utils.SendPengReq(deskData, 0))
	// 0玩家碰牌成功进入自询
	assert.Nil(t, utils.WaitZixunNtf(deskData, 0))
	// 0 出 14
	assert.Nil(t, utils.SendChupaiReq(deskData, 0, 14))

	// 1玩家能碰14
	assert.Nil(t, utils.WaitChupaiWenxunNtf(deskData, 1, true, true, false))
	// 1玩家发送碰请求14
	assert.Nil(t, utils.SendPengReq(deskData, 1))
	// 1玩家碰牌成功进入自询
	assert.Nil(t, utils.WaitZixunNtf(deskData, 1))
	// 1 出 13
	assert.Nil(t, utils.SendChupaiReq(deskData, 1, 13))

	// 0玩家能碰13
	assert.Nil(t, utils.WaitChupaiWenxunNtf(deskData, 0, true, false, false))
	// 0玩家发送碰请求13
	assert.Nil(t, utils.SendPengReq(deskData, 0))
	// 0玩家碰牌成功进入自询
	assert.Nil(t, utils.WaitZixunNtf(deskData, 0))
	// 0 出 25
	assert.Nil(t, utils.SendChupaiReq(deskData, 0, 25))
	// 1玩家摸牌15
	utils.CheckZixunNtf(t, deskData, 1, false, false, true)
	return deskData
}

//TestFan_Jingoudiao_Dianpao 金钩钓立即点炮自摸测试
// 1玩家放弃自摸胡15,出15,庄家点炮14
//期望赢分:4
func TestFan_Jingoudiao_Dianpao(t *testing.T) {
	deskData := jingoudiao(t)
	// 1玩家发送弃自摸,胡15
	utils.SendQiReq(deskData, 1)
	// 1 出 15
	assert.Nil(t, utils.SendChupaiReq(deskData, 1, 15))
	//0玩家能点炮胡15
	assert.Nil(t, utils.WaitChupaiWenxunNtf(deskData, 0, false, true, false))
	// 0 号玩家发送胡请求
	assert.Nil(t, utils.SendHuReq(deskData, 0))

	// 检测所有玩家收到点炮通知x
	utils.CheckHuNotify(t, deskData, []int{0}, 1, 15, room.HuType_HT_DIANPAO)
	// 检测金钩钓点炮分数
	utils.CheckInstantSettleScoreNotify(t, deskData, 0, 4)
}

//TestFan_Jingoudiao_Zimo 金钩钓立即结算自摸测试
// 1玩家自摸胡15
//期望赢分:24 = 4 * 2 *3
func TestFan_Jingoudiao_Zimo(t *testing.T) {
	deskData := jingoudiao(t)
	// 1玩家发送胡,胡15
	utils.SendHuReq(deskData, 1)
	// 1 号玩家发送胡请求
	assert.Nil(t, utils.SendHuReq(deskData, 1))

	// 检测所有玩家收到自摸通知
	utils.CheckHuNotify(t, deskData, []int{1}, 1, 15, room.HuType_HT_ZIMO)

	// 检测金钩钓自摸分数,金钩钓4倍*自摸2倍
	winScro := 4 * 2 * (len(deskData.Players) - 1)
	utils.CheckInstantSettleScoreNotify(t, deskData, 1, int64(winScro))
}
