package recovertests

import (
	"steve/client_pb/msgid"
	"steve/client_pb/room"
	"steve/simulate/config"
	"steve/simulate/connect"
	"steve/simulate/global"
	"steve/simulate/utils"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Test_DisconnectRecover 退出后再进入恢复牌局
// step1: 开局 0号玩家为庄家
// step2: 0号玩家在收到 自询 后，记录桌面信息， 断开连接
// step3: 3号玩家收到 出牌问询通知 后，0号玩家发送 登录请求，
// step4: 0号玩家收到 收到 登录应答 后发送 恢复牌局请求
// step5: 0号玩家收到 恢复牌局应答 判断数据的正确性
func Test_DisconnectRecover(t *testing.T) {
	params := global.NewCommonStartGameParams()
	params.WallCards = []uint32{31, 31, 31, 31, 32, 32, 32, 32}
	disconnectSeat := params.BankerSeat
	mopaiSeat := (disconnectSeat + 1) % len(params.Cards)
	// step 1
	deskData, err := utils.StartGame(params)
	assert.Nil(t, err)
	assert.NotNil(t, deskData)
	// step 2
	disconnectPlayer := utils.GetDeskPlayerBySeat(disconnectSeat, deskData)
	expector, _ := disconnectPlayer.Expectors[msgid.MsgID_ROOM_ZIXUN_NTF]
	ntf1 := &room.RoomZixunNtf{}
	assert.Nil(t, expector.Recv(global.DefaultWaitMessageTime, ntf1))
	// 出牌，时间好掌握
	assert.Nil(t, utils.SendChupaiReq(deskData, disconnectSeat, params.Cards[disconnectSeat][0]))
	// 发牌后睡眠，在关闭链接前保证发送成功
	time.Sleep(time.Second)
	assert.Nil(t, disconnectPlayer.Player.GetClient().Stop())

	// step 3
	mopaiPlayer := utils.GetDeskPlayerBySeat(mopaiSeat, deskData)
	expector, _ = mopaiPlayer.Expectors[msgid.MsgID_ROOM_CHUPAIWENXUN_NTF]
	ntf2 := &room.RoomChupaiWenxunNtf{}
	assert.Nil(t, expector.Recv(global.DefaultWaitMessageTime, ntf2))
	client := connect.NewTestClient(config.ServerAddr, config.ClientVersion)
	assert.NotNil(t, client)
	player, err := utils.LoginUser(client, disconnectPlayer.Player.GetUsrName())
	assert.Nil(t, err)
	assert.NotNil(t, player)
	assert.Equal(t, disconnectPlayer.Player.GetID(), player.GetID())

	// step 4
	utils.UpdatePlayerClientInfo(client, player, deskData)
	// // 发牌后睡眠，在关闭链接前保证发送成功
	// time.Sleep(time.Second)
	// assert.Nil(t, utils.SendNeedRecoverGameReq(disconnectSeat, deskData))
	// expector, _ = disconnectPlayer.Expectors[msgid.MsgID_ROOM_DESK_NEED_RESUME_RSP]
	// ntf4 := room.RoomDeskNeedReusmeRsp{}
	// assert.Nil(t, expector.Recv(global.DefaultWaitMessageTime, &ntf4))
	// assert.True(t, ntf4.GetIsNeed())
	// assert.Equal(t, 1, ntf4.GetGameId()) // gameid后续修改

	assert.Nil(t, utils.SendRecoverGameReq(disconnectSeat, deskData))
	// step 5
	disconnectPlayer = utils.GetDeskPlayerBySeat(disconnectSeat, deskData)
	expector, _ = disconnectPlayer.Expectors[msgid.MsgID_ROOM_RESUME_GAME_RSP]
	ntf3 := &room.RoomResumeGameRsp{}
	assert.Nil(t, expector.Recv(global.DefaultWaitMessageTime, ntf3))
	assert.Equal(t, room.RoomError_SUCCESS, ntf3.GetResumeRes())
	assert.Equal(t, room.GameStage_GAMESTAGE_PLAYCARD, ntf3.GetGameInfo().GetGameStage())
	// assert.Equal(t, room.GameId_GAMEID_XUELIU, ntf3.GetGameInfo().GetGameId())
}
