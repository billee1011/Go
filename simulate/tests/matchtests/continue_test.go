package matchtests

import (
	"steve/client_pb/common"
	"steve/client_pb/match"
	"steve/client_pb/msgid"
	"steve/client_pb/room"
	"steve/simulate/global"
	"steve/simulate/interfaces"
	"steve/simulate/utils"
	"testing"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
)

// 开始游戏然后结束
// 返回参与游戏的玩家列表
func startAndFinishGame(t *testing.T) []interfaces.ClientPlayer {
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

	players := make([]interfaces.ClientPlayer, 0, len(deskData.Players))
	for _, deskPlayer := range deskData.Players {
		players = append(players, deskPlayer.Player)
	}
	return players
}

// Test_ContinueMajong 测试麻将续局功能
// 开始一局游戏，等游戏结束后，4 个玩家均申请续局
// 4 个玩家均会收到房间创建通知，且每个玩家座位号不变
func Test_ContinueMajong(t *testing.T) {
	players := startAndFinishGame(t)
	time.Sleep(10 * time.Millisecond) // 等待 10ms 确保匹配服已经接收到续局牌桌

	for _, player := range players {
		player.AddExpectors(msgid.MsgID_MATCH_CONTINUE_RSP, msgid.MsgID_ROOM_START_GAME_NTF)
		player.GetClient().SendPackage(utils.CreateMsgHead(msgid.MsgID_MATCH_CONTINUE_REQ), &match.MatchDeskContinueReq{
			GameId: common.GameId_GAMEID_XUELIU.Enum(),
			Cancel: proto.Bool(false),
		})
		expector := player.GetExpector(msgid.MsgID_MATCH_CONTINUE_RSP)
		assert.Nil(t, expector.Recv(global.DefaultWaitMessageTime, nil))
	}
	// 所有玩家收到游戏开始通知
	for _, player := range players {
		expector := player.GetExpector(msgid.MsgID_ROOM_START_GAME_NTF)
		startGameNotify := room.RoomStartGameNtf{}
		assert.Nil(t, expector.Recv(global.DefaultWaitMessageTime, &startGameNotify))
	}
}

// Test_ContinueCancel 测试取消续局
func Test_ContinueCancel(t *testing.T) {
	players := startAndFinishGame(t)
	time.Sleep(10 * time.Millisecond) // 等待 10ms 确保匹配服已经接收到续局牌桌

	for i := 1; i < len(players); i++ {
		players[i].AddExpectors(msgid.MsgID_MATCH_CONTINUE_DESK_DIMISS_NTF)
	}

	players[0].GetClient().SendPackage(utils.CreateMsgHead(msgid.MsgID_MATCH_CONTINUE_REQ), &match.MatchDeskContinueReq{
		GameId: common.GameId_GAMEID_XUELIU.Enum(),
		Cancel: proto.Bool(true),
	})
	for i := 1; i < len(players); i++ {
		expector := players[i].GetExpector(msgid.MsgID_MATCH_CONTINUE_DESK_DIMISS_NTF)
		expector.Recv(global.DefaultWaitMessageTime, nil)
	}
}
