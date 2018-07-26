package recovertests

import (
	"steve/client_pb/msgid"
	"steve/client_pb/room"
	"steve/simulate/global"
	"steve/simulate/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test_QuitRecover 退出后再进入恢复牌局
// step1: 开局 0号玩家为庄家
// step2: 0号玩家在收到 自询 后，记录桌面信息，发送 退出游戏请求
// step3: 3号玩家收到 自询通知 后，0号玩家发送 加入桌面请求
// step4: 0号玩家收到 加入桌面请求应答：有正在进行的游戏，发送 恢复牌局请求
// step5: 0号玩家收到 恢复牌局应答 判断数据的正确性
func Test_QuitRecover(t *testing.T) {
	params := global.NewCommonStartGameParams()
	params.WallCards = []uint32{31, 31, 31, 31, 32, 32, 32, 32}
	quitSeat := params.BankerSeat
	mopaiSeat := (quitSeat + 1) % len(params.Cards)
	// step 1
	deskData, err := utils.StartGame(params)
	assert.Nil(t, err)
	assert.NotNil(t, deskData)
	// step 2
	quitPlayer := utils.GetDeskPlayerBySeat(quitSeat, deskData)
	expector, _ := quitPlayer.Expectors[msgid.MsgID_ROOM_ZIXUN_NTF]
	ntf1 := &room.RoomZixunNtf{}
	assert.Nil(t, expector.Recv(global.DefaultWaitMessageTime, ntf1))
	assert.Nil(t, utils.SendQuitReq(deskData, quitSeat))

	// step 3、4
	mopaiPlayer := utils.GetDeskPlayerBySeat(mopaiSeat, deskData)
	expector, _ = mopaiPlayer.Expectors[msgid.MsgID_ROOM_ZIXUN_NTF]
	ntf2 := &room.RoomZixunNtf{}
	assert.Nil(t, expector.Recv(global.DefaultWaitMessageTime, ntf2))

	// TODO 匹配功能还未完善，先不测这个
	// rsp, err := utils.ApplyJoinDesk(quitPlayer.Player, common.GameId_GAMEID_XUELIU)
	// assert.Nil(t, err)
	// assert.Equal(t, room.RoomError_DESK_GAME_PLAYING, rsp.GetErrCode())
	assert.Nil(t, utils.SendRecoverGameReq(quitSeat, deskData))

	// step 5
	expector, _ = quitPlayer.Expectors[msgid.MsgID_ROOM_RESUME_GAME_RSP]
	ntf3 := &room.RoomResumeGameRsp{}
	assert.Nil(t, expector.Recv(global.DefaultWaitMessageTime, ntf3))
	assert.Equal(t, room.RoomError_SUCCESS, ntf3.GetResumeRes())
	assert.Equal(t, room.GameStage_GAMESTAGE_PLAYCARD, ntf3.GetGameInfo().GetGameStage())
	var player *room.GamePlayerInfo
	for _, player = range ntf3.GetGameInfo().GetPlayers() {
		if player.GetPlayerInfo().GetSeat() == uint32(quitSeat) {
			break
		}
	}
	assert.True(t, player.GetIsTuoguan())
	assert.Equal(t, room.XingPaiState_XP_STATE_NORMAL, player.GetXpState())
	assert.Equal(t, room.GameId_GAMEID_XUELIU, ntf3.GetGameInfo().GetGameId())
}
