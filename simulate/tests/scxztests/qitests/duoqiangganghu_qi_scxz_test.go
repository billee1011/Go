package qitests

import (
	 "steve/client_pb/msgId"
	"steve/client_pb/room"
	"steve/simulate/global"
	"steve/simulate/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test_SCXZ_Duo_Qiangganghu_Qi 测试多家抢杠胡,1个玩家点弃
// 开始游戏，庄家出9万， 1 号玩家可以碰，其他玩家不可以杠和胡
// 1号玩家请求碰。 并且打出6万，没人可以碰杠胡。
// 2号玩家摸8万， 打出9筒， 没人可以碰杠胡
// 3号玩家摸8万，打出9筒，没人可以碰杠胡。
// 0号玩家摸8万，并且打出9筒，没人可以碰杠胡
// 1号玩家摸9万，并且请求执行补杠。 0,2,3号玩家可以抢杠胡
// 期望：
// 1. 所有玩家收到等待抢杠胡通知，杠的玩家为1号玩家， 杠的牌为9W， 并且0,2,3号玩家收到的通知中可以抢杠胡
// 2. 0,3号玩家请求胡,2号玩家选择弃，所有玩家收到胡通知，胡的玩家为0,3号玩家，胡的牌为9W， 胡牌来源是1号玩家，胡类型为抢杠胡
func Test_SCXZ_Duo_Qiangganghu_Qi(t *testing.T) {
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
	// 0,3号玩家发送枪杠胡请求,2号玩家发送弃牌
	assert.Nil(t, utils.SendHuReq(deskData, 0))
	assert.Nil(t, utils.SendQiReq(deskData, 2))
	assert.Nil(t, utils.SendHuReq(deskData, 3))
	// 检测0, 2, 3玩家收到点炮通知
	utils.CheckHuNotify(t, deskData, []int{0, 3}, 1, 19, room.HuType_HT_QIANGGANGHU)
}
