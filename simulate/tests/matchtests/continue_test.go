package matchtests

import (
	"steve/client_pb/common"
	"steve/client_pb/match"
	"steve/client_pb/msgid"
	"steve/client_pb/room"
	"steve/simulate/global"
	"steve/simulate/utils"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Test_ContinueMajong 测试麻将续局功能
// 开始一局游戏，等游戏结束后，4 个玩家均申请续局
// 4 个玩家均会收到房间创建通知，且每个玩家座位号不变
func Test_ContinueMajong(t *testing.T) {
	params := global.NewCommonStartGameParams()
	params.WallCards = []uint32{}
	deskData, err := utils.StartGame(params)
	assert.NotNil(t, deskData)
	assert.Nil(t, err)

	// 庄家出 3万，游戏结束
	utils.SendChupaiReq(deskData, 0, 13)

	// 等待游戏结束通知和总结算通知
	bankerPlayer := utils.GetDeskPlayerBySeat(0, deskData)
	gameOverNtfExpector := bankerPlayer.Expectors[msgid.MsgID_ROOM_GAMEOVER_NTF]
	settleNtfExpector := bankerPlayer.Expectors[msgid.MsgID_ROOM_ROUND_SETTLE]
	assert.Nil(t, gameOverNtfExpector.Recv(global.DefaultWaitMessageTime, nil))
	assert.Nil(t, settleNtfExpector.Recv(global.DefaultWaitMessageTime, nil))

	time.Sleep(100 * time.Millisecond) // 等待 100ms 再申请续局

	for _, deskPlayer := range deskData.Players {
		deskPlayer.Player.AddExpectors(msgid.MsgID_MATCH_CONTINUE_RSP, msgid.MsgID_ROOM_START_GAME_NTF)
		deskPlayer.Player.GetClient().SendPackage(utils.CreateMsgHead(msgid.MsgID_MATCH_CONTINUE_REQ), &match.MatchDeskContinueReq{
			GameId: common.GameId_GAMEID_XUELIU.Enum(),
		})
		expector := deskPlayer.Player.GetExpector(msgid.MsgID_MATCH_CONTINUE_RSP)
		assert.Nil(t, expector.Recv(global.DefaultWaitMessageTime, nil))
	}
	// 所有玩家收到游戏开始通知
	for _, deskPlayer := range deskData.Players {
		expector := deskPlayer.Player.GetExpector(msgid.MsgID_ROOM_START_GAME_NTF)
		startGameNotify := room.RoomStartGameNtf{}
		assert.Nil(t, expector.Recv(global.DefaultWaitMessageTime, &startGameNotify))
	}
}
