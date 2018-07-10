package tests

import (
	msgid "steve/client_pb/msgId"
	"steve/client_pb/room"
	"steve/simulate/global"
	"steve/simulate/utils"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
)

// Test_SCXZ_Mopai 摸牌测试
// 步骤：庄家出1筒
// 期望：所有玩家收到庄家下家摸牌通知
func Test_SCXZ_Mopai(t *testing.T) {
	params := global.NewCommonStartGameParams()
	// 庄家最后一张牌改为1筒
	params.Cards[0][13] = 31
	params.GameID = room.GameId_GAMEID_XUEZHAN // 血战
	params.PeiPaiGame = "scxz"
	deskData, err := utils.StartGame(params)

	assert.NotNil(t, deskData)
	assert.Nil(t, err)

	// 庄家出1筒
	assert.Nil(t, utils.WaitZixunNtf(deskData, params.BankerSeat))
	zjPlayer := utils.GetDeskPlayerBySeat(params.BankerSeat, deskData)
	zjClient := zjPlayer.Player.GetClient()
	zjClient.SendPackage(utils.CreateMsgHead(msgid.MsgID_ROOM_CHUPAI_REQ), &room.RoomChupaiReq{
		Card: proto.Uint32(31),
	})

	mopaiSeat := (params.BankerSeat + 1) % len(deskData.Players)
	mopaiPlayer := utils.GetDeskPlayerBySeat(mopaiSeat, deskData)
	mopaiPlayerID := mopaiPlayer.Player.GetID()
	// 所有玩家收到庄家下家摸牌通知
	for _, deskPlayer := range deskData.Players {
		expector, _ := deskPlayer.Expectors[msgid.MsgID_ROOM_MOPAI_NTF]
		ntf := room.RoomMopaiNtf{}
		assert.Nil(t, expector.Recv(global.DefaultWaitMessageTime, &ntf))
		assert.Equal(t, mopaiPlayerID, ntf.GetPlayer())
		assert.Equal(t, false, ntf.GetBack())
		if deskPlayer.Seat == mopaiSeat {
			assert.Equal(t, uint32(31), ntf.GetCard()) // 摸了一张 1 筒
		}
	}
}

// Test_SCXZ_Hued_NotMopai 测试胡过玩家是否还能在摸牌
// 步骤：庄家天胡自摸，庄下家即1玩家（摸牌，出牌5条），2玩家（摸牌，出牌5条），3玩家（摸牌，出牌5条）
// 期望: 庄没有摸牌，直接跳过，庄下家即1玩家摸牌
func Test_SCXZ_Hued_NotMopai(t *testing.T) {
	params := global.NewCommonStartGameParams()
	params.Cards = [][]uint32{
		{11, 11, 11, 12, 12, 12, 13, 13, 13, 14, 14, 14, 31, 31},
		{15, 15, 15, 25, 26, 16, 16, 16, 17, 17, 17, 17, 18},
		{21, 21, 21, 25, 26, 22, 22, 22, 23, 23, 23, 23, 24},
		{35, 35, 35, 25, 26, 36, 36, 36, 37, 37, 37, 37, 38},
	}
	params.HszCards = [][]uint32{}
	params.GameID = room.GameId_GAMEID_XUEZHAN // 血战
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
	var Int5t uint32 = 25
	assert.Nil(t, utils.WaitZixunNtf(deskData, 1))
	assert.Nil(t, utils.SendChupaiReq(deskData, 1, Int5t))

	// 2号玩家摸牌自询，出牌
	assert.Nil(t, utils.WaitZixunNtf(deskData, 2))
	assert.Nil(t, utils.SendChupaiReq(deskData, 2, Int5t))
	// 3号玩家摸牌自询，出牌
	assert.Nil(t, utils.WaitZixunNtf(deskData, 3))
	assert.Nil(t, utils.SendChupaiReq(deskData, 3, Int5t))

	// 跳过庄家， 庄下家即1玩家摸牌
	assert.Nil(t, utils.WaitZixunNtf(deskData, 1))
}

// Test_GiveUp_MoPai 测试认输玩家是否还能在摸牌
// 步骤：（所有玩家金币数设置为2）
//1.庄家出牌
//2.下家摸牌出牌，庄家明杠下家，下家钱不足认输，
//3.庄家明杠后出牌，其他人没有出牌问询
// 期望:直接跳过下家，由上家摸牌
func Test_GiveUp_MoPai(t *testing.T) {
	params := global.NewCommonStartGameParams()
	params.Cards = [][]uint32{
		{11, 11, 11, 12, 12, 12, 13, 13, 13, 14, 14, 14, 31, 31},
		{15, 15, 15, 11, 26, 16, 16, 16, 17, 17, 17, 17, 18},
		{21, 21, 21, 25, 26, 22, 22, 22, 23, 23, 23, 23, 24},
		{35, 35, 35, 25, 26, 36, 36, 36, 37, 37, 37, 37, 38},
	}
	params.HszCards = [][]uint32{}
	params.GameID = room.GameId_GAMEID_XUEZHAN // 血战
	params.PeiPaiGame = "scxz"
	params.IsHsz = false // 不换三张
	// 根据座位设置玩家金币数
	params.PlayerSeatGold = map[int]uint64{
		0: 12, 1: 1, 2: 1, 3: 1,
	}
	params.WallCards = []uint32{31, 31, 32, 33}
	params.DingqueColor = []room.CardColor{room.CardColor_CC_TONG, room.CardColor_CC_TONG, room.CardColor_CC_TONG, room.CardColor_CC_TONG}
	deskData, err := utils.StartGame(params)
	assert.NotNil(t, deskData)
	assert.Nil(t, err)

	banker := params.BankerSeat
	// 庄家出牌31
	assert.Nil(t, utils.WaitZixunNtf(deskData, banker))
	assert.Nil(t, utils.SendChupaiReq(deskData, banker, 31))
	// 下家出牌11
	assert.Nil(t, utils.WaitZixunNtf(deskData, 1))
	assert.Nil(t, utils.SendChupaiReq(deskData, 1, 11))
	// 庄家收到出牌问询通知，能杠，能碰，并发送杠11请求
	assert.Nil(t, utils.WaitChupaiWenxunNtf(deskData, banker, true, false, true))
	utils.SendGangReq(deskData, banker, 11, room.GangType_MingGang)

	// 检测所有玩家收到杠通知
	bankerGPlayerID := utils.GetDeskPlayerBySeat(banker, deskData).Player.GetID()
	bGPlayerID := utils.GetDeskPlayerBySeat(1, deskData).Player.GetID()
	utils.CheckGangNotify(t, deskData, bankerGPlayerID, bGPlayerID, 11, room.GangType_MingGang)

	// 庄家摸牌31,出牌31
	assert.Nil(t, utils.WaitZixunNtf(deskData, banker))
	assert.Nil(t, utils.SendChupaiReq(deskData, banker, 31))

	// （跳过下家摸牌）对家摸牌32
	assert.Nil(t, utils.WaitZixunNtf(deskData, 2))
}
