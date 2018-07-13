package recovertests

import (
	"steve/client_pb/common"
	msgid "steve/client_pb/msgid"
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
	params.GameID = room.GameId_GAMEID_XUEZHAN // 血战
	params.PeiPaiGame = "scxz"
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

	quitRspExpector := zimoPlayer.Expectors[msgid.MsgID_ROOM_DESK_QUIT_RSP]
	assert.Nil(t, utils.SendQuitReq(deskData, quitSeat))
	// 其他玩家收到该玩家退出通知
	utils.RecvQuitNtf(t, deskData, []int{0, 2, 3})
	quitResponse := room.RoomDeskQuitRsp{}
	assert.Nil(t, quitRspExpector.Recv(global.DefaultWaitMessageTime, &quitResponse))
	assert.Equal(t, room.RoomError_SUCCESS, quitResponse.GetErrCode())

	playerState, err := utils.GetDeskPlayerState(zimoPlayer.Player)
	assert.Nil(t, err)
	assert.Equal(t, playerState, common.PlayerState_PS_IDLE)

	// assert.Nil(t, utils.SendNeedRecoverGameReq(quitSeat, deskData))
	// expector, _ = zimoPlayer.Expectors[msgid.MsgID_ROOM_DESK_NEED_RESUME_RSP]
	// rsp1 := room.RoomDeskNeedReusmeRsp{}
	// assert.Nil(t, expector.Recv(global.DefaultWaitMessageTime, &rsp1))
	// assert.False(t, rsp1.GetIsNeed())

	utils.ApplyJoinDesk(zimoPlayer.Player, room.GameId_GAMEID_XUELIU)

	// 再加入3个玩家凑够4人开局避免影响其他测试用例
	newPlayers, err := utils.CreateAndLoginUsers(3)
	assert.Nil(t, err)
	err = utils.ApplyJoinDeskPlayers(newPlayers, room.GameId_GAMEID_XUELIU)
	assert.Nil(t, err)
	expector, _ = zimoPlayer.Expectors[msgid.MsgID_ROOM_DESK_CREATED_NTF]
	ntf1 := room.RoomDeskCreatedNtf{}
	assert.Nil(t, expector.Recv(global.DefaultWaitMessageTime, &ntf1))
}
