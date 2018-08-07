package matchtests

import (
	"steve/client_pb/common"
	"steve/client_pb/hall"
	"steve/client_pb/msgid"
	"steve/simulate/global"
	"steve/simulate/interfaces"
	"steve/simulate/utils"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func assertPlayerState(t *testing.T, player interfaces.ClientPlayer, stateID common.PlayerState, gameID common.GameId) {
	player.AddExpectors(msgid.MsgID_HALL_GET_PLAYER_STATE_RSP)

	client := player.GetClient()
	client.SendPackage(utils.CreateMsgHead(msgid.MsgID_HALL_GET_PLAYER_STATE_REQ), &hall.HallGetPlayerStateReq{})

	expector := player.GetExpector(msgid.MsgID_HALL_GET_PLAYER_STATE_RSP)
	response := hall.HallGetPlayerStateRsp{}
	expector.Recv(global.DefaultWaitMessageTime, &response)
	assert.Equal(t, response.GetPlayerState(), stateID)
	if stateID == common.PlayerState_PS_GAMEING || stateID == common.PlayerState_PS_MATCHING {
		assert.Equal(t, gameID, response.GetGameId())
	}
}

// Test_PlayerStates 测试玩家空闲状态
func Test_PlayerStates_Idle(t *testing.T) {
	player, err := utils.LoginNewPlayer()
	assert.Nil(t, err)
	assert.NotNil(t, player)

	assertPlayerState(t, player, common.PlayerState_PS_IDLE, common.GameId_GAMEID_XUEZHAN)
}

// Test_PlayerStates_Gameing 测试玩家在游戏中的状态
func Test_PlayerStates_Gameing(t *testing.T) {
	params := global.NewCommonStartGameParams()
	deskData, err := utils.StartGame(params)
	assert.Nil(t, err)

	deskPlayer := utils.GetDeskPlayerBySeat(0, deskData)
	assertPlayerState(t, deskPlayer.Player, common.PlayerState_PS_GAMEING, common.GameId(params.GameID))
}

// Test_PlayerStates_Matching 测试玩家在匹配中的状态
func Test_PlayerStates_Matching(t *testing.T) {
	player, err := utils.LoginNewPlayer()
	assert.Nil(t, err)
	assert.NotNil(t, player)

	_, err = utils.ApplyJoinDesk(player, common.GameId_GAMEID_XUELIU)
	assert.Nil(t, err)

	time.Sleep(10 * time.Millisecond)
	assertPlayerState(t, player, common.PlayerState_PS_MATCHING, common.GameId_GAMEID_XUELIU)
	// 再次登录玩家，取消匹配
	utils.LoginPlayer(player.GetAccountID(), "")
}
