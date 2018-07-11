package recovertests

import (
	msgid "steve/client_pb/msgId"
	"steve/client_pb/room"
	"steve/simulate/global"
	"steve/simulate/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test_SCXZ_Zimo 自摸测试
// 开始游戏后，庄家出1筒，没有人可以碰杠胡。1 号玩家摸 9W 且可以自摸
// 期望：
// 1. 1号玩家收到自询通知，且可以自摸
// 2. 1号玩家发送胡请求后，所有玩家收到胡通知， 胡牌者为1号玩家，胡类型为自摸，胡的牌为9W
func Test_SCXZ_Zimo_Recover(t *testing.T) {
	var Int1B uint32 = 31
	var Int9W uint32 = 19
	params := global.NewCommonStartGameParams()
	params.BankerSeat = 0
	zimoSeat := 1
	quitSeat := zimoSeat
	bankerSeat := params.BankerSeat

	// 庄家的最后一张牌改为 1B
	params.Cards[bankerSeat][13] = 31
	// 1 号玩家最后1张牌改为 9W
	params.Cards[zimoSeat][12] = 19
	// 墙牌改成 9W 。 墙牌有两张，否则就是海底捞了
	params.WallCards = []uint32{19, 31}

	deskData, err := utils.StartGame(params)
	assert.Nil(t, err)
	assert.NotNil(t, deskData)

	assert.Nil(t, utils.SendChupaiReq(deskData, bankerSeat, Int1B))

	// 1 号玩家收到可自摸通知
	zimoPlayer := utils.GetDeskPlayerBySeat(zimoSeat, deskData)
	expector, _ := zimoPlayer.Expectors[msgid.MsgID_ROOM_ZIXUN_NTF]
	ntf := room.RoomZixunNtf{}
	assert.Nil(t, expector.Recv(global.DefaultWaitMessageTime, &ntf))
	assert.True(t, ntf.GetEnableZimo())

	// 发送胡请求
	assert.Nil(t, utils.SendHuReq(deskData, zimoSeat))

	// 检测所有玩家收到自摸通知
	utils.CheckHuNotify(t, deskData, []int{zimoSeat}, zimoSeat, Int9W, room.HuType_HT_DIHU)

	// 检测所有玩家收到自摸结算通知
	utils.CheckZiMoSettleNotify(t, deskData, []int{zimoSeat}, zimoSeat, Int9W, room.HuType_HT_DIHU)

	assert.Nil(t, utils.SendQuitReq(deskData, quitSeat))
	// 其他玩家收到该玩家退出通知
	utils.RecvQuitNtf(t, deskData, []int{0, 2, 3})
	assert.Nil(t, utils.SendNeedRecoverGameReq(quitSeat, deskData))

	// 需要恢复对局
	expector, _ = zimoPlayer.Expectors[msgid.MsgID_ROOM_DESK_NEED_RESUME_RSP]
	rsp1 := room.RoomDeskNeedReusmeRsp{}
	assert.Nil(t, expector.Recv(global.DefaultWaitMessageTime, &rsp1))
	assert.True(t, rsp1.GetIsNeed())

	// 请求加入失败  新架构匹配服务没有识别 在游戏中
	// rsp2, err := utils.ApplyJoinDesk(zimoPlayer.Player, room.GameId_GAMEID_XUEZHAN)
	// assert.Nil(t, err)
	// assert.Equal(t, room.RoomError_DESK_GAME_PLAYING, rsp2.GetErrCode())

	// 请求恢复对局
	assert.Nil(t, utils.SendRecoverGameReq(quitSeat, deskData))
	expector, _ = zimoPlayer.Expectors[msgid.MsgID_ROOM_RESUME_GAME_RSP]
	rsp3 := &room.RoomResumeGameRsp{}
	assert.Nil(t, expector.Recv(global.DefaultWaitMessageTime, rsp3))
	assert.Equal(t, room.RoomError_SUCCESS, rsp3.GetResumeRes())
	assert.Equal(t, room.GameStage_GAMESTAGE_PLAYCARD, rsp3.GetGameInfo().GetGameStage())

	var player *room.GamePlayerInfo
	for _, player = range rsp3.GetGameInfo().GetPlayers() {
		if player.GetPlayerInfo().GetSeat() == uint32(zimoSeat) {
			break
		}
	}
	assert.False(t, player.GetIsTuoguan())
	assert.Equal(t, room.XingPaiState_XP_STATE_HU, player.GetXpState())
}
