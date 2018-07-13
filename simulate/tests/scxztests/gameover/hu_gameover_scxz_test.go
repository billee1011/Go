package tests

import (
	msgid "steve/client_pb/msgid"
	"steve/client_pb/room"
	"steve/simulate/global"
	"steve/simulate/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test_SCXZ_Duo_Dianpao_GameOver 多家点炮胡，正常玩家不足，游戏结束
// 开始游戏后，庄家出9W，其他玩家都可以胡
// 期望：
// 1. 1，2,3号玩家收到出牌问询通知，且可以胡
// 2. 1,2,3号玩家发送胡请求后，所有玩家收到胡通知， 胡牌者为1,2,3号玩家，胡类型为点炮，胡的牌为9W
// 3. 牌墙还剩余2张,正常玩家只剩下1人，游戏结束
func Test_SCXZ_Duo_Dianpao_GameOver(t *testing.T) {
	var Int9W uint32 = 19
	params := global.NewCommonStartGameParams()
	params.GameID = room.GameId_GAMEID_XUEZHAN // 血战
	params.PeiPaiGame = "scxz"
	params.BankerSeat = 0
	params.WallCards = []uint32{31, 33}
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
	// assert.Nil(t, utils.SendQiReq(deskData, hu2Seat))

	// 3 号玩家发送胡请求
	assert.Nil(t, utils.SendHuReq(deskData, hu3Seat))

	// 检测所有玩家收到点炮通知
	utils.CheckHuNotify(t, deskData, []int{hu1Seat, hu2Seat, hu3Seat}, bankerSeat, Int9W, room.HuType_HT_DIANPAO)
	// utils.CheckHuNotify(t, deskData, []int{hu1Seat, hu3Seat}, bankerSeat, Int9W, room.HuType_HT_DIANPAO)

	// 检测0, 2, 3玩家收到点炮结算通知
	utils.CheckDianPaoSettleNotify(t, deskData, []int{hu1Seat, hu2Seat, hu3Seat}, bankerSeat, Int9W, room.HuType_HT_DIANPAO)

	// 游戏结束
	utils.WaitGameOverNtf(t, deskData)
}

// Test_SCXZ_ZiMo_GameOver 测试自摸，正常玩家不足，游戏结束
// 步骤：庄家天胡自摸，庄下家自摸，对家自摸
// 期望:  牌墙还剩余2张,正常玩家只剩下1人，游戏结束
func Test_SCXZ_ZiMo_GameOver(t *testing.T) {
	params := global.NewCommonStartGameParams()
	params.Cards = [][]uint32{
		{11, 11, 11, 11, 12, 12, 12, 12, 13, 13, 13, 13, 14, 14},
		{15, 15, 15, 15, 16, 16, 16, 16, 17, 17, 17, 17, 18},
		{21, 21, 21, 21, 22, 22, 22, 22, 23, 23, 23, 23, 24},
		{25, 25, 25, 25, 26, 26, 26, 26, 27, 27, 27, 27, 28},
	}
	params.HszCards = [][]uint32{}
	params.GameID = room.GameId_GAMEID_XUEZHAN // 血战
	params.PeiPaiGame = "scxz"
	params.IsHsz = false // 不换三张
	// 根据座位设置玩家金币数
	params.PlayerSeatGold = map[int]uint64{
		0: 1000000, 1: 1000000, 2: 1000000, 3: 1000000,
	}
	params.WallCards = []uint32{18, 24, 31, 31}
	params.DingqueColor = []room.CardColor{room.CardColor_CC_TIAO, room.CardColor_CC_TIAO, room.CardColor_CC_TONG, room.CardColor_CC_TIAO}
	deskData, err := utils.StartGame(params)
	assert.NotNil(t, deskData)
	assert.Nil(t, err)

	banker := params.BankerSeat
	// 庄家自摸11
	assert.Nil(t, utils.WaitZixunNtf(deskData, banker))
	assert.Nil(t, utils.SendHuReq(deskData, banker))
	// 检测所有玩家收到天胡通知
	utils.CheckHuNotify(t, deskData, []int{banker}, banker, 11, room.HuType_HT_TIANHU)

	// 庄家下家自摸18
	assert.Nil(t, utils.WaitZixunNtf(deskData, 1))
	assert.Nil(t, utils.SendHuReq(deskData, 1))
	// 检测所有玩家收到地胡胡通知
	utils.CheckHuNotify(t, deskData, []int{1}, 1, 18, room.HuType_HT_DIHU)

	// 庄家对家自摸24
	assert.Nil(t, utils.WaitZixunNtf(deskData, 2))
	assert.Nil(t, utils.SendHuReq(deskData, 2))
	// 检测所有玩家收到地胡通知
	utils.CheckHuNotify(t, deskData, []int{2}, 2, 24, room.HuType_HT_DIHU)

	//游戏结束
	utils.WaitGameOverNtf(t, deskData)

}

// Test_SCXZ_Duo_Qiangganghu_GameOver 测试多家抢杠胡，正常玩家不足，游戏结束
// 开始游戏，庄家出9万， 1 号玩家可以碰，其他玩家不可以杠和胡
// 1号玩家请求碰。 并且打出6万，没人可以碰杠胡。
// 2号玩家摸8万， 打出9筒， 没人可以碰杠胡
// 3号玩家摸8万，打出9筒，没人可以碰杠胡。
// 0号玩家摸8万，并且打出9筒，没人可以碰杠胡
// 1号玩家摸9万，并且请求执行补杠。 0,2,3号玩家可以抢杠胡
// 期望：
// 1. 所有玩家收到等待抢杠胡通知，杠的玩家为1号玩家， 杠的牌为9W， 并且0,2,3号玩家收到的通知中可以抢杠胡
// 2. 0,2,3号玩家请求胡，所有玩家收到胡通知，胡的玩家为0,2,3号玩家，胡的牌为9W， 胡牌来源是1号玩家，胡类型为抢杠胡
// 3. 牌墙有剩余，正常玩家人数不足，游戏结束
func Test_SCXZ_Duo_Qiangganghu_GameOver(t *testing.T) {
	params := global.NewCommonStartGameParams()
	params.GameID = room.GameId_GAMEID_XUEZHAN // 血战
	params.PeiPaiGame = "scxz"
	params.BankerSeat = 0
	// 庄家的初始手牌： 11,11,11,11,12,12,12,12,13,13,13,39,17,19 8w
	params.Cards[0][13] = 39
	params.Cards[0][12] = 17
	params.Cards[0][11] = 19
	// 1 号玩家初始手牌： 15,15,15,15,16,16,16,16,17,27,31,19,19 9w
	params.Cards[1][12] = 19
	params.Cards[1][11] = 19
	params.Cards[1][10] = 31
	params.Cards[1][9] = 27
	// 2 号玩家初始手牌： 21,21,21,21,22,23,22,22,22,23,23,17,39  8w
	params.Cards[2][12] = 17
	params.Cards[2][11] = 39
	// 3 号玩家初始手牌： 25,25,25,25,26,26,26,26,27,27,27,17,39 8w
	params.Cards[3][12] = 39
	params.Cards[3][11] = 17

	// 墙牌改为 8W, 8W, 8W, 9W， 3B
	params.WallCards = []uint32{18, 18, 18, 19, 33}
	// 修改所有定缺颜色为筒
	params.DingqueColor = []room.CardColor{room.CardColor_CC_TONG, room.CardColor_CC_TONG, room.CardColor_CC_TONG, room.CardColor_CC_TONG}
	// 开始游戏
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
	assert.Nil(t, utils.SendGangReq(deskData, 1, 19, room.GangType_BuGang))

	// 所有玩家收到等待抢杠胡通知， 2号玩家可以抢杠胡， 其他玩家不能抢杠胡
	gangPlayer := utils.GetDeskPlayerBySeat(1, deskData)
	for i := 0; i < 4; i++ {
		deskPlayer := utils.GetDeskPlayerBySeat(i, deskData)
		expector, _ := deskPlayer.Expectors[msgid.MsgID_ROOM_WAIT_QIANGGANGHU_NTF]
		ntf := room.RoomWaitQianggangHuNtf{}
		assert.Nil(t, expector.Recv(global.DefaultWaitMessageTime, &ntf))
		assert.Equal(t, uint32(19), ntf.GetCard())
		assert.Equal(t, gangPlayer.Player.GetID(), ntf.GetFromPlayerId())
		if i != 1 {
			assert.Equal(t, true, ntf.GetSelfCan())
		}
	}
	// 0,2,3号玩家发送枪杠胡请求
	assert.Nil(t, utils.SendHuReq(deskData, 0))
	assert.Nil(t, utils.SendHuReq(deskData, 2))
	assert.Nil(t, utils.SendHuReq(deskData, 3))
	// 检测0, 2, 3玩家收到点炮通知
	utils.CheckHuNotify(t, deskData, []int{0, 2, 3}, 1, 19, room.HuType_HT_QIANGGANGHU)

	//等待游戏结束通知等待游戏结束通知
	utils.WaitGameOverNtf(t, deskData)
}
