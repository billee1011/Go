package matchtests

import (
	"steve/client_pb/common"
	"steve/client_pb/match"
	"steve/client_pb/msgid"
	"steve/client_pb/room"
	"steve/simulate/global"
	"steve/simulate/interfaces"
	"steve/simulate/utils"
	"steve/simulate/utils/doudizhu"
	"testing"
	"time"

	"github.com/Sirupsen/logrus"
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

// 开始斗地主然后结束
// 返回参与游戏的玩家列表
func startDDZAndFinishGame(t *testing.T) []interfaces.ClientPlayer {

	// 配牌1
	params := doudizhu.NewStartDDZGameParamsTest2()

	// 开始游戏
	deskData, err := utils.StartDDZGame(params)
	assert.NotNil(t, deskData)
	assert.Nil(t, err)

	// 叫地主用例1
	doudizhu.JiaodizhuTest1(deskData)

	// 加倍用例1
	doudizhu.JiabeiTest1(deskData)

	// 出牌用例1
	doudizhu.PlaycardTest1(deskData)

	// 所有玩家的clientPlayer
	players := make([]interfaces.ClientPlayer, 0, len(deskData.Players))
	for _, deskPlayer := range deskData.Players {
		players = append(players, deskPlayer.Player)
	}

	return players
}

// Test_ContinueDDZ 测试斗地主续局功能
// 开始一局游戏，等游戏结束后，3个玩家均申请续局
// 4 个玩家均会收到房间创建通知，且每个玩家座位号不变
func Test_ContinueDDZ(t *testing.T) {

	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "continue_test.go::Test_ContinueDDZ",
	})

	players := startDDZAndFinishGame(t)
	time.Sleep(10 * time.Millisecond) // 等待 10ms 确保match服已经接收到room服的续局牌桌的请求

	for _, player := range players {

		// 准备：续局匹配的响应通知，斗地主开始游戏的通知
		player.AddExpectors(msgid.MsgID_MATCH_CONTINUE_RSP, msgid.MsgID_ROOM_DDZ_START_GAME_NTF)

		// 发出游戏续局请求
		player.GetClient().SendPackage(utils.CreateMsgHead(msgid.MsgID_MATCH_CONTINUE_REQ), &match.MatchDeskContinueReq{
			GameId: common.GameId_GAMEID_DOUDIZHU.Enum(),
			Cancel: proto.Bool(false), // false表开始续局
		})

		// 续局期待
		expector := player.GetExpector(msgid.MsgID_MATCH_CONTINUE_RSP)

		matchRsp := match.MatchRsp{}

		// 接收续局回应
		assert.Nil(t, expector.Recv(global.DefaultWaitMessageTime, &matchRsp))

		// 不为0时，说明续局的请求失败了，不再继续
		if matchRsp.GetErrCode() != 0 {
			return
		}
	}

	// 所有玩家收到斗地主游戏开始通知
	for _, player := range players {

		// 斗地主开始游戏的通知
		expector := player.GetExpector(msgid.MsgID_ROOM_DDZ_START_GAME_NTF)

		startGameNotify := room.DDZStartGameNtf{}

		assert.Nil(t, expector.Recv(global.DefaultWaitMessageTime, &startGameNotify))

		logEntry.Infof("斗地主续局后收到了开始游戏的通知，playerID = %v, nextStage = %v", startGameNotify.GetPlayerId(), startGameNotify.GetNextStage())
	}
}

// Test_ContinueCancelDDZ 测试斗地主取消续局
func Test_ContinueCancelDDZ(t *testing.T) {
	logEntry := logrus.WithFields(logrus.Fields{
		"func_name": "continue_test.go::Test_ContinueCancelDDZ",
	})

	players := startDDZAndFinishGame(t)
	time.Sleep(10 * time.Millisecond) // 等待 10ms 确保match服已经接收到room服的续局牌桌的请求

	// 准备：续局牌桌解散的期望
	for i := 1; i < len(players); i++ {
		players[i].AddExpectors(msgid.MsgID_MATCH_CONTINUE_DESK_DIMISS_NTF)
	}

	// 由第一个玩家发出游戏取消续局的请求
	players[0].GetClient().SendPackage(utils.CreateMsgHead(msgid.MsgID_MATCH_CONTINUE_REQ), &match.MatchDeskContinueReq{
		GameId: common.GameId_GAMEID_DOUDIZHU.Enum(),
		Cancel: proto.Bool(true), // true表取消续局
	})

	// 剩余两个玩家收到续局牌桌解散的通知
	for i := 1; i < len(players); i++ {
		expector := players[i].GetExpector(msgid.MsgID_MATCH_CONTINUE_DESK_DIMISS_NTF)

		deskMissNtf := match.MatchContinueDeskDimissNtf{}
		assert.Nil(t, expector.Recv(global.DefaultWaitMessageTime, &deskMissNtf))

		logEntry.Infof("斗地主续局取消后收到了续局牌桌解散的通知,reserve = %v", deskMissNtf.GetReserve())
	}
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
		assert.Nil(t, expector.Recv(global.DefaultWaitMessageTime, nil))
	}
}
