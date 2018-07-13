package hutests

import (
	msgid "steve/client_pb/msgid"
	"steve/client_pb/room"
	"steve/simulate/global"
	"steve/simulate/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test_SCXZ_Qiangganghu 测试抢杠胡
// 开始游戏，庄家出9万， 1 号玩家可以碰，其他玩家不可以杠和胡
// 1号玩家请求碰。 并且打出6万，没人可以碰杠胡。
// 2号玩家摸8万， 打出9筒， 没人可以碰杠胡
// 3号玩家摸8万，打出9筒，没人可以碰杠胡。
// 0号玩家摸8万，并且打出9筒，没人可以碰杠胡
// 1号玩家摸9万，并且请求执行补杠。 2号玩家可以抢杠胡
// 期望：
// 1. 所有玩家收到等待抢杠胡通知，杠的玩家为1号玩家， 杠的牌为9W， 并且2号玩家收到的通知中可以抢杠胡
// 2. 2号玩家请求胡，所有玩家收到胡通知，胡的玩家为2号玩家，胡的牌为9W， 胡牌来源是1号玩家，胡类型为抢杠胡
func Test_SCXZ_Qiangganghu(t *testing.T) {
	params := global.NewCommonStartGameParams()
	params.GameID = room.GameId_GAMEID_XUEZHAN // 血战
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

	assert.Nil(t, utils.SendHuReq(deskData, 2))
	utils.CheckHuNotify(t, deskData, []int{2}, 1, 19, room.HuType_HT_QIANGGANGHU)
}

// makeRoomCards(Card1W, Card1W, Card1W, Card1W, Card2W, Card2W, Card2W, Card2W, Card3W, Card3W, Card3W, Card3W, Card4W, Card4W),
// makeRoomCards(Card5W, Card5W, Card5W, Card5W, Card6W, Card6W, Card6W, Card6W, Card7W, Card7W, Card7W, Card7W, Card8W),
// makeRoomCards(Card1T, Card1T, Card1T, Card1T, Card2T, Card2T, Card2T, Card2T, Card3T, Card3T, Card3T, Card3T, Card4T),
// makeRoomCards(Card5T, Card5T, Card5T, Card5T, Card6T, Card6T, Card6T, Card6T, Card7T, Card7T, Card7T, Card7T, Card8T),

//Test_SCXZ_Hued_NotQiangGangHu 测试胡过玩家能不能在提示抢杠胡
// 步骤：胡牌后抢杠胡（东打5万，南碰打4万，东胡，南摸5万补杠，东可以抢杠胡）
// 期望：南补杠成功，东提示抢杠胡
func Test_SCXZ_Hued_NotQiangGangHu(t *testing.T) {
	params := global.NewCommonStartGameParams()
	params.Cards = [][]uint32{
		{11, 11, 11, 11, 12, 12, 12, 12, 13, 13, 13, 23, 23, 24},
		{19, 19, 19, 19, 18, 18, 18, 18, 17, 17, 27, 27, 26},
		{21, 21, 21, 21, 22, 22, 22, 22, 23, 23, 13, 14, 15},
		{29, 29, 29, 29, 28, 28, 28, 28, 27, 27, 14, 15, 15},
	}
	params.GameID = room.GameId_GAMEID_XUEZHAN // 血战
	params.PeiPaiGame = "scxz"
	params.WallCards = []uint32{15, 37, 39, 39, 39, 38, 38}
	params.IsHsz = true // 换三张
	params.HszCards = [][]uint32{
		{23, 23, 24},
		{27, 27, 26},
		{13, 14, 15},
		{14, 15, 15},
	}
	params.HszDir = room.Direction_Opposite
	params.DingqueColor = []room.CardColor{room.CardColor_CC_TONG, room.CardColor_CC_TIAO, room.CardColor_CC_TIAO, room.CardColor_CC_TIAO}
	deskData, err := utils.StartGame(params)
	assert.NotNil(t, deskData)
	assert.Nil(t, err)

	banker := params.BankerSeat
	// 庄家出5万
	assert.Nil(t, utils.WaitZixunNtf(deskData, banker))
	assert.Nil(t, utils.SendChupaiReq(deskData, banker, 15))

	// 1玩家能碰5万，并碰5万
	assert.Nil(t, utils.WaitChupaiWenxunNtf(deskData, 1, true, false, false))
	assert.Nil(t, utils.SendPengReq(deskData, 1))

	// 1玩家出牌4万
	assert.Nil(t, utils.WaitZixunNtf(deskData, 1))
	assert.Nil(t, utils.SendChupaiReq(deskData, 1, 14))

	// 庄家胡4万
	assert.Nil(t, utils.WaitChupaiWenxunNtf(deskData, banker, false, true, false))
	assert.Nil(t, utils.SendHuReq(deskData, banker))

	// 1玩家摸5万补杠
	utils.CheckZixunNtf(t, deskData, 1, true, true, false)
	assert.Nil(t, utils.SendGangReq(deskData, 1, 15, room.GangType_BuGang))

	// 杠成功，收到杠通知
	xiajiaID := utils.GetDeskPlayerBySeat(1, deskData).Player.GetID()
	utils.CheckGangNotify(t, deskData, uint64(xiajiaID), uint64(xiajiaID), uint32(15), room.GangType_BuGang)
}
