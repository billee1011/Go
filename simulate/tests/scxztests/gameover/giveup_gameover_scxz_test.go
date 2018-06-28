package tests

import (
	"steve/client_pb/room"
	"steve/simulate/global"
	"steve/simulate/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

//Test_SCXZ_AnGang_GiveUp_GameOver 测试暗杠后，其他玩家钱不够都认输，认输后，正常状态玩家不足，游戏结束
//步骤：所有玩家金币数只有1,庄家起手暗杠，其他玩家钱不足都认输
//期望：暗杠后，正常状态玩家不足，游戏结束
func Test_SCXZ_AnGang_GiveUp_GameOver(t *testing.T) {
	params := global.NewCommonStartGameParams()
	params.HszCards = [][]uint32{}
	params.GameID = room.GameId_GAMEID_XUEZHAN // 血战
	params.PeiPaiGame = "scxz"
	params.IsHsz = false // 不换三张
	// 根据座位设置玩家金币数
	params.PlayerSeatGold = map[int]uint64{
		0: 8, 1: 2, 2: 2, 3: 2,
	}
	params.DingqueColor = []room.CardColor{room.CardColor_CC_TIAO, room.CardColor_CC_TIAO, room.CardColor_CC_TIAO, room.CardColor_CC_TIAO}
	deskData, err := utils.StartGame(params)
	assert.NotNil(t, deskData)
	assert.Nil(t, err)
	banker := params.BankerSeat
	// 庄家暗杠
	assert.Nil(t, utils.WaitZixunNtf(deskData, banker))
	assert.Nil(t, utils.SendGangReq(deskData, banker, 11, room.GangType_AnGang))

	playeID := utils.GetDeskPlayerBySeat(banker, deskData).Player.GetID()
	// 其他玩家收到暗杠通知
	utils.CheckGangNotify(t, deskData, playeID, playeID, uint32(11), room.GangType_AnGang)

	// 游戏结束
	utils.WaitGameOverNtf(t, deskData)
}

//Test_SCXZ_BuGang_GiveUp_GameOver 测试补杠后，其他玩家钱不够都认输，认输后，正常状态玩家不足，游戏结束
//步骤： 所有玩家金币数只有1
// 1.庄家出12
// 2.下家碰12,下家打出17
// 3.对家摸到18,对家打出22
// 4.尾家摸到18,尾家打出27
// 5.庄家摸到18,庄家打出16
// 6.下家摸到19,此时下家可以补杠12
// 7.下家选择补杠12,其他玩家钱不足都认输
//期待:所有人收到下家补杠的广播后，正常状态玩家不足，游戏结束
func Test_SCXZ_BuGang_GiveUp_GameOver(t *testing.T) {
	param := global.NewCommonStartGameParams()
	param.GameID = room.GameId_GAMEID_XUEZHAN // 血战
	param.PeiPaiGame = "scxz"
	param.BankerSeat = 0
	// 根据座位设置玩家金币数
	param.PlayerSeatGold = map[int]uint64{
		0: 1, 1: 8, 2: 1, 3: 1,
	}
	param.Cards[0][4] = 16
	param.Cards[0][5] = 16
	param.Cards[0][6] = 16
	param.Cards[1][4] = 12
	param.Cards[1][5] = 12
	param.Cards[1][6] = 12
	param.WallCards = []uint32{18, 18, 18, 19, 33, 39, 38}
	deskData, err := utils.StartGame(param)
	assert.NotNil(t, deskData)
	assert.Nil(t, err)
	utils.WaitZixunNtf(deskData, deskData.BankerSeat)
	//庄家出牌,出12
	utils.SendChupaiReq(deskData, deskData.BankerSeat, uint32(12))
	//检查出牌响应
	utils.CheckChuPaiNotify(t, deskData, uint32(12), deskData.BankerSeat)
	//下家请求碰12
	xjSeat := (deskData.BankerSeat + 1) % len(deskData.Players)
	utils.SendPengReq(deskData, xjSeat)
	//检查碰的通知
	utils.CheckPengNotify(t, deskData, xjSeat, 12)
	//碰成功后收到自询通知
	utils.CheckZixunNotify(t, deskData, xjSeat)
	//下家出牌请求
	utils.SendChupaiReq(deskData, xjSeat, 17)
	//下家出牌响应
	utils.CheckChuPaiNotify(t, deskData, 17, xjSeat)
	//对家摸牌(自寻)响应
	djSeat := (xjSeat + 1) % len(deskData.Players)
	utils.CheckMoPaiNotify(t, deskData, djSeat, 18)
	//对家出牌请求
	utils.SendChupaiReq(deskData, djSeat, 22)
	//对家出牌响应
	utils.CheckChuPaiNotify(t, deskData, 22, djSeat)
	//尾家摸牌(自寻)响应
	wjSeat := (djSeat + 1) % len(deskData.Players)
	utils.CheckMoPaiNotify(t, deskData, wjSeat, 18)
	//尾家出牌请求
	utils.SendChupaiReq(deskData, wjSeat, 27)
	//尾家出牌响应
	utils.CheckChuPaiNotify(t, deskData, 27, wjSeat)
	//庄家摸牌(自寻)响应
	utils.CheckMoPaiNotify(t, deskData, deskData.BankerSeat, 18)
	//庄家出牌请求
	utils.SendChupaiReq(deskData, deskData.BankerSeat, 16)
	//庄家出牌响应
	utils.CheckChuPaiNotify(t, deskData, 16, deskData.BankerSeat)
	//下家摸牌(自寻)响应
	utils.CheckMoPaiNotify(t, deskData, xjSeat, 19)
	//下家补杠
	utils.SendGangReq(deskData, xjSeat, 12, room.GangType_BuGang)
	//下家补杠响应
	player := utils.GetDeskPlayerBySeat(xjSeat, deskData)
	utils.CheckGangNotify(t, deskData, player.Player.GetID(), player.Player.GetID(), 12, room.GangType_BuGang)

	// 游戏结束
	utils.WaitGameOverNtf(t, deskData)
}

//Test_SCXZ_MingGang_GiveUp_GameOver 测试明杠后，其他玩家钱不够都认输，认输后，正常状态玩家不足，游戏结束
//步骤：所有玩家金币数只有1
//1.庄家出牌15,下家明杠15，庄家钱不足认输
//2.对家出牌16,下家明杠16，对家钱不足认输
//3.上家出牌17,下家明杠17，上家钱不足认输
//期望：正常状态玩家不足，游戏结束
func Test_SCXZ_MingGang_GiveUp_GameOver(t *testing.T) {
	param := global.NewCommonStartGameParams()
	param.GameID = room.GameId_GAMEID_XUEZHAN // 血战
	param.PeiPaiGame = "scxz"
	param.BankerSeat = 0
	param.IsHsz = false // 不换三张
	// 根据座位设置玩家金币数
	param.PlayerSeatGold = map[int]uint64{
		0: 1, 1: 9, 2: 1, 3: 1,
	}
	param.Cards = [][]uint32{
		{11, 11, 11, 15, 12, 12, 12, 12, 39, 13, 13, 13, 14, 14},
		{15, 15, 15, 11, 16, 16, 16, 21, 17, 17, 17, 25, 39},
		{21, 21, 21, 16, 22, 22, 22, 22, 39, 23, 23, 23, 24},
		{25, 25, 25, 17, 26, 26, 26, 26, 39, 27, 27, 27, 28},
	}
	param.DingqueColor = []room.CardColor{room.CardColor_CC_TONG, room.CardColor_CC_TONG, room.CardColor_CC_TONG, room.CardColor_CC_TONG}
	param.WallCards = []uint32{33, 33, 34, 34, 35, 36, 37, 31, 31, 32, 32, 32, 35, 35, 38, 38}
	deskData, err := utils.StartGame(param)
	assert.NotNil(t, deskData)
	assert.Nil(t, err)
	// 庄家
	bankerPlayerID := utils.GetDeskPlayerBySeat(param.BankerSeat, deskData).Player.GetID()
	//庄家自询，出牌
	utils.WaitZixunNtf(deskData, deskData.BankerSeat)
	utils.SendChupaiReq(deskData, deskData.BankerSeat, uint32(15))
	//检查出牌响应
	utils.CheckChuPaiNotify(t, deskData, uint32(15), deskData.BankerSeat)

	// 下家明杠庄家15
	xiaPlayerID := utils.GetDeskPlayerBySeat(1, deskData).Player.GetID()
	// 1 号玩家收到出牌问询通知， 可以杠
	assert.Nil(t, utils.WaitChupaiWenxunNtf(deskData, 1, true, false, true))
	// 1 号玩家发送杠请求
	assert.Nil(t, utils.SendGangReq(deskData, 1, 15, room.GangType_MingGang))
	// 检测所有玩家收到杠通知,并等待自询
	utils.CheckGangNotify(t, deskData, xiaPlayerID, bankerPlayerID, 15, room.GangType_MingGang)
	// 摸牌33
	utils.WaitZixunNtf(deskData, 1)
	// 1 号玩家出牌33
	utils.SendChupaiReq(deskData, 1, uint32(33))
	//检查出牌响应
	utils.CheckChuPaiNotify(t, deskData, uint32(33), 1)

	// 对家摸牌33
	duiPlayerID := utils.GetDeskPlayerBySeat(2, deskData).Player.GetID()
	utils.WaitZixunNtf(deskData, 2)
	// 2 号玩家出牌16
	utils.SendChupaiReq(deskData, 2, uint32(16))
	//检查出牌响应
	utils.CheckChuPaiNotify(t, deskData, uint32(16), 2)

	// 下家明杠对家16
	// 1 号玩家收到出牌问询通知， 可以杠
	assert.Nil(t, utils.WaitChupaiWenxunNtf(deskData, 1, true, false, true))
	// 1 号玩家发送杠请求
	assert.Nil(t, utils.SendGangReq(deskData, 1, 16, room.GangType_MingGang))
	// 检测所有玩家收到杠通知,并等待自询
	utils.CheckGangNotify(t, deskData, xiaPlayerID, duiPlayerID, 16, room.GangType_MingGang)
	// 摸牌34
	utils.WaitZixunNtf(deskData, 1)
	// 1 号玩家出牌34
	utils.SendChupaiReq(deskData, 1, uint32(34))
	//检查出牌响应
	utils.CheckChuPaiNotify(t, deskData, uint32(34), 1)

	// 上家摸牌34
	shangPlayerID := utils.GetDeskPlayerBySeat(3, deskData).Player.GetID()
	utils.WaitZixunNtf(deskData, 3)
	// 3 号玩家出牌17
	utils.SendChupaiReq(deskData, 3, uint32(17))
	//检查出牌响应
	utils.CheckChuPaiNotify(t, deskData, uint32(17), 3)

	// 下家明杠上家17
	// 1 号玩家收到出牌问询通知， 可以杠
	assert.Nil(t, utils.WaitChupaiWenxunNtf(deskData, 1, true, false, true))
	// 1 号玩家发送杠请求
	assert.Nil(t, utils.SendGangReq(deskData, 1, 17, room.GangType_MingGang))
	// 检测所有玩家收到杠通知,并等待自询
	utils.CheckGangNotify(t, deskData, xiaPlayerID, shangPlayerID, 17, room.GangType_MingGang)

	// 游戏结束
	utils.WaitGameOverNtf(t, deskData)
}
