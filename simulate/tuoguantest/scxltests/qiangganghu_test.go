package tuoguantest

import (
	 "steve/client_pb/msgId"
	"steve/client_pb/room"
	"steve/simulate/global"
	"steve/simulate/utils"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

//Test_Qiangganghu01 庄家打出9w,下家碰并且打出6w,接下来对家,尾家,庄家分别摸牌并打出9b,下家摸牌,并且请求补杠,此时对家可以抢杠胡
// 期待: 對家托管,默认选过
func Test_Qiangganghu01(t *testing.T) {
	params := global.NewCommonStartGameParams()

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
	assert.Nil(t, utils.SendGangReq(deskData, 1, 19, room.GangType_BuGang))

	// 所有玩家收到等待抢杠胡通知， 2号玩家可以抢杠胡， 其他玩家不能抢杠胡
	gangPlayer := utils.GetDeskPlayerBySeat(1, deskData)
	for i := 0; i < 4; i++ {
		deskPlayer := utils.GetDeskPlayerBySeat(i, deskData)
		expector, _ := deskPlayer.Expectors[msgId.MsgID_ROOM_WAIT_QIANGGANGHU_NTF]
		ntf := room.RoomWaitQianggangHuNtf{}
		assert.Nil(t, expector.Recv(global.DefaultWaitMessageTime, &ntf))
		assert.Equal(t, uint32(19), ntf.GetCard())
		assert.Equal(t, gangPlayer.Player.GetID(), ntf.GetFromPlayerId())
		assert.Equal(t, i == 2, ntf.GetSelfCan())
	}
	assert.Nil(t, utils.SendQuitReq(deskData, 2))
	player := utils.GetDeskPlayerBySeat(1, deskData)
	expector, _ := player.Expectors[msgId.MsgID_ROOM_GANG_NTF]
	ntf := room.RoomGangNtf{}
	assert.Nil(t, expector.Recv(time.Second*10, &ntf))
	assert.Equal(t, uint32(19), ntf.GetCard())
	assert.Equal(t, room.GangType_BuGang, ntf.GetGangType())
	assert.Equal(t, player.Player.GetID(), ntf.GetFromPlayerId())
	assert.Equal(t, player.Player.GetID(), ntf.GetToPlayerId())
}

//Test_Qiangganghu02 庄家打出9w,下家碰并且打出1b,接下来对家自摸,尾家,庄家分别摸牌并打出9b,下家摸牌,并且请求补杠,此时对家可以抢杠胡
// 期待: 對家托管,默认选抢杠胡
func Test_Qiangganghu02(t *testing.T) {
	params := global.NewCommonStartGameParams()

	params.BankerSeat = 0
	// 庄家的初始手牌： 11,11,11,11,12,12,12,12,13,13,13,39,31,19
	params.Cards[0][13] = 39
	params.Cards[0][12] = 31
	params.Cards[0][11] = 19
	// 1 号玩家初始手牌： 15,15,15,15,16,16,16,16,17,17,29,19,19
	params.Cards[1][12] = 19
	params.Cards[1][11] = 19
	params.Cards[1][10] = 29
	params.Cards[1][5] = 31
	// 2 号玩家初始手牌： 21,21,21,21,22,22,22,22,23,23,23,17,18
	params.Cards[2][11] = 17
	params.Cards[2][12] = 18
	// 3 号玩家初始手牌： 25,25,25,25,26,26,26,26,27,27,27,27,39
	params.Cards[3][12] = 39

	params.DingqueColor[2] = room.CardColor_CC_TONG

	// 墙牌改为 8W, 8W, 8W, 9W， 3B
	params.WallCards = []uint32{16, 18, 18, 19, 33}

	deskData, err := utils.StartGame(params)
	assert.Nil(t, err)
	assert.NotNil(t, deskData)
	// 庄家出 9W
	assert.Nil(t, utils.WaitZixunNtf(deskData, 0))
	assert.Nil(t, utils.SendChupaiReq(deskData, 0, 19))
	// 1 号玩家等可碰通知， 然后请求碰， 再打出1b,2号玩家可胡,选择弃
	assert.Nil(t, utils.WaitChupaiWenxunNtf(deskData, 1, true, false, false))
	assert.Nil(t, utils.SendPengReq(deskData, 1))
	assert.Nil(t, utils.SendQiReq(deskData, 2))
	assert.Nil(t, utils.WaitZixunNtf(deskData, 1))
	assert.Nil(t, utils.SendChupaiReq(deskData, 1, 31))
	// 2 号玩家等待自询通知， 然后自摸6w
	assert.Nil(t, utils.WaitZixunNtf(deskData, 2))
	assert.Nil(t, utils.SendHuReq(deskData, 2))
	utils.CheckHuNotify(t, deskData, []int{2}, 2, uint32(16), room.HuType_HT_DIHU)
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
		expector, _ := deskPlayer.Expectors[msgId.MsgID_ROOM_WAIT_QIANGGANGHU_NTF]
		ntf := room.RoomWaitQianggangHuNtf{}
		assert.Nil(t, expector.Recv(global.DefaultWaitMessageTime, &ntf))
		assert.Equal(t, uint32(19), ntf.GetCard())
		assert.Equal(t, gangPlayer.Player.GetID(), ntf.GetFromPlayerId())
		assert.Equal(t, i == 2, ntf.GetSelfCan())
	}
	//玩家2之前开过胡,托管默认给胡
	assert.Nil(t, utils.SendQuitReq(deskData, 2))
	utils.CheckHuNotifyBySeats(t, deskData, []int{2}, 1, uint32(19), room.HuType_HT_QIANGGANGHU, []int{0, 1, 3})
}
