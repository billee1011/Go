package qitests

import (
	"steve/client_pb/common"
	"steve/client_pb/msgid"
	"steve/client_pb/room"
	"steve/simulate/global"
	"steve/simulate/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test_SCXZ_Qiangganghu_qi 枪杠胡弃测试
// 开始游戏后，庄家出9W，1 号玩家请求碰再打出6W，继续行牌，所有玩家出9B，直到1号玩家摸到9W请求补杠9W
// 期望：
// 1号玩家点补杠后，2号玩家可以抢杠胡9W，2号玩家点弃，1号玩家补杠成功，
// 1号玩家将收到杠通知，杠后摸牌后将收到自询通知
func Test_SCXZ_Qiangganghu_qi(t *testing.T) {
	params := global.NewCommonStartGameParams()
	params.GameID = common.GameId_GAMEID_XUEZHAN // 血战
	params.PeiPaiGame = "scxz"
	params.BankerSeat = 0
	// 庄家的初始手牌： 11,11,11,11,12,12,12,12,13,13,13,39,31,19
	params.Cards[0][13] = 39
	params.Cards[0][12] = 31
	params.Cards[0][11] = 19
	// 1 号玩家初始手牌： 15,15,15,15,16,16,16,16,17,17,29,19,19
	params.Cards[1][12] = 19
	params.Cards[1][11] = 19
	params.Cards[1][10] = 29
	// 2 号玩家初始手牌： 21,21,21,21,22,22,22,22,23,23,23,17,39
	params.Cards[2][12] = 17
	params.Cards[2][11] = 39
	// 3 号玩家初始手牌： 25,25,25,25,26,26,26,26,27,27,27,27,39
	params.Cards[3][12] = 39
	params.DingqueColor[2] = room.CardColor_CC_TONG
	// 墙牌改为 8W, 8W, 8W, 9W， 3B
	params.WallCards = []uint32{18, 18, 18, 19, 33}

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
	// 发送补杠请求
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
		assert.Equal(t, i == 2, ntf.GetSelfCan())
	}
	// 2号玩家发送弃请求
	assert.Nil(t, utils.SendQiReq(deskData, 2))
	//检查1号家补杠的通知
	utils.CheckGangNotify(t, deskData, gangPlayer.Player.GetID(), gangPlayer.Player.GetID(), uint32(19), room.GangType_BuGang)
	// 1号玩家等待接收到自询
	utils.WaitZixunNtf(deskData, 1)
}
